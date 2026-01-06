package ruby

import (
	"regexp"
	"strings"

	"github.com/bearer/bearer/pkg/detectors/ruby/datatype"
	"github.com/bearer/bearer/pkg/util/file"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/interfacedetector"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/report/variables"
	"github.com/bearer/bearer/pkg/util/regex"

	"github.com/smacker/go-tree-sitter/ruby"
)

var (
	language = ruby.GetLanguage()

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(element_reference
			object: (constant) @object
			(string) @key) @node

		(call
			receiver: (constant) @object
			method: (identifier) @method
			arguments: (argument_list . (string) @key)) @node
	`)

	ignoredFilenames = []*regexp.Regexp{
		regexp.MustCompile(`(^|/)Gemfile$`),
		regexp.MustCompile(`(^|/)db/migrate/`),
		regexp.MustCompile(`\.rake$`),
	}
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
	if file.Language != "Ruby" {
		return false, nil
	}

	if regex.AnyMatch(ignoredFilenames, file.RelativePath) {
		return true, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	if err := Annotate(tree); err != nil {
		return false, err
	}

	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorRuby,
		PathAllowed:      false,
	}); err != nil {
		return false, err
	}

	datatypes := datatype.Discover(tree.RootNode(), detector.idGenerator)
	report.AddDataType(detections.TypeSchema, detectors.DetectorRuby, detector.idGenerator, datatypes, nil)

	return true, nil
}

func Annotate(tree *parser.Tree) error {
	if err := annotateEnvironmentVariables(tree); err != nil {
		return err
	}

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "interpolation":
			// FirstChild() uses NamedChild(0) internally, getting the actual expression
			if firstChild := node.FirstChild(); firstChild != nil {
				value.Append(firstChild.Value())
			}

			return
		case "binary":
			if unnamed := node.FirstUnnamedChild(); unnamed != nil && unnamed.Content() == "+" {
				if left := node.ChildByFieldName("left"); left != nil {
					value.Append(left.Value())
				}
				if right := node.ChildByFieldName("right"); right != nil {
					value.Append(right.Value())
				}

				return
			}
		case "identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return
		case "instance_variable":
			value.AppendVariableReference(variables.VariableName, stripVariableName(node.Content()))

			return
		case "string":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				return nil
			}, func(child *parser.Node) error {
				value.Append(child.Value())

				return nil
			})

			return
		case "string_content":
			value.AppendString(node.Content())

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree) error {
	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		node := captures["node"]
		key := stripQuotes(captures["key"].Content())
		object := captures["object"].Content()
		methodNode := captures["method"]

		if object != "ENV" || (methodNode != nil && methodNode.Content() != "fetch") {
			return nil
		}

		value := values.New()
		value.AppendVariableReference(variables.VariableEnvironment, key)
		node.SetValue(value)

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `'"`)
}

func stripVariableName(value string) string {
	return strings.Trim(value, `@`)
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {

		case "pair":
			if parent.ChildByFieldName("key").Equal(lastNode) {
				// { "ignored.domain" => "..." }
				return false
			}

		case "element_reference":
			// something["ignored.domain"]
			return false
		}

		lastNode = parent
	}

	return true
}
