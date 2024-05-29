package knex

import (
	"sort"

	"github.com/bearer/bearer/pkg/detectors/javascript/util"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	reportknex "github.com/bearer/bearer/pkg/report/frameworks/knex"
	sitter "github.com/smacker/go-tree-sitter"
)

type FunctionType struct {
	node         *parser.Node
	terminating  bool
	containsKnex bool
	types        []string
}

func detectFunctionTypes(report report.Report, tree *parser.Tree, language *sitter.Language, knexImports []util.Import) {
	knexFunctionTypes := make(map[parser.NodeID]FunctionType)
	helperNodes := make(map[parser.NodeID]FunctionType)
	tree.WalkBottomUp(func(node *parser.Node) error { //nolint:all,errcheck
		functionType := FunctionType{
			containsKnex: false,
			terminating:  false,
			node:         node,
		}

		// knex detection
		if node.Type() == "identifier" {
			functionType.containsKnex = isTypeKnexUsed(node, knexImports)
		}

		// child property push up
		for i := 0; i < node.ChildCount(); i++ {
			child := node.Child(i)

			if helperNodes[child.ID()].containsKnex && !helperNodes[child.ID()].terminating {
				functionType.containsKnex = true
			}

			if !helperNodes[child.ID()].terminating {
				functionType.types = append(functionType.types, helperNodes[child.ID()].types...)
			}
		}
		// add types
		if node.Type() == "type_identifier" {
			functionType.types = append(functionType.types, node.Content())
		}

		// add terminating
		nodeType := node.Type()
		if nodeType == "expression_statement" || nodeType == "lexical_declaration" || nodeType == "statement_block" {
			functionType.terminating = true
		}

		if functionType.containsKnex && functionType.terminating && len(functionType.types) > 0 {
			knexFunctionTypes[node.ID()] = functionType
		}

		helperNodes[node.ID()] = functionType

		return nil
	})

	// sort by line number
	var sortedTypes []FunctionType

	for _, functionType := range knexFunctionTypes {
		sortedTypes = append(sortedTypes, functionType)
	}

	sort.Slice(sortedTypes, func(i, j int) bool {
		lineNumberA := sortedTypes[i].node.Source(false).StartLineNumber
		lineNumberB := sortedTypes[j].node.Source(false).StartLineNumber
		return *lineNumberA < *lineNumberB
	})

	for _, functionType := range sortedTypes {
		for _, dataType := range functionType.types {
			report.AddFramework(detectors.DetectorTypescript, reportknex.TypeFunction, reportknex.Function{DataType: dataType}, functionType.node.Source(false))
		}
	}
}

const knexModule = "knex"

func isTypeKnexUsed(node *parser.Node, imports []util.Import) bool {
	if node.Type() == "identifier" {
		parent := node.Parent()

		if parent == nil {
			return false
		}

		if parent.Type() != "member_expression" && parent.Type() != "call_expression" {
			return false
		}

		toMatch := node.Content()

		for _, imported := range imports {
			if imported.IsMaster {
				if imported.Alias == toMatch {
					return true
				}
			}

			if !imported.IsMaster && !imported.IsRoot && imported.Name == knexModule {
				if imported.Alias == knexModule {
					return true
				}
			}
		}
	}
	return false
}
