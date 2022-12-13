package flag

import (
	"time"
)

var (
	WorkersFlag = Flag{
		Name:       "workers",
		ConfigName: "worker.workers",
		Value:      1,
		Usage:      "The number of processing workers to spawn.",
	}
	TimeoutFlag = Flag{
		Name:       "timeout",
		ConfigName: "worker.timeout",
		Value:      10 * time.Minute,
		Usage:      "The maximum time alloted to complete the scan.",
	}
	TimeoutFileMinimumFlag = Flag{
		Name:       "timeout-file-min",
		ConfigName: "worker.timeout-file-min",
		Value:      5 * time.Second,
		Usage:      "Minimum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes.",
	}
	TimeoutFileMaximumFlag = Flag{
		Name:       "timeout-file-max",
		ConfigName: "worker.timeout-file-max",
		Value:      30 * time.Second,
		Usage:      "Maximum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes.",
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
		Usage:      "Maximum time to wait for a worker process to come online.",
	}
	FileSizeMaximumFlag = Flag{
		Name:       "file-size-max",
		ConfigName: "worker.file-size-max",
		Value:      2 * 1000 * 1000, // 2MB
		Usage:      "Ignore files larger than the specified value.",
	}
	FilesToBatchFlag = Flag{
		Name:       "files-to-batch",
		ConfigName: "worker.files-to-batch",
		Value:      1,
		Usage:      "Specify the number of files to batch per worker.",
	}
	MemoryMaximumFlag = Flag{
		Name:       "memory-max",
		ConfigName: "worker.memory-max",
		Value:      800 * 1000 * 1000, // 800 MB
		Usage:      "If the memory needed to scan a file surpasses the specified limit, skip the file.",
	}
	ExistingWorkerFlag = Flag{
		Name:       "existing-worker",
		ConfigName: "worker.existing-worker",
		Value:      "",
		Usage:      "Specify the URL of an existing worker.",
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
	Workers                   int           `mapstructure:"workers" json:"workers" yaml:"workers"`
	Timeout                   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	TimeoutFileMinimum        time.Duration `mapstructure:"timeout-file-min" json:"timeout-file-min" yaml:"timeout-file-min"`
	TimeoutFileMaximum        time.Duration `mapstructure:"timeout-file-max"  json:"timeout-file-max" yaml:"timeout-file-max"`
	TimeoutFileSecondPerBytes int           `mapstructure:"timeout-file-second-per-bytes" json:"timeout-file-second-per-bytes" yaml:"timeout-file-second-per-bytes"`
	TimeoutWorkerOnline       time.Duration `mapstructure:"timeout-worker-online" json:"timeout-worker-online" yaml:"timeout-worker-online"`
	FileSizeMaximum           int           `mapstructure:"file-size-max" json:"file-size-max" yaml:"file-size-max"`
	FilesToBatch              int           `mapstructure:"files-to-batch" json:"files-to-batch" yaml:"files-to-batch"`
	MemoryMaximum             int           `mapstructure:"memory-max" json:"memory-max" yaml:"memory-max"`
	ExistingWorker            string        `mapstructure:"existing-worker" json:"existing-worker" yaml:"existing-worker"`
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
