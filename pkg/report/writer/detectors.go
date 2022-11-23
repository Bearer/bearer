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
	Value        schema.Schema
	Source       *source.Source
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
	StoredSchemas SchemaGroup
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

func (report *Detectors) AddSchema(
	detectorType detectors.Type,
	schema schema.Schema,
	source source.Source,
) {
	// @todo FIXME: Add classification here

	report.AddDetection(detections.TypeSchema, detectorType, source, schema)
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

func (report *Detectors) AddFillerLine() {
	report.Add(&detections.Detection{
		Type: detections.TypeFiller,
	})
}

func (report *Detectors) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
