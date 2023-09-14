package javascript

import (
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/smacker/go-tree-sitter/javascript"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/detectors/javascript/datatype"
	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/interfacedetector"
	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/values"
	"github.com/bearer/bearer/internal/report/variables"
	"github.com/bearer/bearer/internal/util/file"
)

var (
	environmentVariableQuery = parser.QueryMustCompile(javascript.GetLanguage(), `
	(member_expression
		object: (member_expression) @object
		property: (property_identifier) @key) @node

	(subscript_expression
		object: (member_expression) @object
		index: (string) @key) @node

	(variable_declarator
		name: (object_pattern (shorthand_property_identifier_pattern) @key @node)
			value: (member_expression) @object)
`)
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

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Language != "JavaScript" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, javascript.GetLanguage())
	if err != nil {
		return false, err
	}
	defer tree.Close()

	return detector.processTree(tree, report)
}

func ProcessRaw(
	raw_string string,
	report report.Report,
	f *file.FileInfo,
	offset int,
	idGenerator nodeid.Generator,
) (bool, error) {
	fileInfo := &file.FileInfo{
		LanguageType: enry.Programming,
		Language:     "JavaScript",
	}

	tree, err := parser.ParseBytes(fileInfo, f.Path, []byte(raw_string), javascript.GetLanguage(), offset)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	jsDetector := New(idGenerator).(*detector)
	_, err = jsDetector.processTree(tree, report)

	return true, err
}

func (detector *detector) processTree(tree *parser.Tree, report report.Report) (bool, error) {
	if err := annotate(tree, environmentVariableQuery); err != nil {
		return false, err
	}

	datatype.Discover(report, tree, detector.idGenerator)

	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorJavascript,
		PathAllowed:      false,
	}); err != nil {
		return false, err
	}

	return true, nil
}

func annotate(tree *parser.Tree, environmentVariablesQuery *sitter.Query) error {
	if err := annotateEnvironmentVariables(tree, environmentVariablesQuery); err != nil {
		return err
	}

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "template_substitution":
			value.Append(node.FirstChild().Value())

			return
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "identifier", "property_identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return

		case "string", "template_string":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				value.AppendString(text)

				return nil
			}, func(child *parser.Node) error {
				value.Append(child.Value())

				return nil
			})

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree, query *sitter.Query) error {
	return tree.Query(query, func(captures parser.Captures) error {
		if captures["object"].Content() != "process.ENV" {
			return nil
		}

		node := captures["node"]
		keyNode := captures["key"]
		key := stripQuotes(keyNode.Content())

		value := values.New()
		value.AppendVariableReference(variables.VariableEnvironment, key)
		node.SetValue(value)

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `'"`)
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "decorator":
			// @MyDecorator("ignored")
			return false
		case "pair":
			if parent.ChildByFieldName("key").Equal(lastNode) {
				// { 'ignored.domain': '...' }
				return false
			}
		case "import_statement":
			// import * from 'ignored'
			return false
		case "subscript_expression":
			// something['ignored.domain']
			return false
		case "jsx_element", "jsx_self_closing_element":
			// <img src="ignored.domain"/>
			return false
		}

		lastNode = parent
	}

	return true
}
