package testhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"

	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
	"github.com/bearer/curio/pkg/report/detections"
	reportdetectors "github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/report/secret"
	"github.com/bearer/curio/pkg/report/source"

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
		t.Errorf("report has errored %e", err)
	}

	err = detectors.ExtractWithDetectors(path, files, &report, registrations)
	if len(report.Errors) > 0 {
		t.Errorf("report has some errors %#v", report.Errors)
	}

	if !assert.Nil(t, err) {
		t.Errorf("report has errored %e", err)
	}

	return &report
}

func RegistrationFor(detectorType reportdetectors.Type) []detectors.InitializedDetector {
	for _, registration := range detectors.Registrations() {
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

func (report *InMemoryReport) AddCreateView(
	detectorType reportdetectors.Type,
	createview createview.View,
) {

	report.CreateView = append(report.CreateView, &detections.Detection{
		DetectorType: detectorType,
		Value:        createview,
		Source:       createview.Source,
		Type:         detections.TypeCreateView,
	})
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
	dependency dependencies.Dependency,
	source source.Source,
) {

	report.Dependencies = append(report.Dependencies, &detections.Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: detections.TypeDependency})
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
