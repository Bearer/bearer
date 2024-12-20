package datatype

import (
	"sort"
	"strings"

	classificationschema "github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/util/classify"
	"github.com/bearer/bearer/pkg/util/normalize_key"
	"github.com/bearer/bearer/pkg/util/pluralize"
)

type ReportDataType interface {
	AddDataType(detectionType detections.DetectionType, detectorType detectors.Type, generator nodeid.Generator, values map[parser.NodeID]*DataType, parent *parser.Node)
}

type DataType struct {
	Node       *parser.Node
	Name       string
	Type       string
	TextType   string
	Properties map[string]DataTypable
	IsHelper   bool // helper dataTypes and their child datatypes don't get exported
	UUID       string
}

type ClassifiedDatatype struct {
	*DataType
	Classification *classificationschema.ClassifiedDatatype
	Properties     map[string]DataTypable
}

func (datatype *ClassifiedDatatype) GetClassification() interface{} {
	return (*datatype.Classification).Classification
}

func (datatype *ClassifiedDatatype) GetProperties() map[string]DataTypable {
	return datatype.Properties
}

func BuildClassifiedDatatype(datatype *DataType, classification *classificationschema.ClassifiedDatatype) *ClassifiedDatatype {
	properties := make(map[string]DataTypable)

	for key, propertyDatatypable := range datatype.Properties {
		propertyDatatype, ok := propertyDatatypable.(*DataType)
		if !ok {
			continue
		}

		for _, classificationProperty := range classification.Properties {
			if key == classificationProperty.Name {
				properties[key] = &ClassifiedDatatype{
					DataType:       propertyDatatype,
					Classification: classificationProperty,
					Properties:     make(map[string]DataTypable),
				}
			}
		}
	}
	result := &ClassifiedDatatype{
		DataType:       datatype,
		Classification: classification,
		Properties:     properties,
	}

	return result
}

func (datatype *DataType) GetClassification() interface{} {
	return nil
}

func (datatype *DataType) SetName(name string) {
	datatype.Name = name
}

func (datatype *DataType) GetName() string {
	return datatype.Name
}

func (datatype *DataType) GetNormalizedName() string {
	return normalize_key.Normalize(datatype.Name)
}

func (datatype *DataType) GetNode() *parser.Node {
	return datatype.Node
}

func (datatype *DataType) GetProperties() map[string]DataTypable {
	return datatype.Properties
}

func (datatype *DataType) SetProperty(key string, property DataTypable) {
	datatype.Properties[key] = property
}

func (datatype *DataType) GetUUID() string {
	return datatype.UUID
}

func (datatype *DataType) GetIsHelper() bool {
	return datatype.IsHelper
}

func (datatype *DataType) GetTextType() string {
	return datatype.TextType
}

func (datatype *DataType) GetType() string {
	return datatype.Type
}

func (datatype *DataType) SetUUID(UUID string) {
	datatype.UUID = UUID
}

func (datatype *DataType) DeleteProperty(name string) {
	delete(datatype.Properties, name)
}

func (datatype *DataType) CreateProperties() {
	datatype.Properties = make(map[string]DataTypable)
}

func (datatype *DataType) Clone() DataTypable {
	cloned := &DataType{
		Node:       datatype.Node,
		Name:       datatype.Name,
		Type:       datatype.Type,
		TextType:   datatype.TextType,
		Properties: make(map[string]DataTypable),
		IsHelper:   datatype.IsHelper,
		UUID:       datatype.UUID,
	}

	for nodeID, child := range datatype.Properties {
		cloned.Properties[nodeID] = child.Clone()
	}

	return cloned
}

func (datatype *DataType) ToClassificationRequestDetection() *classificationschema.ClassificationRequestDetection {
	detection := &classificationschema.ClassificationRequestDetection{
		Name:       datatype.Name,
		SimpleType: datatype.Type,
	}

	for _, datatypeProperty := range datatype.Properties {
		detection.Properties = append(detection.Properties, datatypeProperty.ToClassificationRequestDetection())
	}

	return detection
}

type DataTypable interface {
	DeleteProperty(name string)
	GetClassification() interface{}
	SetUUID(string)
	GetUUID() string
	GetIsHelper() bool
	GetTextType() string
	GetType() string
	GetName() string
	GetNormalizedName() string
	SetName(string)
	GetNode() *parser.Node
	GetProperties() map[string]DataTypable
	Clone() DataTypable
	SetProperty(string, DataTypable)
	CreateProperties()
	ToClassificationRequestDetection() *classificationschema.ClassificationRequestDetection
}

func ExportClassified[D DataTypable](report detections.ReportDetection, detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]D, parent *parser.Node) {
	exportSchemas(report, detectionType, detectorType, idGenerator, values, parent)
}

func exportSchemas[D DataTypable](report detections.ReportDetection, detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]D, parent *parser.Node) {
	sortedDataTypes := SortParserMap(values)

	parentName := ""
	parentUUID := idGenerator.GenerateId()
	for _, value := range sortedDataTypes {
		dataTypeToSchema(report, detectionType, detectorType, idGenerator, value, parentName, parentUUID, false, parent)
	}
}

func dataTypeToSchema[D DataTypable](report detections.ReportDetection, detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, dataType D, parentName string, parentUUID string, shouldExport bool, parent *parser.Node) {
	if dataType.GetIsHelper() {
		return
	}

	if classification, ok := dataType.GetClassification().(classificationschema.Classification); ok && classification.Decision.State != classify.Valid {
		return
	}

	selfUUID := dataType.GetUUID()
	if selfUUID == "" {
		selfUUID = idGenerator.GenerateId()
	}

	selfName := dataType.GetName()

	if shouldExport {
		var sourceSchema *schema.Source

		if parent != nil {
			parentContent := parent.Content()
			sourceSchema = &schema.Source{
				Content:           parentContent,
				StartLineNumber:   parent.StartLineNumber(),
				StartColumnNumber: parent.StartColumnNumber(),
				EndLineNumber:     parent.EndLineNumber(),
				EndColumnNumber:   parent.EndColumnNumber(),
			}
		}

		normalizedObjectName := pluralize.Singular(strings.ToLower(parentName))
		normalizedFieldName := pluralize.Singular(strings.ToLower(selfName))

		report.AddDetection(
			detectionType,
			detectorType,
			dataType.GetNode().Source(false),
			schema.Schema{
				ObjectName:           parentName,
				FieldName:            selfName,
				ObjectUUID:           parentUUID,
				FieldUUID:            selfUUID,
				FieldType:            dataType.GetTextType(),
				SimpleFieldType:      dataType.GetType(),
				Classification:       dataType.GetClassification(),
				Source:               sourceSchema,
				NormalizedObjectName: normalizedObjectName,
				NormalizedFieldName:  normalizedFieldName,
			},
		)
	}

	sortedProperties := SortStringMap(dataType.GetProperties())

	for _, property := range sortedProperties {
		dataTypeToSchema(report, detectionType, detectorType, idGenerator, property, selfName, selfUUID, true, parent)
	}
}

func SortStringMap[D DataTypable](input map[string]D) []D {
	var sortedDataTypes []D
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return SortSlice(sortedDataTypes)
}

func SortParserMap[D DataTypable](input map[parser.NodeID]D) []D {
	var sortedDataTypes []D
	for _, value := range input {
		sortedDataTypes = append(sortedDataTypes, value)
	}
	return SortSlice(sortedDataTypes)
}

func SortSlice[D DataTypable](input []D) []D {
	sort.Slice(input, func(i, j int) bool {
		lineNumberA := input[i].GetNode().Source(false).StartLineNumber
		lineNumberB := input[j].GetNode().Source(false).StartLineNumber

		if *lineNumberA != *lineNumberB {
			return *lineNumberA < *lineNumberB
		}

		columnNumberA := input[i].GetNode().Source(false).StartColumnNumber
		columnNumberB := input[j].GetNode().Source(false).StartColumnNumber

		if *columnNumberA != *columnNumberB {
			return *columnNumberA < *columnNumberB
		}

		if strings.Compare(input[i].GetName(), input[j].GetName()) == -1 {
			return true
		}

		return false
	})
	return input
}
