package writer

import (
	"fmt"
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

type DataFlow struct {
	Blamer     blamer.Blamer
	Classifier *classsification.Classifier
	File       io.Writer
}

func (report *DataFlow) AddInterface(
	detectorType detectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeInterface}
	classifiedDetection, err := report.Classifier.Interfaces.Classify(*detection)
	if err != nil {
		report.AddError(detection.Source.Filename, fmt.Errorf("classification interfaces error: %s", err))
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = detections.TypeInterfaceClassified
	report.Add(classifiedDetection)
}

func (report *DataFlow) AddDataType(detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType) {
	classifiedDatatypes := make(map[parser.NodeID]*classsificationschema.ClassifiedDatatype, 0)
	for nodeID, target := range values {
		classified, err := report.Classifier.Schema.Classify(classsificationschema.DataTypeDetection{
			Value:        target,
			Filename:     target.GetNode().Source(false).Filename,
			DetectorType: detectorType,
		})
		if err != nil {
			report.AddError(target.GetNode().Source(false).Filename, fmt.Errorf("classification datatypes error: %s", err))
		}

		classifiedDatatypes[nodeID] = classified
	}

	if detectionType == detections.TypeCustom {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, false, values)
	} else {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, true, values)
	}
}

func (report *DataFlow) AddDependency(
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

func (report *DataFlow) AddSchema(
	detectorType detectors.Type,
	schema schema.Schema,
	source source.Source,
) {
	report.AddDetection(detections.TypeSchema, detectorType, source, schema)
}

func (report *DataFlow) AddOperation(
	detectorType detectors.Type,
	operation operations.Operation,
	source source.Source,
) {

}

func (report *DataFlow) AddCreateView(
	detectorType detectors.Type,
	createview createview.View,
) {

}

func (report *DataFlow) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
}

func (report *DataFlow) AddDetection(detectionType detections.DetectionType, detectorType detectors.Type, source source.Source, value interface{}) {

}

func (report *DataFlow) AddFramework(
	detectorType detectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {

}

func (report *DataFlow) AddError(filePath string, err error) {

}

func (report *DataFlow) AddFillerLine() {
	report.Add(&detections.Detection{
		Type: detections.TypeFiller,
	})
}

func (report *DataFlow) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
