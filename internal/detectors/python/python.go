package python

import (
	"strings"

	"github.com/smacker/go-tree-sitter/python"

	"github.com/bearer/bearer/internal/detectors/python/datatype"
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

const environItem = "environ"

var (
	language = python.GetLanguage()

	importsQuery = parser.QueryMustCompile(language, `
		(import_statement name: (dotted_name) @module)

		(import_statement name: (aliased_import
			name: (dotted_name) @module
			alias: (identifier) @alias))

		(import_from_statement
			module_name: (dotted_name) @module
			name: (dotted_name) @item)
	`)

	environmentVariableQuery = parser.QueryMustCompile(language, `
		(subscript value: (attribute) @item subscript: (string) @key) @detection

		(call
			function: (attribute) @item
			arguments: (argument_list . (string) @key)) @detection
	`)
)

type detector struct {
	idGenerator nodeid.Generator
}
type pyImport struct {
	module string
	alias  string
	item   string
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
	if file.Language != "Python" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	datatype.Discover(report, tree, detector.idGenerator)

	return processTree(tree, report)
}

func ProcessRaw(file *file.FileInfo, report report.Report, input []byte, offset int, idGenerator nodeid.Generator) (bool, error) {
	tree, err := parser.ParseBytes(file, file.Path, input, language, offset)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	datatype.Discover(report, tree, idGenerator)

	return processTree(tree, report)
}

func processTree(tree *parser.Tree, report report.Report) (bool, error) {
	if err := annotate(tree); err != nil {
		return false, err
	}
	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorPython,
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
		case "binary_operator":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "string":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				value.AppendString(text)

				return nil
			}, func(child *parser.Node) error {
				value.Append(child.Value())

				return nil
			})

			return
		case "concatenated_string":
			for i := 0; i < node.ChildCount(); i++ {
				value.Append(node.Child(i).Value())
			}
		case "interpolation":
			value.Append(node.FirstChild().Value())

			return
		case "identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree) error {
	imports, err := getImports(tree)
	if err != nil {
		return err
	}

	environImports := getEnvironImports(imports)
	if len(environImports) == 0 {
		return nil
	}

	return tree.Query(environmentVariableQuery, func(captures parser.Captures) error {
		if !isEnvironmentVarCall(environImports, captures["item"].Content()) {
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

func isEnvironmentVarCall(environImports []string, item string) bool {
	for _, envImport := range environImports {
		if item == envImport || item == envImport+".get" || item == envImport+".pop" {
			return true
		}
	}

	return false
}

func stripQuotes(value string) string {
	return strings.Trim(value, `'"`)
}

func getImports(tree *parser.Tree) ([]pyImport, error) {
	var imports []pyImport
	err := tree.Query(importsQuery, func(captures parser.Captures) error {
		alias := ""
		if aliasNode := captures["alias"]; aliasNode != nil {
			alias = aliasNode.Content()
		}

		item := ""
		if itemNode := captures["item"]; itemNode != nil {
			item = itemNode.Content()
		}

		imports = append(imports, pyImport{
			module: captures["module"].Content(),
			alias:  alias,
			item:   item,
		})

		return nil
	})

	return imports, err
}

func getEnvironImports(imports []pyImport) []string {
	var result []string
	for _, pyImport := range imports {
		if pyImport.module != "os" {
			continue
		}

		switch {
		case pyImport.item == environItem:
			result = append(result, environItem)
		case pyImport.alias != "":
			result = append(result, pyImport.alias+"."+environItem)
		default:
			result = append(result, "os."+environItem)
		}
	}

	return result
}

func acceptExpression(node *parser.Node) bool {
	if node.Type() == "string" {
		quotes := node.FirstUnnamedChild().Content()
		if quotes == `'''` || quotes == `"""` {
			return false
		}
	}

	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "decorator":
			// @something("ignored.domain")
			return false
		case "pair":
			if parent.ChildByFieldName("key").Equal(lastNode) {
				// { "ignored.domain" => "..." }
				return false
			}
		case "subscript":
			// something["ignored.domain"]
			return false
		}

		lastNode = parent
	}

	return true
}
