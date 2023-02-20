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

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"golang.org/x/xerrors"

	"github.com/bearer/curio/cmd/curio/build"
	"github.com/bearer/curio/pkg/commands/process/balancer"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/flag"
	reportoutput "github.com/bearer/curio/pkg/report/output"
	outputhandler "github.com/bearer/curio/pkg/util/output"

	"github.com/bearer/curio/pkg/types"
)

// TargetKind represents what kind of artifact curio scans
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
	balancer       *balancer.Monitor
	reportPath     string
	reuseDetection bool
}

// NewRunner initializes Runner that provides scanning functionalities.
func NewRunner(ctx context.Context, scanSettings settings.Config) Runner {
	r := &runner{}

	r.balancer = balancer.New(scanSettings)

	scanID, err := buildScanID(scanSettings)
	if err != nil {
		log.Error().Msgf("failed to build scan id for caching %e", err)
	}

	path := os.TempDir() + "/curio" + scanID
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
	detectorsHashBuilder := md5.New()
	var ruleNames []string
	for key := range scanSettings.Rules {
		ruleNames = append(ruleNames, key)
	}
	sort.Strings(ruleNames)

	for _, ruleName := range ruleNames {
		_, err := detectorsHashBuilder.Write([]byte(ruleName))
		if err != nil {
			return "", err
		}
		detectorContent, err := json.Marshal(scanSettings.Rules[ruleName])
		if err != nil {
			return "", err
		}
		_, err = detectorsHashBuilder.Write(detectorContent)
		if err != nil {
			return "", err
		}
	}

	configHash := hex.EncodeToString(detectorsHashBuilder.Sum(nil)[:])

	// we want curio sha as it might change detections
	curioSHA := build.CommitSHA

	scanID := strings.TrimSuffix(string(sha), "\n") + "-" + curioSHA + "-" + configHash + ".jsonl"

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
		task := r.balancer.ScheduleTask(work.ProcessRequest{
			Repository: work.Repository{
				Dir:               opts.Target,
				PreviousCommitSHA: "",
				CommitSHA:         "",
			},
			ReportPath: r.reportPath,
		})
		result := <-task.Done

		if result.Error != nil {
			return types.Report{}, result.Error
		}
	}

	return types.Report{
		Path: r.reportPath,
	}, nil
}

// Run performs artifact scanning
func Run(ctx context.Context, opts flag.Options, targetKind TargetKind) (err error) {
	scanSettings, err := settings.FromOptions(opts)
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

	r := NewRunner(ctx, scanSettings)
	defer r.Close(ctx)

	var report types.Report
	switch targetKind {
	case TargetFilesystem:
		if report, err = r.ScanFilesystem(ctx, opts); err != nil {
			if errors.Is(err, balancer.ErrFileListEmpty) {
				outputhandler.StdOutLogger().Msgf("directory empty: %s", err)
				os.Exit(0)
				return
			}

			return xerrors.Errorf("filesystem scan error: %w", err)
		}
	}

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

	if config.Report.Format == flag.FormatEmpty {
		if config.Report.Report == flag.ReportSummary {
			// for policy report, default report format is NOT JSON
			reportPassed, err := reportoutput.ReportSummary(report, logger, config)
			if err != nil {
				return false, fmt.Errorf("error generating report %w", err)
			}
			outputCachedDataWarning(cacheUsed, config.Scan.Quiet)
			return reportPassed, nil
		}
		if config.Report.Report == flag.ReportPrivacy {
			// for privacy report, default report format is CSV
			err := reportoutput.ReportCSV(report, logger, config)
			if err != nil {
				return false, fmt.Errorf("error generating report %w", err)
			}
			outputCachedDataWarning(cacheUsed, config.Scan.Quiet)
			return true, nil
		}
	}

	switch config.Report.Format {
	case flag.FormatEmpty, flag.FormatJSON:
		// default report format for is JSON
		err := reportoutput.ReportJSON(report, logger, config)
		if err != nil {
			return false, fmt.Errorf("error generating report %w", err)
		}
	case flag.FormatYAML:
		err := reportoutput.ReportYAML(report, logger, config)
		if err != nil {
			return false, fmt.Errorf("error generating report %w", err)
		}
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
