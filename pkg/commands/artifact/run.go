package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bearer/bearer/pkg/report/coverage"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/artifact/scanid"
	"github.com/bearer/bearer/pkg/commands/process/filelist"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	settingsloader "github.com/bearer/bearer/pkg/commands/process/settings/loader"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/report/output/stats"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	scannerstats "github.com/bearer/bearer/pkg/scanner/stats"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/ignore"
	ignoretypes "github.com/bearer/bearer/pkg/util/ignore/types"
	outputhandler "github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/bearer/bearer/pkg/version_check"

	"github.com/bearer/bearer/pkg/types"
)

// TargetKind represents what kind of artifact bearer scans
type TargetKind string

const (
	TargetFilesystem TargetKind = "fs"
	TargetRepository TargetKind = "repo"
)

type ScannerConfig struct {
	Target   string
	Artifact types.Artifact
}

type Runner interface {
	// Cached returns true if cached data was used in scan
	CacheUsed() bool
	// ReportPath returns the filename of the report
	ReportPath() string
	// Scan gathers the findings
	Scan(ctx context.Context, opts flagtypes.Options) ([]files.File, *basebranchfindings.Findings, error)
	// Report a writes a report
	Report(files []files.File, baseBranchFindings *basebranchfindings.Findings) (bool, error)
}

type runner struct {
	targetPath,
	reportPath string
	reuseDetection bool
	goclocResult   *gocloc.Result
	scanSettings   settings.Config
	stats          *scannerstats.Stats
	gitContext     *gitrepository.Context
	engine         engine.Engine
	discovery      coverage.Discovery
}

// NewRunner initializes Runner that provides scanning functionalities.
func NewRunner(
	ctx context.Context,
	scanSettings settings.Config,
	gitContext *gitrepository.Context,
	targetPath string,
	goclocResult *gocloc.Result,
	stats *scannerstats.Stats,
	engine engine.Engine,
) (Runner, error) {
	r := &runner{
		scanSettings: scanSettings,
		targetPath:   targetPath,
		goclocResult: goclocResult,
		stats:        stats,
		gitContext:   gitContext,
		engine:       engine,
	}

	scanID, err := scanid.Build(scanSettings, gitContext)
	if err != nil {
		return nil, fmt.Errorf("failed to build scan id for caching: %w", err)
	}

	path := os.TempDir() + "/bearer" + scanID
	completedPath := strings.Replace(path, ".jsonl", "-completed.jsonl", 1)

	r.reportPath = path

	log.Debug().Msgf("creating report %s", path)

	if _, err := os.Stat(completedPath); err == nil {
		// diff can't use the cache because the base branch scan data is not in the report
		if !scanSettings.Scan.Force && !scanSettings.Scan.Diff {
			// force is not set, and we are not running a diff scan
			r.reuseDetection = true
			log.Debug().Msgf("reuse detection for %s", path)
			r.reportPath = completedPath

			return r, nil
		} else {
			if _, err = os.Stat(path); err == nil {
				err := os.Remove(path)
				if err != nil {
					log.Error().Msgf("couldn't remove report path %s, %s", path, err.Error())
				}
			}

			err = os.Remove(completedPath)
			if err != nil {
				log.Error().Msgf("couldn't remove completed report path %s, %s", completedPath, err.Error())
			}
		}
	}

	pathCreated, err := os.Create(path)

	log.Debug().Msgf("successfully created reportPath %s", path)

	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	if err := pathCreated.Close(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *runner) CacheUsed() bool {
	return r.reuseDetection
}

func (r *runner) Scan(ctx context.Context, opts flagtypes.Options) (result []files.File, baseResult *basebranchfindings.Findings, retErr error) {
	defer func() {
		if retErr != nil {
			retErr = errors.Join(retErr, r.closeReportInputs())
		}
	}()
	if r.reuseDetection {
		return nil, nil, nil
	}

	if !opts.Quiet {
		outputhandler.StdErrLog(fmt.Sprintf("Scanning target %s", opts.Target))
	}

	repository, err := gitrepository.New(ctx, r.scanSettings, r.targetPath, r.gitContext)
	if err != nil {
		return nil, nil, fmt.Errorf("git repository error: %w", err)
	}

	fileList, err := filelist.Discover(repository, r.targetPath, r.goclocResult, r.scanSettings)
	if err != nil {
		return nil, nil, err
	}

	r.discovery = fileList.Coverage
	if err := r.validateDiscovery(opts.Diff); err != nil {
		return nil, nil, err
	}
	var baseBranchFindings *basebranchfindings.Findings
	if err := repository.WithBaseBranch(func() error {
		if !opts.Quiet {
			outputhandler.StdErrLog(fmt.Sprintf("\nScanning base branch %s", r.gitContext.BaseBranch))
		}

		baseBranchFindings, err = r.scanBaseBranch(fileList)
		if err != nil {
			return err
		}

		if !opts.Quiet {
			outputhandler.StdErrLog("\nScanning current branch")
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	// The current Engine interface cannot cancel a scan without racing its
	// lazy initialization. Do not concurrently call Close here. Check cancellation
	// before and after processing; the caller must still bound the process lifetime.
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	scanInputs := fileList.Files
	if opts.Diff {
		// Diff discovery preserves the base/head mapping. Evaluate eligibility
		// separately in each checked-out tree; this mode still has no coverage.
		scanInputs, err = r.configuredInputs(scanInputs)
		if err != nil {
			return nil, nil, err
		}
	}
	scanErr := r.scanWithInputs(r.reportPath, scanInputs)
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	if err := scanErr; err != nil {
		return nil, nil, err
	}

	return scanInputs, baseBranchFindings, nil
}

func (r *runner) scanBaseBranch(fileList *files.List) (resultFindings *basebranchfindings.Findings, retErr error) {
	defer func() { retErr = errors.Join(retErr, r.closeReportInputs()) }()
	result := basebranchfindings.New(fileList)

	eligibleBase, err := r.configuredInputs(fileList.BaseFiles)
	if err != nil {
		return nil, err
	}
	if len(eligibleBase) == 0 {
		return result, nil
	}

	if err := r.scanWithInputs(r.reportPath+".base", eligibleBase); err != nil {
		return nil, err
	}

	report := types.Report{
		Path:        r.reportPath + ".base",
		Inputgocloc: r.goclocResult,
		HasFiles:    len(eligibleBase) != 0,
	}

	reportData, err := reportoutput.GetData(report, r.scanSettings, r.gitContext, nil)
	if err != nil {
		return nil, err
	}

	for _, findings := range reportData.FindingsBySeverity {
		for _, finding := range findings {
			result.Add(finding.Rule.Id, finding.Filename, finding.Sink.Start, finding.Sink.End)
		}
	}

	return result, nil
}

func getIgnoredFingerprints(settings settings.Config) (
	useCloudIgnores bool,
	ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint,
	staleIgnoredFingerprintIds []string,
	err error,
) {
	localIgnoredFingerprints, _, _, err := ignore.GetIgnoredFingerprints(settings.IgnoreFile, &settings.Target)
	if err != nil {
		return useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, err
	}

	return false, localIgnoredFingerprints, []string{}, nil
}

type ReportFailedError int

func (exitcode ReportFailedError) Error() string {
	return "Report failed with exitcode"
}

// Run performs artifact scanning
func Run(ctx context.Context, opts flagtypes.Options, engine engine.Engine) (err error) {
	targetPath, err := file.CanonicalPath(opts.Target)
	if err != nil {
		return fmt.Errorf("failed to get absolute target: %w", err)
	}

	if err := validateLanguages(engine, opts.Language); err != nil {
		return err
	}

	inputgocloc, err := stats.GoclocDetectorOutput(targetPath, opts)
	if err != nil {
		log.Debug().Msgf("Error in line of code output %s", err)
		return err
	}

	if allowed := allowedGoclocLanguages(engine, opts.Language); allowed != nil {
		filterGoclocResult(inputgocloc, allowed)
	}

	foundLanguageIDs := GetFoundLanguageIDs(engine, inputgocloc.Languages)

	// set used language list for external rules to empty if we dont use them
	metaLanguageList := foundLanguageIDs
	if opts.RuleOptions.DisableDefaultRules {
		metaLanguageList = make([]string, 0)
	}

	versionMeta, err := version_check.GetScanVersionMeta(ctx, opts, metaLanguageList)
	if err != nil {
		log.Debug().Msgf("failed: %s", err)
	} else {
		version_check.DisplayBinaryVersionWarning(versionMeta, opts.ScanOptions.Quiet)
	}

	gitContext, err := gitrepository.NewContext(&opts)
	if err != nil {
		return fmt.Errorf("failed to get git context: %w", err)
	}

	if opts.Diff && gitContext == nil {
		return errors.New("--diff option requires a git repository")
	}

	if !opts.Quiet {
		outputhandler.StdErrLog("Loading rules")
	}

	if err := engine.Initialize(opts.LogLevel); err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	scanSettings, err := settingsloader.FromOptions(opts, versionMeta, engine, foundLanguageIDs)
	scanSettings.Target = opts.Target
	if err != nil {
		return err
	}
	scanSettings.CloudIgnoresUsed, scanSettings.IgnoredFingerprints, scanSettings.StaleIgnoredFingerprintIds, err = getIgnoredFingerprints(
		scanSettings,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, scanSettings.Worker.Timeout)
	defer cancel()

	defer func() {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn().Msg("Increase --timeout value")
		}
	}()

	var stats *scannerstats.Stats
	if scanSettings.Debug {
		stats = scannerstats.New()
	}

	r, err := NewRunner(ctx, scanSettings, gitContext, targetPath, inputgocloc, stats, engine)
	if err != nil {
		return err
	}

	// Cover cancellation between Scan and Report as well as report failures.
	defer func() {
		if owned, ok := r.(*runner); ok {
			err = errors.Join(err, owned.closeReportInputs())
		}
	}()
	files, baseBranchFindings, err := r.Scan(ctx, opts)
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	reportFailed, err := r.Report(files, baseBranchFindings)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("report error: %w", err)
	} else {
		reportPath := r.ReportPath()
		if !strings.HasSuffix(reportPath, "-completed.jsonl") {
			newPath := strings.Replace(reportPath, ".jsonl", "-completed.jsonl", 1)
			log.Debug().Msgf("renaming report %s -> %s", reportPath, newPath)
			err := os.Rename(reportPath, newPath)
			if err != nil {
				return fmt.Errorf("failed to rename report file %s -> %s: %w", reportPath, newPath, err)
			}
		}
	}

	if stats != nil {
		outputhandler.StdErrLog(fmt.Sprintf("=====================================\n\nProfile\n\n%s", stats.String()))
	}

	if reportFailed {
		if scanSettings.Scan.ExitCode == -1 {
			return ReportFailedError(1)
		} else {
			return ReportFailedError(scanSettings.Scan.ExitCode)
		}
	}

	return nil
}

func (r *runner) Report(
	files []files.File,
	baseBranchFindings *basebranchfindings.Findings,
) (failed bool, retErr error) {
	defer func() { retErr = errors.Join(retErr, r.closeReportInputs()) }()
	startTime := time.Now()
	cacheUsed := r.CacheUsed()
	r.scanSettings.CacheUsed = cacheUsed

	report := types.Report{
		Path:        r.reportPath,
		Inputgocloc: r.goclocResult,
		HasFiles:    r.CacheUsed() || len(files) != 0,
	}

	// Validate the entire input stream before formatting. Malformed or
	// truncated records must not become a successful cached or fresh report.
	raw, err := os.ReadFile(r.reportPath)
	if err != nil {
		return false, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	fileLists := 0
	var contextReceipts []*file.ContextReceipt
	var reportInputNames []string
	for {
		var rawRecord json.RawMessage
		err = decoder.Decode(&rawRecord)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if err = coverage.CheckJSON(rawRecord); err != nil {
			return false, err
		}
		var record map[string]any
		recordDecoder := json.NewDecoder(bytes.NewReader(rawRecord))
		// The envelope is typed below; opaque values must not pass through float64.
		recordDecoder.UseNumber()
		if err = recordDecoder.Decode(&record); err != nil {
			return false, err
		}
		if record == nil {
			return false, errors.New("null report record")
		}
		kind, ok := record["type"].(string)
		if !ok {
			return false, errors.New("report record type required")
		}
		if err := rejectLegacyAliases(record, []string{"type", "detector_type", "detector_language", "commit_sha", "framework_detection_type", "source", "value", "filenames", "coverage", "message", "file", "file_size", "timeout_duration", "error"}, false); err != nil {
			return false, err
		}
		switch kind {
		case "file_list":
			if err := rejectLegacyAliases(record, []string{"type", "filenames"}, true); err != nil {
				return false, err
			}
			fileLists++
			names, ok := record["filenames"].([]any)
			if !ok || fileLists != 1 {
				return false, errors.New("invalid/duplicate report file list")
			}
			seen := map[string]bool{}
			for _, value := range names {
				name, ok := value.(string)
				if !ok || name == "" || seen[name] {
					return false, errors.New("invalid report filename")
				}
				seen[name] = true
				reportInputNames = append(reportInputNames, name)
			}
		case file.ContextReceiptType:
			receipt, e := file.DecodeContextReceipt(rawRecord)
			if e != nil {
				return false, e
			}
			contextReceipts = append(contextReceipts, receipt)
		case "processing_coverage":
			payload, ok := record["coverage"]
			if !ok {
				return false, errors.New("missing coverage record")
			}
			encoded, err := json.Marshal(payload)
			if err != nil {
				return false, err
			}
			if _, err := coverage.Decode(encoded); err != nil {
				return false, err
			}
		case "error", "file_error":
			return false, errors.New("report records processing failure")
		case "dependency", "dependency_classified", "interface", "interface_classified", "schema", "schema_classified", "create_view", "framework", "framework_classified", "filler", "secret_leak", "custom", "custom_classified", "custom_risk", "expected_detection", "operation":
			if err := validateLegacyDetection(rawRecord); err != nil {
				return false, err
			}

		default:
			return false, errors.New("unknown report record type")
		}
	}
	if fileLists != 1 {
		return false, errors.New("report file list unavailable")
	}

	seenContext := map[string]bool{}
	listedContext := map[string]bool{}
	for _, name := range reportInputNames {
		listedContext[name] = true
	}
	for _, receipt := range contextReceipts {
		if seenContext[receipt.Input.Relative] || !listedContext[receipt.Input.Relative] {
			return false, errors.New("duplicate/unlisted context handoff")
		}
		seenContext[receipt.Input.Relative] = true
	}
	if !cacheUsed {
		if len(contextReceipts) != len(reportInputNames) {
			return false, errors.New("missing context handoff")
		}
		scope := r.scanSettings.ReportInputs
		if scope == nil {
			return false, errors.New("fresh context handoff scope unavailable")
		}
		seen := map[string]bool{}
		for _, receipt := range contextReceipts {
			if seen[receipt.Input.Relative] {
				return false, errors.New("duplicate context handoff")
			}
			seen[receipt.Input.Relative] = true
			if e := scope.ImportContextReceipt(receipt, nil); e != nil {
				return false, e
			}
		}
		if e := scope.RequireContextReceipts(reportInputNames); e != nil {
			return false, e
		}
	}
	if err := r.prepareCachedStatisticsInputs(cacheUsed); err != nil {
		return false, err
	}

	if cacheUsed && !r.scanSettings.Scan.Quiet {
		// output cached data warning at start of report
		outputhandler.StdErrLog("Using cached data")
	}

	reportData, err := reportoutput.GetData(report, r.scanSettings, r.gitContext, baseBranchFindings)
	if err != nil {
		return false, err
	}

	if err := r.bindReportCoverage(reportData, cacheUsed); err != nil {
		return false, err
	}
	endTime := time.Now()

	reportSupported := anySupportedLanguagesPresent(r.engine, report.Inputgocloc, r.scanSettings)
	if !reportSupported && r.scanSettings.Report.Report != flag.ReportPrivacy && !r.scanSettings.Scan.Quiet {
		var placeholderStr *strings.Builder
		placeholderStr, err = getPlaceholderOutput(reportData, report, r.scanSettings, report.Inputgocloc)
		if err != nil {
			return false, err
		}

		if err := r.acceptReport(func() error { return deliverReport(r.scanSettings.Report.Output, []byte(placeholderStr.String())) }); err != nil {
			return false, err
		}
		return true, nil
	}

	formatStr, err := reportoutput.FormatOutput(
		reportData,
		r.scanSettings,
		r.engine,
		report.Inputgocloc,
		startTime,
		endTime,
	)
	if err != nil {
		return false, fmt.Errorf("error generating report %s", err)
	}

	if !r.scanSettings.Scan.Quiet {
		// add cached data warning message
		if cacheUsed {
			outputhandler.StdErrLog("Cached data used (no code changes detected). Unexpected? Use --force to force a re-scan.\n")
		}
		// add cloud info message
		if r.scanSettings.Client != nil {
			if r.scanSettings.Client.Error == nil {
				outputhandler.StdErrLog("Data successfully sent to Bearer Cloud.")
			} else {
				// client error
				outputhandler.StdErrLog(fmt.Sprintf("Failed to send data to Bearer Cloud. %s ", *r.scanSettings.Client.Error))
			}
		}
	}

	if r.scanSettings.LoadedRuleCount == 0 && slices.Contains(r.scanSettings.Scan.Scanner, flag.ScannerSAST) && r.scanSettings.Report.Report == flag.ReportSecurity {
		return false, fmt.Errorf("%d rules found for supported language, default rules could not be downloaded or possibly disabled without using --external-rule-dir", len(r.scanSettings.Rules))
	}

	if err := r.acceptReport(func() error { return deliverReport(r.scanSettings.Report.Output, []byte(formatStr)) }); err != nil {
		return false, err
	}
	return reportData.ReportFailed, nil
}

func (r *runner) ReportPath() string {
	return r.reportPath
}

func anySupportedLanguagesPresent(engine engine.Engine, inputgocloc *gocloc.Result, config settings.Config) bool {
	if inputgocloc == nil {
		return true
	}

	ruleLanguages := make(map[string]bool)
	for _, rule := range config.Rules {
		for _, language := range rule.Languages {
			ruleLanguages[language] = true
		}
	}

	for _, goclocLanguage := range inputgocloc.Languages {
		for _, language := range engine.GetLanguages() {
			if slices.Contains(language.GoclocLanguages(), goclocLanguage.Name) {
				return true
			}
		}
	}

	log.Debug().Msg("No language found for which rules are applicable")
	return false
}

func getPlaceholderOutput(reportData *outputtypes.ReportData, report types.Report, config settings.Config, inputgocloc *gocloc.Result) (outputStr *strings.Builder, err error) {
	if err := reportoutput.GetDataflow(reportData, report, config, true); err != nil {
		return nil, err
	}

	return stats.GetPlaceholderOutput(reportData, inputgocloc, config)
}

func GetFoundLanguageIDs(engine engine.Engine, goclocLanguages map[string]*gocloc.Language) []string {
	foundLanguages := set.New[string]()

	for _, goclocLanguage := range goclocLanguages {
		for _, language := range engine.GetLanguages() {
			if slices.Contains(language.GoclocLanguages(), goclocLanguage.Name) {
				foundLanguages.Add(language.ID())
			}
		}
	}

	keys := foundLanguages.Items()
	sort.Strings(keys)

	return keys
}

func validateLanguages(engine engine.Engine, requested []string) error {
	if len(requested) == 0 {
		return nil
	}

	valid := set.New[string]()
	for _, language := range engine.GetLanguages() {
		valid.Add(language.ID())
	}

	var unknown []string
	for _, id := range requested {
		if !valid.Has(id) {
			unknown = append(unknown, id)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	validIDs := valid.Items()
	sort.Strings(validIDs)

	return fmt.Errorf(
		"invalid language argument: %s; supported values: %s",
		strings.Join(unknown, ", "),
		strings.Join(validIDs, ", "),
	)
}

func filterGoclocResult(result *gocloc.Result, allowed set.Set[string]) {
	if result == nil {
		return
	}

	for path, file := range result.Files {
		if !allowed.Has(file.Lang) {
			delete(result.Files, path)
		}
	}

	for name := range result.Languages {
		if !allowed.Has(name) {
			delete(result.Languages, name)
		}
	}

	total := &gocloc.Language{Name: "TOTAL"}
	for _, language := range result.Languages {
		total.Code += language.Code
		total.Comments += language.Comments
		total.Blanks += language.Blanks
		total.Total += language.Total
		total.Files = append(total.Files, language.Files...)
	}
	result.Total = total
}

func allowedGoclocLanguages(engine engine.Engine, requested []string) set.Set[string] {
	if len(requested) == 0 {
		return nil
	}

	enabled := set.New[string]()
	enabled.AddAll(requested)

	allowed := set.New[string]()
	for _, language := range engine.GetLanguages() {
		if enabled.Has(language.ID()) {
			allowed.AddAll(language.GoclocLanguages())
		}
	}

	return allowed
}

// Only the explicit unsupported diff inventory is compatible with a legacy
// findings report. Other incomplete discovery must never enter this route.
func (r *runner) unsupportedDiffCoverage() bool {
	d := r.discovery
	return r.scanSettings.Scan.Diff && r.gitContext != nil && r.gitContext.BaseCommitHash != "" &&
		!d.Complete && len(d.Files) == 0 && len(d.Excluded) == 0 && len(d.Errors) == 1 &&
		d.Errors[0] == "git diff discovery does not expose a supported same-run terminal file inventory"
}

func (r *runner) validateDiscovery(diff bool) error {
	if diff != r.scanSettings.Scan.Diff {
		return errors.New("inconsistent diff mode")
	}
	if diff && r.unsupportedDiffCoverage() {
		return nil
	}
	if !r.discovery.Complete || len(r.discovery.Errors) > 0 {
		return errors.New("incomplete input discovery")
	}
	return nil
}

func (r *runner) bindReportCoverage(data *outputtypes.ReportData, cached bool) error {
	if cached && r.CacheUsed() && !r.scanSettings.Scan.Force && !r.scanSettings.Scan.Diff {
		// NewRunner selected an existing completed cache. Its findings may render,
		// but even valid prior-run coverage cannot describe processing in this run.
		if data.Dataflow == nil || len(data.Dataflow.Errors) != 0 {
			return errors.New("cached processing diagnostics incomplete")
		}
		if data.Coverage != nil {
			if err := data.Coverage.Validate(data.Coverage.Discovery.Files); err != nil {
				return err
			}
		}
		clearUnavailableCoverage(data)
		return nil
	}

	if !r.unsupportedDiffCoverage() {
		return bindCoverage(data, r.discovery, cached)
	}
	// Preserve processing/diagnostic failures even though discovery coverage for
	// this mode is unavailable. A successful findings report is not coverage proof.
	if cached {
		return errors.New("cached scan has no fresh same-run discovery")
	}
	if data.Coverage == nil {
		return errors.New("processing coverage unavailable")
	}
	if err := data.Coverage.Validate(data.Coverage.Discovery.Files); err != nil {
		return err
	}
	if data.Dataflow == nil || len(data.Dataflow.Errors) != 0 {
		return errors.New("processing diagnostics incomplete")
	}
	clearUnavailableCoverage(data)
	return nil
}

// bindCoverage joins the actual discovery and worker records. It does not claim
// coverage of a caller's external promise; that must be reconciled independently.
func bindCoverage(data *outputtypes.ReportData, discovery coverage.Discovery, cached bool) error {
	if cached {
		return errors.New("cached scan has no fresh same-run discovery")
	}
	if data.Coverage == nil {
		return errors.New("processing coverage unavailable")
	}
	if err := data.Coverage.Validate(discovery.Files); err != nil {
		return err
	}
	if len(data.Coverage.Discovery.Files) != len(discovery.Files) {
		return errors.New("discovery/processing inventory mismatch")
	}
	if data.Dataflow == nil || len(data.Dataflow.Errors) != 0 {
		return errors.New("processing diagnostics incomplete")
	}
	bound := *data.Coverage
	bound.Discovery = discovery
	if err := bound.Validate(discovery.Files); err != nil {
		return err
	}
	// Detector output retains raw records independently of the security coverage
	// field. Bind both representations only after all reconciliation succeeds.
	encoded, err := json.Marshal(&bound)
	if err != nil {
		return err
	}
	var rawCoverage map[string]any
	if err := json.Unmarshal(encoded, &rawCoverage); err != nil {
		return err
	}
	data.Coverage = &bound
	for _, record := range data.Detectors {
		if raw, ok := record.(map[string]any); ok && raw["type"] == "processing_coverage" {
			raw["coverage"] = rawCoverage
		}
	}
	return nil
}

type deliveryFile interface {
	io.Writer
	Sync() error
	Close() error
	Name() string
}
type deliveryOps struct {
	create func(string) (deliveryFile, error)
	rename func(string, string) error
	remove func(string) error
}

func deliverReport(destination string, data []byte) error {
	if destination == "" {
		return coverage.WriteFull(os.Stdout, data)
	}
	return publishReport(destination, data, deliveryOps{
		create: func(dir string) (deliveryFile, error) { return os.CreateTemp(dir, ".bearer-report-*") }, rename: os.Rename, remove: os.Remove,
	})
}
func publishReport(destination string, data []byte, ops deliveryOps) (retErr error) {
	f, err := ops.create(filepath.Dir(destination))
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			retErr = errors.Join(retErr, f.Close())
		}
		_ = ops.remove(f.Name())
	}()
	if err = coverage.WriteFull(f, data); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	err = f.Close()
	closed = true
	if err != nil {
		return err
	}
	if err = ops.rename(f.Name(), destination); err != nil {
		return err
	}
	return nil
}

// rejectLegacyAliases preserves unrelated legacy extensions, but never permits
// an alternate spelling that encoding/json would fold onto a declared field.
func rejectLegacyAliases(record map[string]any, keys []string, exact bool) error {
	for key := range record {
		known := false
		for _, canonical := range keys {
			if key == canonical {
				known = true
				break
			}
			if strings.EqualFold(key, canonical) {
				return errors.New("aliased legacy report field")
			}
		}
		if exact && !known {
			return errors.New("unknown legacy report field")
		}
	}
	return nil
}

func validateLegacyDetection(raw json.RawMessage) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	checkString := func(fields map[string]json.RawMessage, key string, required bool) error {
		raw, present := fields[key]
		if !present {
			if required {
				return fmt.Errorf("missing legacy %s", key)
			}
			return nil
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return fmt.Errorf("null legacy %s", key)
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("invalid legacy %s", key)
		}
		return nil
	}
	for _, key := range []string{"type", "detector_type"} {
		if err := checkString(fields, key, true); err != nil {
			return err
		}
	}
	for _, key := range []string{"detector_language", "commit_sha", "framework_detection_type"} {
		if err := checkString(fields, key, false); err != nil {
			return err
		}
	}
	if _, ok := fields["value"]; !ok {
		return errors.New("missing legacy detection value")
	}
	var sourceFields map[string]json.RawMessage
	if err := json.Unmarshal(fields["source"], &sourceFields); err != nil || sourceFields == nil {
		return errors.New("invalid legacy detection source")
	}
	sourceKeys := make(map[string]any, len(sourceFields))
	for key := range sourceFields {
		sourceKeys[key] = nil
	}
	if err := rejectLegacyAliases(sourceKeys, []string{"filename", "full_filename", "language", "language_type", "start_line_number", "start_column_number", "end_line_number", "end_column_number", "text"}, false); err != nil {
		return err
	}
	if err := checkString(sourceFields, "filename", true); err != nil {
		return err
	}
	for _, key := range []string{"full_filename", "language", "language_type"} {
		if err := checkString(sourceFields, key, false); err != nil {
			return err
		}
	}
	for _, key := range []string{"start_line_number", "start_column_number", "end_line_number", "end_column_number"} {
		if raw, present := sourceFields[key]; present {
			var value *int
			if err := json.Unmarshal(raw, &value); err != nil {
				return fmt.Errorf("invalid legacy source %s", key)
			}
		}
	}
	if raw, present := sourceFields["text"]; present {
		var value *string
		if err := json.Unmarshal(raw, &value); err != nil {
			return errors.New("invalid legacy source text")
		}
	}
	return nil
}

// Keep unavailable coverage out of both public representations. A non-nil raw
// slice also prevents a later formatter/placeholder from reloading the cache.
func clearUnavailableCoverage(data *outputtypes.ReportData) {
	data.Coverage = nil
	records := make([]any, 0, len(data.Detectors))
	for _, record := range data.Detectors {
		if value, ok := record.(map[string]any); ok && value["type"] == "processing_coverage" {
			continue
		}
		records = append(records, record)
	}
	data.Detectors = records
}

// Select scheduled inputs for the current checked-out view without mutating the
// diff mapping or any independently fixed external coverage promise.
func (r *runner) configuredInputs(inputs []files.File) ([]files.File, error) {
	eligible := make([]files.File, 0, len(inputs))
	for _, input := range inputs {
		reason, err := file.ConfiguredExclusion(r.targetPath, input.FilePath, r.scanSettings.Scan.SkipTest, r.scanSettings.Scan.SkipGitIgnore)
		if err != nil {
			return nil, err
		}
		if reason == "" {
			eligible = append(eligible, input)
		}
	}
	return eligible, nil
}

func (r *runner) scanWithInputs(reportPath string, inputs []files.File) error {
	if r.scanSettings.ReportInputs != nil {
		return errors.New("previous report input scope still owned")
	}
	scope, err := file.OpenInputScope(r.targetPath)
	if err != nil {
		return err
	}
	r.scanSettings.ReportInputs = scope
	return r.engine.Scan(&r.scanSettings, r.stats, reportPath, r.targetPath, inputs)
}
func (r *runner) validateReportInputs() error {
	if r.scanSettings.ReportInputs != nil {
		return r.scanSettings.ReportInputs.Validate()
	}
	return nil
}
func (r *runner) closeReportInputs() error {
	scope := r.scanSettings.ReportInputs
	r.scanSettings.ReportInputs = nil
	if scope == nil {
		return nil
	}
	return scope.Close()
}

// Reconcile after delivery as well as before it. A delivered partial report is
// not accepted when final input/descriptor close fails.
func (r *runner) acceptReport(deliver func() error) (retErr error) {
	defer func() { retErr = errors.Join(retErr, r.closeReportInputs()) }()
	if err := r.validateReportInputs(); err != nil {
		return err
	}
	return deliver()
}

// Cached statistics and finding excerpts describe current display context,
// not processing by the cached scan. Only a genuinely selected, structurally validated cache may
// obtain this scope. Fresh scans must retain their original processing scope.
// Report's error defer and acceptReport own it through delivery/final close.
func (r *runner) prepareCachedStatisticsInputs(cached bool) error {
	if !cached || r.scanSettings.ReportInputs != nil {
		return nil
	}
	scope, err := file.OpenInputScope(r.targetPath)
	if err != nil {
		return err
	}
	r.scanSettings.ReportInputs = scope
	return nil
}
