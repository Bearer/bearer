package orchestrator

import (
	"errors"
	"io"
	"os"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/bearer/bearer/pkg/commands/process/orchestrator/filelist"
	"github.com/bearer/bearer/pkg/commands/process/repo_info"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/pool"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/util/jsonlines"
	"github.com/bearer/bearer/pkg/util/output"
	bearerprogress "github.com/bearer/bearer/pkg/util/progressbar"
	"github.com/bearer/bearer/pkg/util/tmpfile"
)

var ErrFileListEmpty = errors.New("We couldn't find any files to scan in the specified directory.")

type orchestrator struct {
	repository          work.Repository
	config              settings.Config
	reportFile          *os.File
	files               []work.File
	maxWorkersSemaphore chan struct{}
	done                chan struct{}
	pool                *pool.Pool
	progressBar         *progressbar.ProgressBar
	waitGroup           sync.WaitGroup
}

func newOrchestrator(
	repository work.Repository,
	config settings.Config,
	reportPath string,
) (*orchestrator, error) {
	reportFile, err := os.Create(reportPath)
	if err != nil {
		return nil, err
	}

	files, err := filelist.Discover(config.Scan.Target, config)
	if err != nil {
		reportFile.Close()
		return nil, err
	}

	if len(files) == 0 {
		reportFile.Close()
		return nil, ErrFileListEmpty
	}

	return &orchestrator{
		repository:          repository,
		config:              config,
		reportFile:          reportFile,
		files:               files,
		maxWorkersSemaphore: make(chan struct{}, config.Scan.Parallel),
		done:                make(chan struct{}),
		pool:                pool.New(config),
		progressBar:         bearerprogress.GetProgressBar(len(files), config, "files"),
	}, nil
}

func (orchestrator *orchestrator) Scan() error {
	if err := repo_info.ReportRepositoryInfo(orchestrator.reportFile, orchestrator.repository, nil); err != nil {
		return err
	}

	for _, file := range orchestrator.files {
		select {
		case <-orchestrator.done:
			log.Debug().Msgf("orchestrator stopping early due to close")
			return nil
		default:
		}

		orchestrator.waitGroup.Add(1)
		go orchestrator.scanFile(file)
	}

	orchestrator.waitGroup.Wait()

	return nil
}

func (orchestrator *orchestrator) scanFile(file work.File) {
	orchestrator.maxWorkersSemaphore <- struct{}{}
	tmpReportFile := tmpfile.Create(".jsonl")

	defer func() {
		if err := orchestrator.progressBar.Add(1); err != nil {
			log.Debug().Msgf("failed to write progress bar for %s", file.FilePath)
		}

		<-orchestrator.maxWorkersSemaphore
		orchestrator.waitGroup.Done()
		os.RemoveAll(tmpReportFile)
	}()

	if err := orchestrator.pool.Scan(work.ProcessRequest{
		Repository: orchestrator.repository,
		File:       file,
		ReportPath: tmpReportFile,
	}); err != nil {
		log.Debug().Msgf("error processing %s: %s", file.FilePath, err)
		orchestrator.writeFileError(file, err)
		return
	}

	orchestrator.writeFileResult(tmpReportFile)
}

func (orchestrator *orchestrator) Close() {
	close(orchestrator.done)
	orchestrator.reportFile.Close()
	orchestrator.progressBar.Close()
	orchestrator.pool.Close()
}

func (orchestrator *orchestrator) writeFileResult(reportPath string) {
	reportFile, err := os.Open(reportPath)
	if err != nil {
		log.Error().Msgf("failed to open tmp report file %s: %s", reportPath, err)
		return
	}
	defer reportFile.Close()

	reportBytes, err := io.ReadAll(reportFile)
	if err != nil {
		log.Error().Msgf("failed to read tmp report file %s: %s", reportPath, err)
		return
	}

	orchestrator.reportFile.Write(reportBytes) //nolint:all,errcheck
}

func (orchestrator *orchestrator) writeFileError(file work.File, fileErr error) {
	fullPath := path.Join(orchestrator.config.Scan.Target, file.FilePath)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Debug().Msgf("failed to stat file %s: %s", fullPath, err)
		return
	}

	detections := []detections.FileFailedDetection{{
		Type:     detections.TypeFileFailed,
		File:     file.FilePath,
		FileSize: int(fileInfo.Size()),
		Timeout:  file.Timeout,
		Error:    fileErr.Error(),
	}}

	if err := jsonlines.Encode(orchestrator.reportFile, &detections); err != nil {
		log.Error().Msgf("failed to encode error for %s: %s", fullPath, err)
	}
}

func Scan(repository work.Repository, config settings.Config, reportPath string) error {
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Scanning target %s", config.Scan.Target)
	}

	orchestrator, err := newOrchestrator(repository, config, reportPath)
	if err != nil {
		return err
	}

	err = orchestrator.Scan()
	orchestrator.Close()
	return err
}
