package flag

import (
	"errors"
)

type Severity int

var (
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatEmpty = ""

	ReportDetectors = "detectors" // nodoc: internal report type
	ReportDataFlow  = "dataflow"
	ReportSummary   = "summary"
	ReportStats     = "stats"
)

var ErrInvalidFormat = errors.New("invalid format argument; supported values: json, yaml")
var ErrInvalidReport = errors.New("invalid report argument; supported values: summary, dataflow, stats")

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml)",
	}
	ReportFlag = Flag{
		Name:       "report",
		ConfigName: "report.report",
		Value:      ReportSummary,
		Usage:      "Specify the type of report (summary, dataflow, stats).",
	}
	OutputFlag = Flag{
		Name:       "output",
		ConfigName: "report.output",
		Value:      "",
		Usage:      "Specify the output path for the report.",
	}
)

type ReportFlagGroup struct {
	Format *Flag
	Report *Flag
	Output *Flag
}

type ReportOptions struct {
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	Report string `mapstructure:"report" json:"report" yaml:"report"`
	Output string `mapstructure:"output" json:"output" yaml:"output"`
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

func (f *ReportFlagGroup) ToOptions() (ReportOptions, error) {
	format := getString(f.Format)
	switch format {
	case FormatYAML:
	case FormatJSON:
	case FormatEmpty:
	default:
		return ReportOptions{}, ErrInvalidFormat
	}

	report := getString(f.Report)
	switch report {
	case ReportDetectors:
	case ReportDataFlow:
	case ReportSummary:
	case ReportStats:
	default:
		return ReportOptions{}, ErrInvalidReport
	}

	return ReportOptions{
		Format: format,
		Report: report,
		Output: getString(f.Output),
	}, nil
}
