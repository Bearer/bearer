package datatype

import (
	"errors"
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
)

type Finder struct {
	tree      *parser.Tree
	values    map[parser.NodeID]*DataType
	parseNode func(finder *Finder, node *parser.Node, value *DataType) bool
}

type DataType struct {
	Node       *parser.Node
	Name       string
	Type       string
	TextType   string
	Properties map[string]*DataType
	IsHelper   bool // helper dataTypes and their child datatypes don't get exported
	UUID       string
}

func NewDataType() DataType {
	return DataType{
		Properties: make(map[string]*DataType),
	}
}

func NewFinder(tree *parser.Tree, parseNode func(finder *Finder, node *parser.Node, value *DataType) bool) *Finder {
	return &Finder{
		tree:      tree,
		parseNode: parseNode,
		values:    make(map[parser.NodeID]*DataType),
	}
}

func NewExport(report report.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*DataType) {
	finder := &Finder{
		values: values,
	}
	finder.ExportSchemas(report, detectorType, idGenerator, true)
}

func NewCompleteExport(report report.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*DataType) {
	finder := &Finder{
		values: values,
	}
	finder.ExportSchemas(report, detectorType, idGenerator, false)
}

func (finder *Finder) Find() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck

		value := &DataType{
			Properties: make(map[string]*DataType),
		}
		found := finder.parseNode(finder, child, value)
		if found {
			finder.values[child.ID()] = value
		}

		return nil
	})
}

func (finder *Finder) GetValues() map[parser.NodeID]*DataType {
	return finder.values
}

func (finder *Finder) ExportSchemas(report report.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, ignoreFirst bool) {
	sortedDataTypes := sortParserDataTypeMap(finder.values)

	parentName := ""
	parentUUID := idGenerator.GenerateId()
	for _, value := range sortedDataTypes {
		dataTypeToSchema(report, detectorType, idGenerator, value, parentName, parentUUID, !ignoreFirst)
	}
}

func dataTypeToSchema(report report.SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, dataType *DataType, parentName string, parentUUID string, shouldExport bool) {
	if dataType.IsHelper {
		return
	}

	selfUUID := dataType.UUID
	if dataType.UUID == "" {
		selfUUID = idGenerator.GenerateId()
	}

	selfName := dataType.Name

	if shouldExport {
		report.AddSchema(detectorType, schema.Schema{
			ObjectName:      parentName,
			FieldName:       selfName,
			ObjectUUID:      parentUUID,
			FieldUUID:       selfUUID,
			FieldType:       dataType.TextType,
			SimpleFieldType: dataType.Type,
		}, dataType.Node.Source(false))
	}

	sortedProperties := sortStringDataTypeMap(dataType.Properties)

	for _, property := range sortedProperties {
		dataTypeToSchema(report, detectorType, idGenerator, property, selfName, selfUUID, true)
	}
}

func sortStringDataTypeMap(input map[string]*DataType) []*DataType {
	var sortedDataTypes []*DataType
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return sortSliceDataType(sortedDataTypes)
}

func sortParserDataTypeMap(input map[parser.NodeID]*DataType) []*DataType {
	var sortedDataTypes []*DataType
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return sortSliceDataType(sortedDataTypes)
}

func sortSliceDataType(input []*DataType) []*DataType {
	sort.Slice(input, func(i, j int) bool {
		lineNumberA := input[i].Node.Source(false).LineNumber
		lineNumberB := input[j].Node.Source(false).LineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := input[i].Node.Source(false).ColumnNumber
		columnNumberB := input[j].Node.Source(false).ColumnNumber

		if *columnNumberA != *columnNumberB {
			return *columnNumberA < *columnNumberB
		}

		if strings.Compare(input[i].Name, input[j].Name) == -1 {
			return true
		}

		return false
	})
	return input
}

func DeepestSingleChild(datatype *DataType) (*DataType, error) {
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

func PruneMap(datatypes map[parser.NodeID]*DataType) {
	for key, datatype := range datatypes {
		if Prune(datatype) {
			delete(datatypes, key)
		}
	}
}

func Prune(datatype *DataType) bool {
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
