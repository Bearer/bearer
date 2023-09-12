package artifact

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"golang.org/x/exp/maps"

	evalstats "github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/pkg/commands/artifact/scanid"
	"github.com/bearer/bearer/pkg/commands/process/filelist"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/github_api"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/report/output/stats"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	outputhandler "github.com/bearer/bearer/pkg/util/output"

	"github.com/bearer/bearer/pkg/types"
)

var ErrFileListEmpty = errors.New("couldn't find any files to scan in the specified directory")

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
	Scan(ctx context.Context, opts flag.Options) ([]files.File, *basebranchfindings.Findings, error)
	// Report a writes a report
	Report(files []files.File, baseBranchFindings *basebranchfindings.Findings) (bool, error)
	// Close closes runner
	Close(ctx context.Context) error
}

type runner struct {
	reportPath     string
	reuseDetection bool
	goclocResult   *gocloc.Result
	scanSettings   settings.Config
	stats          *evalstats.Stats
}

// NewRunner initializes Runner that provides scanning functionalities.
func NewRunner(
	ctx context.Context,
	scanSettings settings.Config,
	goclocResult *gocloc.Result,
	stats *evalstats.Stats,
) Runner {
	r := &runner{
		scanSettings: scanSettings,
		goclocResult: goclocResult,
		stats:        stats,
	}

	scanID, err := scanid.Build(scanSettings)
	if err != nil {
		log.Error().Msgf("failed to build scan id for caching %s", err)
	}

	path := os.TempDir() + "/bearer" + scanID
	completedPath := strings.Replace(path, ".jsonl", "-completed.jsonl", 1)

	r.reportPath = path

	log.Debug().Msgf("creating report %s", path)

	if _, err := os.Stat(completedPath); err == nil {
		if !scanSettings.Scan.Force {
			r.reuseDetection = true
			log.Debug().Msgf("reuse detection for %s", path)
			r.reportPath = completedPath

			return r
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

	return r
}

func (r *runner) CacheUsed() bool {
	return r.reuseDetection
}

// Close closes everything
func (r *runner) Close(ctx context.Context) error {
	return nil
}

func (r *runner) Scan(ctx context.Context, opts flag.Options) ([]files.File, *basebranchfindings.Findings, error) {
	if r.reuseDetection {
		return nil, nil, nil
	}

	if !opts.Quiet {
		outputhandler.StdErrLog(fmt.Sprintf("Scanning target %s", opts.Target))
	}

	targetPath, err := filepath.Abs(opts.Target)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute target: %w", err)
	}

	repository, err := gitrepository.New(ctx, r.scanSettings, targetPath, opts.DiffBaseBranch)
	if err != nil {
		return nil, nil, fmt.Errorf("git repository error: %w", err)
	}

	fileList, err := filelist.Discover(repository, targetPath, r.goclocResult, r.scanSettings)
	if err != nil {
		return nil, nil, err
	}

	if len(fileList.Files) == 0 {
		return nil, nil, ErrFileListEmpty
	}

	orchestrator, err := orchestrator.New(
		work.Repository{Dir: opts.Target},
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
			outputhandler.StdErrLog(fmt.Sprintf("\nScanning base branch %s", opts.DiffBaseBranch))
		}

		if err := orchestrator.Scan(r.reportPath+".base", fileList.BaseFiles); err != nil {
			return err
		}

		report := types.Report{Path: r.reportPath + ".base", Inputgocloc: r.goclocResult}

		reportData, err := reportoutput.GetData(report, r.scanSettings, nil)
		if err != nil {
			return err
		}

		baseBranchFindings = buildBaseBranchFindings(reportData, fileList)

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

// Run performs artifact scanning
func Run(ctx context.Context, opts flag.Options) (err error) {
	if !opts.Quiet {
		outputhandler.StdErrLog("Loading rules")
	}

	github_api.VersionCheck(ctx, opts.GeneralOptions.DisableVersionCheck, opts.ScanOptions.Quiet)

	inputgocloc, err := stats.GoclocDetectorOutput(opts.ScanOptions.Target)
	if err != nil {
		log.Debug().Msgf("Error in line of code output %s", err)
		return err
	}
	scanSettings, err := settings.FromOptions(opts, FormatFoundLanguages(inputgocloc.Languages))
	scanSettings.Target = opts.Target
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

	var stats *evalstats.Stats
	if scanSettings.Debug {
		stats = evalstats.New()
	}

	gitrepository.ConfigureGithubAuth(scanSettings.Scan.GithubToken)

	r := NewRunner(ctx, scanSettings, inputgocloc, stats)
	defer r.Close(ctx)

	if !r.CacheUsed() && scanSettings.CacheUsed {
		// re-cache rules
		if opts.ScanOptions.Force && !opts.ScanOptions.Quiet {
			outputhandler.StdOutLog("Caching rules")
		}
		if err = settings.RefreshRules(scanSettings, opts.ExternalRuleDir, opts.RuleOptions, FormatFoundLanguages(inputgocloc.Languages)); err != nil {
			return err
		}
	}

	files, baseBranchFindings, err := r.Scan(ctx, opts)
	if err != nil {
		if errors.Is(err, ErrFileListEmpty) {
			outputhandler.StdOutLog(err.Error())
			os.Exit(0)
			return
		}

		return fmt.Errorf("scan error: %w", err)
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

	report := types.Report{Path: r.reportPath, Inputgocloc: r.goclocResult}

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

	reportData, err := reportoutput.GetData(report, r.scanSettings, baseBranchFindings)
	if err != nil {
		return false, err
	}

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

	logger(*formatStr)

	outputCachedDataWarning(cacheUsed, r.scanSettings.Scan.Quiet)

	return reportData.ReportFailed, nil
}

func (r *runner) ReportPath() string {
	return r.reportPath
}

func outputCachedDataWarning(cacheUsed bool, quietMode bool) {
	if quietMode || !cacheUsed {
		return
	}

	outputhandler.StdErrLog("Cached data used (no code changes detected). Unexpected? Use --force to force a re-scan.\n")
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

		if supportedLangPresent {
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

func buildBaseBranchFindings(reportData *outputtypes.ReportData, fileList *files.List) *basebranchfindings.Findings {
	result := basebranchfindings.New(fileList)

	for _, findings := range reportData.FindingsBySeverity {
		for _, finding := range findings {
			result.Add(
				finding.Rule.Id,
				finding.Filename,
				finding.Sink.Start,
				finding.Sink.End,
			)
		}
	}

	return result
}
