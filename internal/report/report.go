package report

import (
	"github.com/bearer/bearer/internal/report/dependencies"
	"github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/frameworks"
	"github.com/bearer/bearer/internal/report/interfaces"
	"github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/report/schema/datatype"

	"github.com/bearer/bearer/internal/report/secret"
	"github.com/bearer/bearer/internal/report/source"
)

type Report interface {
	detections.ReportDetection
	schema.ReportSchema
	datatype.ReportDataType
	AddInterface(detectorType detectors.Type, data interfaces.Interface, source source.Source)
	AddFramework(detectorType detectors.Type, frameworkType frameworks.Type, data interface{}, source source.Source)
	AddDependency(detectorType detectors.Type, detectorLanguage detectors.Language, dependency dependencies.Dependency, source source.Source)
	AddSecretLeak(secret secret.Secret, source source.Source)
	AddError(filePath string, err error)
}
