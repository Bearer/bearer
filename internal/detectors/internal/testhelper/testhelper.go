package testhelper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/internal/detectors"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/nodeid"

	"github.com/bearer/bearer/internal/report/dependencies"
	"github.com/bearer/bearer/internal/report/detections"
	reportdetectors "github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/frameworks"
	"github.com/bearer/bearer/internal/report/interfaces"
	"github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/report/schema/datatype"
	"github.com/bearer/bearer/internal/report/secret"
	"github.com/bearer/bearer/internal/report/source"

	"github.com/stretchr/testify/assert"
)

func Extract(
	t assert.TestingT,
	path string,
	registrations []detectors.InitializedDetector,
	detectorType reportdetectors.Type,
) *InMemoryReport {
	report := InMemoryReport{}

	var files []string

	err := filepath.Walk(path,
		func(fullPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			file := strings.TrimPrefix(fullPath, path)

			files = append(files, file)
			return nil
		})
	if !assert.Nil(t, err) {
		t.Errorf("report has errored %s", err)
	}

	for _, filename := range files {
		err = detectors.ExtractWithDetectors(context.Background(), path, filename, &report, nil, registrations, nil)
		if !assert.Nil(t, err) {
			t.Errorf("report has errored %s", err)
		}
	}

	if len(report.Errors) > 0 {
		t.Errorf("report has some errors %#v", report.Errors)
	}

	return &report
}

func RegistrationFor(detectorType reportdetectors.Type) []detectors.InitializedDetector {
	scanners := []string{"sast", "secrets"}
	for _, registration := range detectors.Registrations(scanners) {
		if registration.Type == detectorType {
			return []detectors.InitializedDetector{registration}
		}
	}

	panic(fmt.Sprintf("missing registration for '%s'", detectorType))
}

type InMemoryReport struct {
	CustomDetections        []detections.Detection
	Detections              []*detections.Detection
	Dependencies            []*detections.Detection
	Frameworks              []*detections.FrameworkDetection
	Errors                  []*detections.ErrorDetection
	SecretLeaks             []*detections.Detection
	CreateView              []*detections.Detection
	SchemaGroupDetectorType reportdetectors.Type
	SchemaGroupObjectName   string
}

func (report *InMemoryReport) AddDetection(
	detectionType detections.DetectionType,
	detectorType reportdetectors.Type,
	source source.Source,
	value interface{},
) {
	detection := &detections.Detection{
		Type:         detectionType,
		DetectorType: detectorType,
		Source:       source,
		Value:        value,
	}
	if detectionType == detections.TypeCustom || detectionType == detections.TypeCustomRisk {
		report.CustomDetections = append(report.CustomDetections, *detection)
	} else {
		report.Detections = append(report.Detections, detection)
	}
}

func (report *InMemoryReport) SchemaGroupBegin(detectorType reportdetectors.Type, node *parser.Node, schema schema.Schema, source *source.Source, parent *parser.Node) {
	report.SchemaGroupDetectorType = detectorType
	report.SchemaGroupObjectName = schema.ObjectName
}

func (report *InMemoryReport) SchemaGroupIsOpen() bool {
	return report.SchemaGroupDetectorType != ""
}

func (report *InMemoryReport) SchemaGroupShouldClose(tableName string) bool {
	return tableName != report.SchemaGroupObjectName
}

func (report *InMemoryReport) SchemaGroupAddItem(node *parser.Node, schema schema.Schema, source *source.Source) {
	report.Detections = append(report.Detections, &detections.Detection{
		DetectorType: report.SchemaGroupDetectorType,
		Value:        schema,
		Source:       *source,
		Type:         detections.TypeSchema,
	})
}

func (report *InMemoryReport) SchemaGroupEnd(idGenerator nodeid.Generator) {
	report.SchemaGroupDetectorType = ""
	report.SchemaGroupObjectName = ""
}

func (report *InMemoryReport) AddDataType(
	detectionType detections.DetectionType,
	detectorType reportdetectors.Type,
	idGenerator nodeid.Generator,
	values map[parser.NodeID]*datatype.DataType,
	parent *parser.Node,
) {

	datatype.ExportClassified(report, detectionType, detectorType, idGenerator, values, nil)
}

func (report *InMemoryReport) AddInterface(
	detectorType reportdetectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	report.Detections = append(report.Detections, &detections.Detection{
		DetectorType: detectorType,
		Value:        data,
		Source:       source,
		Type:         detections.TypeInterface,
	})
}

func (report *InMemoryReport) AddFramework(
	detectorType reportdetectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {

	report.Frameworks = append(report.Frameworks, &detections.FrameworkDetection{
		Type:          detections.TypeFramework,
		DetectorType:  detectorType,
		FrameworkType: frameworkType,
		Source:        source,
		Value:         data,
	})
}

func (report *InMemoryReport) AddDependency(
	detectorType reportdetectors.Type,
	detectorLanguage reportdetectors.Language,
	dependency dependencies.Dependency,
	source source.Source,
) {

	report.Dependencies = append(report.Dependencies, &detections.Detection{
		DetectorType:     detectorType,
		DetectorLanguage: detectorLanguage,
		Value:            dependency,
		Source:           source,
		Type:             detections.TypeDependency,
	})
}

func (report *InMemoryReport) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {

	report.SecretLeaks = append(report.SecretLeaks, &detections.Detection{
		DetectorType: reportdetectors.DetectorGitleaks,
		Value:        secret,
		Source:       source,
		Type:         detections.TypeSecretleak,
	})
}

func (report *InMemoryReport) AddError(filePath string, err error) {
	report.Errors = append(report.Errors, &detections.ErrorDetection{
		Type:    detections.TypeError,
		Message: err.Error(),
		File:    filePath,
	})
}
