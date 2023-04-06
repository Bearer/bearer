package flag

import (
	"errors"

	"github.com/bearer/bearer/pkg/types"
)

var (
	FormatSarif = "sarif"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatEmpty = ""

	ReportPrivacy   = "privacy"
	ReportSecurity  = "security"
	ReportDataFlow  = "dataflow"
	ReportDetectors = "detectors" // nodoc: internal report type
	ReportSaaS      = "saas"      // nodoc: internal report type
	ReportStats     = "stats"     // nodoc: internal report type

	DefaultSeverity = "critical,high,medium,low,warning"
)

var ErrInvalidFormat = errors.New("invalid format argument; supported values: json, yaml, sarif")
var ErrInvalidReport = errors.New("invalid report argument; supported values: security, privacy")
var ErrInvalidSeverity = errors.New("invalid severity argument; supported values: critical, high, medium, low, warning")

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml, sarif)",
	}
	ReportFlag = Flag{
		Name:       "report",
		ConfigName: "report.report",
		Value:      ReportSecurity,
		Usage:      "Specify the type of report (security, privacy, dataflow).",
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
	report := getString(f.Report)
	switch report {
	case ReportPrivacy:
	case ReportSecurity:
	case ReportDataFlow:
	// hidden flags for development use
	case ReportDetectors:
	case ReportSaaS:
	case ReportStats:
	default:
		return ReportOptions{}, ErrInvalidReport
	}

	format := getString(f.Format)
	switch format {
	case FormatYAML:
	case FormatJSON:
	case FormatEmpty:
	case FormatSarif:
		if report != ReportSecurity {
			return ReportOptions{}, ErrInvalidFormat
		}
	default:
		return ReportOptions{}, ErrInvalidFormat
	}

	severity := getStringSlice(f.Severity)
	severityMapping := make(map[string]bool)

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
		case types.LevelWarning:
			severityMapping[types.LevelWarning] = true
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
