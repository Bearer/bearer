package flag

import (
	"time"
)

var (
	WorkersFlag = Flag{
		Name:       "workers",
		ConfigName: "worker.workers",
		Value:      1,
		Usage:      "number of processing workers to spawn",
	}
	TimeoutFlag = Flag{
		Name:       "timeout",
		ConfigName: "worker.timeout",
		Value:      10 * time.Minute,
		Usage:      "time allowed to complete scan",
	}
	TimeoutFileMinimumFlag = Flag{
		Name:       "timeout-file-min",
		ConfigName: "worker.timeout-file-min",
		Value:      5 * time.Second,
		Usage:      "minimum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes",
	}
	TimeoutFileMaximumFlag = Flag{
		Name:       "timeout-file-max",
		ConfigName: "worker.timeout-file-max",
		Value:      30 * time.Second,
		Usage:      "maximum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes",
	}
	TimeoutFileSecondPerBytesFlag = Flag{
		Name:       "timeout-file-second-per-bytes",
		ConfigName: "worker.timeout-file-second-per-bytes",
		Value:      10 * 1000, // 10kb/s
		Usage:      "number of file size bytes producing a second of timeout assigned to scanning a file",
	}
	TimeoutWorkerOnlineFlag = Flag{
		Name:       "timeout-worker-online",
		ConfigName: "worker.timeout-worker-online",
		Value:      60 * time.Second,
		Usage:      "maximum time for worker process to come online",
	}
	FileSizeMaximumFlag = Flag{
		Name:       "file-size-max",
		ConfigName: "worker.file-size-max",
		Value:      100 * 1000, // 100 KB
		Usage:      "ignore files with file size larger than this config",
	}
	FilesToBatchFlag = Flag{
		Name:       "files-to-batch",
		ConfigName: "worker.files-to-batch",
		Value:      1,
		Usage:      "number of files to batch to worker",
	}
	MemoryMaximumFlag = Flag{
		Name:       "memory-max",
		ConfigName: "worker.memory-max",
		Value:      800 * 1000 * 1000, // 800 MB
		Usage:      "if memory needed to scan a file surpasses this limit, skip the file",
	}
	ExistingWorkerFlag = Flag{
		Name:       "existing-worker",
		ConfigName: "worker.existing-worker",
		Value:      "",
		Usage:      "URL of an existing worker",
	}
)

type WorkerFlagGroup struct {
	Workers                   *Flag
	Timeout                   *Flag
	TimeoutFileMinimum        *Flag
	TimeoutFileMaximum        *Flag
	TimeoutFileSecondPerBytes *Flag
	TimeoutWorkerOnline       *Flag
	FileSizeMaximum           *Flag
	FilesToBatch              *Flag
	MemoryMaximum             *Flag
	ExistingWorker            *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type WorkerOptions struct {
	Workers                   int           `json:"workers" yaml:"workers"`
	Timeout                   time.Duration `json:"timeout" yaml:"timeout"`
	TimeoutFileMinimum        time.Duration `json:"timeout_file_minimum" yaml:"timeout_file_minimum"`
	TimeoutFileMaximum        time.Duration `json:"timeout_file_maximum" yaml:"timeout_file_maximum"`
	TimeoutFileSecondPerBytes int           `json:"timeout_file_second_per_bytes" yaml:"timeout_file_second_per_bytes"`
	TimeoutWorkerOnline       time.Duration `json:"timeout_worker_online" yaml:"timeout_worker_online"`
	FileSizeMaximum           int           `json:"file_size_maximum" yaml:"file_size_maximum"`
	FilesToBatch              int           `json:"files_to_batch" yaml:"files_to_batch"`
	MemoryMaximum             int           `json:"memory_maximum" yaml:"memory_maximum"`
	ExistingWorker            string        `json:"existing_worker" yaml:"existing_worker"`
}

func NewWorkerFlagGroup() *WorkerFlagGroup {
	return &WorkerFlagGroup{
		Workers:                   &WorkersFlag,
		Timeout:                   &TimeoutFlag,
		TimeoutFileMinimum:        &TimeoutFileMinimumFlag,
		TimeoutFileMaximum:        &TimeoutFileMaximumFlag,
		TimeoutFileSecondPerBytes: &TimeoutFileSecondPerBytesFlag,
		TimeoutWorkerOnline:       &TimeoutWorkerOnlineFlag,
		FileSizeMaximum:           &FileSizeMaximumFlag,
		FilesToBatch:              &FilesToBatchFlag,
		MemoryMaximum:             &MemoryMaximumFlag,
		ExistingWorker:            &ExistingWorkerFlag,
	}
}

func (f *WorkerFlagGroup) Name() string {
	return "Worker"
}

func (f *WorkerFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.Workers,
		f.Timeout,
		f.TimeoutFileMinimum,
		f.TimeoutFileMaximum,
		f.TimeoutFileSecondPerBytes,
		f.TimeoutWorkerOnline,
		f.FileSizeMaximum,
		f.FilesToBatch,
		f.MemoryMaximum,
		f.ExistingWorker,
	}
}

func (f *WorkerFlagGroup) ToOptions() WorkerOptions {
	return WorkerOptions{
		Workers:                   getInt(f.Workers),
		Timeout:                   getDuration(f.Timeout),
		TimeoutFileMinimum:        getDuration(f.TimeoutFileMinimum),
		TimeoutFileMaximum:        getDuration(f.TimeoutFileMaximum),
		TimeoutFileSecondPerBytes: getInt(f.TimeoutFileSecondPerBytes),
		TimeoutWorkerOnline:       getDuration(f.TimeoutWorkerOnline),
		FilesToBatch:              getInt(f.FilesToBatch),
		FileSizeMaximum:           getInt(f.FileSizeMaximum),
		MemoryMaximum:             getInt(f.MemoryMaximum),
		ExistingWorker:            getString(f.ExistingWorker),
	}
}
