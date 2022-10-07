package artifact

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"golang.org/x/xerrors"

	"github.com/bearer/curio/pkg/commands/process/balancer"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/tmpfile"

	"github.com/bearer/curio/pkg/scanner"
	"github.com/bearer/curio/pkg/types"
)

// TargetKind represents what kind of artifact curio scans
type TargetKind string

const (
	TargetFilesystem TargetKind = "fs"
	TargetRepository TargetKind = "repo"
)

// InitializeScanner defines the initialize function signature of scanner
type InitializeScanner func(context.Context, ScannerConfig) (scanner.Scanner, func(), error)

type ScannerConfig struct {
	Target   string
	Artifact types.Artifact
}

type Runner interface {
	// ScanFilesystem scans a filesystem
	ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanRepository scans repository
	ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error)
	// Filter filter a report
	Filter(ctx context.Context, opts flag.Options, report types.Report) (types.Report, error)
	// Report a writes a report
	Report(opts flag.Options, report types.Report) error
	// Close closes runner
	Close(ctx context.Context) error
}

type runner struct {
}

type runnerOption func(*runner)

// NewRunner initializes Runner that provides scanning functionalities.
// It is possible to return SkipScan and it must be handled by caller.
func NewRunner(ctx context.Context, cliOptions flag.Options, opts ...runnerOption) (Runner, error) {
	r := &runner{}
	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Close closes everything
func (r *runner) Close(ctx context.Context) error {
	return nil
}

func (r *runner) ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error) {
	return r.scanFS(ctx, opts)
}

func (r *runner) scanFS(ctx context.Context, opts flag.Options) (types.Report, error) {
	var s InitializeScanner
	// Scan filesystem

	return r.scanArtifact(ctx, opts, s)
}

func (r *runner) ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error) {

	return r.scanArtifact(ctx, opts, repositoryStandaloneScanner)
}

func (r *runner) scanArtifact(ctx context.Context, opts flag.Options, scanner InitializeScanner) (types.Report, error) {
	reportpath := tmpfile.Create(os.TempDir(), ".jsonl")
	balancer := balancer.New(settings.WorkerSettings{
		FilesToBatch:          1,
		Count:                 1,
		Memory:                800 * 1024 * 1024,
		ProcessOnlineTimeout:  60 * time.Second,
		TimeoutSecondPerBytes: 10 * 1000,
		TimeoutMinimum:        5 * time.Second,
		TimeoutMaximum:        300 * time.Second,
		MaximumFileSize:       25 * 1000 * 1000,
	})
	task := balancer.ScheduleTask(work.ProcessRequest{
		Repository: work.Repository{
			Dir:               opts.Target,
			PreviousCommitSHA: "",
			CommitSHA:         "",
		},
		FilePath:             reportpath,
		CustomDetectorConfig: nil,
	})
	result := <-task.Done

	if result.Error != nil {
		return types.Report{}, result.Error
	}

	log.Debug().Msgf("report is %s", reportpath)

	return types.Report{
		Path: reportpath,
	}, nil
}

func (r *runner) Filter(ctx context.Context, opts flag.Options, report types.Report) (types.Report, error) {

	return report, nil
}

func (r *runner) Report(opts flag.Options, report types.Report) error {

	return nil
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

	if opts.GenerateDefaultConfig {
		log.Info().Msg("Writing the default config to curio-default.yaml...")
		return viper.SafeWriteConfigAs("curio-default.yaml")
	}

	r, err := NewRunner(ctx, opts)
	if err != nil {
		return xerrors.Errorf("init error: %w", err)
	}
	defer r.Close(ctx)

	var report types.Report
	switch targetKind {
	case TargetFilesystem:
		if report, err = r.ScanFilesystem(ctx, opts); err != nil {
			return xerrors.Errorf("filesystem scan error: %w", err)
		}
	}

	report, err = r.Filter(ctx, opts, report)
	if err != nil {
		return xerrors.Errorf("filter error: %w", err)
	}

	if err = r.Report(opts, report); err != nil {
		return xerrors.Errorf("report error: %w", err)
	}

	Exit(opts, report.Failed())

	return nil
}

func initScannerConfig(opts flag.Options) (ScannerConfig, types.ScanOptions, error) { //nolint:all,unused
	target := opts.Target

	scanOptions := types.ScanOptions{
		// SecurityChecks: opts.SecurityChecks,
		// FilePatterns:   opts.FilePatterns,
	}

	return ScannerConfig{
		Target: target,
	}, scanOptions, nil
}

func scan(ctx context.Context, opts flag.Options, initializeScanner InitializeScanner) ( //nolint:all,unused
	types.Report, error) {

	scannerConfig, scanOptions, err := initScannerConfig(opts)
	if err != nil {
		return types.Report{}, err
	}

	s, cleanup, err := initializeScanner(ctx, scannerConfig)
	if err != nil {
		return types.Report{}, xerrors.Errorf("unable to initialize a scanner: %w", err)
	}
	defer cleanup()

	report, err := s.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return types.Report{}, xerrors.Errorf("scan failed: %w", err)
	}
	return report, nil
}

func Exit(opts flag.Options, failedResults bool) {
	if opts.ExitCode != 0 && failedResults {
		os.Exit(opts.ExitCode)
	}
}
