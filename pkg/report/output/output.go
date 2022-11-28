package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"golang.org/x/exp/maps"

	"github.com/bearer/curio/pkg/report/output/detectors"
	"github.com/bearer/curio/pkg/report/output/policies"
	"github.com/bearer/curio/pkg/report/output/stats"
	"github.com/bearer/curio/pkg/types"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
)

var ErrUndefinedFormat = errors.New("undefined output format")
var underline = color.New(color.Underline).SprintFunc()
var severityColorFns = map[string]func(x ...interface{}) string{
	settings.LevelCritical: color.New(color.FgRed).SprintFunc(),
	settings.LevelHigh:     color.New(color.FgHiRed).SprintFunc(),
	settings.LevelMedium:   color.New(color.FgYellow).SprintFunc(),
	settings.LevelLow:      color.New(color.FgBlue).SprintFunc(),
}

func ReportPolicies(report types.Report, output *zerolog.Event, config settings.Config) error {
	policyResults, err := getPolicyReportOutput(report, config)
	if err != nil {
		return err
	}

	reportStr := &strings.Builder{}
	reportStr.WriteString("\n\nPolicy Report\n")
	reportStr.WriteString("\n=====================================")

	writePolicyListToString(reportStr, config.Policies)

	breachedPolicies := map[string]map[string]bool{
		settings.LevelCritical: make(map[string]bool),
		settings.LevelHigh:     make(map[string]bool),
		settings.LevelMedium:   make(map[string]bool),
		settings.LevelLow:      make(map[string]bool),
	}

	for _, policyLevel := range []string{
		settings.LevelCritical,
		settings.LevelHigh,
		settings.LevelMedium,
		settings.LevelLow,
	} {
		for _, policyBreach := range policyResults[policyLevel] {
			breachedPolicies[policyLevel][policyBreach.PolicyName] = true
			writePolicyBreachToString(reportStr, policyBreach, policyLevel)
		}
	}

	writeSummaryToString(reportStr, policyResults, len(config.Policies), breachedPolicies)

	output.Msg(reportStr.String())

	return nil
}

func ReportJSON(report types.Report, output *zerolog.Event, config settings.Config) error {
	outputDetections, err := getReportOutput(report, config)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
}

func ReportYAML(report types.Report, output *zerolog.Event, config settings.Config) error {
	ouputDetections, err := getReportOutput(report, config)
	if err != nil {
		return err
	}

	jsonBytes, err := yaml.Marshal(&ouputDetections)
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
}

func getReportOutput(report types.Report, config settings.Config) (any, error) {
	switch config.Report.Report {
	case flag.ReportDetectors:
		return detectors.GetOutput(report)
	case flag.ReportDataFlow:
		return getDataflow(report, config, false)
	case flag.ReportPolicies:
		dataflow, err := getDataflow(report, config, true)
		if err != nil {
			return nil, err
		}

		return policies.GetOutput(dataflow, config)
	case flag.ReportStats:
		lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
		if err != nil {
			return nil, err
		}

		detectorsOutput, err := detectors.GetOutput(report)
		if err != nil {
			return nil, err
		}

		dataflowOutput, err := dataflow.GetOutput(detectorsOutput, config, true)
		if err != nil {
			return nil, err
		}

		return stats.GetOutput(lineOfCodeOutput, dataflowOutput, config)
	}

	return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getPolicyReportOutput(report types.Report, config settings.Config) (map[string][]policies.PolicyResult, error) {
	detections, err := detectors.GetOutput(report)
	if err != nil {
		return nil, err
	}

	dataflow, err := dataflow.GetOutput(detections, config, true)
	if err != nil {
		return nil, err
	}

	return policies.GetOutput(dataflow, config)
}

func getDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, error) {
	reportedDetections, err := detectors.GetOutput(report)
	if err != nil {
		return nil, err
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func writePolicyListToString(reportStr *strings.Builder, policies map[string]*settings.Policy) {
	// list policies that were run
	reportStr.WriteString("\nPolicy list: \n\n")
	for key := range policies {
		policy := policies[key]
		reportStr.WriteString(color.HiBlackString("- " + policy.Name + "\n"))
	}
}

func writeSummaryToString(
	reportStr *strings.Builder,
	policyResults map[string][]policies.PolicyResult,
	policyCount int, breachedPolicies map[string]map[string]bool,
) {
	reportStr.WriteString("\n=====================================")

	// give summary including counts
	if len(policyResults) == 0 {
		reportStr.WriteString("\n\n")
		reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
		reportStr.WriteString(fmt.Sprint(policyCount) + " policies were run and no breaches were detected.\n\n")
		return
	}

	criticalCount := len(policyResults[settings.LevelCritical])
	highCount := len(policyResults[settings.LevelHigh])
	mediumCount := len(policyResults[settings.LevelMedium])
	lowCount := len(policyResults[settings.LevelLow])

	totalCount := criticalCount + highCount + mediumCount + lowCount

	reportStr.WriteString("\n\n")
	reportStr.WriteString(color.RedString("Policy breaches detected\n\n"))
	reportStr.WriteString(fmt.Sprint(policyCount) + " policies were run ")
	reportStr.WriteString("and " + fmt.Sprint(totalCount) + " breaches were detected.\n\n")

	// critical count
	reportStr.WriteString(formatSeverity(settings.LevelCritical) + fmt.Sprint(criticalCount))
	if len(breachedPolicies[settings.LevelCritical]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(breachedPolicies[settings.LevelCritical]), ", ") + ")")
	}
	// high count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelHigh) + fmt.Sprint(highCount))
	if len(breachedPolicies[settings.LevelHigh]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(breachedPolicies[settings.LevelHigh]), ", ") + ")")
	}
	// medium count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelMedium) + fmt.Sprint(mediumCount))
	if len(breachedPolicies[settings.LevelMedium]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(breachedPolicies[settings.LevelMedium]), ", ") + ")")
	}
	// low count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelLow) + fmt.Sprint(lowCount))
	if len(breachedPolicies[settings.LevelLow]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(breachedPolicies[settings.LevelLow]), ", ") + ")")
	}

	reportStr.WriteString("\n\n")
}

func writePolicyBreachToString(reportStr *strings.Builder, policyBreach policies.PolicyResult, policySeverity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(policySeverity))
	reportStr.WriteString(policyBreach.PolicyName + " policy breach with " + policyBreach.CategoryGroup + "\n")
	reportStr.WriteString(color.HiBlackString(policyBreach.PolicyDescription + "\n"))
	reportStr.WriteString("\n")
	reportStr.WriteString(color.HiBlueString("File: " + underline(policyBreach.Filename+":"+fmt.Sprint(policyBreach.LineNumber)) + "\n"))
	reportStr.WriteString("\n")
	reportStr.WriteString(highlightCodeExtract(policyBreach.LineNumber, policyBreach.ParentLineNumber, policyBreach.ParentContent))
}

func formatSeverity(policySeverity string) string {
	severityColorFn, ok := severityColorFns[policySeverity]
	if !ok {
		return strings.ToUpper(policySeverity)
	}
	return severityColorFn(strings.ToUpper(policySeverity + ": "))
}

func highlightCodeExtract(lineNumber int, extractStartLineNumber int, extract string) string {
	result := ""
	targetIndex := lineNumber - extractStartLineNumber
	for index, line := range strings.Split(extract, "\n") {
		if index == targetIndex {
			result += color.MagentaString(" " + fmt.Sprint(extractStartLineNumber+index) + " ")
			result += color.MagentaString(line) + "\n"
		} else {
			result += " " + fmt.Sprint(extractStartLineNumber+index) + " "
			result += line + "\n"
		}
	}

	return result
}
