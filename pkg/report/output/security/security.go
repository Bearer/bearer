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
	"github.com/rodaine/table"
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

type RuleCounter struct {
	DefaultRuleCount int
	CustomRuleCount  int
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
	LineNumber       int         `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	FullFilename     string      `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	Filename         string      `json:"filename,omitempty" yaml:"filename,omitempty"`
	DataType         *DataType   `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	CategoryGroups   []string    `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Source           Source      `json:"source,omitempty" yaml:"source,omitempty"`
	Sink             Sink        `json:"sink,omitempty" yaml:"sink,omitempty"`
	ParentLineNumber int         `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string      `json:"snippet,omitempty" yaml:"snippet,omitempty"`
	Fingerprint      string      `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	OldFingerprint   string      `json:"old_fingerprint,omitempty" yaml:"old_fingerprint,omitempty"`
	DetailedContext  string      `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	CodeExtract      string      `json:"code_extract,omitempty" yaml:"code_extract,omitempty"`
	RawCodeExtract   []file.Line `json:"-" yaml:"-"`
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
		output.StdErrLog("Evaluating rules")
	}

	builtInFingerprints, err := evaluateRules(summaryResults, config.BuiltInRules, config, dataflow, true)
	if err != nil {
		return nil, err
	}
	fingerprints, err := evaluateRules(summaryResults, config.Rules, config, dataflow, false)
	if err != nil {
		return nil, err
	}

	if !config.Scan.Quiet {
		fingerprints = append(fingerprints, builtInFingerprints...)
		unusedFingerprints := removeUnusedFingerprints(fingerprints, config.Report.ExcludeFingerprint)
		if len(unusedFingerprints) > 0 {
			output.StdErrLog("\n=====================================\n")
			output.StdErrLog(fmt.Sprintf("%d excluded fingerprints present in your Bearer configuration file are no longer detected:", len(unusedFingerprints)))
			for _, fingerprint := range unusedFingerprints {
				output.StdErrLog(fmt.Sprintf("  - %s", fingerprint))
			}
			output.StdErrLog("\n=====================================")
		}
	}

	return &summaryResults, nil
}

func evaluateRules(
	summaryResults Results,
	rules map[string]*settings.Rule,
	config settings.Config,
	dataflow *dataflow.DataFlow,
	builtIn bool,
) ([]string, error) {
	outputResults := map[string][]Result{}

	var bar *progressbar.ProgressBar
	if !builtIn {
		bar = bearerprogressbar.GetProgressBar(len(rules), config, "rules")
	}

	var fingerprints []string

	for _, rule := range maputil.ToSortedSlice(rules) {
		if !builtIn {
			err := bar.Add(1)
			if err != nil {
				output.StdErrLog(fmt.Sprintf("Rule %s failed to write progress bar %s", rule.Id, err))
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
			return fingerprints, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return fingerprints, err
			}

			var results map[string][]Output
			err = json.Unmarshal(jsonRes, &results)
			if err != nil {
				return fingerprints, err
			}

			ruleSummary := &Rule{
				Title:            rule.Description,
				Description:      rule.RemediationMessage,
				Id:               rule.Id,
				CWEIDs:           rule.CWEIDs,
				DocumentationUrl: rule.DocumentationUrl,
			}

			instanceCount := make(map[string]int)
			policyFailures := results["policy_failure"]
			sortByLineNumber(policyFailures)

			for i, output := range policyFailures {
				fingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.Filename)
				oldFingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.FullFilename)
				fingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(fingerprintId)), instanceCount[output.Filename])
				oldFingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(oldFingerprintId)), i)
				instanceCount[output.Filename]++
				fingerprints = append(fingerprints, fingerprint)
				if config.Report.ExcludeFingerprint[fingerprint] {
					// skip finding - fingerprint is in exclude list
					log.Debug().Msgf("Excluding finding with fingerprint %s", fingerprint)
					continue
				}

				rawCodeExtract := codeExtract(output.FullFilename, output.Source, output.Sink)
				codeExtract := getExtract(rawCodeExtract)

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
					CodeExtract:      codeExtract,
					RawCodeExtract:   rawCodeExtract,
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

	return fingerprints, nil
}

func removeUnusedFingerprints(detectedFingerprints []string, excludeFingerprints map[string]bool) []string {
	filteredFingerprints := []string{}

	for fingerprint := range excludeFingerprints {
		if !slice.Contains(detectedFingerprints, fingerprint) {
			filteredFingerprints = append(filteredFingerprints, fingerprint)
		}
	}

	return filteredFingerprints
}

func getExtract(rawCodeExtract []file.Line) string {
	var parts []string
	for _, line := range rawCodeExtract {
		parts = append(parts, line.Extract)
	}

	return strings.Join(parts, "\n")
}

func BuildReportString(config settings.Config, results *Results, lineOfCodeOutput *gocloc.Result, dataflow *dataflow.DataFlow) (*strings.Builder, bool) {
	severityForFailure := config.Report.Severity
	reportStr := &strings.Builder{}

	reportStr.WriteString("\n\nSecurity Report\n")
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

	if rulesAvailableCount == 0 {
		return reportStr, false
	}

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

	if config.Client == nil {
		reportStr.WriteString("\nReady to take the next step, learn more about Bearer Cloud https://www.bearer.com/bearer-cloud\n")
	}

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
	ruleCountPerLang, totalRuleCount, defaultRulesUsed := countRules(rules, languages, config, false)
	builtInRuleCountPerLang, totalBuiltInRuleCount, builtInRulesUsed := countRules(builtInRules, languages, config, true)

	// combine default and built-in rules per lang
	for _, lang := range maps.Keys(builtInRuleCountPerLang) {
		if ruleCount, ok := ruleCountPerLang[lang]; ok {
			ruleCount.DefaultRuleCount += builtInRuleCountPerLang[lang].DefaultRuleCount
			ruleCountPerLang[lang] = ruleCount
		} else {
			ruleCountPerLang[lang] = builtInRuleCountPerLang[lang]
		}
	}

	totalRuleCount += totalBuiltInRuleCount

	if totalRuleCount == 0 {
		reportStr.WriteString("\n\nZero rules found. A security report requires rules to function. Please check configuration.\n")
		return 0
	}
	reportStr.WriteString("\n\nRules: \n")

	if defaultRulesUsed || builtInRulesUsed {
		reportStr.WriteString(color.HiBlackString(fmt.Sprintf("https://docs.bearer.com/reference/rules [%s]\n\n", config.BearerRulesVersion)))
	}

	tbl := table.New("Language", "Default Rules", "Custom Rules").WithWriter(reportStr)
	for _, lang := range maputil.SortedStringKeys(ruleCountPerLang) {
		ruleCount := ruleCountPerLang[lang]
		tbl.AddRow(lang, ruleCount.DefaultRuleCount, ruleCount.CustomRuleCount)
	}
	tbl.Print()

	return totalRuleCount
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
) (
	ruleCountPerLang map[string]RuleCounter,
	totalRuleCount int,
	defaultRulesUsed bool,
) {
	ruleCountPerLang = make(map[string]RuleCounter)
	totalRuleCount = 0

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

		// increase total count by 1
		totalRuleCount += 1

		defaultRule := strings.HasPrefix(rule.DocumentationUrl, "https://docs.bearer.com") || builtIn
		if ruleCount, ok := ruleCountPerLang[rule.Language()]; ok {
			if defaultRule {
				if !defaultRulesUsed {
					defaultRulesUsed = true
				}
				ruleCount.DefaultRuleCount += 1
				ruleCountPerLang[rule.Language()] = ruleCount
			} else {
				ruleCount.CustomRuleCount += 1
				ruleCountPerLang[rule.Language()] = ruleCount
			}
		} else {
			if defaultRule {
				if !defaultRulesUsed {
					defaultRulesUsed = true
				}
				ruleCountPerLang[rule.Language()] = RuleCounter{
					CustomRuleCount:  0,
					DefaultRuleCount: 1,
				}
			} else {
				ruleCountPerLang[rule.Language()] = RuleCounter{
					CustomRuleCount:  1,
					DefaultRuleCount: 0,
				}
			}
		}
	}

	return ruleCountPerLang, totalRuleCount, defaultRulesUsed
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
	reportStr.WriteString(HighlightCodeExtract(result))
}

func formatSeverity(severity string) string {
	severityColorFn, ok := severityColorFns[severity]
	if !ok {
		return strings.ToUpper(severity)
	}
	return severityColorFn(strings.ToUpper(severity + ": "))
}

func HighlightCodeExtract(record Result) string {
	result := ""
	for _, line := range record.RawCodeExtract {
		if line.Strip {
			result += color.HiBlackString(
				fmt.Sprintf(" %s %s", strings.Repeat(" ", iterativeDigitsCount(line.LineNumber)), line.Extract),
			)
		} else {
			result += color.HiMagentaString(fmt.Sprintf(" %d ", line.LineNumber))
			if line.LineNumber == record.Source.Start && line.LineNumber == record.Source.End {
				for i, char := range line.Extract {
					if i >= record.Source.Column.Start-1 && i < record.Source.Column.End-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber == record.Source.Start && line.LineNumber <= record.Source.End {
				for i, char := range line.Extract {
					if i >= record.Source.Column.Start-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber == record.Source.End && line.LineNumber >= record.Source.Start {
				for i, char := range line.Extract {
					if i <= record.Source.Column.End-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber > record.Source.Start && line.LineNumber < record.Source.End {
				result += color.MagentaString("%s", line.Extract)
			} else {
				result += color.HiMagentaString(line.Extract)
			}
		}

		if line.LineNumber != record.Sink.End {
			result += "\n"
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

func sortByLineNumber(outputs []Output) {
	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].LineNumber < outputs[j].LineNumber
	})
}

func codeExtract(filename string, Source Source, Sink Sink) []file.Line {
	code, err := file.ReadFileSinkLines(
		filename,
		Sink.Start,
		Sink.End,
		Source.Start,
		Source.End,
		settings.CodeExtractBuffer,
	)

	if err != nil {
		return []file.Line{}
	}

	return code
}
