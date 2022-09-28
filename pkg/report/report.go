package report

import (
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"

	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/util/blamer"

	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/secret"
	"github.com/bearer/curio/pkg/report/source"
)

type DetectionType string

var TypeDependency DetectionType = "dependency"
var TypeInterface DetectionType = "interface"
var TypeSchema DetectionType = "schema"
var TypeCreateView DetectionType = "create_view"
var TypeOperation DetectionType = "operation"
var TypeFramework DetectionType = "framework"
var TypeFiller DetectionType = "filler"
var TypeError DetectionType = "error"
var TypeFileFailed DetectionType = "file_error"
var TypeSecretleak DetectionType = "secret_leak"
var TypeCustom DetectionType = "custom"

type Report interface {
	SchemaReport
	CustomReport
	AddCreateView(detectorType detectors.Type, createView createview.View)
	AddOperation(detectorType detectors.Type, operation operations.Operation, source source.Source)
	AddInterface(detectorType detectors.Type, data interfaces.Interface, source source.Source)
	AddFramework(detectorType detectors.Type, frameworkType frameworks.Type, data interface{}, source source.Source)
	AddDependency(detectorType detectors.Type, dependency dependencies.Dependency, source source.Source)
	AddSecretLeak(secret secret.Secret, source source.Source)
	AddFillerLine()
	AddError(err error)
}

type SchemaReport interface {
	AddSchema(detectorType detectors.Type, schema schema.Schema, source source.Source)
}

// broker writes those for files that scanner fails to proccess
type FileFailedDetection struct {
	Type     DetectionType `json:"type"`
	File     string        `json:"file"`
	FileSize int           `json:"file_size"`
	Timeout  time.Duration `json:"timeout_duration"`
	Error    string        `json:"error"`
}

type ErrorDetection struct {
	Type    DetectionType `json:"type"`
	Message string        `json:"message"`
}

type FrameworkDetection struct {
	Type          DetectionType   `json:"type"`
	DetectorType  detectors.Type  `json:"detector_type"`
	FrameworkType frameworks.Type `json:"framework_detection_type"`
	CommitSHA     string          `json:"commit_sha"`
	Source        source.Source   `json:"source"`
	Value         interface{}     `json:"value"`
}

type Detection struct {
	Type         DetectionType  `json:"type"`
	DetectorType detectors.Type `json:"detector_type"`
	CommitSHA    string         `json:"commit_sha"`
	Source       source.Source  `json:"source"`
	Value        interface{}    `json:"value"`
}

type JsonLinesReport struct {
	Blamer blamer.Blamer
	File   io.Writer
}

func (report *JsonLinesReport) AddDependency(
	detectorType detectors.Type,
	dependency dependencies.Dependency,
	source source.Source,
) {

	report.addDetection(&Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: TypeDependency})
}

func (report *JsonLinesReport) AddInterface(
	detectorType detectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	report.addDetection(&Detection{DetectorType: detectorType, Value: data, Source: source, Type: TypeInterface})
}

func (report *JsonLinesReport) AddOperation(
	detectorType detectors.Type,
	operation operations.Operation,
	source source.Source,
) {
	report.addDetection(&Detection{DetectorType: detectorType, Value: operation, Source: source, Type: TypeOperation})
}

func (report *JsonLinesReport) AddCreateView(
	detectorType detectors.Type,
	createview createview.View,
) {
	for _, field := range createview.Fields {
		field.CommitSHA = report.Blamer.SHAForLine(field.Source.Filename, *field.Source.LineNumber)
	}

	for _, field := range createview.From {
		field.CommitSHA = report.Blamer.SHAForLine(field.Source.Filename, *field.Source.LineNumber)
	}

	report.addDetection(&Detection{DetectorType: detectorType, Value: createview, Source: createview.Source, Type: TypeCreateView})
}

func (report *JsonLinesReport) AddSchema(
	detectorType detectors.Type,
	schema schema.Schema,
	source source.Source,
) {

	report.addDetection(&Detection{DetectorType: detectorType, Value: schema, Source: source, Type: TypeSchema})
}

func (report *JsonLinesReport) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
	report.addDetection(&Detection{DetectorType: detectors.DetectorGitleaks, Value: secret, Source: source, Type: TypeSecretleak})
}

func (report *JsonLinesReport) addDetection(data *Detection) {
	if data.Source.LineNumber != nil {
		data.CommitSHA = report.Blamer.SHAForLine(data.Source.Filename, *data.Source.LineNumber)
	}

	detectionsToAdd := []*Detection{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}

func (report *JsonLinesReport) AddFramework(
	detectorType detectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {
	var commitSHA string
	if source.LineNumber != nil {
		commitSHA = report.Blamer.SHAForLine(source.Filename, *source.LineNumber)
	}

	detectionsToAdd := []*FrameworkDetection{{
		Type:          TypeFramework,
		DetectorType:  detectorType,
		FrameworkType: frameworkType,
		CommitSHA:     commitSHA,
		Source:        source,
		Value:         data,
	}}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}

func (report *JsonLinesReport) AddError(err error) {
	data := []*ErrorDetection{{
		Type:    TypeError,
		Message: err.Error(),
	}}

	errorEncoding := jsonlines.Encode(report.File, &data)
	if errorEncoding != nil {
		log.Printf("failed to encode data line %e", err)
	}
}

func (report *JsonLinesReport) AddFillerLine() {
	data := []*Detection{{
		Type: TypeFiller,
	}}

	errorEncoding := jsonlines.Encode(report.File, &data)
	if errorEncoding != nil {
		log.Printf("failed to encode data line %e", errorEncoding)
	}
}
