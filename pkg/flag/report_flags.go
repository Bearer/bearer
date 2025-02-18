package flag

import (
	"errors"
	"strings"

	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	globaltypes "github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/set"
	sliceutil "github.com/bearer/bearer/pkg/util/slices"
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

type reportFlagGroup struct{ flagGroupBase }

var ReportFlagGroup = &reportFlagGroup{flagGroupBase{name: "Report"}}

var (
	FormatFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      FormatEmpty,
		Usage:      "Specify report format (json, yaml, sarif, gitlab-sast, rdjson, html)",
	})
	ReportFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "report",
		ConfigName: "report.report",
		Value:      ReportSecurity,
		Usage:      "Specify the type of report (security, privacy, dataflow).",
	})
	OutputFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "output",
		ConfigName: "report.output",
		Value:      "",
		Usage:      "Specify the output path for the report.",
	})
	SeverityFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "severity",
		ConfigName: "report.severity",
		Value:      strings.Join(globaltypes.Severities, ","),
		Usage:      "Specify which severities are included in the report.",
	})
	FailOnSeverityFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "fail-on-severity",
		ConfigName: "report.fail-on-severity",
		Value:      strings.Join(sliceutil.Except(globaltypes.Severities, globaltypes.LevelWarning), ","),
		Usage:      "Specify which severities cause the report to fail. Works in conjunction with --exit-code.",
	})
	ExcludeFingerprintFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:            "exclude-fingerprint",
		ConfigName:      "report.exclude-fingerprint",
		Value:           []string{},
		Usage:           "Specify the comma-separated fingerprints of the findings you would like to exclude from the report.",
		DisableInConfig: true,
		Hide:            true,
		Deprecated:      true,
	})
	NoExtractFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "no-extract",
		ConfigName: "report.no-extract",
		Value:      false,
		Usage:      "Do not include code extract in report.",
	})
	NoRuleMetaFlag = ReportFlagGroup.add(flagtypes.Flag{
		Name:       "no-rule-meta",
		ConfigName: "report.no-rule-meta",
		Value:      false,
		Usage:      "Do not include rule description content.",
	})
)

type ReportOptions struct {
	Format             string          `mapstructure:"format" json:"format" yaml:"format"`
	Report             string          `mapstructure:"report" json:"report" yaml:"report"`
	Output             string          `mapstructure:"output" json:"output" yaml:"output"`
	Severity           set.Set[string] `mapstructure:"severity" json:"severity" yaml:"severity"`
	FailOnSeverity     set.Set[string] `mapstructure:"fail-on-severity" json:"fail-on-severity" yaml:"fail-on-severity"`
	ExcludeFingerprint map[string]bool `mapstructure:"exclude_fingerprints" json:"exclude_fingerprints" yaml:"exclude_fingerprints"`
}

func (reportFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	invalidFormat := ErrInvalidFormatDefault
	report := getString(ReportFlag)
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
		return ErrInvalidReport
	}

	format := getString(FormatFlag)
	switch format {
	case FormatYAML:
	case FormatJSON:
	case FormatEmpty:
	case FormatHTML:
		if report != ReportPrivacy && report != ReportSecurity {
			return invalidFormat
		}
	case FormatCSV:
		if report != ReportPrivacy {
			return invalidFormat
		}
	case FormatSarif, FormatGitLabSast, FormatReviewDog, FormatJSONV2:
		if report != ReportSecurity {
			return invalidFormat
		}
	default:
		return invalidFormat
	}

	severity := getSeverities(SeverityFlag)
	if severity == nil {
		return ErrInvalidSeverity
	}
	failOnSeverity := getSeverities(FailOnSeverityFlag)
	if failOnSeverity == nil {
		return ErrInvalidFailOnSeverity
	}

	// turn string slice into map for ease of access
	excludeFingerprints := getStringSlice(ExcludeFingerprintFlag)
	excludeFingerprintsMapping := make(map[string]bool)
	for _, fingerprint := range excludeFingerprints {
		excludeFingerprintsMapping[fingerprint] = true
	}

	options.ReportOptions = flagtypes.ReportOptions{
		Format:             format,
		Report:             report,
		Output:             getString(OutputFlag),
		Severity:           severity,
		FailOnSeverity:     failOnSeverity,
		ExcludeFingerprint: excludeFingerprintsMapping,
		NoExtract:          getBool(NoExtractFlag),
		NoRuleMeta:         getBool(NoRuleMetaFlag),
	}

	return nil
}
