package security

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
	"github.com/schollz/progressbar/v3"
	"github.com/ssoroka/slice"
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
	IsLocal          *bool    `json:"is_local,omitempty" yaml:"is_local,omitempty"`
	ParentLineNumber int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	LineNumber       int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename         string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups   []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Severity         string   `json:"severity,omitempty" yaml:"severity,omitempty"`
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
	summaryResults := make(map[string][]Result)
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Evaluating rules")
	}

	err := evaluateRules(summaryResults, config.BuiltInRules, config, dataflow, true)
	if err != nil {
		return nil, err
	}
	err = evaluateRules(summaryResults, config.Rules, config, dataflow, false)
	if err != nil {
		return nil, err
	}

	return summaryResults, nil
}

func evaluateRules(
	summaryResults map[string][]Result,
	rules map[string]*settings.Rule,
	config settings.Config,
	dataflow *dataflow.DataFlow,
	builtIn bool,
) error {
	var bar *progressbar.ProgressBar
	if !builtIn {
		bar = output.GetProgressBar(len(rules), config, "rules")
	}

	for _, rule := range maputil.ToSortedSlice(rules) {
		if !builtIn {
			err := bar.Add(1)
			if err != nil {
				output.StdErrLogger().Msgf("Rule %s failed to write progress bar %s", rule.Id, err)
			}
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
			return err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return err
			}

			var results map[string][]Output
			err = json.Unmarshal(jsonRes, &results)
			if err != nil {
				return err
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

				severity := CalculateSeverity(result.CategoryGroups, rule.Severity, output.IsLocal != nil && *output.IsLocal)

				if config.Report.Severity[severity] {
					summaryResults[severity] = append(summaryResults[severity], result)
				}
			}
		}
	}

	return nil
}

func BuildReportString(config settings.Config, results map[string][]Result, lineOfCodeOutput *gocloc.Result, dataflow *dataflow.DataFlow) (*strings.Builder, bool) {
	rules := config.Rules
	builtInRules := config.BuiltInRules

	withoutColor := config.Report.Output != ""
	severityForFailure := config.Report.Severity
	reportStr := &strings.Builder{}

	reportStr.WriteString("\n\nSummary Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if withoutColor && !initialColorSetting {
		color.NoColor = true
	}

	rulesAvailableCount := writeRuleListToString(reportStr, rules, builtInRules, lineOfCodeOutput.Languages, config)

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
		reportStr.WriteString("\nNeed to add your own custom rule? Check out the guide: https://docs.bearer.com/guides/custom-rule\n")
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

func CalculateSeverity(groups []string, severity string, hasLocalDataTypes bool) string {
	if severity == types.LevelWarning {
		return types.LevelWarning
	}

	// highest sensitive data category
	sensitiveDataCategoryWeighting := 0
	if slices.Contains(groups, "PHI") {
		sensitiveDataCategoryWeighting = 3
	} else if slices.Contains(groups, "Personal Data (Sensitive)") {
		sensitiveDataCategoryWeighting = 3
	} else if slices.Contains(groups, "Personal Data") {
		sensitiveDataCategoryWeighting = 2
	} else if slices.Contains(groups, "PII") {
		sensitiveDataCategoryWeighting = 1
	}

	var ruleSeverityWeighting int
	switch severity {
	case types.LevelCritical:
		ruleSeverityWeighting = 8
	case types.LevelHigh:
		ruleSeverityWeighting = 5
	case types.LevelMedium:
		ruleSeverityWeighting = 3
	default:
		ruleSeverityWeighting = 2 // low weighting as default
	}

	triggerWeighting := 1
	if hasLocalDataTypes {
		triggerWeighting = 2
	}

	switch finalWeighting := ruleSeverityWeighting + (sensitiveDataCategoryWeighting * triggerWeighting); {
	case finalWeighting >= 8:
		return types.LevelCritical
	case finalWeighting >= 5:
		return types.LevelHigh
	case finalWeighting >= 3:
		return types.LevelMedium
	}

	return types.LevelLow
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
	builtInRules map[string]*settings.Rule,
	languages map[string]*gocloc.Language,
	config settings.Config,
) int {
	// list rules that were run
	reportStr.WriteString("\n\nRules: \n")

	defaultRuleCount, customRuleCount := countRules(rules, languages, config, false)
	builtInCount, _ := countRules(builtInRules, languages, config, true)
	defaultRuleCount = defaultRuleCount + builtInCount

	reportStr.WriteString(fmt.Sprintf(" - %d default rules applied ", defaultRuleCount))
	reportStr.WriteString(color.HiBlackString("(https://docs.bearer.com/reference/rules)\n"))
	if customRuleCount > 0 {
		reportStr.WriteString(fmt.Sprintf(" - %d custom rules applied", customRuleCount))
	}

	return defaultRuleCount + customRuleCount
}

func countRules(
	rules map[string]*settings.Rule,
	languages map[string]*gocloc.Language,
	config settings.Config,
	builtIn bool,
) (int, int) {
	defaultRuleCount := 0
	customRuleCount := 0

	for key := range rules {
		rule := rules[key]
		if !rule.PolicyType() {
			continue
		}

		var shouldCount bool

		if rule.Language() == "secret" {
			shouldCount = slice.Contains(config.Scan.Scanner, "secrets")
		} else {
			shouldCount = languages[rule.Language()] != nil && slice.Contains(config.Scan.Scanner, "sast")
		}

		if !shouldCount {
			continue
		}

		if strings.HasPrefix(rule.DocumentationUrl, "https://docs.bearer.com") || builtIn {
			defaultRuleCount++
		} else {
			customRuleCount++
		}
	}

	return defaultRuleCount, customRuleCount
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
