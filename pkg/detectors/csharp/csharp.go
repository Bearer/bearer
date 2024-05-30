package csharp

import (
	"regexp"

	"github.com/smacker/go-tree-sitter/csharp"

	"github.com/bearer/bearer/pkg/detectors/csharp/datatype"
	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/interfacedetector"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/report/variables"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

var (
	language = csharp.GetLanguage()

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(invocation_expression
			(member_access_expression) @function
			(argument_list . (argument .
				[(string_literal) (verbatim_string_literal)] @key))) @detection
	`)

	trimQuotesPattern = regexp.MustCompile(`@?"(.*)"`)

	ignoredBinaryOperators = []string{"!=", "=="}
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
	if file.Language != "C#" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	datatype.Discover(report, tree, detector.idGenerator)

	if err := annotate(tree); err != nil {
		return false, err
	}
	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorCSharp,
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
		case "interpolated_string_text", "interpolated_verbatim_string_text":
			value.AppendString(node.Content())

			return
		case "interpolation":
			value.Append(node.FirstChild().Value())

			return
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "interpolated_string_expression", "string_literal", "verbatim_string_literal":
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

func annotateEnvironmentVariables(tree *parser.Tree) error {
	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		function := captures["function"].Content()
		if function != "System.Environment.GetEnvironmentVariable" && function != "Environment.GetEnvironmentVariable" {
			return nil
		}

		detectionNode := captures["detection"]
		key := stripQuotes(captures["key"].Content())

		value := values.New()
		value.AppendVariableReference(variables.VariableEnvironment, key)
		detectionNode.SetValue(value)

		return nil
	})
}

func stripQuotes(value string) string {
	return trimQuotesPattern.ReplaceAllString(value, "$1")
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "initializer_expression":
			if grandParent := parent.Parent(); grandParent != nil &&
				grandParent.Type() == "initializer_expression" &&
				parent.Child(0).Equal(lastNode) {

				// This probably matches other things too but we're aiming for:
				// 	 new Dictionary<...>{ {"ignored.domain", "..."}, ... }
				return false
			}
		case "element_binding_expression":
			// { ["ignored.domain"] = "..." }
			return false
		case "element_access_expression":
			// something["ignored.domain"]
			return false
		case "binary_expression":
			// Ignore certain boolean expressions. eg `"a" == "b"`
			if stringutil.SliceContains(ignoredBinaryOperators, parent.FirstUnnamedChild().Content()) {
				return false
			}
		}

		lastNode = parent
	}

	return true
}
