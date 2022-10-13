package java

import (
	"strings"

	"github.com/smacker/go-tree-sitter/java"

	"github.com/bearer/curio/pkg/detectors/java/datatype"
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/interfacedetector"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/report/variables"
	"github.com/bearer/curio/pkg/util/file"
)

var (
	language = java.GetLanguage()

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(method_invocation
			object: (identifier) @object
			name: (identifier) @method
			arguments: (argument_list . (string_literal) @key)) @node
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
	if file.Language != "Java" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	if err := annotate(tree); err != nil {
		return false, err
	}

	datatype.Discover(report, tree, detector.idGenerator)

	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		DetectorType:     detectors.DetectorJava,
		AcceptExpression: acceptExpression,
		PathAllowed:      false,
	}); err != nil {
		return false, err
	}

	return true, nil
}

func annotate(tree *parser.Tree) error {
	if err := annotateEnvironmentVariables(tree); err != nil {
		return err
	}

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return
		case "object_creation_expression":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				return nil
			}, func(child *parser.Node) error {
				value.Append(child.Value())

				return nil
			})

			return
		case "string_literal":
			value.AppendString(stripQuotes(node.Content()))

			return
		case "program", "block", "class_declaration", "class_body", "variable_declarator", "local_variable_declaration", "type_identifier", "method_declaration", "throws":
			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree) error {
	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		object := captures["object"].Content()
		method := captures["method"].Content()
		if object != "System" || method != "getenv" {
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
	return strings.Trim(value, `"`)
}

func acceptExpression(node *parser.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		// something["ignored.domain"]
		if parent.Type() == "array_access" {
			return false
		}
	}

	return true
}
