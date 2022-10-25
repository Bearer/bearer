package writer

import (
	"io"
	"log"

	classsification "github.com/bearer/curio/pkg/classification"
	classsificationschema "github.com/bearer/curio/pkg/classification/schema"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"

	reporttypes "github.com/bearer/curio/pkg/report"
	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
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
	detection := &reporttypes.Detection{DetectorType: detectorType, Value: data, Source: source, Type: reporttypes.TypeInterface}
	classifiedDetection, err := report.Classifier.Interfaces.Classify(*detection)
	if err != nil {
		report.AddError(err)
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = reporttypes.TypeInterfaceClassified
	report.Add(classifiedDetection)
}

func (report *JSONLines) AddOperation(
	detectorType detectors.Type,
	operation operations.Operation,
	source source.Source,
) {
	report.addDetection(&reporttypes.Detection{DetectorType: detectorType, Value: operation, Source: source, Type: reporttypes.TypeOperation})
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

	report.addDetection(&reporttypes.Detection{DetectorType: detectorType, Value: createview, Source: createview.Source, Type: reporttypes.TypeCreateView})
}

func (report *JSONLines) AddDatatype(detectorType detectors.Type, idGenerator nodeid.Generator, ignorefirst bool, values map[parser.NodeID]*datatype.DataType) {
	var classifiedDatatypes []*classsificationschema.ClassifiedDatatype
	for _, target := range values {
		classified, err := report.Classifier.Schema.Classify(classsificationschema.DataTypeDetection{
			Value:        *target,
			Filename:     target.Node.Source(false).Filename,
			DetectorType: detectorType,
		})
		if err != nil {
			report.AddError(err)
		}

		classifiedDatatypes = append(classifiedDatatypes, classified)
	}

	datatype.ExportSchemas(report, detectorType, idGenerator, ignorefirst, classifiedDatatypes)
}

func (report *JSONLines) AddSchema(
	detectorType detectors.Type,
	schema schema.Schema,
	source source.Source,
) {
	report.addDetection(&reporttypes.Detection{DetectorType: detectorType, Value: schema, Source: source, Type: reporttypes.TypeSchema})
}

func (report *JSONLines) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
	report.addDetection(&reporttypes.Detection{DetectorType: detectors.DetectorGitleaks, Value: secret, Source: source, Type: reporttypes.TypeSecretleak})
}

func (report *JSONLines) addDetection(data *reporttypes.Detection) {
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

	detection := &reporttypes.Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: reporttypes.TypeDependency}
	classifiedDetection, err := report.Classifier.Dependencies.Classify(*detection)
	if err != nil {
		report.AddError(err)
		return
	}

	if classifiedDetection.Source.LineNumber != nil {
		classifiedDetection.CommitSHA = report.Blamer.SHAForLine(classifiedDetection.Source.Filename, *classifiedDetection.Source.LineNumber)
	}

	classifiedDetection.Type = reporttypes.TypeDependencyClassified
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

	report.Add(&reporttypes.FrameworkDetection{
		Type:          reporttypes.TypeFramework,
		DetectorType:  detectorType,
		FrameworkType: frameworkType,
		CommitSHA:     commitSHA,
		Source:        source,
		Value:         data,
	})
}

func (report *JSONLines) AddError(err error) {
	report.Add(&reporttypes.ErrorDetection{
		Type:    reporttypes.TypeError,
		Message: err.Error(),
	})
}

func (report *JSONLines) AddFillerLine() {
	report.Add(&reporttypes.Detection{
		Type: reporttypes.TypeFiller,
	})
}

func (report *JSONLines) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
