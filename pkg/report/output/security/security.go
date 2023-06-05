package security

import (
	"crypto/md5"
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
	bearerprogressbar "github.com/bearer/bearer/pkg/util/progressbar"
	"github.com/bearer/bearer/pkg/util/rego"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/fatih/color"
	"github.com/hhatto/gocloc"
	log "github.com/rs/zerolog/log"
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
var orderedSeverityLevels = []string{
	types.LevelCritical,
	types.LevelHigh,
	types.LevelMedium,
	types.LevelLow,
	types.LevelWarning,
}

type Results = map[string][]Result

type Input struct {
	RuleId         string             `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule     `json:"rule" yaml:"rule"`
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type Location struct {
	Start  int    `json:"start" yaml:"start"`
	End    int    `json:"end" yaml:"end"`
	Column Column `json:"column" yaml:"column"`
}

type Source struct {
	*Location
}

type Column struct {
	Start int `json:"start" yaml:"start"`
	End   int `json:"end" yaml:"end"`
}

type Sink struct {
	*Location
	Content string `json:"content" yaml:"content"`
}

type Output struct {
	IsLocal         *bool     `json:"is_local,omitempty" yaml:"is_local,omitempty"`
	Source          Source    `json:"source,omitempty" yaml:"source,omitempty"`
	Sink            Sink      `json:"sink,omitempty" yaml:"sink,omitempty"`
	LineNumber      int       `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename        string    `json:"filename,omitempty" yaml:"filename,omitempty"`
	FullFilename    string    `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	CategoryGroups  []string  `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	DataType        *DataType `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	Severity        string    `json:"severity,omitempty" yaml:"severity,omitempty"`
	DetailedContext string    `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type DataType struct {
	CategoryUUID string `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
}

type Result struct {
	*Rule
	LineNumber       int       `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	FullFilename     string    `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	Filename         string    `json:"filename,omitempty" yaml:"filename,omitempty"`
	DataType         *DataType `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	CategoryGroups   []string  `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Source           Source    `json:"source,omitempty" yaml:"source,omitempty"`
	Sink             Sink      `json:"sink,omitempty" yaml:"sink,omitempty"`
	ParentLineNumber int       `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string    `json:"snippet,omitempty" yaml:"snippet,omitempty"`
	Fingerprint      string    `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	OldFingerprint   string    `json:"old_fingerprint,omitempty" yaml:"old_fingerprint,omitempty"`
	DetailedContext  string    `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type Rule struct {
	CWEIDs           []string `json:"cwe_ids" yaml:"cwe_ids"`
	Id               string   `json:"id" yaml:"id"`
	Title            string   `json:"title" yaml:"title"`
	Description      string   `json:"description" yaml:"description"`
	DocumentationUrl string   `json:"documentation_url" yaml:"documentation_url"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (*Results, error) {
	summaryResults := make(Results)
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

	return &summaryResults, nil
}

func evaluateRules(
	summaryResults Results,
	rules map[string]*settings.Rule,
	config settings.Config,
	dataflow *dataflow.DataFlow,
	builtIn bool,
) error {
	outputResults := map[string][]Result{}

	var bar *progressbar.ProgressBar
	if !builtIn {
		bar = bearerprogressbar.GetProgressBar(len(rules), config, "rules")
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

			for i, output := range results["policy_failure"] {
				ruleSummary := &Rule{
					Title:            rule.Description,
					Description:      rule.RemediationMessage,
					Id:               rule.Id,
					CWEIDs:           rule.CWEIDs,
					DocumentationUrl: rule.DocumentationUrl,
				}

				fingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.Filename)
				oldFingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.FullFilename)
				fingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(fingerprintId)), i)
				oldFingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(oldFingerprintId)), i)

				if config.Report.ExcludeFingerprint[fingerprint] {
					// skip finding - fingerprint is in exclude list
					log.Debug().Msgf("Excluding finding with fingerprint %s", fingerprint)
					continue
				}

				result := Result{
					Rule:             ruleSummary,
					FullFilename:     output.FullFilename,
					Filename:         output.Filename,
					LineNumber:       output.LineNumber,
					CategoryGroups:   output.CategoryGroups,
					DataType:         output.DataType,
					Source:           output.Source,
					Sink:             output.Sink,
					ParentLineNumber: output.Sink.Start,
					ParentContent:    output.Sink.Content,
					DetailedContext:  output.DetailedContext,
					Fingerprint:      fingerprint,
					OldFingerprint:   oldFingerprint,
				}

				severity := CalculateSeverity(result.CategoryGroups, rule.Severity, output.IsLocal != nil && *output.IsLocal)

				if config.Report.Severity[severity] {
					outputResults[severity] = append(outputResults[severity], result)
				}
			}
		}
	}

	outputResults = removeDuplicates(outputResults)

	for _, resultsSlice := range outputResults {
		sortResult(resultsSlice)
	}

	for severity, resultSlice := range outputResults {
		summaryResults[severity] = append(summaryResults[severity], resultSlice...)
	}

	return nil
}

func BuildReportString(config settings.Config, results *Results, lineOfCodeOutput *gocloc.Result, dataflow *dataflow.DataFlow) (*strings.Builder, bool) {
	severityForFailure := config.Report.Severity
	reportStr := &strings.Builder{}

	reportStr.WriteString("\n\nSummary Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if config.NoColor && !initialColorSetting {
		color.NoColor = true
	}

	rulesAvailableCount := writeRuleListToString(
		reportStr,
		config.Rules,
		config.BuiltInRules,
		lineOfCodeOutput.Languages,
		config,
	)

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
		if severityLevel != types.LevelWarning && len((*results)[severityLevel]) != 0 {
			// fail the report if we have failures above the severity threshold
			reportPassed = false
		}

		for _, failure := range (*results)[severityLevel] {
			for i := 0; i < len(failure.CWEIDs); i++ {
				failures[severityLevel]["CWE-"+failure.CWEIDs[i]] = true
			}
			writeFailureToString(reportStr, failure, severityLevel)
		}
	}

	if reportPassed {
		reportStr.WriteString("\nNeed to add your own custom rule? Check out the guide: https://docs.bearer.com/guides/custom-rule\n")
	}

	noFailureSummary := checkAndWriteFailureSummaryToString(reportStr, *results, rulesAvailableCount, failures, severityForFailure)

	if noFailureSummary {
		writeSuccessToString(rulesAvailableCount, reportStr)
		writeStatsToString(reportStr, config, lineOfCodeOutput, dataflow)
	}

	writeApiClientResultToString(reportStr, config)

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
	statistics, _, err := stats.GetOutput(lineOfCodeOutput, dataflow, config)
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

func writeApiClientResultToString(
	reportStr *strings.Builder,
	config settings.Config,
) {
	if config.Client != nil {
		if config.Client.Error == nil {
			reportStr.WriteString("\nData successfully sent to Bearer Cloud.\n")
		} else {
			reportStr.WriteString(fmt.Sprintf("\nFailed to send data to Bearer Cloud. %s \n", *config.Client.Error))
		}
	}
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
		} else if slice.Contains(config.Scan.Scanner, "sast") {
			if rule.Language() == "JavaScript" {
				shouldCount = languages["JavaScript"] != nil || languages["TypeScript"] != nil
			} else {
				shouldCount = languages[rule.Language()] != nil
			}
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
	results Results,
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
	reportStr.WriteString(color.RedString(fmt.Sprint(ruleCount) + " checks, " + fmt.Sprint(failureCount+warningCount) + " findings\n\n"))

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
	reportStr.WriteString(result.Title)
	cweCount := len(result.CWEIDs)
	if cweCount > 0 {
		var displayCWEList = []string{}
		for i := 0; i < cweCount; i++ {
			displayCWEList = append(displayCWEList, "CWE-"+result.CWEIDs[i])
		}
		reportStr.WriteString(" [" + strings.Join(displayCWEList, ", ") + "]")
	}
	reportStr.WriteString("\n")

	if result.DocumentationUrl != "" {
		reportStr.WriteString(color.HiBlackString(result.DocumentationUrl + "\n"))
	}

	reportStr.WriteString(color.HiBlackString("To exclude this finding, use the flag --exclude-fingerprint=" + result.Fingerprint + "\n"))
	reportStr.WriteString("\n")
	if result.DetailedContext != "" {
		reportStr.WriteString("Detected: " + result.DetailedContext + "\n\n")
	}
	reportStr.WriteString(color.HiBlueString("File: " + underline(result.FullFilename+":"+fmt.Sprint(result.LineNumber)) + "\n"))

	reportStr.WriteString("\n")
	reportStr.WriteString(highlightCodeExtract(result.FullFilename, result.LineNumber, result.Sink.Start, result.Sink.Content, result))
}

func formatSeverity(severity string) string {
	severityColorFn, ok := severityColorFns[severity]
	if !ok {
		return strings.ToUpper(severity)
	}
	return severityColorFn(strings.ToUpper(severity + ": "))
}

func highlightCodeExtract(fileName string, lineNumber int, extractStartLineNumber int, extract string, record Result) string {
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
			for i, char := range line {
				if i >= record.Source.Column.Start-1 && i < record.Source.Column.End-1 {
					result += color.BlueString(fmt.Sprintf("%c", char))
				} else {
					result += color.MagentaString(fmt.Sprintf("%c", char))
				}
			}
			result += "\n"
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

// removeDuplicates removes detections for same detector with same line number by keeping only a single highest severity detection
func removeDuplicates(data map[string][]Result) map[string][]Result {
	filteredData := map[string][]Result{}

	type Key struct {
		LineNumber int
		FileName   string
		Detector   string
	}

	reportedDetections := set.Set[Key]{}

	// filter duplicates
	for _, severity := range orderedSeverityLevels {
		resultsSlice, hasSeverity := data[severity]
		if !hasSeverity {
			continue
		}

		for _, result := range resultsSlice {
			key := Key{
				LineNumber: result.LineNumber,
				FileName:   result.Filename,
				Detector:   result.Rule.Id,
			}
			if reportedDetections.Add(key) {
				filteredData[severity] = append(filteredData[severity], result)
			}
		}
	}

	return filteredData
}

func sortResult(data []Result) {
	sort.Slice(data, func(i, j int) bool {
		vulnerabilityA := data[i]
		vulnerabilityB := data[j]

		if vulnerabilityA.Rule.Id < vulnerabilityB.Rule.Id {
			return true
		}
		if vulnerabilityA.Rule.Id > vulnerabilityB.Rule.Id {
			return false
		}

		if vulnerabilityA.Filename < vulnerabilityB.Filename {
			return true
		}
		if vulnerabilityA.Filename > vulnerabilityB.Filename {
			return false
		}

		if vulnerabilityA.LineNumber < vulnerabilityB.LineNumber {
			return true
		}
		if vulnerabilityA.LineNumber > vulnerabilityB.LineNumber {
			return false
		}

		if vulnerabilityA.ParentLineNumber < vulnerabilityB.ParentLineNumber {
			return true
		}
		if vulnerabilityA.ParentLineNumber > vulnerabilityB.ParentLineNumber {
			return false
		}

		if vulnerabilityA.ParentContent < vulnerabilityB.ParentContent {
			return true
		}

		return false
	})
}
