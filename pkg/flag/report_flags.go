package flag

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/xerrors"

	"github.com/bearer/curio/pkg/report/output"
)

const (
	SeverityUnknown Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

var (
	SeverityNames = []string{
		"UNKNOWN",
		"LOW",
		"MEDIUM",
		"HIGH",
		"CRITICAL",
	}
)

type Severity int

// e.g. config yaml:
//
//	format: table
//	dependency-tree: true
//	severity: HIGH,CRITICAL
var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "format",
		Shorthand:  "f",
		Value:      output.TypeJSONLines,
		Usage:      "format (table, json, jsonline)",
	}
	// ReportFormatFlag = Flag{
	// 	Name:       "report",
	// 	ConfigName: "report",
	// 	Value:      "all",
	// 	Usage:      "specify a report format for the output. (all,summary)",
	// }
	IgnoreFileFlag = Flag{
		Name:       "ignorefile",
		ConfigName: "ignorefile",
		Value:      "",
		Usage:      "specify .curioignore file",
	}
	// IgnorePolicyFlag = Flag{
	// 	Name:       "ignore-policy",
	// 	ConfigName: "ignore-policy",
	// 	Value:      "",
	// 	Usage:      "specify the Rego file path to evaluate each vulnerability",
	// }
	ExitCodeFlag = Flag{
		Name:       "exit-code",
		ConfigName: "exit-code",
		Value:      0,
		Usage:      "specify exit code when any security issues are found",
	}
	OutputFlag = Flag{
		Name:       "output",
		ConfigName: "output",
		Shorthand:  "o",
		Value:      "",
		Usage:      "output file name",
	}
)

// ReportFlagGroup composes common printer flag structs
// used for commands requiring reporting logic.
type ReportFlagGroup struct {
	Format *Flag
	// ReportFormat   *Flag
	// IgnorePolicy   *Flag
	ExitCode   *Flag
	Output     *Flag
	IgnoreFile *Flag
	// Severity *Flag
}

type ReportOptions struct {
	Format string
	// ReportFormat   string
	IgnoreFile string
	ExitCode   int
	// IgnorePolicy string
	Output io.Writer
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format: &FormatFlag,
		// ReportFormat:   &ReportFormatFlag,
		// IgnorePolicy:   &IgnorePolicyFlag,
		ExitCode: &ExitCodeFlag,
		Output:   &OutputFlag,
	}
}

func NewSeverity(severity string) (Severity, error) {
	for i, name := range SeverityNames {
		if severity == name {
			return Severity(i), nil
		}
	}
	return SeverityUnknown, fmt.Errorf("unknown severity: %s", severity)
}

func (f *ReportFlagGroup) Name() string {
	return "Report"
}

func (f *ReportFlagGroup) Flags() []*Flag {
	return []*Flag{f.Format,
		// f.ReportFormat,
		f.ExitCode, f.Output, f.IgnoreFile,
	}
}

func (f *ReportFlagGroup) ToOptions(out io.Writer) (ReportOptions, error) {
	format := getString(f.Format)
	output := getString(f.Output)

	if output != "" {
		var err error
		if out, err = os.Create(output); err != nil {
			return ReportOptions{}, xerrors.Errorf("failed to create an output file: %w", err)
		}
	}

	return ReportOptions{
		Format: format,
		// ReportFormat:   getString(f.ReportFormat),
		IgnoreFile: getString(f.IgnoreFile),
		ExitCode:   getInt(f.ExitCode),
		Output:     out,
	}, nil
}

// func splitSeverity(severity []string) []Severity {
// 	switch {
// 	case len(severity) == 0:
// 		return nil
// 	case len(severity) == 1 && strings.Contains(severity[0], ","): // get severities from flag
// 		severity = strings.Split(severity[0], ",")
// 	}

// 	var severities []Severity
// 	for _, s := range severity {
// 		sev, err := NewSeverity(strings.ToUpper(s))
// 		if err != nil {
// 			// log.Logger.Warnf("unknown severity option: %s", err)
// 			continue
// 		}
// 		severities = append(severities, sev)
// 	}
// 	// log.Logger.Debugf("Severities: %q", severities)
// 	return severities
// }
