package testhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/report"

	reporttypes "github.com/bearer/curio/pkg/report"
	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/dependencies"
	reportdetectors "github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/schema"
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
	CustomDetections []*reporttypes.CustomDetection
	Detections       []*reporttypes.Detection
	Dependencies     []*reporttypes.Detection
	Frameworks       []*reporttypes.FrameworkDetection
	Errors           []*reporttypes.ErrorDetection
	SecretLeaks      []*reporttypes.Detection
	CreateView       []*reporttypes.Detection
}

func (report *InMemoryReport) AddInterface(
	detectorType reportdetectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	report.Detections = append(report.Detections, &reporttypes.Detection{
		DetectorType: detectorType,
		Value:        data,
		Source:       source,
		Type:         reporttypes.TypeInterface,
	})
}

func (inMemReport *InMemoryReport) AddCustomDetection(ruleName string, source source.Source, value schema.Schema) {
	data := &report.CustomDetection{
		Type:         report.TypeCustom,
		DetectorType: reportdetectors.Type(ruleName),
		Source:       source,
		Value:        value,
	}
	inMemReport.CustomDetections = append(inMemReport.CustomDetections, data)
}

func (report *InMemoryReport) AddOperation(
	detectorType reportdetectors.Type,
	operation operations.Operation,
	source source.Source,
) {
	report.Detections = append(report.Detections, &reporttypes.Detection{
		DetectorType: detectorType,
		Value:        operation,
		Source:       source,
		Type:         reporttypes.TypeOperation,
	})
}

func (report *InMemoryReport) AddSchema(
	detectorType reportdetectors.Type,
	schema schema.Schema,
	source source.Source,
) {

	report.Detections = append(report.Detections, &reporttypes.Detection{
		DetectorType: detectorType,
		Value:        schema,
		Source:       source,
		Type:         reporttypes.TypeSchema,
	})
}

func (report *InMemoryReport) AddCreateView(
	detectorType reportdetectors.Type,
	createview createview.View,
) {

	report.CreateView = append(report.CreateView, &reporttypes.Detection{
		DetectorType: detectorType,
		Value:        createview,
		Source:       createview.Source,
		Type:         reporttypes.TypeCreateView,
	})
}

func (report *InMemoryReport) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {

	report.SecretLeaks = append(report.SecretLeaks, &reporttypes.Detection{
		DetectorType: reportdetectors.DetectorGitleaks,
		Value:        secret,
		Source:       source,
		Type:         reporttypes.TypeSecretleak,
	})
}

func (report *InMemoryReport) AddDependency(
	detectorType reportdetectors.Type,
	dependency dependencies.Dependency,
	source source.Source,
) {

	report.Dependencies = append(report.Dependencies, &reporttypes.Detection{DetectorType: detectorType, Value: dependency, Source: source, Type: reporttypes.TypeDependency})
}

func (report *InMemoryReport) AddFramework(
	detectorType reportdetectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {

	report.Frameworks = append(report.Frameworks, &reporttypes.FrameworkDetection{
		Type:          reporttypes.TypeFramework,
		DetectorType:  detectorType,
		FrameworkType: frameworkType,
		Source:        source,
		Value:         data,
	})
}

func (report *InMemoryReport) AddError(err error) {
	report.Errors = append(report.Errors, &reporttypes.ErrorDetection{
		Type:    reporttypes.TypeError,
		Message: err.Error(),
	})
}

func (report *InMemoryReport) AddFillerLine() {

}
