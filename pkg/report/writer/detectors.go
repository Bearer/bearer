package writer

import (
	"fmt"
	"io"
	"log"

	classification "github.com/bearer/curio/pkg/classification"
	classificationschema "github.com/bearer/curio/pkg/classification/schema"
	zerolog "github.com/rs/zerolog/log"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"

	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/report/secret"
	"github.com/bearer/curio/pkg/report/source"

	"github.com/bearer/curio/pkg/util/blamer"
	"github.com/wlredeye/jsonlines"
)

type StoredSchema struct {
	Value  schema.Schema
	Source *source.Source
	Parent *parser.Node
}

type StoredSchemaNodes = map[*parser.Node]*StoredSchema

type SchemaGroup struct {
	Node         *parser.Node
	ParentSchema StoredSchema
	DetectorType detectors.Type
	Schemas      StoredSchemaNodes
}

type Detectors struct {
	Blamer        blamer.Blamer
	Classifier    *classification.Classifier
	File          io.Writer
	StoredSchemas *SchemaGroup
}

func (report *Detectors) AddInterface(
	detectorType detectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeInterface}
	classifiedDetection, err := report.Classifier.Interfaces.Classify(*detection)
	if err != nil {
		zerolog.Debug().Msgf("classification interfaces error from %s: %s", detection.Source.Filename, err)
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeInterfaceClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddCreateView(
	detectorType detectors.Type,
	createview createview.View,
) {
	for _, field := range createview.Fields {
		field.CommitSHA = report.Blamer.SHAForLine(field.Source.Filename, *field.Source.LineNumber)
	}

	for _, field := range createview.From {
		field.CommitSHA = report.Blamer.SHAForLine(field.Source.Filename, *field.Source.LineNumber)
	}

	report.AddDetection(detections.TypeCreateView, detectorType, createview.Source, createview)
}

func (report *Detectors) AddDataType(detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType, parent *parser.Node) {
	classifiedDatatypes := make(map[parser.NodeID]*classificationschema.ClassifiedDatatype, 0)
	for nodeID, target := range values {
		classified := report.Classifier.Schema.Classify(classificationschema.DataTypeDetection{
			Value:        target,
			Filename:     target.GetNode().Source(false).Filename,
			DetectorType: detectorType,
		})

		classifiedDatatypes[nodeID] = classified
	}

	if detectionType == detections.TypeCustom {
		datatype.ExportClassified(report, detections.TypeCustomClassified, detectorType, idGenerator, classifiedDatatypes, parent)
	} else {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, classifiedDatatypes, nil)
	}
}

func (report *Detectors) SchemaGroupBegin(detectorType detectors.Type, node *parser.Node, schema schema.Schema, source *source.Source, parent *parser.Node) {
	if report.SchemaGroupIsOpen() {
		zerolog.Warn().Msg("schema group already open")
	}
	report.StoredSchemas = &SchemaGroup{
		Node: node,
		ParentSchema: StoredSchema{
			Value:  schema,
			Source: source,
			Parent: parent,
		},
		DetectorType: detectorType,
		Schemas:      make(StoredSchemaNodes),
	}
}

func (report *Detectors) SchemaGroupIsOpen() bool {
	return report.StoredSchemas != nil
}

func (report *Detectors) SchemaGroupShouldClose(tableName string) bool {
	if report.StoredSchemas == nil {
		return false
	}
	return tableName != report.StoredSchemas.ParentSchema.Value.ObjectName
}

func (report *Detectors) SchemaGroupAddItem(node *parser.Node, schema schema.Schema, source *source.Source) {
	report.StoredSchemas.Schemas[node] = &StoredSchema{Value: schema, Source: source, Parent: report.StoredSchemas.ParentSchema.Parent}
}

func (report *Detectors) SchemaGroupEnd(idGenerator nodeid.Generator) {
	// Build child data types
	childDataTypes := map[string]datatype.DataTypable{}
	for node, storedSchema := range report.StoredSchemas.Schemas {
		schema := storedSchema.Value

		childName := schema.FieldName
		childDataTypes[childName] = &datatype.DataType{
			Node:       node,
			Name:       childName,
			Type:       schema.SimpleFieldType,
			TextType:   schema.FieldType,
			Properties: map[string]datatype.DataTypable{},
			UUID:       schema.FieldUUID,
		}
	}

	// Build parent data type
	parentDataType := &datatype.DataType{
		Node:       report.StoredSchemas.Node,
		Name:       report.StoredSchemas.ParentSchema.Value.ObjectName,
		Type:       "",
		TextType:   "",
		Properties: childDataTypes,
		UUID:       report.StoredSchemas.ParentSchema.Value.ObjectUUID,
	}

	classifiedDatatypes := make(map[parser.NodeID]*classificationschema.ClassifiedDatatype, 0)

	// Classify child data types
	for node, storedSchema := range report.StoredSchemas.Schemas {
		source := storedSchema.Source
		schema := storedSchema.Value

		childName := schema.FieldName
		value := childDataTypes[childName]
		detection := classificationschema.DataTypeDetection{DetectorType: report.StoredSchemas.DetectorType, Value: value, Filename: source.Filename}
		classifiedDatatype := report.Classifier.Schema.Classify(detection)

		classifiedDatatypes[node.ID()] = classifiedDatatype
	}

	// Classify parent data type
	parentDetection := classificationschema.DataTypeDetection{DetectorType: report.StoredSchemas.DetectorType, Value: parentDataType, Filename: report.StoredSchemas.ParentSchema.Source.Filename}
	classifiedParentDatatype := report.Classifier.Schema.Classify(parentDetection)
	classifiedDatatypes[report.StoredSchemas.Node.ID()] = classifiedParentDatatype

	// Export classified data types
	datatype.ExportClassified(report, detections.TypeSchemaClassified, report.StoredSchemas.DetectorType, idGenerator, classifiedDatatypes, report.StoredSchemas.ParentSchema.Parent)

	// Clear the map of stored schema detection information
	report.StoredSchemas = nil
}

func (report *Detectors) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
	report.AddDetection(detections.TypeSecretleak, detectors.DetectorGitleaks, source, secret)
}

func (report *Detectors) AddDetection(detectionType detections.DetectionType, detectorType detectors.Type, source source.Source, value interface{}) {
	data := &detections.Detection{
		Type:         detectionType,
		DetectorType: detectorType,
		Source:       source,
		Value:        value,
	}

	if data.Source.LineNumber != nil {
		data.CommitSHA = report.Blamer.SHAForLine(data.Source.Filename, *data.Source.LineNumber)
	}

	report.Add(data)
}

func (report *Detectors) AddDependency(
	detectorType detectors.Type,
	dependency dependencies.Dependency,
	source source.Source,
) {

	detection := &detections.Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: detections.TypeDependency}
	classifiedDetection, err := report.Classifier.Dependencies.Classify(*detection)
	if err != nil {
		report.AddError(source.Filename, fmt.Errorf("classification dependencies error: %s", err))
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeDependencyClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddFramework(
	detectorType detectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeFramework}
	classifiedDetection, err := report.Classifier.Frameworks.Classify(*detection)
	if err != nil {
		report.AddError(source.Filename, fmt.Errorf("classification frameworks error: %s", err))
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeFrameworkClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddError(filePath string, err error) {
	report.Add(&detections.ErrorDetection{
		Type:    detections.TypeError,
		Message: err.Error(),
		File:    filePath,
	})
}

func (report *Detectors) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
