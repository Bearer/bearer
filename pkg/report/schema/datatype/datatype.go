package datatype

import (
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
)

type DataType struct {
	Node       *parser.Node
	Name       string
	Type       string
	TextType   string
	Properties map[string]*DataType
	IsHelper   bool // helper dataTypes and their child datatypes don't get exported
	UUID       string
}

type SchemaReport interface {
	AddSchema(detectorType detectors.Type, schema schema.Schema, source source.Source)
}

func ExportSchemas(report SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, ignoreFirst bool, values map[parser.NodeID]*DataType) {
	sortedDataTypes := SortParserMap(values)

	parentName := ""
	parentUUID := idGenerator.GenerateId()
	for _, value := range sortedDataTypes {
		dataTypeToSchema(report, detectorType, idGenerator, value, parentName, parentUUID, !ignoreFirst)
	}
}

func dataTypeToSchema(report SchemaReport, detectorType detectors.Type, idGenerator nodeid.Generator, dataType *DataType, parentName string, parentUUID string, shouldExport bool) {
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

	sortedProperties := SortStringMap(dataType.Properties)

	for _, property := range sortedProperties {
		dataTypeToSchema(report, detectorType, idGenerator, property, selfName, selfUUID, true)
	}
}

func SortStringMap(input map[string]*DataType) []*DataType {
	var sortedDataTypes []*DataType
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return SortSlice(sortedDataTypes)
}

func SortParserMap(input map[parser.NodeID]*DataType) []*DataType {
	var sortedDataTypes []*DataType
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return SortSlice(sortedDataTypes)
}

func SortSlice(input []*DataType) []*DataType {
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
