package php

import (
	"os"
	"regexp"

	"github.com/bearer/bearer/internal/detectors/html"
	"github.com/bearer/bearer/internal/detectors/php/context"
	"github.com/bearer/bearer/internal/detectors/php/datatype"
	"github.com/bearer/bearer/internal/parser/interfacedetector"
	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/parser/schema"
	php "github.com/bearer/bearer/internal/parser/sitter/php2"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	reporttypes "github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/values"
	"github.com/bearer/bearer/internal/report/variables"
)

var (
	language = php.GetLanguage()

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(subscript_expression (variable_name) @variable . [(encapsed_string) (string)] @key) @node
	`)
	queryText = parser.QueryMustCompile(language, `
		(text) @param_text
	`)

	haltCompiler = regexp.MustCompile(`(?i)__HALT_COMPILER\(\)`)
)

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report reporttypes.Report) (bool, error) {
	if file.Language != "PHP" {
		return false, nil
	}

	fileContent, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return false, nil
	}
	if haltCompiler.Match(fileContent) {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	datatype.Discover(report, tree, detector.idGenerator)

	schemaFinder := schema.New(tree, annotateSchemaVariables)
	schemaFinder.Annotate()
	tree.SetValues(schemaFinder.ToVariableValues())

	contextResolver := context.FindContext(tree)

	tree.SetValues(make(map[*sitter.Node]*values.Value))
	if err := annotate(tree); err != nil {
		return false, err
	}

	if err := parseEmbedHtml(file, report, tree, detector.idGenerator); err != nil {
		return false, err
	}

	return true, interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorPHP,
		PathAllowed:      false,
		ContextResolver:  contextResolver,
	})

}

func parseEmbedHtml(file *file.FileInfo, report reporttypes.Report, tree *parser.Tree, idGenerator nodeid.Generator) error {
	captures := tree.QueryMustPass(queryText)

	for _, capture := range captures {
		_, err := html.ProcessRaw(file, report, []byte(capture["param_text"].Content()), capture["param_text"].StartLineNumber(), idGenerator)

		if err != nil {
			return err
		}
	}

	return nil
}

func annotate(tree *parser.Tree) error {
	if err := annotateEnvironmentVariables(tree); err != nil {
		return err
	}

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "." {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}

		case "encapsed_string":
			node.EachPart( //nolint:all,errcheck
				func(text string) error {
					value.AppendString(stringutil.StripQuotes(text))

					return nil
				},
				func(child *parser.Node) error {
					value.Append(child.Value())

					return nil
				})

			return

		case "string":
			nodeContent := stringutil.StripQuotes(node.Content())
			value.AppendString(nodeContent)

			return

		case "variable_name":
			value.AppendVariableReference(variables.VariableName, node.FirstChild().Content())

			return

		case "name":
			if node.Parent().Type() == "member_access_expression" {
				value.AppendVariableReference(variables.VariableName, stringutil.StripQuotes(node.Content()))
			}

			return

		case "dynamic_variable_name":
			if node.Value() != nil {
				value.Append(node.Value())
			}

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree) error {
	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		if captures["variable"].Content() != "$_ENV" {
			return nil
		}

		node := captures["node"]
		key := stringutil.StripQuotes(captures["key"].Content())

		value := values.New()
		value.AppendVariableReference(variables.VariableEnvironment, key)
		node.SetValue(value)

		return nil
	})
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "array_element_initializer":
			if parent.Child(0).Equal(lastNode) {
				// array( "ignored.domain" => "..." )
				return false
			}
		case "subscript_expression":
			// something["ignored.domain"]
			return false
		case "require_once_expression":
			return false
		case "require_expression":
			return false
		}

		lastNode = parent
	}

	return true
}

func annotateSchemaVariables(finder *schema.Finder, node *parser.Node, value *schema.Node) {
	value.Terminating = true

	content := schema.Variable(stringutil.StripQuotes(node.Content()))
	switch node.Type() {
	case "name":
		if node.Parent().Type() == "member_access_expression" {
			value.Variables = append(value.Variables, &content)
			value.Terminating = false
			return
		}

		if node.Parent().Type() == "member_call_expression" {
			value.Variables = append(value.Variables, &content)
			value.Terminating = false
			return
		}

		if node.Parent().Type() == "class_constant_access_expression" {
			value.Variables = append(value.Variables, &content)
			value.Terminating = false
			return
		}

	case "variable_name":
		value.Variables = append(value.Variables, &content)
		value.Terminating = false
		return
	case "member_call_expression":
		if node.Parent().Type() == "member_access_expression" || node.Parent().Type() == "class_constant_access_expression" {
			value.Terminating = false
		}
	case "member_access_expression":
		value.Terminating = false
	case "class_constant_access_expression":
		value.Terminating = false
	}

	parts := finder.NonTerminatingValues(node)
	value.Variables = append(value.Variables, parts...)
}
