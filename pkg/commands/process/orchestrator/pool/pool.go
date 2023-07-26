package pool

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/rs/zerolog/log"
)

type Pool struct {
	processOptions ProcessOptions
	stats          *stats.Stats
	mutex          sync.Mutex
	nextId         int
	closed         bool
	available      []*Process
}

func New(config settings.Config, stats *stats.Stats) *Pool {
	executable, err := os.Executable()
	if err != nil {
		output.Fatal(fmt.Sprintf("failed to get current command executable %s", err))
	}

	baseArguments := []string{"processing-worker", "--log-level", config.Scan.LogLevel}
	if config.DebugProfile {
		baseArguments = append(baseArguments, "--debug-profile")
	}

	return &Pool{
		processOptions: ProcessOptions{
			executable:    executable,
			baseArguments: baseArguments,
			config:        config,
		},
		stats: stats,
	}
}

func (pool *Pool) Scan(request work.ProcessRequest) error {
	process, err := pool.get()
	if err != nil {
		return err
	}

	startTime := time.Now()
	log.Debug().Msgf("processing file %s using %s", request.File.FilePath, process.id)

	response, err := process.Scan(request)
	if err != nil {
		process.Close()
		pool.stats.FileFailed(
			request.File.FilePath,
			translateErrorForStats(err.Error(), true),
			startTime,
			process.memoryUsage,
		)
		return err
	}

	pool.stats.AddFileStats(response.FileStats)

	pool.mutex.Lock()
	pool.available = append(pool.available, process)
	pool.mutex.Unlock()

	if response.Error != "" {
		pool.stats.FileFailed(
			request.File.FilePath,
			translateErrorForStats(response.Error, false),
			startTime,
			0,
		)
		return errors.New(response.Error)
	}

	duration := pool.stats.File(request.File.FilePath, startTime)
	log.Debug().Msgf(
		"processing complete for %s [%s]",
		request.File.FilePath,
		duration.Truncate(time.Millisecond),
	)

	return nil
}

func (pool *Pool) get() (*Process, error) {
	pool.mutex.Lock()

	if pool.closed {
		pool.mutex.Unlock()
		return nil, errors.New("pool closed")
	}

	if len(pool.available) != 0 {
		process := pool.available[len(pool.available)-1]
		pool.available = pool.available[:len(pool.available)-1]
		pool.mutex.Unlock()
		return process, nil
	}

	id := fmt.Sprintf("worker-%d", pool.nextId)
	pool.nextId++

	pool.mutex.Unlock()

	process, err := newProcess(&pool.processOptions, id)
	if err != nil {
		return nil, fmt.Errorf("error spawning %s: %w", id, err)
	}

	return process, nil
}

func (pool *Pool) Close() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.closed {
		return
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(pool.available))

	for _, process := range pool.available {
		go func(process *Process) {
			process.Close()
			waitGroup.Done()
		}(process)
	}

	waitGroup.Wait()
	pool.closed = true
}

func translateErrorForStats(message string, processKilled bool) stats.FailedReason {
	switch message {
	case ErrorOutOfMemory.Error():
		return stats.MemoryLimitReason
	case worker.ErrorTimeoutReached.Error():
		if processKilled {
			return stats.KilledAfterTimeoutReason
		}

		return stats.TimeoutReason
	default:
		return stats.UnexpectedReason
	}
}
