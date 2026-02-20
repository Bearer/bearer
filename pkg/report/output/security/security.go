package security

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/hhatto/gocloc"
	"github.com/rodaine/table"
	"github.com/schollz/progressbar/v3"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	"github.com/bearer/bearer/pkg/scanner/language"
	globaltypes "github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/file"
	ignoretypes "github.com/bearer/bearer/pkg/util/ignore/types"
	"github.com/bearer/bearer/pkg/util/maputil"
	"github.com/bearer/bearer/pkg/util/output"
	bearerprogressbar "github.com/bearer/bearer/pkg/util/progressbar"
	"github.com/bearer/bearer/pkg/util/rego"
	"github.com/bearer/bearer/pkg/util/set"

	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	types "github.com/bearer/bearer/pkg/report/output/security/types"
	stats "github.com/bearer/bearer/pkg/report/output/stats"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
)

var underline = color.New(color.Underline).SprintFunc()
var severityColorFns = map[string]func(x ...interface{}) string{
	globaltypes.LevelCritical: color.New(color.FgRed).SprintFunc(),
	globaltypes.LevelHigh:     color.New(color.FgHiRed).SprintFunc(),
	globaltypes.LevelMedium:   color.New(color.FgYellow).SprintFunc(),
	globaltypes.LevelLow:      color.New(color.FgBlue).SprintFunc(),
	globaltypes.LevelWarning:  color.New(color.FgCyan).SprintFunc(),
}

type ExpectedDetections = []types.ExpectedDetection
type RawFindings = []types.RawFinding
type Findings = map[string][]types.Finding
type IgnoredFindings = map[string][]types.IgnoredFinding

type Input struct {
	RuleId         string                `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule        `json:"rule" yaml:"rule"`
	Dataflow       *outputtypes.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory     `json:"data_categories" yaml:"data_categories"`
}

type RuleCounter struct {
	DefaultRuleCount int
	CustomRuleCount  int
}

type Output struct {
	IsLocal         *bool           `json:"is_local,omitempty" yaml:"is_local,omitempty"`
	Source          types.Source    `json:"source,omitempty" yaml:"source,omitempty"`
	Sink            types.Sink      `json:"sink,omitempty" yaml:"sink,omitempty"`
	LineNumber      int             `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename        string          `json:"filename,omitempty" yaml:"filename,omitempty"`
	FullFilename    string          `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	CategoryGroups  []string        `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	DataType        *types.DataType `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	Severity        string          `json:"severity,omitempty" yaml:"severity,omitempty"`
	DetailedContext string          `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

func AddReportData(
	reportData *outputtypes.ReportData,
	config settings.Config,
	baseBranchFindings *basebranchfindings.Findings,
	hasFiles bool,
) error {
	summaryFindings := make(Findings)
	ignoredSummaryFindings := make(IgnoredFindings)
	reportData.FindingsBySeverity = summaryFindings
	reportData.IgnoredFindingsBySeverity = ignoredSummaryFindings

	if !hasFiles {
		return nil
	}

	dataflow := reportData.Dataflow
	if !config.Scan.Quiet {
		output.StdErrLog("Evaluating rules")
	}

	builtInFingerprints, builtInFailed, err := evaluateRules(summaryFindings, ignoredSummaryFindings, config.BuiltInRules, config, dataflow, baseBranchFindings, true)
	if err != nil {
		return err
	}
	fingerprints, failed, err := evaluateRules(summaryFindings, ignoredSummaryFindings, config.Rules, config, dataflow, baseBranchFindings, false)
	if err != nil {
		return err
	}

	for severity, findingsSlice := range summaryFindings {
		for _, finding := range findingsSlice {
			reportData.RawFindings = append(reportData.RawFindings, finding.ToRawFinding(severity))
		}
	}

	for _, expectedDetectionPerRule := range dataflow.ExpectedDetections {
		for _, location := range expectedDetectionPerRule.Locations {
			reportData.ExpectedDetections = append(reportData.ExpectedDetections, types.ExpectedDetection{
				RuleID: expectedDetectionPerRule.DetectorID,
				Location: types.Location{
					Start: location.Source.StartLineNumber,
					End:   location.Source.EndLineNumber,
					Column: types.Column{
						Start: location.Source.StartColumnNumber,
						End:   location.Source.EndColumnNumber,
					},
				},
			})
		}
	}

	if !config.Scan.Quiet {
		fingerprintOutput(
			append(fingerprints, builtInFingerprints...),
			config.CloudIgnoresUsed,
			config.Report.ExcludeFingerprint,
			config.IgnoredFingerprints,
			config.StaleIgnoredFingerprintIds,
			config.Scan.Diff,
		)
	}

	reportData.ReportFailed = builtInFailed || failed
	return nil
}

func evaluateRules(
	summaryFindings Findings,
	ignoredSummaryFindings IgnoredFindings,
	rules map[string]*settings.Rule,
	config settings.Config,
	dataflow *outputtypes.DataFlow,
	baseBranchFindings *basebranchfindings.Findings,
	builtIn bool,
) ([]string, bool, error) {
	outputFindings := map[string][]types.Finding{}
	ignoredOutputFindings := map[string][]types.IgnoredFinding{}

	var bar *progressbar.ProgressBar
	if !builtIn {
		bar = bearerprogressbar.GetProgressBar(len(rules), &config)
	}

	var fingerprints []string
	failed := false

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
			return fingerprints, false, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return fingerprints, false, err
			}

			var results map[string][]Output
			err = json.Unmarshal(jsonRes, &results)
			if err != nil {
				return fingerprints, false, err
			}
			var ruleSummary *types.Rule
			if config.Report.NoRuleMeta {
				ruleSummary = &types.Rule{
					Title:  rule.Description,
					Id:     rule.Id,
					CWEIDs: rule.CWEIDs,
				}
			} else {
				ruleSummary = &types.Rule{
					Title:            rule.Description,
					Description:      rule.RemediationMessage,
					Id:               rule.Id,
					CWEIDs:           rule.CWEIDs,
					DocumentationUrl: rule.DocumentationUrl,
				}
			}

			instanceCount := make(map[string]int)
			policyFailures := results["policy_failure"]
			sortByLineNumber(policyFailures)

			for i, output := range policyFailures {
				instanceID := instanceCount[output.Filename]
				instanceCount[output.Filename]++

				if baseBranchFindings != nil &&
					baseBranchFindings.Consume(rule.Id, output.Filename, output.Sink.Start, output.Sink.End) {
					continue
				}

				fingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.Filename)
				oldFingerprintId := fmt.Sprintf("%s_%s", rule.Id, output.FullFilename)
				fingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(fingerprintId)), instanceID)
				oldFingerprint := fmt.Sprintf("%x_%d", md5.Sum([]byte(oldFingerprintId)), i)
				fingerprints = append(fingerprints, fingerprint)

				rawCodeExtract := []file.Line{}
				if !config.Report.NoExtract {
					rawCodeExtract = codeExtract(output.FullFilename, output.Source, output.Sink)
				}
				codeExtract := getExtract(rawCodeExtract)

				finding := types.Finding{
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

				ignoredFingerprint, ignored := config.IgnoredFingerprints[fingerprint]
				if !ignored && !config.CloudIgnoresUsed {
					// check for legacy excluded fingerprint
					ignored = config.Report.ExcludeFingerprint[fingerprint]
				}

				severityMeta := CalculateSeverity(finding.CategoryGroups, rule.GetSeverity(), output.IsLocal != nil && *output.IsLocal)
				severity := severityMeta.DisplaySeverity

				if config.Report.Severity.Has(severity) {
					finding.SeverityMeta = severityMeta
					if ignored {
						ignoredOutputFindings[severity] = append(ignoredOutputFindings[severity], types.IgnoredFinding{Finding: finding, IgnoreMeta: ignoredFingerprint})
					} else {
						outputFindings[severity] = append(outputFindings[severity], finding)

						if config.Report.FailOnSeverity.Has(severity) {
							failed = true
						}
					}
				}
			}
		}
	}

	sortFindingsBySeverity(summaryFindings, outputFindings)
	sortFindingsBySeverity(ignoredSummaryFindings, ignoredOutputFindings)

	return fingerprints, failed, nil
}

func sortFindingsBySeverity[F types.GenericFinding](findingsBySeverity map[string][]F, outputFindings map[string][]F) {
	outputFindings = removeDuplicates(outputFindings)

	for severity, findingsSlice := range outputFindings {
		sortFindings(findingsSlice)
		findingsBySeverity[severity] = append(findingsBySeverity[severity], findingsSlice...)
	}
}

func fingerprintOutput(
	fingerprints []string,
	cloudIgnoresUsed bool,
	legacyExcludedFingerprints map[string]bool,
	ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint,
	staleFingerprints []string,
	diffScan bool,
) {
	if cloudIgnoresUsed {
		if len(ignoredFingerprints) > 0 || len(staleFingerprints) > 0 {
			output.StdErrLog("\n=====================================\n")
			if len(ignoredFingerprints) > 0 {
				output.StdErrLog(fmt.Sprintf("%d findings have been ignored from Bearer Cloud", len(ignoredFingerprints)))
			}

			if len(staleFingerprints) > 0 {
				// ignore file entries that have been e.g. re-opened in the Cloud
				output.StdErrLog(fmt.Sprintf("%d fingerprints present in your ignore file are stale and have not been applied", len(staleFingerprints)))
				for _, fingerprintId := range staleFingerprints {
					output.StdErrLog(fmt.Sprintf("  - %s", fingerprintId))
					output.StdErrLog(color.HiBlackString("\tTo remove this fingerprint from your ignore file, run: bearer ignore remove " + fingerprintId))
				}
			}
			output.StdErrLog("\n=====================================\n")
		}
		return
	}
	unusedFingerprints, unusedLegacyFingerprints := removeUnusedFingerprints(
		fingerprints,
		legacyExcludedFingerprints,
		ignoredFingerprints,
	)
	if len(legacyExcludedFingerprints) > 0 || len(unusedFingerprints) > 0 || len(unusedLegacyFingerprints) > 0 {
		output.StdErrLog("\n=====================================\n")
		// legacy
		if len(legacyExcludedFingerprints) > 0 {
			output.StdErrLog(color.HiYellowString("Note: exclude-fingerprints is being replaced by bearer ignore. To use the new ignore functionality, run bearer ignore migrate. See https://docs.bearer.com/reference/commands/#ignore_migrate."))
		}

		if !diffScan { // stale ignored fingerprint warning is misleading for diff scans
			output.StdErrLog("\n")
			if len(unusedLegacyFingerprints) > 0 {
				output.StdErrLog(fmt.Sprintf("%d ignored fingerprints present in your Bearer Configuration file are no longer detected:", len(unusedLegacyFingerprints)))
				for _, fingerprint := range unusedLegacyFingerprints {
					output.StdErrLog(fmt.Sprintf("  - %s", fingerprint))
				}
			}
			// end legacy

			if len(unusedFingerprints) > 0 {
				output.StdErrLog(fmt.Sprintf("%d ignored fingerprints present in your ignore file are no longer detected:", len(unusedFingerprints)))
				for _, fingerprintId := range unusedFingerprints {
					fingerprint, ok := ignoredFingerprints[fingerprintId]
					if !ok {
						// fingerprint will always be found, but if not let's not blow up the scan
						continue
					}

					if fingerprint.Comment == nil {
						output.StdErrLog(fmt.Sprintf("  - %s", fingerprintId))
					} else {
						output.StdErrLog(fmt.Sprintf("  - %s (%s)", fingerprintId, *fingerprint.Comment))
					}
					output.StdErrLog(color.HiBlackString("\tTo remove this fingerprint from your ignore file, run: bearer ignore remove " + fingerprintId))
				}
			}
		}

		output.StdErrLog("\n=====================================")
	}
}

func removeUnusedFingerprints(
	detectedFingerprints []string,
	excludeFingerprints map[string]bool,
	ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint) ([]string, []string) {

	filteredBearerIgnoreFingerprints := make(map[string]bool)
	for fingerprint := range ignoredFingerprints {
		if !slices.Contains(detectedFingerprints, fingerprint) {
			filteredBearerIgnoreFingerprints[fingerprint] = true
		}
	}

	// legacy
	filteredExcludeFingerprints := make(map[string]bool)
	for fingerprint := range excludeFingerprints {
		if !slices.Contains(detectedFingerprints, fingerprint) {
			filteredExcludeFingerprints[fingerprint] = true
		}
	}
	// end legacy

	return slices.Collect(maps.Keys(filteredBearerIgnoreFingerprints)), slices.Collect(maps.Keys(filteredExcludeFingerprints))
}

func getExtract(rawCodeExtract []file.Line) string {
	var parts []string
	for _, line := range rawCodeExtract {
		parts = append(parts, line.Extract)
	}

	return strings.Join(parts, "\n")
}

func BuildReportString(
	reportData *outputtypes.ReportData,
	config settings.Config,
	engine engine.Engine,
	lineOfCodeOutput *gocloc.Result,
) *strings.Builder {
	reportStr := &strings.Builder{}

	if len(reportData.Files) == 0 {
		reportStr.WriteString(
			"\ncouldn't find any files to scan in the specified directory, " +
				"for diff scans this can mean the compared branches were identical",
		)

		return reportStr
	}

	reportStr.WriteString("\n\nSecurity Report\n")
	reportStr.WriteString("\n=====================================")

	initialColorSetting := color.NoColor
	if config.NoColor && !initialColorSetting {
		color.NoColor = true
	}

	rulesAvailableCount := writeRuleListToString(
		reportStr,
		engine,
		config.Rules,
		config.BuiltInRules,
		reportData.Dataflow.Dependencies,
		lineOfCodeOutput.Languages,
		config,
	)

	if rulesAvailableCount == 0 {
		return reportStr
	}

	failures := map[string]map[string]bool{
		globaltypes.LevelCritical: make(map[string]bool),
		globaltypes.LevelHigh:     make(map[string]bool),
		globaltypes.LevelMedium:   make(map[string]bool),
		globaltypes.LevelLow:      make(map[string]bool),
		globaltypes.LevelWarning:  make(map[string]bool),
	}

	for _, severityLevel := range globaltypes.Severities {
		for _, failure := range reportData.FindingsBySeverity[severityLevel] {
			for i := 0; i < len(failure.CWEIDs); i++ {
				failures[severityLevel]["CWE-"+failure.CWEIDs[i]] = true
			}
			writeFailureToString(reportStr, failure, severityLevel)
		}
	}

	if !reportData.ReportFailed {
		reportStr.WriteString("\nNeed to add your own custom rule? Check out the guide: https://docs.bearer.com/guides/custom-rule\n")
	}

	noFailureSummary := checkAndWriteFailureSummaryToString(reportStr, reportData.FindingsBySeverity, rulesAvailableCount, failures, config.Report.Severity)

	if noFailureSummary {
		writeSuccessToString(rulesAvailableCount, reportStr)
		writeStatsToString(reportData, reportStr, config, lineOfCodeOutput)
	}

	color.NoColor = initialColorSetting

	return reportStr
}

func CalculateSeverity(groups []string, severity string, hasLocalDataTypes bool) types.SeverityMeta {
	if severity == globaltypes.LevelWarning {
		return types.SeverityMeta{
			RuleSeverity:    severity,
			DisplaySeverity: globaltypes.LevelWarning,
		}
	}

	return types.SeverityMeta{
		RuleSeverity:    severity,
		DisplaySeverity: severity,
	}
}

func writeStatsToString(
	reportData *outputtypes.ReportData,
	reportStr *strings.Builder,
	config settings.Config,
	lineOfCodeOutput *gocloc.Result,
) {
	if err := stats.AddReportData(reportData, lineOfCodeOutput, config); err != nil {
		return
	}
	if stats.AnythingFoundFor(reportData.Stats) {
		reportStr.WriteString("\nBearer found:\n")
		stats.WriteStatsToString(reportStr, reportData.Stats)
		reportStr.WriteString("\n")
	}
}

func writeRuleListToString(
	reportStr *strings.Builder,
	engine engine.Engine,
	rules map[string]*settings.Rule,
	builtInRules map[string]*settings.Rule,
	reportedDependencies []dataflowtypes.Dependency,
	languages map[string]*gocloc.Language,
	config settings.Config,
) int {
	ruleCountPerLang, totalRuleCount, defaultRulesUsed := countRules(engine, rules, languages, config, false)
	builtInRuleCountPerLang, totalBuiltInRuleCount, builtInRulesUsed := countRules(engine, builtInRules, languages, config, true)

	// combine default and built-in rules per lang
	for lang := range builtInRuleCountPerLang {
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

	tbl := table.New("Language", "Default Rules", "Custom Rules", "Files").WithWriter(reportStr)

	unsupportedLanguages := make(map[string]bool)
	for _, languageFiles := range getFileCountByLanguage(engine, languages) {
		languageID := languageFiles.language.ID()
		displayName := languageFiles.language.DisplayName()
		if ruleCount, ok := ruleCountPerLang[languageID]; ok {
			tbl.AddRow(
				displayName,
				ruleCount.DefaultRuleCount,
				ruleCount.CustomRuleCount,
				languageFiles.fileCount,
			)
		} else {
			for _, reportedDependency := range reportedDependencies {
				if reportedDependency.DetectorLanguage == languageID {
					unsupportedLanguages[languageID] = true
					tbl.AddRow(displayName, 0, 0, languageFiles.fileCount)
					break
				}
			}
		}
	}

	tbl.Print()

	if len(unsupportedLanguages) > 0 {
		sortedUnsupportedLanguages := slices.Sorted(maps.Keys(unsupportedLanguages))
		reportStr.WriteString(fmt.Sprintf(
			"\nWarning: Only partial support is offered for %s.\n",
			strings.Join(sortedUnsupportedLanguages, ", "),
		))
		reportStr.WriteString(color.HiBlackString(
			"For more information, see https://docs.bearer.com/reference/supported-languages\n",
		))
	}

	return totalRuleCount
}

type languageFiles struct {
	language  language.Language
	fileCount int
}

func getFileCountByLanguage(engine engine.Engine, languages map[string]*gocloc.Language) []languageFiles {
	filesByLanguage := make(map[language.Language]int)

	for _, goclocLanguage := range languages {
		language := getLanguageByGocloc(engine, goclocLanguage)
		if language != nil {
			filesByLanguage[language] += len(goclocLanguage.Files)
		}
	}

	result := make([]languageFiles, 0, len(filesByLanguage))
	for language, fileCount := range filesByLanguage {
		result = append(result, languageFiles{
			language:  language,
			fileCount: fileCount,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].fileCount > result[j].fileCount
	})

	return result
}

func countRules(
	engine engine.Engine,
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
		var languageID string

		if rule.IsSecrets() {
			shouldCount = slices.Contains(config.Scan.Scanner, "secrets")
		} else if slices.Contains(config.Scan.Scanner, "sast") {
			languageID = rule.Languages[0]
			language := engine.GetLanguageById(languageID)
			for _, name := range language.GoclocLanguages() {
				if languages[name] != nil {
					shouldCount = true
				}
			}
		}

		if !shouldCount {
			continue
		}

		// increase total count by 1
		totalRuleCount += 1

		defaultRule := strings.HasPrefix(rule.DocumentationUrl, "https://docs.bearer.com") || builtIn
		if ruleCount, ok := ruleCountPerLang[languageID]; ok {
			if defaultRule {
				if !defaultRulesUsed {
					defaultRulesUsed = true
				}
				ruleCount.DefaultRuleCount += 1
				ruleCountPerLang[languageID] = ruleCount
			} else {
				ruleCount.CustomRuleCount += 1
				ruleCountPerLang[languageID] = ruleCount
			}
		} else {
			if defaultRule {
				if !defaultRulesUsed {
					defaultRulesUsed = true
				}
				ruleCountPerLang[languageID] = RuleCounter{
					CustomRuleCount:  0,
					DefaultRuleCount: 1,
				}
			} else {
				ruleCountPerLang[languageID] = RuleCounter{
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
	findings Findings,
	ruleCount int,
	failures map[string]map[string]bool,
	reportedSeverity set.Set[string],
) bool {
	reportStr.WriteString("\n=====================================")

	if len(findings) == 0 {
		return true
	}

	// give summary including counts
	failureCount := 0
	warningCount := 0
	for _, severityLevel := range globaltypes.Severities {
		if severityLevel == globaltypes.LevelWarning {
			warningCount += len(findings[severityLevel])
			continue
		}
		failureCount += len(findings[severityLevel])
	}

	if failureCount == 0 && warningCount == 0 {
		return true
	}

	reportStr.WriteString("\n\n")
	reportStr.WriteString(color.RedString(fmt.Sprint(ruleCount) + " checks, " + fmt.Sprint(failureCount+warningCount) + " findings\n"))

	for _, severityLevel := range globaltypes.Severities {
		if !reportedSeverity.Has(severityLevel) {
			continue
		}
		reportStr.WriteString("\n" + formatSeverity(severityLevel) + fmt.Sprint(len(findings[severityLevel])))
		if len(failures[severityLevel]) > 0 {
			ruleIds := slices.Collect(maps.Keys(failures[severityLevel]))
			sort.Strings(ruleIds)
			if len(ruleIds) > 0 {
				reportStr.WriteString(" (" + strings.Join(ruleIds, ", ") + ")")
			}
		}
	}

	reportStr.WriteString("\n")

	return false
}

func writeFailureToString(reportStr *strings.Builder, finding types.Finding, severity string) {
	reportStr.WriteString("\n\n")
	reportStr.WriteString(formatSeverity(severity))
	reportStr.WriteString(finding.Title)
	cweCount := len(finding.CWEIDs)
	if cweCount > 0 {
		var displayCWEList = []string{}
		for i := 0; i < cweCount; i++ {
			displayCWEList = append(displayCWEList, "CWE-"+finding.CWEIDs[i])
		}
		reportStr.WriteString(" [" + strings.Join(displayCWEList, ", ") + "]")
	}
	reportStr.WriteString("\n")

	if finding.DocumentationUrl != "" {
		reportStr.WriteString(color.HiBlackString(finding.DocumentationUrl + "\n"))
	}

	reportStr.WriteString(color.HiBlackString("To ignore this finding, run: bearer ignore add " + finding.Fingerprint + "\n"))
	reportStr.WriteString("\n")
	if finding.DetailedContext != "" {
		reportStr.WriteString("Detected: " + finding.DetailedContext + "\n\n")
	}
	reportStr.WriteString(color.HiBlueString("File: " + underline(finding.FullFilename+":"+fmt.Sprint(finding.LineNumber)) + "\n"))

	reportStr.WriteString("\n")
	reportStr.WriteString(finding.HighlightCodeExtract())
}

func formatSeverity(severity string) string {
	severityColorFn, ok := severityColorFns[severity]
	if !ok {
		return strings.ToUpper(severity)
	}
	return severityColorFn(strings.ToUpper(severity + ": "))
}

type key struct {
	LineNumber int
	FileName   string
	Detector   string
}

// removeDuplicates removes detections for same detector with same line number by keeping only a single highest severity detection
func removeDuplicates[F types.GenericFinding](data map[string][]F) map[string][]F {
	filteredData := map[string][]F{}

	reportedDetections := set.Set[key]{}

	// filter duplicates
	for _, severity := range globaltypes.Severities {
		findingsSlice, hasSeverity := data[severity]
		if !hasSeverity {
			continue
		}

		for _, genericFinding := range findingsSlice {
			finding := genericFinding.GetFinding()
			key := key{
				LineNumber: finding.LineNumber,
				FileName:   finding.Filename,
				Detector:   finding.Rule.Id,
			}
			if reportedDetections.Add(key) {
				filteredData[severity] = append(filteredData[severity], genericFinding)
			}
		}
	}

	return filteredData
}

func sortFindings[F types.GenericFinding](data []F) {
	sort.Slice(data, func(i, j int) bool {
		vulnerabilityA := data[i].GetFinding()
		vulnerabilityB := data[j].GetFinding()

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

func codeExtract(filename string, Source types.Source, Sink types.Sink) []file.Line {
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

func getLanguageByGocloc(engine engine.Engine, goclocLanguage *gocloc.Language) language.Language {
	for _, language := range engine.GetLanguages() {
		if slices.Contains(language.GoclocLanguages(), goclocLanguage.Name) {
			return language
		}
	}

	return nil
}
