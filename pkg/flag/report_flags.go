package flag

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		Usage:      "format (table, json, sarif, template, cyclonedx, spdx, spdx-json, github, cosign-vuln)",
	}
	ReportFormatFlag = Flag{
		Name:       "report",
		ConfigName: "report",
		Value:      "all",
		Usage:      "specify a report format for the output. (all,summary)",
	}
	TemplateFlag = Flag{
		Name:       "template",
		ConfigName: "template",
		Shorthand:  "t",
		Value:      "",
		Usage:      "output template",
	}
	DependencyTreeFlag = Flag{
		Name:       "dependency-tree",
		ConfigName: "dependency-tree",
		Value:      false,
		Usage:      "show dependency origin tree (EXPERIMENTAL)",
	}
	ListAllPkgsFlag = Flag{
		Name:       "list-all-pkgs",
		ConfigName: "list-all-pkgs",
		Value:      false,
		Usage:      "enabling the option will output all packages regardless of vulnerability",
	}
	IgnoreFileFlag = Flag{
		Name:       "ignorefile",
		ConfigName: "ignorefile",
		Value:      "",
		Usage:      "specify .trivyignore file",
	}
	IgnorePolicyFlag = Flag{
		Name:       "ignore-policy",
		ConfigName: "ignore-policy",
		Value:      "",
		Usage:      "specify the Rego file path to evaluate each vulnerability",
	}
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
	SeverityFlag = Flag{
		Name:       "severity",
		ConfigName: "severity",
		Shorthand:  "s",
		Usage:      "severities of security issues to be displayed (comma separated)",
	}
)

// ReportFlagGroup composes common printer flag structs
// used for commands requiring reporting logic.
type ReportFlagGroup struct {
	Format         *Flag
	ReportFormat   *Flag
	Template       *Flag
	DependencyTree *Flag
	ListAllPkgs    *Flag
	IgnoreFile     *Flag
	IgnorePolicy   *Flag
	ExitCode       *Flag
	Output         *Flag
	Severity       *Flag
}

type ReportOptions struct {
	Format         string
	ReportFormat   string
	Template       string
	DependencyTree bool
	ListAllPkgs    bool
	IgnoreFile     string
	ExitCode       int
	IgnorePolicy   string
	Output         io.Writer
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format:         &FormatFlag,
		ReportFormat:   &ReportFormatFlag,
		Template:       &TemplateFlag,
		DependencyTree: &DependencyTreeFlag,
		ListAllPkgs:    &ListAllPkgsFlag,
		IgnoreFile:     &IgnoreFileFlag,
		IgnorePolicy:   &IgnorePolicyFlag,
		ExitCode:       &ExitCodeFlag,
		Output:         &OutputFlag,
		Severity:       &SeverityFlag,
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
	return []*Flag{f.Format, f.ReportFormat, f.Template, f.DependencyTree, f.ListAllPkgs, f.IgnoreFile,
		f.IgnorePolicy, f.ExitCode, f.Output, f.Severity}
}

func (f *ReportFlagGroup) ToOptions(out io.Writer) (ReportOptions, error) {
	format := getString(f.Format)
	template := getString(f.Template)
	dependencyTree := getBool(f.DependencyTree)
	listAllPkgs := getBool(f.ListAllPkgs)
	output := getString(f.Output)

	if template != "" {
		if format == "" {
			// log.Logger.Warn("'--template' is ignored because '--format template' is not specified. Use '--template' option with '--format template' option.")
		} else if format != "template" {
			// log.Logger.Warnf("'--template' is ignored because '--format %s' is specified. Use '--template' option with '--format template' option.", format)
		}
	} else {
		// if format == report.FormatTemplate {
		// 	// log.Logger.Warn("'--format template' is ignored because '--template' is not specified. Specify '--template' option when you use '--format template'.")
		// }
	}

	// "--list-all-pkgs" option is unavailable with "--format table".
	// If user specifies "--list-all-pkgs" with "--format table", we should warn it.
	// if listAllPkgs && format == report.FormatTable {
	// 	// log.Logger.Warn(`"--list-all-pkgs" cannot be used with "--format table". Try "--format json" or other formats.`)
	// }

	// "--dependency-tree" option is available only with "--format table".
	if dependencyTree {
		// log.Logger.Infof(`"--dependency-tree" only shows dependencies for "package-lock.json" files`)
		// if format != report.FormatTable {
		// 	// log.Logger.Warn(`"--dependency-tree" can be used only with "--format table".`)
		// }
	}

	// Enable '--list-all-pkgs' if needed
	if f.forceListAllPkgs(format, listAllPkgs, dependencyTree) {
		listAllPkgs = true
	}

	if output != "" {
		var err error
		if out, err = os.Create(output); err != nil {
			return ReportOptions{}, xerrors.Errorf("failed to create an output file: %w", err)
		}
	}

	return ReportOptions{
		Format:         format,
		ReportFormat:   getString(f.ReportFormat),
		Template:       template,
		DependencyTree: dependencyTree,
		ListAllPkgs:    listAllPkgs,
		IgnoreFile:     getString(f.IgnoreFile),
		ExitCode:       getInt(f.ExitCode),
		IgnorePolicy:   getString(f.IgnorePolicy),
		Output:         out,
	}, nil
}

func (f *ReportFlagGroup) forceListAllPkgs(format string, listAllPkgs, dependencyTree bool) bool {
	// if slices.Contains(report.SupportedSBOMFormats, format) && !listAllPkgs {
	// 	// log.Logger.Debugf("%q automatically enables '--list-all-pkgs'.", report.SupportedSBOMFormats)
	// 	return true
	// }
	if dependencyTree && !listAllPkgs {
		// log.Logger.Debugf("'--dependency-tree' enables '--list-all-pkgs'.")
		return true
	}
	return false
}

func splitSeverity(severity []string) []Severity {
	switch {
	case len(severity) == 0:
		return nil
	case len(severity) == 1 && strings.Contains(severity[0], ","): // get severities from flag
		severity = strings.Split(severity[0], ",")
	}

	var severities []Severity
	for _, s := range severity {
		sev, err := NewSeverity(strings.ToUpper(s))
		if err != nil {
			// log.Logger.Warnf("unknown severity option: %s", err)
			continue
		}
		severities = append(severities, sev)
	}
	// log.Logger.Debugf("Severities: %q", severities)
	return severities
}
