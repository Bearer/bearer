package report

import (
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
)

type Report interface {
	detections.ReportDetection
	schema.ReportSchema
	datatype.ReportDataType
	AddCreateView(detectorType detectors.Type, createView createview.View)
	AddInterface(detectorType detectors.Type, data interfaces.Interface, source source.Source)
	AddFramework(detectorType detectors.Type, frameworkType frameworks.Type, data interface{}, source source.Source)
	AddDependency(detectorType detectors.Type, dependency dependencies.Dependency, source source.Source)
	AddSecretLeak(secret secret.Secret, source source.Source)
	AddError(filePath string, err error)
}

type ReportOutput interface {
	Output(filePath string, outputType string)
}
