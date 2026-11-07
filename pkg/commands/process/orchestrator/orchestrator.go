package orchestrator

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/bearer/bearer/pkg/report/coverage"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"io"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/scanner/stats"
	bearerprogress "github.com/bearer/bearer/pkg/util/progressbar"

	"github.com/bearer/bearer/pkg/commands/process/orchestrator/pool"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
)

type Orchestrator struct {
	repository          work.Repository
	config              *settings.Config
	maxWorkersSemaphore chan struct{}
	done                chan struct{}
	pool                interface {
		Scan(work.ProcessRequest) error
		Close()
	}
	closeOnce   sync.Once
	reportMutex sync.Mutex
	inputScope  *inputfile.InputScope
}

func New(
	repository work.Repository,
	config *settings.Config,
	stats *stats.Stats,
	estimatedFileCount int,
) (*Orchestrator, error) {
	parallel := getParallel(estimatedFileCount, config)
	log.Debug().Msgf("number of workers: %d", parallel)

	return &Orchestrator{
		repository:          repository,
		config:              config,
		maxWorkersSemaphore: make(chan struct{}, parallel),
		done:                make(chan struct{}),
		pool:                pool.New(config, stats),
	}, nil
}

func (orchestrator *Orchestrator) Scan(reportPath string, inputs []files.File) (retErr error) {
	scope := orchestrator.config.ReportInputs
	ownScope := scope == nil
	var err error
	if ownScope {
		scope, err = inputfile.OpenInputScope(orchestrator.repository.Dir)
	}
	if err != nil {
		return err
	}
	orchestrator.inputScope = scope
	defer func() {
		if ownScope {
			retErr = errors.Join(retErr, scope.Close())
		} else {
			retErr = errors.Join(retErr, scope.Validate())
		}
		orchestrator.inputScope = nil
	}()
	fileComplete := make(chan coverage.Outcome, len(inputs))
	reportFile, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, reportFile.Close()) }()
	configBytes, err := json.Marshal(orchestrator.config)
	if err != nil {
		return err
	}
	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return err
	}
	names := make([]string, 0, len(inputs))
	for _, f := range inputs {
		names = append(names, f.FilePath)
	}
	record := coverage.Record{Schema: coverage.Schema, RunID: hex.EncodeToString(nonce), ConfigurationSHA256: coverage.Digest(configBytes), Discovery: coverage.Discovery{Files: names, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{}}
	// Reject duplicate/invalid work before starting any worker.
	shape := record
	shape.Complete = true
	for _, f := range names {
		shape.Outcomes = append(shape.Outcomes, coverage.Outcome{Path: f, SHA256: coverage.Digest(nil), Status: "success"})
	}
	if err = shape.Validate(names); err != nil {
		return err
	}
	var workers sync.WaitGroup
	// Join workers before closing their report; cancellation is never completion.
	defer workers.Wait()
	for _, f := range inputs {
		select {
		case <-orchestrator.done:
			return errors.New("processing cancelled")
		default:
		}
		workers.Add(1)
		go func(f files.File) { defer workers.Done(); fileComplete <- orchestrator.scanFile(reportFile, f) }(f)
	}
	outcomes, err := orchestrator.waitForScan(fileComplete, len(inputs))
	workers.Wait()
	err = errors.Join(err, scope.Validate())
	record.Outcomes = outcomes
	if err != nil {
		return err
	}
	record.Complete = true
	if err = record.Validate(names); err != nil {
		return err
	}
	select {
	case <-orchestrator.done:
		return errors.New("processing cancelled")
	default:
	}
	if err = orchestrator.writeFileList(reportFile, inputs); err != nil {
		return err
	}
	encoded, err := json.Marshal(detections.CoverageDetection{Type: detections.TypeCoverage, Coverage: record})
	if err != nil {
		return err
	}
	if err = coverage.WriteFull(reportFile, append(encoded, '\n')); err != nil {
		return err
	}
	return reportFile.Sync()
}

func (orchestrator *Orchestrator) waitForScan(ch <-chan coverage.Outcome, total int) ([]coverage.Outcome, error) {
	progressBar := bearerprogress.GetProgressBar(total, orchestrator.config)
	defer progressBar.Close()
	outcomes := make([]coverage.Outcome, 0, total)
	for len(outcomes) < total {
		select {
		case <-orchestrator.done:
			return outcomes, errors.New("processing cancelled")
		case outcome, ok := <-ch:
			if !ok {
				return outcomes, errors.New("missing terminal outcomes")
			}
			outcomes = append(outcomes, outcome)
			_ = progressBar.Add(1)
		}
	}
	select {
	case <-orchestrator.done:
		return outcomes, errors.New("processing cancelled")
	default:
	}
	return outcomes, nil
}

func (orchestrator *Orchestrator) scanFile(reportFile io.Writer, file files.File) (out coverage.Outcome) {
	out = coverage.Outcome{Path: file.FilePath, Status: "error", Error: "processing did not complete"}
	select {
	case <-orchestrator.done:
		out.Status = "cancelled"
		return
	case orchestrator.maxWorkersSemaphore <- struct{}{}:
	}
	defer func() { <-orchestrator.maxWorkersSemaphore }()
	select {
	case <-orchestrator.done:
		out.Status = "cancelled"
		return
	default:
	}
	fail := func(err error) coverage.Outcome { out.Error = err.Error(); return out }
	scope := orchestrator.inputScope
	ownScope := false
	if scope == nil {
		var err error
		scope, err = inputfile.OpenInputScope(orchestrator.repository.Dir)
		if err != nil {
			return fail(err)
		}
		ownScope = true
	}
	bound, err := scope.Acquire(file.FilePath, nil)
	if err != nil {
		if ownScope {
			err = errors.Join(err, scope.Close())
		}
		return fail(err)
	}
	defer func() {
		err := bound.Close()
		if ownScope {
			err = errors.Join(err, scope.Close())
		}
		if err != nil {
			out.Status = "error"
			out.SHA256 = ""
			out.Error = err.Error()
		}
	}()
	expectation := bound.Expectation()
	temp, err := os.CreateTemp("", "bearer-outcome-*.jsonl")
	if err != nil {
		return fail(err)
	}
	tmpPath := temp.Name()
	defer os.Remove(tmpPath)
	if err = temp.Close(); err != nil {
		return fail(err)
	}
	if err = orchestrator.pool.Scan(work.ProcessRequest{Repository: orchestrator.repository, File: file, ReportPath: tmpPath, InputExpectation: expectation}); err != nil {
		writeErr := orchestrator.writeFileError(reportFile, file, err)
		return fail(errors.Join(err, writeErr))
	}
	select {
	case <-orchestrator.done:
		out.Status = "cancelled"
		return
	default:
	}
	if err = bound.Validate(); err != nil {
		return fail(err)
	}
	if err = bound.Close(); err != nil {
		return fail(err)
	}
	if err = orchestrator.writeFileResult(reportFile, tmpPath, expectation); err != nil {
		return fail(err)
	}
	out.Status = "success"
	out.Error = ""
	out.SHA256 = expectation.SHA256
	return
}

func (orchestrator *Orchestrator) Close() {
	orchestrator.closeOnce.Do(func() { close(orchestrator.done); orchestrator.pool.Close() })
}

func (orchestrator *Orchestrator) writeFileList(reportFile io.Writer, inputs []files.File) error {
	names := make([]string, len(inputs))
	for i, f := range inputs {
		names[i] = f.FilePath
	}
	b, err := json.Marshal(detections.FileListDetection{Type: detections.TypeFileList, Filenames: names})
	if err != nil {
		return err
	}
	return coverage.WriteFull(reportFile, append(b, '\n'))
}

func (orchestrator *Orchestrator) writeFileResult(reportFile io.Writer, tmpPath string, expected ...*inputfile.InputExpectation) error {
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	var contextReceipt *inputfile.ContextReceipt
	for {
		var raw json.RawMessage
		if err = decoder.Decode(&raw); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err = coverage.CheckJSON(raw); err != nil {
			return err
		}
		var value map[string]any
		if err = json.Unmarshal(raw, &value); err != nil {
			return err
		}
		kind, ok := value["type"].(string)
		if !ok || kind == "" {
			return errors.New("malformed processing diagnostic")
		}
		if kind == inputfile.ContextReceiptType {
			if contextReceipt != nil {
				return errors.New("duplicate context handoff")
			}
			contextReceipt, err = inputfile.DecodeContextReceipt(raw)
			if err != nil {
				return err
			}
			continue
		}
		switch detections.DetectionType(kind) {
		case detections.TypeError, detections.TypeFileFailed, detections.TypeCoverage, detections.TypeFileList:
			return errors.New("error or unexpected terminal diagnostic in worker output")
		}
	}
	if len(expected) > 0 {
		if len(expected) != 1 || expected[0] == nil || contextReceipt == nil || orchestrator.inputScope == nil {
			return errors.New("missing context handoff")
		}
		if err := orchestrator.inputScope.ImportContextReceipt(contextReceipt, expected[0]); err != nil {
			return err
		}
	}
	orchestrator.reportMutex.Lock()
	defer orchestrator.reportMutex.Unlock()
	// Ensure adjacent workers cannot concatenate valid JSON objects on one line.
	if len(data) > 0 && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	return coverage.WriteFull(reportFile, data)
}

func (orchestrator *Orchestrator) writeFileError(reportFile io.Writer, file files.File, fileErr error) error {
	size := -1
	if info, err := os.Stat(path.Join(orchestrator.repository.Dir, file.FilePath)); err == nil {
		size = int(info.Size())
	}
	b, err := json.Marshal(detections.FileFailedDetection{Type: detections.TypeFileFailed, File: file.FilePath, FileSize: size, Timeout: file.Timeout, Error: fileErr.Error()})
	if err != nil {
		return err
	}
	orchestrator.reportMutex.Lock()
	defer orchestrator.reportMutex.Unlock()
	return coverage.WriteFull(reportFile, append(b, '\n'))
}

func getParallel(fileCount int, config *settings.Config) int {
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
