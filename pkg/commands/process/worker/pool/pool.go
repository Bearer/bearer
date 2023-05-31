package pool

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
)

type Pool struct {
	processOptions ProcessOptions
	mutex          sync.Mutex
	nextId         int
	closed         bool
	available      []*Process
}

func New(config settings.Config) *Pool {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal().Msgf("failed to get current command executable %s", err)
	}

	baseArguments := []string{"processing-worker"}
	if config.Scan.Debug {
		baseArguments = append(baseArguments, "--debug")
	}

	return &Pool{
		processOptions: ProcessOptions{
			executable:    executable,
			baseArguments: baseArguments,
			config:        config,
		},
	}
}

func (pool *Pool) Scan(request work.ProcessRequest) error {
	process, err := pool.get()
	if err != nil {
		return err
	}

	response, err := process.Scan(request)
	if err != nil {
		process.Close()
		return err
	}

	pool.mutex.Lock()
	pool.available = append(pool.available, process)
	pool.mutex.Unlock()

	if response.Error != "" {
		return errors.New(response.Error)
	}

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
