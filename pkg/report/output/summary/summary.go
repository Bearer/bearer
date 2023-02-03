package summary

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/types"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/bearer/curio/pkg/util/rego"
	"github.com/fatih/color"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/bearer/curio/pkg/report/output/dataflow"
)

var underline = color.New(color.Underline).SprintFunc()
var severityColorFns = map[string]func(x ...interface{}) string{
	types.LevelCritical: color.New(color.FgRed).SprintFunc(),
	types.LevelHigh:     color.New(color.FgHiRed).SprintFunc(),
	types.LevelMedium:   color.New(color.FgYellow).SprintFunc(),
	types.LevelLow:      color.New(color.FgBlue).SprintFunc(),
	types.LevelWarning:  color.New(color.FgCyan).SprintFunc(),
}

type PolicyInput struct {
	PolicyId       string             `json:"policy_id" yaml:"policy_id"`
	RuleId         string             `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule     `json:"rule" yaml:"rule"`
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
	DetailedContext  string   `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type Result struct {
	PolicyName        string   `json:"policy_name" yaml:"policy_name"`
	PolicyDSRID       string   `json:"policy_dsrid" yaml:"policy_dsrid"`
	PolicyDisplayId   string   `json:"policy_display_id" yaml:"policy_display_id"`
	PolicyDescription string   `json:"policy_description" yaml:"policy_description"`
	LineNumber        int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename          string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups    []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	ParentLineNumber  int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent     string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	OmitParent        bool     `json:"omit_parent,omitempty" yaml:"omit_parent,omitempty"`
	DetailedContext   string   `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (map[string][]Result, error) {
	// policy results grouped by severity (critical, high, ...)
	result := make(map[string][]Result)

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Evaluating rules")
	}

	bar := output.GetProgressBar(len(config.Rules), config, "rules")

	for _, rule := range config.Rules {
		err := bar.Add(1)
		if err != nil {
			output.StdErrLogger().Msgf("Policy %s failed to write progress bar %e", rule.Id, err)
		}

		if !rule.PolicyType() {
			continue
		}

		policy := config.Policies[rule.Type]

		// Create a prepared query that can be evaluated.
		rs, err := rego.RunQuery(policy.Query,
			PolicyInput{
				RuleId:         rule.Id,
				Rule:           rule,
				Dataflow:       dataflow,
				DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
			},
			// TODO: perf question: can we do this once?
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
				policyResult := Result{
					PolicyDescription: rule.Description,
					PolicyDisplayId:   rule.Id,
					PolicyDSRID:       rule.DSRID,
					Filename:          policyOutput.Filename,
					LineNumber:        policyOutput.LineNumber,
					CategoryGroups:    policyOutput.CategoryGroups,
					OmitParent:        rule.OmitParent,
					ParentLineNumber:  policyOutput.ParentLineNumber,
					ParentContent:     policyOutput.ParentContent,
					DetailedContext:   policyOutput.DetailedContext,
				}

				severity := FindHighestSeverity(policyOutput.CategoryGroups, rule.Severity)

				result[severity] = append(result[severity], policyResult)
			}
		}
	}

	return result, nil
}

func BuildReportString(rules map[string]*settings.Rule, results map[string][]Result, severityForFailure map[string]bool, withoutColor bool) (*strings.Builder, bool) {
	reportStr := &strings.Builder{}
	reportStr.WriteString("\n\nSummary Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if withoutColor && !initialColorSetting {
		color.NoColor = true
	}

	writeRuleListToString(reportStr, rules)

	failures := map[string]map[string]bool{
		types.LevelCritical: make(map[string]bool),
		types.LevelHigh:     make(map[string]bool),
		types.LevelMedium:   make(map[string]bool),
		types.LevelLow:      make(map[string]bool),
	}

	reportPassed := true
	for _, severityLevel := range []string{
		types.LevelCritical,
		types.LevelHigh,
		types.LevelMedium,
		types.LevelLow,
	} {
		if severityForFailure[severityLevel] && len(results[severityLevel]) != 0 {
			// fail the report if we have failures above the severity threshold
			reportPassed = false
		}

		for _, failure := range results[severityLevel] {
			failures[severityLevel][failure.PolicyDSRID] = true
			if severityForFailure[severityLevel] {
				writeFailureToString(reportStr, failure, severityLevel)
			}
		}
	}

	writeSummaryToString(reportStr, results, len(rules), failures, severityForFailure)

	color.NoColor = initialColorSetting

	return reportStr, reportPassed
}

func FindHighestSeverity(groups []string, severity map[string]string) string {
	var severities []string
	for _, group := range groups {
		severities = append(severities, severity[group])
	}

	if slices.Contains(severities, "critical") {
		return types.LevelCritical
	} else if slices.Contains(severities, "high") {
		return types.LevelHigh
	} else if slices.Contains(severities, "medium") {
		return types.LevelMedium
	} else if slices.Contains(severities, "low") {
		return types.LevelLow
	} else if slices.Contains(severities, "warning") {
		return types.LevelWarning
	}

	return severity["default"]
}

func writeRuleListToString(
	reportStr *strings.Builder,
	rules map[string]*settings.Rule) {
	// list rules that were run
	reportStr.WriteString("\nChecks: \n\n")
	ruleList := []string{}
	for key := range rules {
		rule := rules[key]
		if !rule.PolicyType() {
			continue
		}
		ruleList = append(ruleList, color.HiBlackString("- "+rule.Description+" - "+key+" ["+rule.DSRID+"/"+rule.Id+"]\n"))
	}

	sort.Strings(ruleList)
	reportStr.WriteString(strings.Join(ruleList, ""))
}

func writeSummaryToString(
	reportStr *strings.Builder,
	policyResults map[string][]Result,
	policyCount int, policyFailures map[string]map[string]bool,
	severityForFailure map[string]bool,
) {
	reportStr.WriteString("\n=====================================")

	if len(policyResults) == 0 {
		reportStr.WriteString("\n\n")
		reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
		reportStr.WriteString(fmt.Sprint(policyCount) + " checks were run and no failures were detected.\n\n")
		return
	}

	// give summary including counts
	failureCount := 0
	for _, severityLevel := range maps.Keys(severityForFailure) {
		if severityForFailure[severityLevel] {
			failureCount += len(policyResults[severityLevel])
			continue
		}
	}

	if failureCount == 0 {
		reportStr.WriteString("\n\n")
		reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
		reportStr.WriteString(fmt.Sprint(policyCount) + " checks were run and no failures were detected.\n\n")
		return
	}

	reportStr.WriteString("\n\n")

	reportStr.WriteString(color.RedString(fmt.Sprint(policyCount) + " checks, " + fmt.Sprint(failureCount) + " failures\n\n"))

	for i, severityLevel := range maps.Keys(severityForFailure) {
		if !severityForFailure[severityLevel] {
			continue
		}

		if i > 0 {
			reportStr.WriteString("\n")
		}
		reportStr.WriteString(formatSeverity(severityLevel) + fmt.Sprint(len(policyResults[severityLevel])))
		if len(policyFailures[severityLevel]) > 0 {
			policyIds := maps.Keys(policyFailures[severityLevel])
			sort.Strings(policyIds)
			reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
		}
	}
	// warning count
	warningCount := 0
	reportStr.WriteString("\n" + formatSeverity(types.LevelWarning) + fmt.Sprint(warningCount))
	if len(policyFailures[types.LevelWarning]) > 0 {
		policyIds := maps.Keys(policyFailures[types.LevelWarning])
		sort.Strings(policyIds)
		reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
	}

	reportStr.WriteString("\n")
}

func writeFailureToString(reportStr *strings.Builder, result Result, policySeverity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(policySeverity))
	reportStr.WriteString(result.PolicyDescription + " [" + result.PolicyDSRID + "]" + "\n")
	reportStr.WriteString(color.HiBlackString("https://curio.sh/reference/rules/" + result.PolicyDisplayId + "\n"))
	reportStr.WriteString(color.HiBlackString("To skip this rule, use the flag --skip-rule=" + result.PolicyDisplayId + "\n"))
	reportStr.WriteString("\n")
	if result.DetailedContext != "" {
		reportStr.WriteString("Detected: " + result.DetailedContext + "\n")
	}
	reportStr.WriteString(color.HiBlueString("File: " + underline(result.Filename+":"+fmt.Sprint(result.LineNumber)) + "\n"))

	if result.DetailedContext == "" {
		reportStr.WriteString("\n")
		reportStr.WriteString(highlightCodeExtract(result.Filename, result.LineNumber, result.ParentLineNumber, result.ParentContent))
	}
}

func formatSeverity(policySeverity string) string {
	severityColorFn, ok := severityColorFns[policySeverity]
	if !ok {
		return strings.ToUpper(policySeverity)
	}
	return severityColorFn(strings.ToUpper(policySeverity + ": "))
}

func highlightCodeExtract(fileName string, lineNumber int, extractStartLineNumber int, extract string) string {
	result := ""
	targetIndex := lineNumber - extractStartLineNumber
	beforeOrAfterDetectionLinesAllowed := 3
	beforeDetection := targetIndex - beforeOrAfterDetectionLinesAllowed
	afterDetection := targetIndex + beforeOrAfterDetectionLinesAllowed
	items := strings.Split(extract, "\n")
	for index, line := range items {
		if index == 0 {
			var err error
			line, err = file.ReadFileSingleLine(fileName, extractStartLineNumber)
			if err != nil {
				break
			}
		}

		if index == targetIndex {
			result += color.MagentaString(fmt.Sprintf(" %d ", extractStartLineNumber+index))
			result += color.MagentaString(line) + "\n"
		} else if index == 0 || len(items)-1 == index {
			result += fmt.Sprintf(" %d ", extractStartLineNumber+index)
			result += fmt.Sprintf("%s\n", line)
		} else if index >= beforeDetection && index <= afterDetection {
			if index == beforeDetection {
				result += fmt.Sprintf("%s\t...\n", strings.Repeat(" ", iterativeDigitsCount(extractStartLineNumber)))
			}
			result += fmt.Sprintf(" %d ", extractStartLineNumber+index)
			result += fmt.Sprintf("%s\n", line)

			if index == afterDetection {
				result += fmt.Sprintf("%s\t...\n", strings.Repeat(" ", iterativeDigitsCount(extractStartLineNumber)))
			}
		}
	}

	return result
}

func iterativeDigitsCount(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}

	return count
}
