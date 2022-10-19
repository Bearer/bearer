package flag

import (
	"io"
)

type Severity int

var (
	FormatJSON      = "json"
	FormatJSONLines = "jsonlines"

	ReportDetectors = "detectors"
	ReportDataFlow  = "dataflow"
)

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      FormatJSON,
		Usage:      "format (json)",
	}
	ReportFlag = Flag{
		Name:       "report",
		ConfigName: "report.report",
		Value:      ReportDetectors,
		Usage:      "specify the kind of report (detectors)",
	}
)

type ReportFlagGroup struct {
	Format *Flag
	Report *Flag
}

type ReportOptions struct {
	Format string
	Report string
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format: &FormatFlag,
		Report: &ReportFlag,
	}
}

func (f *ReportFlagGroup) Name() string {
	return "Report"
}

func (f *ReportFlagGroup) Flags() []*Flag {
	return []*Flag{f.Format, f.Report}
}

func (f *ReportFlagGroup) ToOptions(out io.Writer) ReportOptions {
	return ReportOptions{
		Format: getString(f.Format),
		Report: getString(f.Report),
	}
}
