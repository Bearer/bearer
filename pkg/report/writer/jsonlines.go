package writer

import (
	"io"
	"log"

	classsification "github.com/bearer/curio/pkg/classification"
	classsificationschema "github.com/bearer/curio/pkg/classification/schema"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"

	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/report/secret"
	"github.com/bearer/curio/pkg/report/source"

	"github.com/bearer/curio/pkg/util/blamer"
	"github.com/wlredeye/jsonlines"
)

type JSONLines struct {
	Blamer     blamer.Blamer
	Classifier *classsification.Classifier
	File       io.Writer
}

func (report *JSONLines) AddInterface(
	detectorType detectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeInterface}
	classifiedDetection, err := report.Classifier.Interfaces.Classify(*detection)
	if err != nil {
		report.AddError(detection.Source.Filename, err)
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeInterfaceClassified
	report.Add(classifiedDetection)
}

func (report *JSONLines) AddOperation(
	detectorType detectors.Type,
	operation operations.Operation,
	source source.Source,
) {
	report.AddDetection(detections.TypeOperation, detectorType, source, operation)
}

func (report *JSONLines) AddCreateView(
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

func (report *JSONLines) AddDataType(detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType) {
	classifiedDatatypes := make(map[parser.NodeID]*classsificationschema.ClassifiedDatatype, 0)
	for nodeID, target := range values {
		classified, err := report.Classifier.Schema.Classify(classsificationschema.DataTypeDetection{
			Value:        target,
			Filename:     target.GetNode().Source(false).Filename,
			DetectorType: detectorType,
		})
		if err != nil {
			report.AddError(target.GetNode().Source(false).Filename, err)
		}

		classifiedDatatypes[nodeID] = classified
	}

	if detectionType == detections.TypeCustom {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, false, values)
	} else {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, true, values)
	}
}

func (report *JSONLines) AddSchema(
	detectorType detectors.Type,
	schema schema.Schema,
	source source.Source,
) {
	report.AddDetection(detections.TypeSchema, detectorType, source, schema)
}

func (report *JSONLines) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
	report.AddDetection(detections.TypeSecretleak, detectors.DetectorGitleaks, source, secret)
}

func (report *JSONLines) AddDetection(detectionType detections.DetectionType, detectorType detectors.Type, source source.Source, value interface{}) {
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

func (report *JSONLines) AddDependency(
	detectorType detectors.Type,
	dependency dependencies.Dependency,
	source source.Source,
) {

	detection := &detections.Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: detections.TypeDependency}
	classifiedDetection, err := report.Classifier.Dependencies.Classify(*detection)
	if err != nil {
		report.AddError(source.Filename, err)
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeDependencyClassified
	report.Add(classifiedDetection)
}

func (report *JSONLines) AddFramework(
	detectorType detectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {
	var commitSHA string
	if source.LineNumber != nil {
		commitSHA = report.Blamer.SHAForLine(source.Filename, *source.LineNumber)
	}

	report.Add(&detections.FrameworkDetection{
		Type:          detections.TypeFramework,
		DetectorType:  detectorType,
		FrameworkType: frameworkType,
		CommitSHA:     commitSHA,
		Source:        source,
		Value:         data,
	})
}

func (report *JSONLines) AddError(file string, err error) {
	report.Add(&detections.ErrorDetection{
		Type:    detections.TypeError,
		Message: err.Error(),
		File:    file,
	})
}

func (report *JSONLines) AddFillerLine() {
	report.Add(&detections.Detection{
		Type: detections.TypeFiller,
	})
}

func (report *JSONLines) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
