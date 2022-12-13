package policies

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/bearer/curio/pkg/util/rego"
	"github.com/fatih/color"
	"golang.org/x/exp/maps"

	"github.com/bearer/curio/pkg/report/output/dataflow"
)

var underline = color.New(color.Underline).SprintFunc()
var severityColorFns = map[string]func(x ...interface{}) string{
	settings.LevelCritical: color.New(color.FgRed).SprintFunc(),
	settings.LevelHigh:     color.New(color.FgHiRed).SprintFunc(),
	settings.LevelMedium:   color.New(color.FgYellow).SprintFunc(),
	settings.LevelLow:      color.New(color.FgBlue).SprintFunc(),
}

type PolicyInput struct {
	PolicyId       string             `json:"policy_id" yaml:"policy_id"`
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type PolicyOutput struct {
	ParentLineNumber int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	LineNumber       int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename         string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups   []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Severity         string   `json:"severity,omitempty" yaml:"severity,omitempty"`
	OmitParent       bool     `json:"omit_parent" yaml:"omit_parent"`
}

type PolicyResult struct {
	PolicyName        string   `json:"policy_name" yaml:"policy_name"`
	PolicyDescription string   `json:"policy_description" yaml:"policy_description"`
	LineNumber        int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename          string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups    []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	ParentLineNumber  int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent     string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	OmitParent        bool     `json:"omit_parent" yaml:"omit_parent"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (map[string][]PolicyResult, error) {
	output.StdErrLogger().Msgf("Processing Policies")
	// policy results grouped by severity (critical, high, ...)
	result := make(map[string][]PolicyResult)

	// Ensure a deterministic order
	keys := maps.Keys(config.Policies)
	sort.Strings(keys)

	for _, key := range keys {
		policy := config.Policies[key]
		output.StdErrLogger().Msgf("Processing policy %s", policy.Name)

		// Create a prepared query that can be evaluated.
		rs, err := rego.RunQuery(policy.Query,
			PolicyInput{
				PolicyId:       policy.Id,
				Dataflow:       dataflow,
				DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
			},
			policy.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return nil, err
			}

			var policyResults map[string][]PolicyOutput
			err = json.Unmarshal(jsonRes, &policyResults)
			if err != nil {
				return nil, err
			}

			for _, policyOutput := range policyResults["policy_failure"] {
				policyResult := PolicyResult{
					PolicyName:        policy.Name,
					PolicyDescription: policy.Description,
					Filename:          policyOutput.Filename,
					LineNumber:        policyOutput.LineNumber,
					CategoryGroups:    policyOutput.CategoryGroups,
					OmitParent:        policyOutput.OmitParent,
					ParentLineNumber:  policyOutput.ParentLineNumber,
					ParentContent:     policyOutput.ParentContent,
				}

				result[policyOutput.Severity] = append(result[policyOutput.Severity], policyResult)
			}
		}

		output.StdErrLogger().Msgf("Finished processing policy %s", policy.Name)
	}

	output.StdErrLogger().Msgf("Finished processing policies")
	return result, nil
}

func BuildReportString(policyResults map[string][]PolicyResult, policies map[string]*settings.Policy, withoutColor bool) *strings.Builder {
	reportStr := &strings.Builder{}
	reportStr.WriteString("\n\nPolicy Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if withoutColor && !initialColorSetting {
		color.NoColor = true
	}

	writePolicyListToString(reportStr, policies)

	policyFailures := map[string]map[string]bool{
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
		for _, policyFailure := range policyResults[policyLevel] {
			policyFailures[policyLevel][policyFailure.PolicyName] = true
			writePolicyFailureToString(reportStr, policyFailure, policyLevel)
		}
	}

	writeSummaryToString(reportStr, policyResults, len(policies), policyFailures)

	color.NoColor = initialColorSetting

	return reportStr
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
	policyResults map[string][]PolicyResult,
	policyCount int, policyFailures map[string]map[string]bool,
) {
	reportStr.WriteString("\n=====================================")

	// give summary including counts
	if len(policyResults) == 0 {
		reportStr.WriteString("\n\n")
		reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
		reportStr.WriteString(fmt.Sprint(policyCount) + " policies were run and no failures were detected.\n\n")
		return
	}

	criticalCount := len(policyResults[settings.LevelCritical])
	highCount := len(policyResults[settings.LevelHigh])
	mediumCount := len(policyResults[settings.LevelMedium])
	lowCount := len(policyResults[settings.LevelLow])

	totalCount := criticalCount + highCount + mediumCount + lowCount

	reportStr.WriteString("\n\n")
	reportStr.WriteString(color.RedString("Policy failures detected\n\n"))
	reportStr.WriteString(fmt.Sprint(policyCount) + " policies were run ")
	reportStr.WriteString("and " + fmt.Sprint(totalCount) + " failures were detected.\n\n")

	// critical count
	reportStr.WriteString(formatSeverity(settings.LevelCritical) + fmt.Sprint(criticalCount))
	if len(policyFailures[settings.LevelCritical]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(policyFailures[settings.LevelCritical]), ", ") + ")")
	}
	// high count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelHigh) + fmt.Sprint(highCount))
	if len(policyFailures[settings.LevelHigh]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(policyFailures[settings.LevelHigh]), ", ") + ")")
	}
	// medium count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelMedium) + fmt.Sprint(mediumCount))
	if len(policyFailures[settings.LevelMedium]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(policyFailures[settings.LevelMedium]), ", ") + ")")
	}
	// low count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelLow) + fmt.Sprint(lowCount))
	if len(policyFailures[settings.LevelLow]) > 0 {
		reportStr.WriteString(" (" + strings.Join(maps.Keys(policyFailures[settings.LevelLow]), ", ") + ")")
	}

	reportStr.WriteString("\n\n")
}

func writePolicyFailureToString(reportStr *strings.Builder, policyFailure PolicyResult, policySeverity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(policySeverity))
	reportStr.WriteString(policyFailure.PolicyName + " policy failure with " + strings.Join(policyFailure.CategoryGroups, ", ") + "\n")
	reportStr.WriteString(color.HiBlackString(policyFailure.PolicyDescription + "\n"))
	reportStr.WriteString("\n")
	reportStr.WriteString(color.HiBlueString("File: " + underline(policyFailure.Filename+":"+fmt.Sprint(policyFailure.LineNumber)) + "\n"))

	reportStr.WriteString("\n")
	reportStr.WriteString(highlightCodeExtract(policyFailure.Filename, policyFailure.LineNumber, policyFailure.ParentLineNumber, policyFailure.ParentContent, policyFailure.OmitParent))
}

func formatSeverity(policySeverity string) string {
	severityColorFn, ok := severityColorFns[policySeverity]
	if !ok {
		return strings.ToUpper(policySeverity)
	}
	return severityColorFn(strings.ToUpper(policySeverity + ": "))
}

func highlightCodeExtract(fileName string, lineNumber int, extractStartLineNumber int, extract string, singleLine bool) string {
	result := ""
	targetIndex := lineNumber - extractStartLineNumber
	for index, line := range strings.Split(extract, "\n") {
		if index == 0 {
			var err error
			line, err = file.ReadFileSingleLine(fileName, extractStartLineNumber)
			if err != nil {
				break
			}
		}
		if index == targetIndex {
			result += color.MagentaString(" " + fmt.Sprint(extractStartLineNumber+index) + " ")
			result += color.MagentaString(line) + "\n"
			if singleLine {
				break
			}
		} else if !singleLine {
			result += " " + fmt.Sprint(extractStartLineNumber+index) + " "
			result += line + "\n"
		}
	}

	return result
}
