package balancer

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	config "github.com/bearer/curio/pkg/commands/process/settings"
	workertype "github.com/bearer/curio/pkg/commands/process/worker/work"
)

var ErrorOutOfMemory = errors.New("process went out of memory")
var ErrorProcessCrashed = errors.New("process crashed due to panic or similar")
var ErrorTimeoutReached = errors.New("process hasnt completed in a given time")
var ErrorClosing = errors.New("process hasnt completed but a shutdown was recieved")
var ErrorProcessNotSpawned = errors.New("process wasn't able to come online in a given time")

type Monitor struct {
	workStack []*Task
	workers   []*Worker

	taskRequest  chan *Task
	taskComplete chan *Worker

	context       context.Context
	contextCancel context.CancelFunc

	config config.Config
}

type Task struct {
	Definition workertype.ProcessRequest
	Done       chan *workertype.ProcessResponse
	worker     *Worker
}

func New(settings config.Config) *Monitor {
	monitorCtx, monitorCancel := context.WithCancel(context.Background())
	monitor := &Monitor{
		context:       monitorCtx,
		contextCancel: monitorCancel,

		taskRequest:  make(chan *Task),
		taskComplete: make(chan *Worker, settings.Worker.Workers),

		config: settings,
	}
	go monitor.monitor()

	return monitor
}

func (monitor *Monitor) Close() {
	log.Debug().Msgf("closing")
	monitor.contextCancel()
}

func (monitor *Monitor) ScheduleTask(scanRequest workertype.ProcessRequest) *Task {
	scanRequest.CustomDetectorConfig = monitor.config.CustomDetector.RulesConfig
	task := &Task{
		Definition: scanRequest,
		Done:       make(chan *workertype.ProcessResponse, 1),
	}
	go func() {
		monitor.taskRequest <- task
	}()

	return task
}

func (monitor *Monitor) monitor() {
	log.Debug().Msgf("balancer spawning monitor.....")
	for {
		select {
		case <-monitor.context.Done():
			log.Debug().Msgf("killing workers")
			for _, worker := range monitor.workers {
				log.Debug().Msgf("balancer killing worker %s", worker.uuid)
				worker.kill()
			}
			return
		case task := <-monitor.taskRequest:
			log.Debug().Msgf("balancer got task request")
			if len(monitor.workers) > monitor.config.Worker.Workers {
				log.Debug().Msgf("balancer adding work to stack")

				monitor.workStack = append(monitor.workStack, task)
			} else {
				log.Debug().Msgf("balancer directly offloading work to worker")
				log.Debug().Msgf("balancer got %d number of workers online", len(monitor.workers))
				monitor.spawnWorker(task)
			}
		case worker := <-monitor.taskComplete:
			for i, monitoringWorkers := range monitor.workers {
				if worker == monitoringWorkers {
					log.Debug().Msgf("balancer removing worker from stack")
					monitor.workers = append(monitor.workers[:i], monitor.workers[i+1:]...)
				}
			}
			log.Debug().Msgf("balancer killing previous worker %s", worker.uuid)
			worker.kill()
			if len(monitor.workStack) > 0 {
				log.Debug().Msgf("balancer spawning new worker")

				task := monitor.workStack[0]
				monitor.workStack = monitor.workStack[1:]
				monitor.spawnWorker(task)
			}
		}
	}
}

func (monitor *Monitor) spawnWorker(task *Task) *Worker {
	log.Debug().Msgf("balancer spawning worker")
	ctx, ctxCancel := context.WithCancel(context.Background())
	worker := &Worker{
		context: ctx,
		kill:    ctxCancel,

		taskComplete: monitor.taskComplete,

		port: GetFreePort(),

		chunkDone:      make(chan *workertype.ProcessResponse, 1),
		processErrored: make(chan *workertype.ProcessResponse, 1),

		uuid: uuid.NewString(),

		config: monitor.config,
	}

	task.worker = worker
	worker.task = task
	monitor.workers = append(monitor.workers, worker)

	go worker.Start()
	return worker
}

func GetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to resolve localhost %w", err)).Send()
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to resolve address %w", err)).Send()
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
