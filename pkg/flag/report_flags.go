package flag

import (
	"errors"

	"github.com/bearer/curio/pkg/types"
)

var (
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatEmpty = ""

	ReportPrivacy   = "privacy"
	ReportSummary   = "summary"
	ReportDetectors = "detectors" // nodoc: internal report type
	ReportDataFlow  = "dataflow"  // nodoc: internal report type
	ReportStats     = "stats"     // nodoc: internal report type

	DefaultSeverity = "critical,high,medium,low"
)

var ErrInvalidFormat = errors.New("invalid format argument; supported values: json, yaml")
var ErrInvalidReport = errors.New("invalid report argument; supported values: summary, privacy")
var ErrInvalidSeverity = errors.New("invalid severity argument; supported values: critical, high, medium, low")

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
		Usage:      "Specify the type of report (summary, privacy).",
	}
	OutputFlag = Flag{
		Name:       "output",
		ConfigName: "report.output",
		Value:      "",
		Usage:      "Specify the output path for the report.",
	}
	SeverityFlag = Flag{
		Name:       "severity",
		ConfigName: "report.severity",
		Value:      DefaultSeverity,
		Usage:      "Specify which severities are included in the report.",
	}
)

type ReportFlagGroup struct {
	Format   *Flag
	Report   *Flag
	Output   *Flag
	Severity *Flag
}

type ReportOptions struct {
	Format   string          `mapstructure:"format" json:"format" yaml:"format"`
	Report   string          `mapstructure:"report" json:"report" yaml:"report"`
	Output   string          `mapstructure:"output" json:"output" yaml:"output"`
	Severity map[string]bool `mapstructure:"severity" json:"severity" yaml:"severity"`
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format:   &FormatFlag,
		Report:   &ReportFlag,
		Output:   &OutputFlag,
		Severity: &SeverityFlag,
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
		f.Severity,
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
	case ReportPrivacy:
	case ReportSummary:
	// hidden flags for development use
	case ReportDetectors:
	case ReportDataFlow:
	case ReportStats:
	default:
		return ReportOptions{}, ErrInvalidReport
	}

	severity := getStringSlice(f.Severity)
	// pre-define mapping to ensure ordering
	severityMapping := map[string]bool{
		types.LevelCritical: false,
		types.LevelHigh:     false,
		types.LevelMedium:   false,
		types.LevelLow:      false,
	}

	for _, severityLevel := range severity {
		switch severityLevel {
		case types.LevelCritical:
			severityMapping[types.LevelCritical] = true
		case types.LevelHigh:
			severityMapping[types.LevelHigh] = true
		case types.LevelMedium:
			severityMapping[types.LevelMedium] = true
		case types.LevelLow:
			severityMapping[types.LevelLow] = true
		default:
			return ReportOptions{}, ErrInvalidSeverity
		}
	}

	return ReportOptions{
		Format:   format,
		Report:   report,
		Output:   getString(f.Output),
		Severity: severityMapping,
	}, nil
}
