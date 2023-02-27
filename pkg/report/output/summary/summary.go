package summary

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/maputil"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/rego"
	"github.com/fatih/color"
	"github.com/hhatto/gocloc"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/pkg/report/output/dataflow"
	stats "github.com/bearer/bearer/pkg/report/output/stats"
)

var underline = color.New(color.Underline).SprintFunc()
var severityColorFns = map[string]func(x ...interface{}) string{
	types.LevelCritical: color.New(color.FgRed).SprintFunc(),
	types.LevelHigh:     color.New(color.FgHiRed).SprintFunc(),
	types.LevelMedium:   color.New(color.FgYellow).SprintFunc(),
	types.LevelLow:      color.New(color.FgBlue).SprintFunc(),
	types.LevelWarning:  color.New(color.FgCyan).SprintFunc(),
}
var orderedSeverityLevels = [5]string{
	types.LevelCritical,
	types.LevelHigh,
	types.LevelMedium,
	types.LevelLow,
	types.LevelWarning,
}

type Input struct {
	RuleId         string             `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule     `json:"rule" yaml:"rule"`
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type Output struct {
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
	Rule             *RuleResultSummary `json:"rule" yaml:"rule"`
	LineNumber       int                `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename         string             `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups   []string           `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	ParentLineNumber int                `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string             `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`

	DetailedContext string `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type RuleResultSummary struct {
	CWEIDs           []string `json:"cwe_ids" yaml:"cwe_ids"`
	Id               string   `json:"id" yaml:"id"`
	Description      string   `json:"description" yaml:"description"`
	DocumentationUrl string   `json:"documentation_url" yaml:"documentation_url"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (map[string][]Result, error) {
	// results grouped by severity (critical, high, ...)
	summaryResults := make(map[string][]Result)

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Evaluating rules")
	}

	bar := output.GetProgressBar(len(config.Rules), config, "rules")

	for _, rule := range maputil.ToSortedSlice(config.Rules) {
		err := bar.Add(1)
		if err != nil {
			output.StdErrLogger().Msgf("Rule %s failed to write progress bar %e", rule.Id, err)
		}

		if !rule.PolicyType() {
			continue
		}

		policy := config.Policies[rule.Type]

		// Create a prepared query that can be evaluated.
		rs, err := rego.RunQuery(policy.Query,
			Input{
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

			var results map[string][]Output
			err = json.Unmarshal(jsonRes, &results)
			if err != nil {
				return nil, err
			}

			for _, output := range results["policy_failure"] {
				ruleSummary := &RuleResultSummary{
					Description:      rule.Description,
					Id:               rule.Id,
					CWEIDs:           rule.CWEIDs,
					DocumentationUrl: rule.DocumentationUrl,
				}
				result := Result{
					Rule:             ruleSummary,
					Filename:         output.Filename,
					LineNumber:       output.LineNumber,
					CategoryGroups:   output.CategoryGroups,
					ParentLineNumber: output.ParentLineNumber,
					ParentContent:    output.ParentContent,
					DetailedContext:  output.DetailedContext,
				}

				severity := FindHighestSeverity(result.CategoryGroups, rule.Severity)

				if config.Report.Severity[severity] {
					summaryResults[severity] = append(summaryResults[severity], result)
				}
			}
		}
	}

	return summaryResults, nil
}

func BuildReportString(config settings.Config, results map[string][]Result, lineOfCodeOutput *gocloc.Result, dataflow *dataflow.DataFlow) (*strings.Builder, bool) {
	rules := config.Rules
	withoutColor := config.Report.Output != ""
	severityForFailure := config.Report.Severity
	reportStr := &strings.Builder{}

	reportStr.WriteString("\n\nSummary Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if withoutColor && !initialColorSetting {
		color.NoColor = true
	}

	rulesAvailableCount := writeRuleListToString(reportStr, rules, lineOfCodeOutput.Languages)

	failures := map[string]map[string]bool{
		types.LevelCritical: make(map[string]bool),
		types.LevelHigh:     make(map[string]bool),
		types.LevelMedium:   make(map[string]bool),
		types.LevelLow:      make(map[string]bool),
		types.LevelWarning:  make(map[string]bool),
	}

	reportPassed := true
	for _, severityLevel := range orderedSeverityLevels {
		if !severityForFailure[severityLevel] {
			continue
		}
		if severityLevel != types.LevelWarning && len(results[severityLevel]) != 0 {
			// fail the report if we have failures above the severity threshold
			reportPassed = false
		}

		for _, failure := range results[severityLevel] {
			for i := 0; i < len(failure.Rule.CWEIDs); i++ {
				failures[severityLevel]["CWE-"+failure.Rule.CWEIDs[i]] = true
			}
			writeFailureToString(reportStr, failure, severityLevel)
		}
	}

	if reportPassed {
		reportStr.WriteString("\nNeed to add your own custom rule? Check out the guide: https://docs.bearer.sh/guides/custom-rule\n")
	}

	noFailureSummary := checkAndWriteFailureSummaryToString(reportStr, results, rulesAvailableCount, failures, severityForFailure)

	if noFailureSummary {
		writeSuccessToString(rulesAvailableCount, reportStr)
		writeStatsToString(reportStr, config, lineOfCodeOutput, dataflow)
	}

	reportStr.WriteString("\nNeed help or want to discuss the output? Join the Community https://discord.gg/eaHZBJUXRF\n")

	color.NoColor = initialColorSetting

	return reportStr, reportPassed
}

func FindHighestSeverity(groups []string, severity map[string]string) string {
	var severities []string
	for _, group := range groups {
		severities = append(severities, severity[group])
	}

	if slices.Contains(severities, types.LevelCritical) {
		return types.LevelCritical
	} else if slices.Contains(severities, types.LevelHigh) {
		return types.LevelHigh
	} else if slices.Contains(severities, types.LevelMedium) {
		return types.LevelMedium
	} else if slices.Contains(severities, types.LevelLow) {
		return types.LevelLow
	} else if slices.Contains(severities, types.LevelWarning) {
		return types.LevelWarning
	}

	return severity["default"]
}

func writeStatsToString(
	reportStr *strings.Builder,
	config settings.Config,
	lineOfCodeOutput *gocloc.Result,
	dataflow *dataflow.DataFlow,
) {
	statistics, err := stats.GetOutput(lineOfCodeOutput, dataflow, config)
	if err != nil {
		return
	}
	if stats.AnythingFoundFor(statistics) {
		reportStr.WriteString("\nBearer found:\n")
		stats.WriteStatsToString(reportStr, statistics)
		reportStr.WriteString("\n")
	}
}

func writeRuleListToString(
	reportStr *strings.Builder,
	rules map[string]*settings.Rule,
	languages map[string]*gocloc.Language,
) int {
	// list rules that were run
	reportStr.WriteString("\n\nRules: \n")
	defaultRuleCount := 0
	customRuleCount := 0

	for key := range rules {
		rule := rules[key]
		if !rule.PolicyType() {
			continue
		}

		exists := rule.Language() == "secret" || languages[rule.Language()] != nil

		if !exists {
			continue
		}

		if strings.HasPrefix(rule.DocumentationUrl, "https://docs.bearer.com") {
			defaultRuleCount++
		} else {
			customRuleCount++
		}
	}

	reportStr.WriteString(fmt.Sprintf(" - %d default rules applied ", defaultRuleCount))
	reportStr.WriteString(color.HiBlackString("(https://docs.bearer.com/reference/rules)\n"))
	if customRuleCount > 0 {
		reportStr.WriteString(fmt.Sprintf(" - %d custom rules applied", customRuleCount))
	}

	return defaultRuleCount + customRuleCount
}

func writeSuccessToString(ruleCount int, reportStr *strings.Builder) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(color.HiGreenString("SUCCESS\n\n"))
	reportStr.WriteString(fmt.Sprint(ruleCount) + " checks were run and no failures were detected. Great job! 👏\n")
}

func checkAndWriteFailureSummaryToString(
	reportStr *strings.Builder,
	results map[string][]Result,
	ruleCount int,
	failures map[string]map[string]bool,
	severityForFailure map[string]bool,
) bool {
	reportStr.WriteString("\n=====================================")

	if len(results) == 0 {
		return true
	}

	// give summary including counts
	failureCount := 0
	warningCount := 0
	for _, severityLevel := range maps.Keys(severityForFailure) {
		if !severityForFailure[severityLevel] {
			continue
		}
		if severityLevel == types.LevelWarning {
			warningCount += len(results[severityLevel])
			continue
		}
		failureCount += len(results[severityLevel])
	}

	if failureCount == 0 && warningCount == 0 {
		return true
	}

	reportStr.WriteString("\n\n")
	if failureCount == 0 {
		// only warnings
		reportStr.WriteString(fmt.Sprint(ruleCount) + " checks, " + fmt.Sprint(warningCount) + " warnings\n\n")
	} else {
		reportStr.WriteString(color.RedString(fmt.Sprint(ruleCount) + " checks, " + fmt.Sprint(failureCount) + " failures, " + fmt.Sprint(warningCount) + " warnings\n\n"))
	}

	for i, severityLevel := range orderedSeverityLevels {
		if !severityForFailure[severityLevel] {
			continue
		}
		if i > 0 {
			reportStr.WriteString("\n")
		}
		reportStr.WriteString(formatSeverity(severityLevel) + fmt.Sprint(len(results[severityLevel])))
		if len(failures[severityLevel]) > 0 {
			ruleIds := maps.Keys(failures[severityLevel])
			sort.Strings(ruleIds)
			if len(ruleIds) > 0 {
				reportStr.WriteString(" (" + strings.Join(ruleIds, ", ") + ")")
			}
		}
	}

	reportStr.WriteString("\n")

	return false
}

func writeFailureToString(reportStr *strings.Builder, result Result, severity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(severity))
	reportStr.WriteString(result.Rule.Description)
	cweCount := len(result.Rule.CWEIDs)
	if cweCount > 0 {
		var displayCWEList = []string{}
		for i := 0; i < cweCount; i++ {
			displayCWEList = append(displayCWEList, "CWE-"+result.Rule.CWEIDs[i])
		}
		reportStr.WriteString(" [" + strings.Join(displayCWEList, ", ") + "]")
	}
	reportStr.WriteString("\n")

	if result.Rule.DocumentationUrl != "" {
		reportStr.WriteString(color.HiBlackString(result.Rule.DocumentationUrl + "\n"))
	}

	reportStr.WriteString(color.HiBlackString("To skip this rule, use the flag --skip-rule=" + result.Rule.Id + "\n"))
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

func formatSeverity(severity string) string {
	severityColorFn, ok := severityColorFns[severity]
	if !ok {
		return strings.ToUpper(severity)
	}
	return severityColorFn(strings.ToUpper(severity + ": "))
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
