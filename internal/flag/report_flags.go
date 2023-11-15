package flag

import (
	"errors"
	"strings"

	globaltypes "github.com/bearer/bearer/internal/types"
	"github.com/bearer/bearer/internal/util/set"
	sliceutil "github.com/bearer/bearer/internal/util/slices"
)

var (
	FormatReviewDog  = "rdjson"
	FormatGitLabSast = "gitlab-sast"
	FormatSarif      = "sarif"
	FormatJSON       = "json"
	FormatJSONV2     = "jsonv2"
	FormatYAML       = "yaml"
	FormatHTML       = "html"
	FormatCSV        = "csv"
	FormatEmpty      = ""

	ReportPrivacy   = "privacy"
	ReportSecurity  = "security"
	ReportDataFlow  = "dataflow"
	ReportDetectors = "detectors" // nodoc: internal report type
	ReportSaaS      = "saas"      // nodoc: internal report type
	ReportStats     = "stats"     // nodoc: internal report type
)

var (
	ErrInvalidFormatSecurity = errors.New("invalid format argument for security report; supported values: json, yaml, sarif, gitlab-sast, rdjson, html, jsonv2")
	ErrInvalidFormatPrivacy  = errors.New("invalid format argument for privacy report; supported values: csv, json, yaml, html")
	ErrInvalidFormatDefault  = errors.New("invalid format argument; supported values: json, yaml")
	ErrInvalidReport         = errors.New("invalid report argument; supported values: security, privacy")
	ErrInvalidSeverity       = errors.New("invalid severity argument; supported values: " + strings.Join(globaltypes.Severities, ", "))
	ErrInvalidFailOnSeverity = errors.New("invalid fail-on-severity argument; supported values: " + strings.Join(globaltypes.Severities, ", "))
)

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
		Value:      strings.Join(globaltypes.Severities, ","),
		Usage:      "Specify which severities are included in the report.",
	}
	FailOnSeverityFlag = Flag{
		Name:       "fail-on-severity",
		ConfigName: "report.fail-on-severity",
		Value:      strings.Join(sliceutil.Except(globaltypes.Severities, globaltypes.LevelWarning), ","),
		Usage:      "Specify which severities cause the report to fail. Works in conjunction with --exit-code.",
	}
	ExcludeFingerprintFlag = Flag{
		Name:            "exclude-fingerprint",
		ConfigName:      "report.exclude-fingerprint",
		Value:           []string{},
		Usage:           "Specify the comma-separated fingerprints of the findings you would like to exclude from the report.",
		DisableInConfig: true,
		Hide:            true,
		Deprecated:      true,
	}
)

type ReportFlagGroup struct {
	Format             *Flag
	Report             *Flag
	Output             *Flag
	Severity           *Flag
	FailOnSeverity     *Flag
	ExcludeFingerprint *Flag
}

type ReportOptions struct {
	Format             string          `mapstructure:"format" json:"format" yaml:"format"`
	Report             string          `mapstructure:"report" json:"report" yaml:"report"`
	Output             string          `mapstructure:"output" json:"output" yaml:"output"`
	Severity           set.Set[string] `mapstructure:"severity" json:"severity" yaml:"severity"`
	FailOnSeverity     set.Set[string] `mapstructure:"fail-on-severity" json:"fail-on-severity" yaml:"fail-on-severity"`
	ExcludeFingerprint map[string]bool `mapstructure:"exclude_fingerprints" json:"exclude_fingerprints" yaml:"exclude_fingerprints"`
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format:             &FormatFlag,
		Report:             &ReportFlag,
		Output:             &OutputFlag,
		Severity:           &SeverityFlag,
		FailOnSeverity:     &FailOnSeverityFlag,
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
		f.FailOnSeverity,
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
	case FormatCSV:
		if report != ReportPrivacy {
			return ReportOptions{}, invalidFormat
		}
	case FormatSarif, FormatGitLabSast, FormatReviewDog, FormatJSONV2:
		if report != ReportSecurity {
			return ReportOptions{}, invalidFormat
		}
	default:
		return ReportOptions{}, invalidFormat
	}

	severity := getSeverities(f.Severity)
	if severity == nil {
		return ReportOptions{}, ErrInvalidSeverity
	}
	failOnSeverity := getSeverities(f.FailOnSeverity)
	if failOnSeverity == nil {
		return ReportOptions{}, ErrInvalidFailOnSeverity
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
		Severity:           severity,
		FailOnSeverity:     failOnSeverity,
		ExcludeFingerprint: excludeFingerprintsMapping,
	}, nil
}
