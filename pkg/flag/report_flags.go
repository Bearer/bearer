package flag

import (
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
		Usage:      "specify the kind of report (detectors, dataflow)",
	}
	OutputFlag = Flag{
		Name:       "output",
		ConfigName: "report.output",
		Value:      "",
		Usage:      "path where to save report",
	}
)

type ReportFlagGroup struct {
	Format *Flag
	Report *Flag
	Output *Flag
}

type ReportOptions struct {
	Format string
	Report string
	Output string
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format: &FormatFlag,
		Report: &ReportFlag,
		Output: &OutputFlag,
	}
}

func (f *ReportFlagGroup) Name() string {
	return "Report"
}

func (f *ReportFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.Format,
		f.Report,
		f.Output,
	}
}

func (f *ReportFlagGroup) ToOptions() ReportOptions {
	return ReportOptions{
		Format: getString(f.Format),
		Report: getString(f.Report),
		Output: getString(f.Output),
	}
}
