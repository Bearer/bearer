package pool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/stats"
	"github.com/bearer/bearer/pkg/util/output"
)

type Pool struct {
	processOptions ProcessOptions
	stats          *stats.Stats
	mutex          sync.Mutex
	nextId         int
	closed         bool
	closeDone      chan struct{}
	available      []*Process
	owned          map[*Process]struct{}
	processFactory func(*ProcessOptions, string) (*Process, error)
	context        context.Context
	cancelContext  context.CancelFunc
	starting       map[chan struct{}]struct{}
	cleaning       map[chan struct{}]struct{}
}

func New(config *settings.Config, stats *stats.Stats) *Pool {
	executable, err := os.Executable()
	if err != nil {
		output.Fatal(fmt.Sprintf("failed to get current command executable %s", err))
	}

	baseArguments := []string{"processing-worker", "--log-level", config.LogLevel}
	if config.DebugProfile {
		baseArguments = append(baseArguments, "--debug-profile")
	}

	return &Pool{
		processOptions: ProcessOptions{
			executable:    executable,
			baseArguments: baseArguments,
			config:        config,
		},
		owned: make(map[*Process]struct{}),
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
		pool.remove(process)
		pool.stats.FileFailed(
			request.File.FilePath,
			translateErrorForStats(err.Error(), true),
			startTime,
			process.memoryUsage,
		)
		return err
	}

	if !pool.release(process) {
		return errors.New("pool closed")
	}
	pool.stats.AddFileStats(response.FileStats)

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
		pool.ensureOwnedLocked(process)
		pool.mutex.Unlock()
		return process, nil
	}

	id := fmt.Sprintf("worker-%d", pool.nextId)
	pool.nextId++
	// Own cancellation and completion before calling a factory that may block
	// during initialization. A process returned by that factory is too late.
	if pool.context == nil {
		pool.context, pool.cancelContext = context.WithCancel(context.Background())
	}
	if pool.starting == nil {
		pool.starting = make(map[chan struct{}]struct{})
	}
	started := make(chan struct{})
	pool.starting[started] = struct{}{}
	options := pool.processOptions
	options.parentContext = pool.context
	factory := pool.processFactory
	if factory == nil {
		factory = newProcess
	}
	pool.mutex.Unlock()
	defer func() {
		pool.mutex.Lock()
		delete(pool.starting, started)
		close(started)
		pool.mutex.Unlock()
	}()

	process, err := factory(&options, id)
	if err != nil {
		if process != nil {
			process.Close()
		}
		return nil, fmt.Errorf("error spawning %s: %w", id, err)
	}
	if process == nil {
		return nil, fmt.Errorf("error spawning %s: factory returned no process", id)
	}

	pool.mutex.Lock()
	if pool.closed {
		pool.mutex.Unlock()
		process.Close()
		return nil, errors.New("pool closed")
	}
	pool.ensureOwnedLocked(process)
	pool.mutex.Unlock()

	return process, nil
}

func (pool *Pool) Close() {
	pool.mutex.Lock()
	if pool.closed {
		done := pool.closeDone
		pool.mutex.Unlock()
		if done != nil {
			<-done
		}
		return
	}
	pool.closed = true
	pool.closeDone = make(chan struct{})
	cancel := pool.cancelContext
	starting := make([]chan struct{}, 0, len(pool.starting))
	for done := range pool.starting {
		starting = append(starting, done)
	}
	cleaning := make([]chan struct{}, 0, len(pool.cleaning))
	for done := range pool.cleaning {
		cleaning = append(cleaning, done)
	}
	done := pool.closeDone
	defer close(done)
	processes := make([]*Process, 0, len(pool.owned)+len(pool.available))
	seen := make(map[*Process]struct{}, cap(processes))
	for process := range pool.owned {
		seen[process] = struct{}{}
		processes = append(processes, process)
	}
	for _, process := range pool.available {
		if _, ok := seen[process]; !ok {
			processes = append(processes, process)
		}
	}
	pool.owned = nil
	pool.available = nil
	pool.mutex.Unlock()
	// Cancellation must precede every join, including factories that have not
	// returned yet. Their completion includes failure/late-return cleanup.
	if cancel != nil {
		cancel()
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(processes))

	for _, process := range processes {
		go func(process *Process) {
			process.Close()
			waitGroup.Done()
		}(process)
	}

	waitGroup.Wait()
	for _, done := range starting {
		<-done
	}
	for _, done := range cleaning {
		<-done
	}
}

func (pool *Pool) ensureOwnedLocked(process *Process) {
	if pool.owned == nil {
		pool.owned = make(map[*Process]struct{})
	}
	pool.owned[process] = struct{}{}
}

func (pool *Pool) release(process *Process) bool {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.closed {
		delete(pool.owned, process)
		return false
	}
	pool.ensureOwnedLocked(process)
	pool.available = append(pool.available, process)
	return true
}

// remove transfers a failed scan's process from reusable ownership to a
// joinable cleanup obligation before releasing the mutex. Only this caller
// closes it; a concurrent Pool.Close joins completion instead of closing twice.
func (pool *Pool) remove(process *Process) {
	pool.mutex.Lock()
	if pool.closed {
		// Close has already claimed all remaining owned processes.
		pool.mutex.Unlock()
		return
	}
	if _, ok := pool.owned[process]; !ok {
		pool.mutex.Unlock()
		return
	}
	if pool.cleaning == nil {
		pool.cleaning = make(map[chan struct{}]struct{})
	}
	done := make(chan struct{})
	pool.cleaning[done] = struct{}{}
	delete(pool.owned, process)
	pool.mutex.Unlock()

	process.Close()

	pool.mutex.Lock()
	delete(pool.cleaning, done)
	close(done)
	pool.mutex.Unlock()
}

func translateErrorForStats(message string, processKilled bool) stats.FailedReason {
	if message == ErrorOutOfMemory.Error() {
		return stats.MemoryLimitReason
	}
	// Worker completion may join context/capture/close errors. Preserve the
	// full response text and classify only an exact joined timeout-marker line.
	for _, line := range strings.Split(message, "\n") {
		if line == worker.ErrorTimeoutReached.Error() {
			if processKilled {
				return stats.KilledAfterTimeoutReason
			}
			return stats.TimeoutReason
		}
	}
	return stats.UnexpectedReason
}
