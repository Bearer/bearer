package custom

import (
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
)

type ReportWrapper struct {
	RuleName string
	Report   report.Report
}

func WrapReport(report report.Report, ruleName string) *ReportWrapper {
	return &ReportWrapper{
		RuleName: ruleName,
		Report:   report,
	}
}

func (reportWrapper *ReportWrapper) AddSchema(detectorType detectors.Type, schema schema.Schema, source source.Source) {
	reportWrapper.Report.AddCustomDetection(reportWrapper.RuleName, source, schema)
}
