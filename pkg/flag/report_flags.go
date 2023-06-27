package flag

import (
	"errors"

	"github.com/bearer/bearer/pkg/types"
)

var (
	FormatReviewDog  = "rdjson"
	FormatGitLabSast = "gitlab-sast"
	FormatSarif      = "sarif"
	FormatJSON       = "json"
	FormatYAML       = "yaml"
	FormatHTML       = "html"
	FormatEmpty      = ""

	ReportPrivacy   = "privacy"
	ReportSecurity  = "security"
	ReportDataFlow  = "dataflow"
	ReportDetectors = "detectors" // nodoc: internal report type
	ReportSaaS      = "saas"      // nodoc: internal report type
	ReportStats     = "stats"     // nodoc: internal report type

	DefaultSeverity = "critical,high,medium,low,warning"
)

var ErrInvalidFormatSecurity = errors.New("invalid format argument for security report; supported values: json, yaml, sarif, gitlab-sast, rdjson, html")
var ErrInvalidFormatPrivacy = errors.New("invalid format argument for privacy report; supported values: json, yaml html")
var ErrInvalidFormatDefault = errors.New("invalid format argument; supported values: json, yaml")
var ErrInvalidReport = errors.New("invalid report argument; supported values: security, privacy")
var ErrInvalidSeverity = errors.New("invalid severity argument; supported values: critical, high, medium, low, warning")

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml, sarif, gitlab-sast, rdjson, html)",
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
	ExcludeFingerprintFlag = Flag{
		Name:       "exclude-fingerprint",
		ConfigName: "report.exclude-fingerprint",
		Value:      []string{},
		Usage:      "Specify the comma-separated fingerprints of the findings you would like to exclude from the report.",
	}
)

type ReportFlagGroup struct {
	Format             *Flag
	Report             *Flag
	Output             *Flag
	Severity           *Flag
	ExcludeFingerprint *Flag
}

type ReportOptions struct {
	Format             string          `mapstructure:"format" json:"format" yaml:"format"`
	Report             string          `mapstructure:"report" json:"report" yaml:"report"`
	Output             string          `mapstructure:"output" json:"output" yaml:"output"`
	Severity           map[string]bool `mapstructure:"severity" json:"severity" yaml:"severity"`
	ExcludeFingerprint map[string]bool `mapstructure:"exclude_fingerprints" json:"exclude_fingerprints" yaml:"exclude_fingerprints"`
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format:             &FormatFlag,
		Report:             &ReportFlag,
		Output:             &OutputFlag,
		Severity:           &SeverityFlag,
		ExcludeFingerprint: &ExcludeFingerprintFlag,
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
		f.ExcludeFingerprint,
	}
}

func (f *ReportFlagGroup) ToOptions() (ReportOptions, error) {
	invalidFormat := ErrInvalidFormatDefault
	report := getString(f.Report)
	switch report {
	case ReportPrivacy:
		invalidFormat = ErrInvalidFormatPrivacy
	case ReportSecurity:
		invalidFormat = ErrInvalidFormatSecurity
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
	case FormatHTML:
		if report != ReportPrivacy && report != ReportSecurity {
			return ReportOptions{}, invalidFormat
		}
	case FormatSarif, FormatGitLabSast, FormatReviewDog:
		if report != ReportSecurity {
			return ReportOptions{}, invalidFormat
		}
	default:
		return ReportOptions{}, invalidFormat
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

	// turn string slice into map for ease of access
	excludeFingerprints := getStringSlice(f.ExcludeFingerprint)
	excludeFingerprintsMapping := make(map[string]bool)
	for _, fingerprint := range excludeFingerprints {
		excludeFingerprintsMapping[fingerprint] = true
	}

	return ReportOptions{
		Format:             format,
		Report:             report,
		Output:             getString(f.Output),
		Severity:           severityMapping,
		ExcludeFingerprint: excludeFingerprintsMapping,
	}, nil
}
