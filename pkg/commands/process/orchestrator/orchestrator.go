package orchestrator

import (
	"io"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/pool"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/util/jsonlines"
	bearerprogress "github.com/bearer/bearer/pkg/util/progressbar"
	"github.com/bearer/bearer/pkg/util/tmpfile"
)

type Orchestrator struct {
	repository          work.Repository
	config              settings.Config
	files               []work.File
	maxWorkersSemaphore chan struct{}
	done                chan struct{}
	pool                *pool.Pool
	reportMutex         sync.Mutex
}

func New(
	repository work.Repository,
	config settings.Config,
	files []work.File,
	stats *stats.Stats,
) (*Orchestrator, error) {
	parallel := getParallel(len(files), config)
	log.Debug().Msgf("number of workers: %d", parallel)

	return &Orchestrator{
		repository:          repository,
		config:              config,
		files:               files,
		maxWorkersSemaphore: make(chan struct{}, parallel),
		done:                make(chan struct{}),
		pool:                pool.New(config, stats),
	}, nil
}

func (orchestrator *Orchestrator) Scan(reportPath string) error {
	fileComplete := make(chan struct{}, len(orchestrator.files))
	progressBar := bearerprogress.GetProgressBar(len(orchestrator.files), orchestrator.config, "files")

	reportFile, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer reportFile.Close()

	for _, file := range orchestrator.files {
		select {
		case <-orchestrator.done:
			log.Debug().Msgf("scan stopping early due to close")
			return nil
		default:
		}

		go orchestrator.scanFile(reportFile, fileComplete, file)
	}

	orchestrator.waitForScan(fileComplete, progressBar)

	if err := progressBar.Close(); err != nil {
		log.Debug().Msgf("failed to close progress bar: %s", err)
	}

	return nil
}

func (orchestrator *Orchestrator) waitForScan(fileComplete chan struct{}, progressBar *progressbar.ProgressBar) {
	count := 0

	for {
		select {
		case <-orchestrator.done:
			log.Debug().Msgf("scan stopping early due to close")

			return
		case <-fileComplete:
			count++

			if err := progressBar.Add(1); err != nil {
				log.Debug().Msgf("failed to write progress bar: %s", err)
			}

			if count == len(orchestrator.files) {
				return
			}
		}
	}
}

func (orchestrator *Orchestrator) scanFile(reportFile *os.File, fileComplete chan struct{}, file work.File) {
	orchestrator.maxWorkersSemaphore <- struct{}{}
	tmpReportPath := tmpfile.Create(".jsonl")

	defer func() {
		<-orchestrator.maxWorkersSemaphore
		os.RemoveAll(tmpReportPath)
		fileComplete <- struct{}{}
	}()

	if err := orchestrator.pool.Scan(work.ProcessRequest{
		Repository: orchestrator.repository,
		File:       file,
		ReportPath: tmpReportPath,
	}); err != nil {
		log.Debug().Msgf("error processing %s: %s", file.FilePath, err)
		orchestrator.writeFileError(reportFile, file, err)
		return
	}

	orchestrator.writeFileResult(reportFile, tmpReportPath)
}

func (orchestrator *Orchestrator) Close() {
	close(orchestrator.done)
	orchestrator.pool.Close()
}

func (orchestrator *Orchestrator) writeFileResult(reportFile *os.File, tmpReportPath string) {
	tmpReportFile, err := os.Open(tmpReportPath)
	if err != nil {
		log.Error().Msgf("failed to open tmp report file %s: %s", tmpReportPath, err)
		return
	}
	defer tmpReportFile.Close()

	reportBytes, err := io.ReadAll(tmpReportFile)
	if err != nil {
		log.Error().Msgf("failed to read tmp report file %s: %s", tmpReportPath, err)
		return
	}

	orchestrator.reportMutex.Lock()
	_, err = reportFile.Write(reportBytes)
	if err != nil {
		log.Error().Msgf("failed to write tmp report into main report file %s: %s", tmpReportPath, err)
	}
	orchestrator.reportMutex.Unlock()
}

func (orchestrator *Orchestrator) writeFileError(reportFile *os.File, file work.File, fileErr error) {
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

	orchestrator.reportMutex.Lock()
	if err := jsonlines.Encode(reportFile, &detections); err != nil {
		log.Error().Msgf("failed to encode error for %s: %s", fullPath, err)
	}
	orchestrator.reportMutex.Unlock()
}

func getParallel(fileCount int, config settings.Config) int {
	if config.Scan.Parallel != 0 {
		return config.Scan.Parallel
	}

	result := fileCount / settings.FilesPerWorker

	if result == 0 {
		return 2
	}

	if result > runtime.NumCPU() {
		return runtime.NumCPU()
	}

	return result
}
