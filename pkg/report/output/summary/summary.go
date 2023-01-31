package summary

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
	"golang.org/x/exp/slices"

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

type PolicyResult struct {
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

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (map[string][]PolicyResult, error) {
	// policy results grouped by severity (critical, high, ...)
	result := make(map[string][]PolicyResult)

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
				policyResult := PolicyResult{
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

func BuildReportString(rules map[string]*settings.Rule, policyResults map[string][]PolicyResult, withoutColor bool) *strings.Builder {
	reportStr := &strings.Builder{}
	reportStr.WriteString("\n\nSummary Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if withoutColor && !initialColorSetting {
		color.NoColor = true
	}

	writeRuleListToString(reportStr, rules)

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
			policyFailures[policyLevel][policyFailure.PolicyDSRID] = true
			writePolicyFailureToString(reportStr, policyFailure, policyLevel)
		}
	}

	writeSummaryToString(reportStr, policyResults, len(rules), policyFailures)

	color.NoColor = initialColorSetting

	return reportStr
}

func FindHighestSeverity(groups []string, severity map[string]string) string {
	var severities []string
	for _, group := range groups {
		severities = append(severities, severity[group])
	}

	if slices.Contains(severities, "critical") {
		return settings.LevelCritical
	} else if slices.Contains(severities, "high") {
		return settings.LevelHigh
	} else if slices.Contains(severities, "medium") {
		return settings.LevelMedium
	} else if slices.Contains(severities, "low") {
		return settings.LevelLow
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
	policyResults map[string][]PolicyResult,
	policyCount int, policyFailures map[string]map[string]bool,
) {
	reportStr.WriteString("\n=====================================")

	// give summary including counts
	if len(policyResults) == 0 {
		reportStr.WriteString("\n\n")
		reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
		reportStr.WriteString(fmt.Sprint(policyCount) + " checks were run and no failures were detected.\n\n")
		return
	}

	criticalCount := len(policyResults[settings.LevelCritical])
	highCount := len(policyResults[settings.LevelHigh])
	mediumCount := len(policyResults[settings.LevelMedium])
	lowCount := len(policyResults[settings.LevelLow])

	totalCount := criticalCount + highCount + mediumCount + lowCount

	reportStr.WriteString("\n\n")
	reportStr.WriteString(color.RedString(fmt.Sprint(policyCount) + " checks, " + fmt.Sprint(totalCount) + " failures\n\n"))

	// critical count
	reportStr.WriteString(formatSeverity(settings.LevelCritical) + fmt.Sprint(criticalCount))
	if len(policyFailures[settings.LevelCritical]) > 0 {
		policyIds := maps.Keys(policyFailures[settings.LevelCritical])
		sort.Strings(policyIds)
		reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
	}
	// high count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelHigh) + fmt.Sprint(highCount))
	if len(policyFailures[settings.LevelHigh]) > 0 {
		policyIds := maps.Keys(policyFailures[settings.LevelHigh])
		sort.Strings(policyIds)
		reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
	}
	// medium count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelMedium) + fmt.Sprint(mediumCount))
	if len(policyFailures[settings.LevelMedium]) > 0 {
		policyIds := maps.Keys(policyFailures[settings.LevelMedium])
		sort.Strings(policyIds)
		reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
	}
	// low count
	reportStr.WriteString("\n" + formatSeverity(settings.LevelLow) + fmt.Sprint(lowCount))
	if len(policyFailures[settings.LevelLow]) > 0 {
		policyIds := maps.Keys(policyFailures[settings.LevelLow])
		sort.Strings(policyIds)
		reportStr.WriteString(" (" + strings.Join(policyIds, ", ") + ")")
	}

	reportStr.WriteString("\n")
}

func writePolicyFailureToString(reportStr *strings.Builder, policyFailure PolicyResult, policySeverity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(policySeverity))
	reportStr.WriteString(policyFailure.PolicyDescription + " [" + policyFailure.PolicyDSRID + "]" + "\n")
	reportStr.WriteString(color.HiBlackString("https://curio.sh/reference/rules/" + policyFailure.PolicyDisplayId + "\n"))
	reportStr.WriteString(color.HiBlackString("To skip this rule, use the flag --skip-rule=" + policyFailure.PolicyDisplayId + "\n"))
	reportStr.WriteString("\n")
	if policyFailure.DetailedContext != "" {
		reportStr.WriteString("Detected: " + policyFailure.DetailedContext + "\n")
	}
	reportStr.WriteString(color.HiBlueString("File: " + underline(policyFailure.Filename+":"+fmt.Sprint(policyFailure.LineNumber)) + "\n"))

	if policyFailure.DetailedContext == "" {
		reportStr.WriteString("\n")
		reportStr.WriteString(highlightCodeExtract(policyFailure.Filename, policyFailure.LineNumber, policyFailure.ParentLineNumber, policyFailure.ParentContent))
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
