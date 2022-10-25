package datatype

import (
	"errors"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema/datatype"
)

type Finder struct {
	tree      *parser.Tree
	values    map[parser.NodeID]*datatype.DataType
	parseNode func(finder *Finder, node *parser.Node, value *datatype.DataType) bool
}

func NewDataType() datatype.DataType {
	return datatype.DataType{
		Properties: make(map[string]*datatype.DataType),
	}
}

func NewFinder(tree *parser.Tree, parseNode func(finder *Finder, node *parser.Node, value *datatype.DataType) bool) *Finder {
	return &Finder{
		tree:      tree,
		parseNode: parseNode,
		values:    make(map[parser.NodeID]*datatype.DataType),
	}
}

func NewExport(report datatype.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType) {
	datatype.ExportSchemas(report, detectorType, idGenerator, true, values)
}

func NewCompleteExport(report datatype.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType) {
	datatype.ExportSchemas(report, detectorType, idGenerator, false, values)
}

func (finder *Finder) Find() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck

		value := &datatype.DataType{
			Properties: make(map[string]*datatype.DataType),
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

func DeepestSingleChild(datatype *datatype.DataType) (*datatype.DataType, error) {
	if len(datatype.Properties) > 1 {
		return datatype, errors.New("multiple childs detected")
	}

	if len(datatype.Properties) == 0 {
		return datatype, nil
	}

	if len(datatype.Properties) == 1 {
		for _, child := range datatype.Properties {
			return DeepestSingleChild(child)
		}
	}

	return nil, errors.New("couldn't determine deepest child")
}

func PruneMap(datatypes map[parser.NodeID]*datatype.DataType) {
	for key, datatype := range datatypes {
		if Prune(datatype) {
			delete(datatypes, key)
		}
	}
}

func Prune(datatype *datatype.DataType) bool {
	if len(datatype.Name) < 3 {
		return true
	} else {
		for propertyKey, property := range datatype.Properties {
			if Prune(property) {
				delete(datatype.Properties, propertyKey)
			}
		}
	}
	return false
}
