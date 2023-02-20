package report

import (
	"github.com/bearer/bearer/pkg/report/dependencies"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/frameworks"
	"github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/schema/datatype"

	"github.com/bearer/bearer/pkg/report/secret"
	"github.com/bearer/bearer/pkg/report/source"
)

type Report interface {
	detections.ReportDetection
	schema.ReportSchema
	datatype.ReportDataType
	AddInterface(detectorType detectors.Type, data interfaces.Interface, source source.Source)
	AddFramework(detectorType detectors.Type, frameworkType frameworks.Type, data interface{}, source source.Source)
	AddDependency(detectorType detectors.Type, dependency dependencies.Dependency, source source.Source)
	AddSecretLeak(secret secret.Secret, source source.Source)
	AddError(filePath string, err error)
}

type ReportOutput interface {
	Output(filePath string, outputType string)
}
