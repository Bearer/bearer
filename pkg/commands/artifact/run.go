package artifact

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"golang.org/x/exp/maps"
	"golang.org/x/xerrors"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/github_api"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/report/output/gitlab"
	rdo "github.com/bearer/bearer/pkg/report/output/reviewdog"
	"github.com/bearer/bearer/pkg/report/output/sarif"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/report/output/stats"
	outputhandler "github.com/bearer/bearer/pkg/util/output"

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
	// ScanFilesystem scans a filesystem
	ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanRepository scans repository
	ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error)
	// Report a writes a report
	Report(scanSettings settings.Config, report types.Report) (bool, error)
	// Close closes runner
	Close(ctx context.Context) error
}

type runner struct {
	reportPath     string
	reuseDetection bool
	goclocResult   *gocloc.Result
	scanSettings   settings.Config
}

// NewRunner initializes Runner that provides scanning functionalities.
func NewRunner(ctx context.Context, scanSettings settings.Config, goclocResult *gocloc.Result) Runner {
	r := &runner{scanSettings: scanSettings, goclocResult: goclocResult}

	scanID, err := buildScanID(scanSettings)
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

// buildScanHash builds a hash based on project and settings that latter on gets used for caching scan detections
func buildScanID(scanSettings settings.Config) (string, error) {
	// we want head as project may contain new changes
	cmd := exec.Command("git", "-C", scanSettings.Scan.Target, "rev-parse", "HEAD")
	sha, err := cmd.Output()

	if err != nil {
		log.Debug().Msgf("error getting git sha %s", err.Error())
		sha = []byte(uuid.NewString())
	}

	// we want hash of all active custom detector rules and their content
	hashBuilder := md5.New()
	var ruleNames []string
	for key := range scanSettings.Rules {
		ruleNames = append(ruleNames, key)
	}
	sort.Strings(ruleNames)

	for _, ruleName := range ruleNames {
		_, err := hashBuilder.Write([]byte(ruleName))
		if err != nil {
			return "", err
		}
		detectorContent, err := json.Marshal(scanSettings.Rules[ruleName])
		if err != nil {
			return "", err
		}
		_, err = hashBuilder.Write(detectorContent)
		if err != nil {
			return "", err
		}
	}

	var scanners []string
	scanners = append(scanners, scanSettings.Scan.Scanner...)
	sort.Strings(scanners)

	for _, scanner := range scanners {
		_, err := hashBuilder.Write([]byte(scanner))
		if err != nil {
			return "", err
		}
	}

	configHash := hex.EncodeToString(hashBuilder.Sum(nil)[:])

	// we want sha as it might change detections
	buildSHA := build.CommitSHA
	scanID := strings.TrimSuffix(string(sha), "\n") + "-" + buildSHA + "-" + configHash + ".jsonl"

	return scanID, nil
}

func (r *runner) CacheUsed() bool {
	return r.reuseDetection
}

// Close closes everything
func (r *runner) Close(ctx context.Context) error {
	return nil
}

func (r *runner) ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error) {
	return r.scanFS(ctx, opts)
}

func (r *runner) scanFS(ctx context.Context, opts flag.Options) (types.Report, error) {
	return r.scanArtifact(ctx, opts)
}

func (r *runner) ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error) {
	return r.scanArtifact(ctx, opts)
}

func (r *runner) scanArtifact(ctx context.Context, opts flag.Options) (types.Report, error) {
	if !r.reuseDetection {
		if err := orchestrator.Scan(
			work.Repository{
				Dir:               opts.Target,
				PreviousCommitSHA: "",
				CommitSHA:         "",
			},
			r.scanSettings,
			r.goclocResult,
			r.reportPath,
		); err != nil {
			return types.Report{}, err
		}
	}

	return types.Report{
		Path: r.reportPath,
	}, nil
}

// Run performs artifact scanning
func Run(ctx context.Context, opts flag.Options, targetKind TargetKind) (err error) {
	if !opts.Quiet {
		outputhandler.StdErrLogger().Msg("Loading rules")
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

	r := NewRunner(ctx, scanSettings, inputgocloc)
	defer r.Close(ctx)

	if !r.CacheUsed() && scanSettings.CacheUsed {
		// re-cache rules
		if opts.ScanOptions.Force && !opts.ScanOptions.Quiet {
			outputhandler.StdOutLogger().Msgf("Caching rules")
		}
		if err = settings.RefreshRules(scanSettings, opts.ExternalRuleDir, opts.RuleOptions, FormatFoundLanguages(inputgocloc.Languages)); err != nil {
			return err
		}
	}

	var report types.Report
	switch targetKind {
	case TargetFilesystem:
		if report, err = r.ScanFilesystem(ctx, opts); err != nil {
			if errors.Is(err, orchestrator.ErrFileListEmpty) {
				outputhandler.StdOutLogger().Msgf(err.Error())
				os.Exit(0)
				return
			}

			return xerrors.Errorf("filesystem scan error: %w", err)
		}
	}
	report.Inputgocloc = inputgocloc
	reportPassed, err := r.Report(scanSettings, report)
	if err != nil {
		return xerrors.Errorf("report error: %w", err)
	} else {
		if !strings.HasSuffix(report.Path, "-completed.jsonl") {
			newPath := strings.Replace(report.Path, ".jsonl", "-completed.jsonl", 1)
			log.Debug().Msgf("renaming report %s -> %s", report.Path, newPath)
			err := os.Rename(report.Path, newPath)
			if err != nil {
				return xerrors.Errorf("failed to rename report file %s -> %s: %w", report.Path, newPath, err)
			}
			report.Path = newPath
		}
	}

	if !reportPassed {
		defer os.Exit(1)
	}

	return nil
}

func (r *runner) Report(config settings.Config, report types.Report) (bool, error) {
	startTime := time.Now()
	cacheUsed := r.CacheUsed()
	// if output is defined we want to write only to file
	logger := outputhandler.StdOutLogger()
	if config.Report.Output != "" {
		reportFile, err := os.Create(config.Report.Output)
		if err != nil {
			return false, fmt.Errorf("error creating output file %w", err)
		}
		logger = outputhandler.PlainLogger(reportFile)
	}

	if cacheUsed && !config.Scan.Quiet {
		// output cached data warning at start of report
		outputhandler.StdErrLogger().Msg("Using cached data")
	}

	detections, dataflow, err := reportoutput.GetOutput(report, config)
	if err != nil {
		return false, err
	}

	endTime := time.Now()

	reportSupported, err := anySupportedLanguagesPresent(report.Inputgocloc, config)
	if err != nil {
		return false, err
	}

	if !reportSupported && config.Report.Report != flag.ReportPrivacy {
		var placeholderStr *strings.Builder
		placeholderStr, err = getPlaceholderOutput(report, config, report.Inputgocloc)
		if err != nil {
			return false, err
		}

		logger.Msg(placeholderStr.String())
		return true, nil
	}

	if config.Report.Format == flag.FormatEmpty {
		if config.Report.Report == flag.ReportSecurity {
			// for security report, default report format is Table
			detectionReport := detections.(*security.Results)
			reportStr, reportPassed := security.BuildReportString(config, detectionReport, report.Inputgocloc, dataflow)

			logger.Msg(reportStr.String())

			return reportPassed, nil
		} else if config.Report.Report == flag.ReportPrivacy {
			// for privacy report, default report format is CSV
			content, err := reportoutput.GetPrivacyReportCSVOutput(report, dataflow, config)
			if err != nil {
				return false, fmt.Errorf("error generating report %s", err)
			}

			logger.Msg(*content)

			return true, nil
		}
	}

	switch config.Report.Format {
	case flag.FormatSarif:
		sarifContent, err := sarif.ReportSarif(detections.(*map[string][]security.Result), config.Rules)
		if err != nil {
			return false, fmt.Errorf("error generating sarif report %s", err)
		}
		content, err := outputhandler.ReportJSON(sarifContent)
		if err != nil {
			return false, fmt.Errorf("error generating JSON report %s", err)
		}

		logger.Msg(*content)
	case flag.FormatReviewDog:
		sastContent, err := rdo.ReportReviewdog(detections.(*map[string][]security.Result))
		if err != nil {
			return false, fmt.Errorf("error generating reviewdog report %s", err)
		}
		content, err := reportoutput.ReportJSON(sastContent)
		if err != nil {
			return false, fmt.Errorf("error generating JSON report %s", err)
		}

		logger.Msg(*content)
	case flag.FormatGitLabSast:

		sastContent, err := gitlab.ReportGitLab(detections.(*map[string][]security.Result), startTime, endTime)
		if err != nil {
			return false, fmt.Errorf("error generating gitlab-sast report %s", err)
		}
		content, err := outputhandler.ReportJSON(sastContent)
		if err != nil {
			return false, fmt.Errorf("error generating JSON report %s", err)
		}

		logger.Msg(*content)
	case flag.FormatEmpty, flag.FormatJSON:
		content, err := outputhandler.ReportJSON(detections)
		if err != nil {
			return false, fmt.Errorf("error generating report %s", err)
		}

		logger.Msg(*content)
	case flag.FormatYAML:
		content, err := outputhandler.ReportYAML(detections)
		if err != nil {
			return false, fmt.Errorf("error generating report %s", err)
		}

		logger.Msg(*content)
	}

	outputCachedDataWarning(cacheUsed, config.Scan.Quiet)
	return true, nil
}

func outputCachedDataWarning(cacheUsed bool, quietMode bool) {
	if quietMode || !cacheUsed {
		return
	}

	outputhandler.StdErrLogger().Msg("Cached data used (no code changes detected). Unexpected? Use --force to force a re-scan.\n")
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

	_, javaPresent := foundLanguages["java"]
	if javaPresent {
		return true, nil
	}

	log.Debug().Msg("No language found for which rules are applicable")
	return false, nil
}

func getPlaceholderOutput(report types.Report, config settings.Config, inputgocloc *gocloc.Result) (outputStr *strings.Builder, err error) {
	dataflowOutput, _, err := reportoutput.GetDataflow(report, config, true)
	if err != nil {
		return
	}

	return stats.GetPlaceholderOutput(inputgocloc, dataflowOutput, config)
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
