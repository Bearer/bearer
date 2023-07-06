package orchestrator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
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
	reportMutex         sync.Mutex
}

func newOrchestrator(
	repository work.Repository,
	config settings.Config,
	goclocResult *gocloc.Result,
	reportPath string,
	stats *stats.Stats,
) (*orchestrator, error) {
	reportFile, err := os.Create(reportPath)
	if err != nil {
		return nil, err
	}

	files, err := filelist.Discover(config.Scan.Target, goclocResult, config)
	if err != nil {
		reportFile.Close()
		return nil, err
	}

	if len(files) == 0 {
		reportFile.Close()
		return nil, ErrFileListEmpty
	}

	parallel := getParallel(len(files), config)
	log.Debug().Msgf("number of workers: %d", parallel)

	return &orchestrator{
		repository:          repository,
		config:              config,
		reportFile:          reportFile,
		files:               files,
		maxWorkersSemaphore: make(chan struct{}, parallel),
		done:                make(chan struct{}),
		pool:                pool.New(config, stats),
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

	orchestrator.reportMutex.Lock()
	_, err = orchestrator.reportFile.Write(reportBytes)
	if err != nil {
		log.Error().Msgf("failed to write tmp report into main report file %s: %s", reportPath, err)
	}
	orchestrator.reportMutex.Unlock()
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

	orchestrator.reportMutex.Lock()
	if err := jsonlines.Encode(orchestrator.reportFile, &detections); err != nil {
		log.Error().Msgf("failed to encode error for %s: %s", fullPath, err)
	}
	orchestrator.reportMutex.Unlock()
}

func Scan(
	repository work.Repository,
	config settings.Config,
	goclogResult *gocloc.Result,
	reportPath string,
	stats *stats.Stats,
) error {
	if !config.Scan.Quiet {
		output.StdErrLog(fmt.Sprintf("Scanning target %s", config.Scan.Target))
	}

	orchestrator, err := newOrchestrator(repository, config, goclogResult, reportPath, stats)
	if err != nil {
		return err
	}

	err = orchestrator.Scan()
	orchestrator.Close()
	return err
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
