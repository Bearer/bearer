package golang

import (
	"strings"

	"github.com/go-enry/go-enry/v2/regex"
	"github.com/smacker/go-tree-sitter/golang"

	"github.com/bearer/curio/pkg/detectors/golang/datatype"
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/golang_util"
	"github.com/bearer/curio/pkg/parser/interfacedetector"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/report/variables"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/normalize_key"
	"github.com/bearer/curio/pkg/util/stringutil"
)

var (
	language = golang.GetLanguage()

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(call_expression
			function: [(selector_expression) (identifier)] @function
			arguments: (argument_list
				[(interpreted_string_literal) (raw_string_literal)] @key)) @detection
	`)

	formatStringPattern    = regex.MustCompile(`%([^%]*?[a-zA-Z]|%)`)
	envFunctionNamePattern = regex.MustCompile(`\benv\b`)
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
	if file.Language != "Go" {
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
		DetectorType:     detectors.DetectorGo,
		PathAllowed:      true,
	}); err != nil {
		return false, err
	}

	return true, nil
}

func annotate(tree *parser.Tree) error {
	imports, err := golang_util.GetImports(tree)
	if err != nil {
		return err
	}

	if err := annotateEnvironmentVariables(tree, imports); err != nil {
		return err
	}

	fmtAliases := golang_util.AliasesFor(imports, "fmt")

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "call_expression":
			function := node.ChildByFieldName("function")
			if function.ChildByFieldName("operand") != nil &&
				stringutil.SliceContains(fmtAliases, function.ChildByFieldName("operand").Content()) &&
				function.ChildByFieldName("field").Content() == "Sprintf" {

				populateFormatStringValue(value, node.ChildByFieldName("arguments"))

				return
			}
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return

		case "field_identifier":
			if node.Parent().Type() == "selector_expression" {
				value.AppendVariableReference(variables.VariableName, node.Content())
			}

		case "interpreted_string_literal", "raw_string_literal":
			node.EachPart(func(text string) error {
				value.AppendString(text)

				return nil
			}, func(child *parser.Node) error {
				return nil
			})

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func populateFormatStringValue(value *values.Value, arguments *parser.Node) {
	if arguments.ChildCount() == 0 {
		return
	}

	argumentOffset := 1
	for _, part := range arguments.Child(0).Value().Parts {
		if stringPart, ok := part.(*values.String); ok {
			start := 0
			partValue := stringPart.Value
			matches := formatStringPattern.FindAllStringSubmatchIndex(stringPart.Value, -1)

			for _, match := range matches {
				value.AppendString(partValue[start:match[0]])

				if partValue[match[0]:match[1]] == "%%" {
					value.AppendString("%")
				} else {
					if argument := arguments.Child(argumentOffset); argument != nil {
						value.Append(argument.Value())
					}
					argumentOffset += 1
				}

				start = match[1]
			}

			value.AppendString(partValue[start:])
		} else {
			value.AppendPart(part)
		}
	}
}

func annotateEnvironmentVariables(tree *parser.Tree, imports map[string]string) error {
	osAliases := golang_util.AliasesFor(imports, "os")
	if len(osAliases) == 0 {
		return nil
	}

	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		if !isEnvironmentVarCall(osAliases, captures["function"].Content()) {
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

func isEnvironmentVarCall(osAliases []string, functionName string) bool {
	for _, alias := range osAliases {
		if functionName == alias+".Getenv" ||
			functionName == alias+".LookupEnv" {
			return true
		}
	}

	return envFunctionNamePattern.MatchString(normalize_key.Normalize(functionName))
}

func stripQuotes(value string) string {
	return strings.Trim(value, "`\"")
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "keyed_element":
			if parent.Child(0).Equal(lastNode) {
				// map[string]string{ "ignored.domain": "..." }
				return false
			}
		case "index_expression":
			// something["ignored.domain"]
			return false
		case "import_declaration":
			// import "ignored.domain"
			return false
		}

		lastNode = parent
	}

	return true
}
