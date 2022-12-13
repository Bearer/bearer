package artifact

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	cmd := exec.Command("git", "-C", scanSettings.Scan.Target, "rev-parse", "HEAD")
	sha, err := cmd.Output()

	if err != nil {
		log.Debug().Msgf("error getting git sha %s", err.Error())
		sha = []byte(uuid.NewString())
	}

	path := os.TempDir() + strings.TrimSuffix(string(sha), "\n") + "-" + build.CommitSHA + ".jsonl"
	r.reportPath = path

	if _, err := os.Stat(path); err == nil {
		if !scanSettings.Scan.Force {
			r.reuseDetection = true
			log.Debug().Msgf("reuse detection for %s", path)

			return r
		} else {
			err := os.Remove(path)
			if err != nil {
				log.Error().Msgf("couldn't remove report path %s, %s", path, err.Error())
			}
		}
	}

	log.Debug().Msgf("creating report %s", path)
	pathCreated, err := os.Create(path)

	if err != nil {
		log.Error().Msgf("failed to create path %s, %s, %#v", path, err.Error(), pathCreated)
	}

	log.Debug().Msgf("successfully created reportPath %s", path)

	return r
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
			FilePath: r.reportPath,
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
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	defer func() {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn().Msg("Increase --timeout value")
		}
	}()

	scanSettings, err := settings.FromOptions(opts)
	scanSettings.Target = opts.Target

	if err != nil {
		return err
	}

	r := NewRunner(ctx, scanSettings)
	defer r.Close(ctx)

	var report types.Report
	switch targetKind {
	case TargetFilesystem:
		if report, err = r.ScanFilesystem(ctx, opts); err != nil {
			return xerrors.Errorf("filesystem scan error: %w", err)
		}
	}

	reportPassed, err := r.Report(scanSettings, report)
	if err != nil {
		return xerrors.Errorf("report error: %w", err)
	}

	if !reportPassed {
		defer os.Exit(1)
	}

	return nil
}

func (r *runner) Report(config settings.Config, report types.Report) (bool, error) {
	// if output is defined we want to write only to file
	logger := outputhandler.StdOutLogger()
	if config.Report.Output != "" {
		reportFile, err := os.Create(config.Report.Output)
		if err != nil {
			return false, fmt.Errorf("error creating output file %w", err)
		}
		logger = outputhandler.PlainLogger(reportFile)
	}

	if config.Report.Report == flag.ReportPolicies && config.Report.Format == "" {
		// for policy report, default report format is NOT JSON
		reportPassed, err := reportoutput.ReportPolicies(report, logger, config)
		if err != nil {
			return false, fmt.Errorf("error generating report %w", err)
		}
		return reportPassed, nil
	}

	switch config.Report.Format {
	case "", flag.FormatJSON:
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
	return true, nil
}
