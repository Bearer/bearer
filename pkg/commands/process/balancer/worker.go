package balancer

import (
	"context"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"

	"github.com/bearer/curio/pkg/commands/process/repo_info"
	config "github.com/bearer/curio/pkg/commands/process/settings"
	workertype "github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/git"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/util/tmpfile"
)

type Worker struct {
	FileList []workertype.File

	context context.Context
	kill    context.CancelFunc

	taskComplete chan *Worker

	process *Process

	chunkDone      chan *workertype.ProcessResponse
	processErrored chan *workertype.ProcessResponse

	port int

	task *Task

	uuid string

	config config.Config
}

func (worker *Worker) HasNext() bool {
	return len(worker.FileList) > 0
}

func (worker *Worker) NextChunk() []workertype.File {
	end := worker.config.Worker.FilesToBatch
	if end > len(worker.FileList) {
		end = len(worker.FileList)
	}

	toScan := worker.FileList[:end]

	worker.FileList = worker.FileList[end:]

	return toScan
}

func (worker *Worker) Start() {
	err := worker.SpawnProcess(&worker.task.Definition)
	if err != nil {
		worker.complete(err)
		return
	}

	log.Debug().Msgf("worker uuid %s working on repo %s", worker.uuid, worker.task.Definition.Dir)

	addedFiles := 0
	err = filepath.WalkDir(worker.task.Definition.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file := workertype.File{}

		fileinfo, err := d.Info()
		if err != nil {
			return err
		}
		filesize := fileinfo.Size()

		if filesize > int64(worker.config.Worker.FileSizeMaximum) {
			log.Debug().Msgf("skipping file : %s", path)
			return nil
		}

		file.Timeout = worker.config.Worker.TimeoutFileMinimum
		timeoutFileSize := time.Duration(filesize / int64(worker.config.Worker.TimeoutFileSecondPerBytes) * int64(time.Second))
		if timeoutFileSize > file.Timeout {
			if timeoutFileSize > worker.config.Worker.TimeoutFileMaximum {
				file.Timeout = worker.config.Worker.TimeoutFileMaximum
			} else {
				file.Timeout = timeoutFileSize
			}
		}

		file.FilePath = strings.TrimPrefix(path, worker.task.Definition.Dir)

		addedFiles++
		worker.FileList = append(worker.FileList, file)

		return nil
	})

	if err != nil {
		worker.complete(err)
		return
	}

	reportFile, err := os.Create(worker.task.Definition.FilePath)
	if err != nil {
		worker.complete(err)
		return
	}

	commitList, blameRevisionsFilePath, err := worker.getCommitListAndWriteForBlame()
	if err != nil {
		worker.complete(err)
		return
	}
	defer os.Remove(blameRevisionsFilePath)

	if err := repo_info.ReportRepositoryInfo(
		reportFile,
		worker.task.Definition.Repository,
		commitList,
	); err != nil {
		reportFile.Close()
		worker.complete(err)
		return
	}

	i := 0
	for {
		i++

		tmpReportFile := tmpfile.Create(os.TempDir(), ".jsonl")

		work := worker.NextChunk()

		worker.DoWork(&Task{
			Definition: workertype.ProcessRequest{
				Repository:             worker.task.Definition.Repository,
				Files:                  work,
				FilePath:               tmpReportFile,
				BlameRevisionsFilePath: blameRevisionsFilePath,
			},
			Done: worker.chunkDone,
		})

		var shouldBreak = false
		select {
		case <-worker.context.Done():
			if worker.process != nil {
				worker.process.kill()
				worker.process = nil
			}
			worker.task.Done <- &workertype.ProcessResponse{
				Error: ErrorClosing,
			}
			shouldBreak = true
		case response := <-worker.processErrored:
			// add failed files to report
			log.Debug().Msgf("worker %s got process error %e", worker.uuid, response.Error)
			if worker.process != nil {
				log.Debug().Msgf("process is not nil killing it")
				worker.process.kill()
				worker.process = nil
			}

			worker.logError(reportFile, work, response)

			err := worker.SpawnProcess(&worker.task.Definition)
			if err != nil {
				worker.process.kill()
				worker.process = nil
				worker.complete(err)
				return
			}
		case response := <-worker.chunkDone:
			if response.Error != nil {
				worker.logError(reportFile, work, response)
			}

			// ungzip report and add it to master file
			f, err := os.Open(tmpReportFile)
			if err != nil {
				log.Error().Msgf("worker %s failed to open tmp report chunk file %e", worker.uuid, err)
				worker.complete(err)

				break
			}

			reportBytes, err := ioutil.ReadAll(f)
			if err != nil {
				log.Error().Msgf("worker %s failed to read tmp report chunk file %e", worker.uuid, err)
				worker.complete(err)
				f.Close()
				break
			}

			reportFile.Write(reportBytes) //nolint:all,errcheck
			f.Close()
		}

		os.RemoveAll(tmpReportFile)

		if shouldBreak {
			err := reportFile.Close()
			if err != nil {
				log.Debug().Msgf("worker %s failed to close gzipwriter", worker.uuid)
			}
			err = reportFile.Close()
			if err != nil {
				log.Debug().Msgf("worker %s failed to close reportfile", worker.uuid)
			}
			worker.task.Done <- &workertype.ProcessResponse{
				Error: errors.New("context canceled"),
			}
			break
		}

		if !worker.HasNext() {
			if worker.process != nil {
				log.Debug().Msgf("process is not nil killing it")
				worker.process.kill()
				worker.process = nil
			}

			log.Printf("worker %s closing due to work done", worker.uuid)
			err := reportFile.Close()
			if err != nil {
				log.Debug().Msgf("worker %s failed to close gzipwriter", worker.uuid)
			}
			err = reportFile.Close()
			if err != nil {
				log.Debug().Msgf("worker %s failed to close reportfile", worker.uuid)
			}
			worker.complete(nil)
			break
		}
	}

}

func (worker *Worker) DoWork(task *Task) {
	go worker.process.doTask(task)
}

func (worker *Worker) SpawnProcess(task *workertype.ProcessRequest) error {
	cntx, cntxCancel := context.WithCancel(context.Background())
	worker.process = &Process{
		context:        cntx,
		kill:           cntxCancel,
		chunkDone:      worker.chunkDone,
		processErrored: worker.processErrored,
		port:           worker.port,
		task:           worker.task,

		config: worker.config,

		uuid:       uuid.NewString(),
		workeruuid: worker.uuid,
	}

	return worker.process.StartProcess(task)
}

func (worker *Worker) complete(err error) {
	worker.task.Done <- &workertype.ProcessResponse{Error: err}
	worker.taskComplete <- worker
}

func (worker *Worker) getCommitListAndWriteForBlame() (commitList []git.CommitInfo, blameFilePath string, err error) {
	if worker.task.Definition.PreviousCommitSHA == "" {
		return
	}

	if commitList, err = git.GetCommitList(
		worker.task.Definition.Dir,
		worker.task.Definition.PreviousCommitSHA,
		worker.task.Definition.CommitSHA,
	); err != nil {
		return
	}

	blameRevisionsFile, err := os.CreateTemp("", "blame-revs")
	if err != nil {
		return
	}

	err = git.WriteCommitsForBlame(blameRevisionsFile, commitList)
	blameRevisionsFile.Close()
	if err != nil {
		os.Remove(blameRevisionsFile.Name())
		return
	}

	blameFilePath = blameRevisionsFile.Name()
	return
}

func (worker *Worker) logError(reportFile *os.File, work []workertype.File, response *workertype.ProcessResponse) {
	var errorsToAdd []report.FileFailedDetection
	for _, file := range work {
		fileInfo, err := os.Stat(worker.task.Definition.Dir + "/" + file.FilePath)
		if err != nil {
			log.Debug().Msgf("worker %s failed to stat file %e", worker.uuid, err)
			continue
		}

		errorsToAdd = append(errorsToAdd, report.FileFailedDetection{
			Type:     report.TypeFileFailed,
			File:     file.FilePath,
			FileSize: int(fileInfo.Size()),
			Timeout:  file.Timeout,
			Error:    response.Error.Error(),
		})
	}

	err := jsonlines.Encode(reportFile, &errorsToAdd)
	if err != nil {
		log.Error().Msgf("worker %s failed to encode data line %e", worker.uuid, err)
	}

}
