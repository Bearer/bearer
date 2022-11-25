package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"

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
	outputPolicies, err := getPolicyReportOutput(report, config)
	if err != nil {
		return err
	}

	outputStr := &strings.Builder{}
	outputStr.WriteString("===============================")

	for _, policyBreach := range outputPolicies[settings.LevelCritical] {
		writePolicyBreachToOutput(outputStr, policyBreach, settings.LevelCritical)
	}

	for _, policyBreach := range outputPolicies[settings.LevelHigh] {
		writePolicyBreachToOutput(outputStr, policyBreach, settings.LevelHigh)
	}

	for _, policyBreach := range outputPolicies[settings.LevelMedium] {
		writePolicyBreachToOutput(outputStr, policyBreach, settings.LevelMedium)
	}

	for _, policyBreach := range outputPolicies[settings.LevelLow] {
		writePolicyBreachToOutput(outputStr, policyBreach, settings.LevelLow)
	}

	output.Msg(outputStr.String())

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

func writePolicyBreachToOutput(outputStr *strings.Builder, policyBreach policies.PolicyResult, policySeverity string) {
	outputStr.WriteString("\n\n")
	outputStr.WriteString(formatSeverity(policySeverity))
	outputStr.WriteString(policyBreach.PolicyName + " policy breach with " + policyBreach.CategoryGroup + "\n")
	outputStr.WriteString(color.HiBlackString(policyBreach.PolicyDescription + "\n"))
	outputStr.WriteString("\n")
	outputStr.WriteString(color.HiBlueString("File: " + underline(policyBreach.Filename+":"+fmt.Sprint(policyBreach.LineNumber)) + "\n"))
	outputStr.WriteString("\n")
	outputStr.WriteString(highlightCodeExtract(policyBreach.LineNumber, policyBreach.ParentLineNumber, policyBreach.ParentContent))
	outputStr.WriteString("\n\n")
	outputStr.WriteString("=====================================")
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
