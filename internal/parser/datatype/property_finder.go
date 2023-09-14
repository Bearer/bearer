package datatype

import (
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/report/schema/datatype"
	sitter "github.com/smacker/go-tree-sitter"
)

type PropertyFinder struct {
	Map       map[parser.NodeID]*datatype.DataType
	tree      *parser.Tree
	parseNode func(resolver *PropertyFinder, node *parser.Node)
}

func NewPropertyFinder(tree *parser.Tree, dataTypeMap map[parser.NodeID]*datatype.DataType, parseNode func(resolver *PropertyFinder, node *parser.Node)) *PropertyFinder {
	return &PropertyFinder{
		tree:      tree,
		Map:       dataTypeMap,
		parseNode: parseNode,
	}
}

func (finder *PropertyFinder) Find() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck
		finder.parseNode(finder, child)

		return nil
	})
}

// ResolveContext gets closest context node belongs to
func (finder *PropertyFinder) ResolveClosestDataType(node *sitter.Node) *datatype.DataType {
	currentNode := node
	for {
		if ctx, ok := finder.Map[currentNode]; ok {
			return ctx
		}

		if currentNode.Parent() == nil {
			return nil
		} else {
			currentNode = currentNode.Parent()
		}
	}
}
