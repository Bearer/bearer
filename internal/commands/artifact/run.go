package artifact

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"golang.org/x/exp/maps"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/internal/commands/artifact/scanid"
	"github.com/bearer/bearer/internal/commands/process/filelist"
	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/commands/process/gitrepository"
	"github.com/bearer/bearer/internal/commands/process/orchestrator"
	"github.com/bearer/bearer/internal/commands/process/orchestrator/work"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	flagtypes "github.com/bearer/bearer/internal/flag/types"
	"github.com/bearer/bearer/internal/report/basebranchfindings"
	reportoutput "github.com/bearer/bearer/internal/report/output"
	"github.com/bearer/bearer/internal/report/output/stats"
	outputtypes "github.com/bearer/bearer/internal/report/output/types"
	scannerstats "github.com/bearer/bearer/internal/scanner/stats"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/ignore"
	ignoretypes "github.com/bearer/bearer/internal/util/ignore/types"
	outputhandler "github.com/bearer/bearer/internal/util/output"
	"github.com/bearer/bearer/internal/version_check"

	"github.com/bearer/bearer/internal/types"
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
}

// NewRunner initializes Runner that provides scanning functionalities.
func NewRunner(
	ctx context.Context,
	scanSettings settings.Config,
	gitContext *gitrepository.Context,
	targetPath string,
	goclocResult *gocloc.Result,
	stats *scannerstats.Stats,
) (Runner, error) {
	r := &runner{
		scanSettings: scanSettings,
		targetPath:   targetPath,
		goclocResult: goclocResult,
		stats:        stats,
		gitContext:   gitContext,
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
		log.Error().Msgf("failed to create path %s, %s, %#v", path, err.Error(), pathCreated)
	}

	return r, nil
}

func (r *runner) CacheUsed() bool {
	return r.reuseDetection
}

func (r *runner) Scan(ctx context.Context, opts flagtypes.Options) ([]files.File, *basebranchfindings.Findings, error) {
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

	orchestrator, err := orchestrator.New(
		work.Repository{Dir: r.targetPath},
		r.scanSettings,
		r.stats,
		len(fileList.Files),
	)
	if err != nil {
		return nil, nil, err
	}
	defer orchestrator.Close()

	var baseBranchFindings *basebranchfindings.Findings
	if err := repository.WithBaseBranch(func() error {
		if !opts.Quiet {
			outputhandler.StdErrLog(fmt.Sprintf("\nScanning base branch %s", r.gitContext.BaseBranch))
		}

		baseBranchFindings, err = r.scanBaseBranch(orchestrator, fileList)
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

	if err := orchestrator.Scan(r.reportPath, fileList.Files); err != nil {
		return nil, nil, err
	}

	return fileList.Files, baseBranchFindings, nil
}

func (r *runner) scanBaseBranch(
	orchestrator *orchestrator.Orchestrator,
	fileList *files.List,
) (*basebranchfindings.Findings, error) {
	result := basebranchfindings.New(fileList)

	if len(fileList.BaseFiles) == 0 {
		return result, nil
	}

	if err := orchestrator.Scan(r.reportPath+".base", fileList.BaseFiles); err != nil {
		return nil, err
	}

	report := types.Report{
		Path:        r.reportPath + ".base",
		Inputgocloc: r.goclocResult,
		HasFiles:    len(fileList.BaseFiles) != 0,
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

func getIgnoredFingerprints(client *api.API, settings settings.Config, gitContext *gitrepository.Context, pullRequestNumber string) (
	useCloudIgnores bool,
	ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint,
	staleIgnoredFingerprintIds []string,
	err error,
) {
	localIgnoredFingerprints, _, _, err := ignore.GetIgnoredFingerprints(settings.IgnoreFile, &settings.Target)
	if err != nil {
		return useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, err
	}

	if client != nil && client.Error == nil {
		useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, err = ignore.GetIgnoredFingerprintsFromCloud(
			client,
			gitContext.FullName,
			pullRequestNumber,
			localIgnoredFingerprints,
		)
		if err != nil {
			return useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, err
		}
	}

	if useCloudIgnores {
		return useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, nil
	}

	return false, localIgnoredFingerprints, []string{}, nil
}

// Run performs artifact scanning
func Run(ctx context.Context, opts flagtypes.Options) (err error) {
	targetPath, err := file.CanonicalPath(opts.Target)
	if err != nil {
		return fmt.Errorf("failed to get absolute target: %w", err)
	}

	inputgocloc, err := stats.GoclocDetectorOutput(targetPath, opts)
	if err != nil {
		log.Debug().Msgf("Error in line of code output %s", err)
		return err
	}
	languageList := FormatFoundLanguages(inputgocloc.Languages)

	// set used language list for external rules to empty if we dont use them
	metaLanguageList := languageList
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

	scanSettings, err := settings.FromOptions(opts, versionMeta)
	scanSettings.Target = opts.Target
	if err != nil {
		return err
	}
	scanSettings.CloudIgnoresUsed, scanSettings.IgnoredFingerprints, scanSettings.StaleIgnoredFingerprintIds, err = getIgnoredFingerprints(
		opts.GeneralOptions.Client,
		scanSettings,
		gitContext,
		opts.PullRequestNumber,
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

	r, err := NewRunner(ctx, scanSettings, gitContext, targetPath, inputgocloc, stats)
	if err != nil {
		return err
	}

	files, baseBranchFindings, err := r.Scan(ctx, opts)
	if err != nil {
		return err
	}

	reportFailed, err := r.Report(files, baseBranchFindings)
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
			defer os.Exit(1)
		} else {
			defer os.Exit(scanSettings.Scan.ExitCode)
		}
	}

	return nil
}

func (r *runner) Report(
	files []files.File,
	baseBranchFindings *basebranchfindings.Findings,
) (bool, error) {
	startTime := time.Now()
	cacheUsed := r.CacheUsed()

	report := types.Report{
		Path:        r.reportPath,
		Inputgocloc: r.goclocResult,
		HasFiles:    r.CacheUsed() || len(files) != 0,
	}

	// if output is defined we want to write only to file
	logger := outputhandler.StdOutLog
	if r.scanSettings.Report.Output != "" {
		reportFile, err := os.Create(r.scanSettings.Report.Output)
		if err != nil {
			return false, fmt.Errorf("error creating output file %w", err)
		}
		logger = outputhandler.PlainLogger(reportFile)
	}

	if cacheUsed && !r.scanSettings.Scan.Quiet {
		// output cached data warning at start of report
		outputhandler.StdErrLog("Using cached data")
	}

	reportData, err := reportoutput.GetData(report, r.scanSettings, r.gitContext, baseBranchFindings)
	if err != nil {
		return false, err
	}
	reportoutput.UploadReportToCloud(reportData, r.scanSettings, r.gitContext)

	endTime := time.Now()

	reportSupported, err := anySupportedLanguagesPresent(report.Inputgocloc, r.scanSettings)
	if err != nil {
		return false, err
	}

	if !reportSupported && r.scanSettings.Report.Report != flag.ReportPrivacy {
		var placeholderStr *strings.Builder
		placeholderStr, err = getPlaceholderOutput(reportData, report, r.scanSettings, report.Inputgocloc)
		if err != nil {
			return false, err
		}

		logger(placeholderStr.String())
		return true, nil
	}

	formatStr, err := reportoutput.FormatOutput(
		reportData,
		r.scanSettings,
		report.Inputgocloc,
		startTime,
		endTime,
	)
	if err != nil {
		return false, fmt.Errorf("error generating report %s", err)
	}

	logger(formatStr)

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

	if len(r.scanSettings.Rules) == 0 && slices.Contains(r.scanSettings.Scan.Scanner, flag.ScannerSAST) && r.scanSettings.Report.Report == flag.ReportSecurity {
		return false, fmt.Errorf("%d rules found for supported language, default rules could not be downloaded or possibly disabled without using --external-rule-dir", len(r.scanSettings.Rules))
	}

	return reportData.ReportFailed, nil
}

func (r *runner) ReportPath() string {
	return r.reportPath
}

func anySupportedLanguagesPresent(inputgocloc *gocloc.Result, config settings.Config) (bool, error) {
	if inputgocloc == nil {
		return true, nil
	}

	ruleLanguages := make(map[string]bool)
	for _, rule := range config.Rules {
		for _, language := range rule.Languages {
			ruleLanguages[language] = true
		}
	}

	foundLanguages := make(map[string]bool)
	for _, language := range inputgocloc.Languages {
		foundLanguages[strings.ToLower(language.Name)] = true
	}

	for _, supportedLanguage := range maps.Keys(settings.GetSupportedRuleLanguages()) {
		_, supportedLangPresent := foundLanguages[supportedLanguage]

		if supportedLangPresent && settings.GetSupportedRuleLanguages()[supportedLanguage] {
			return true, nil
		}
	}

	log.Debug().Msg("No language found for which rules are applicable")
	return false, nil
}

func getPlaceholderOutput(reportData *outputtypes.ReportData, report types.Report, config settings.Config, inputgocloc *gocloc.Result) (outputStr *strings.Builder, err error) {
	if err := reportoutput.GetDataflow(reportData, report, config, true); err != nil {
		return nil, err
	}

	return stats.GetPlaceholderOutput(reportData, inputgocloc, config)
}

func FormatFoundLanguages(languages map[string]*gocloc.Language) (foundLanguages []string) {
	var foundLanguagesMap = make(map[string]bool, len(languages))

	for _, language := range languages {
		if language.Name == "TypeScript" {
			foundLanguagesMap["javascript"] = true
		} else {
			foundLanguagesMap[strings.ToLower(language.Name)] = true
		}
	}

	keys := maps.Keys(foundLanguagesMap)
	sort.Strings(keys)

	return keys
}
