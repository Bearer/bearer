package datatype

import (
	"errors"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

type Finder struct {
	tree      *parser.Tree
	values    map[parser.NodeID]*datatype.DataType
	parseNode func(finder *Finder, node *parser.Node, value *datatype.DataType) bool
}

func NewDataType() datatype.DataType {
	return datatype.DataType{
		Properties: make(map[string]datatype.DataTypable),
	}
}

func NewFinder(tree *parser.Tree, parseNode func(finder *Finder, node *parser.Node, value *datatype.DataType) bool) *Finder {
	return &Finder{
		tree:      tree,
		parseNode: parseNode,
		values:    make(map[parser.NodeID]*datatype.DataType),
	}
}

func (finder *Finder) Find() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck

		value := &datatype.DataType{
			Properties: make(map[string]datatype.DataTypable),
		}
		found := finder.parseNode(finder, child, value)
		if found {
			finder.values[child.ID()] = value
		}

		return nil
	})
}

func (finder *Finder) GetValues() map[parser.NodeID]*datatype.DataType {
	return finder.values
}

func DeepestSingleChild(datatype datatype.DataTypable) (datatype.DataTypable, error) {
	if len(datatype.GetProperties()) > 1 {
		return datatype, errors.New("multiple childs detected")
	}

	if len(datatype.GetProperties()) == 0 {
		return datatype, nil
	}

	if len(datatype.GetProperties()) == 1 {
		for _, child := range datatype.GetProperties() {
			return DeepestSingleChild(child)
		}
	}

	return nil, errors.New("couldn't determine deepest child")
}

func PruneMap[D datatype.DataTypable](datatypes map[parser.NodeID]D) {
	for key, datatype := range datatypes {
		if Prune(datatype) {
			delete(datatypes, key)
		}
	}
}

func Prune[D datatype.DataTypable](datatype D) bool {
	if len(datatype.GetName()) < 3 {
		return true
	} else {
		for propertyKey, property := range datatype.GetProperties() {
			if Prune(property) {
				datatype.DeleteProperty(propertyKey)
			}
		}
	}
	return false
}
