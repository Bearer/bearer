package balancer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	config "github.com/bearer/curio/pkg/commands/process/settings"
	workertype "github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/rs/zerolog/log"
	"github.com/struCoder/pidusage"
)

type Process struct {
	context context.Context
	kill    context.CancelFunc

	chunkDone      chan *workertype.ProcessResponse
	processErrored chan *workertype.ProcessResponse

	port int

	workeruuid string
	uuid       string

	task   *Task
	client *http.Client

	config config.Config
}

func (process *Process) StartProcess(task *workertype.ProcessRequest) error {
	port := ":" + strconv.Itoa(process.port)
	var err error
	currentCommand, err := os.Executable()
	if err != nil {
		log.Fatal().Msgf("failed to get current command executable %e", err)
	}

	args := []string{"processing-worker", "--port=" + port}
	if process.config.Scan.Debug {
		args = append(args, "--debug")
	}

	log.Debug().Msgf("spawning on port %s", port)
	cmd := exec.Command(currentCommand, args...)
	cmd.Dir, err = os.Getwd()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("couldn't determine current working dir %w", err)).Send()
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	process.client = &http.Client{
		Timeout: 0,
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to start process %w", err)).Send()
	}

	go process.monitorMemory(cmd.Process.Pid)
	go process.monitorRunning(cmd)

	err = process.WaitForOnline(task)
	if err != nil {
		log.Fatal().Msgf("%s %s failed for worker to become online", process.workeruuid, process.uuid)
		return err
	}

	go func() {
		<-process.context.Done()
		log.Debug().Msgf("%s %s real process handler done trying to kill process", process.workeruuid, process.uuid)
		err = cmd.Process.Kill()
		if err != nil {
			log.Debug().Msgf("failed killing process %e", err)
		}
	}()

	return nil
}

func (process *Process) WaitForOnline(task *workertype.ProcessRequest) error {
	start := time.Now()
	killTime := time.Now().Add(process.config.Worker.Timeout)

	closeCalled := false

	go func() {
		<-process.context.Done()
		closeCalled = true
	}()

	for {
		if closeCalled {
			return nil
		}

		if time.Now().After(killTime) {
			return ErrorProcessNotSpawned
		}

		taskBytes, err := json.Marshal(task)
		if err != nil {
			log.Debug().Msgf("failed to marshall task %e", err)
			continue
		}

		url := "http://localhost:" + strconv.Itoa(process.port) + workertype.RouteStatus

		req, err := http.NewRequestWithContext(process.context, http.MethodPost, url, bytes.NewBuffer(taskBytes))
		if err != nil {
			log.Debug().Msgf("%s %s failed to build status online request %e", process.uuid, process.workeruuid, err)
			continue
		}

		resp, err := process.client.Do(req)

		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			log.Debug().Msgf("%s spawned after %.2f seconds", process.uuid, time.Since(start).Seconds())

			defer resp.Body.Close()

			var result workertype.ProcessResponse
			err := json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				return err
			}

			if result.Error != nil {
				return fmt.Errorf("custom config invalid : %w", result.Error)
			}

			return nil
		}
	}
}

func (process *Process) doTask(task *Task) {
	// add timer for canceling process
	resultChan := make(chan *workertype.ProcessResponse, 1)
	go func() {
		taskBytes, err := json.Marshal(task.Definition)
		if err != nil {
			log.Debug().Msgf("failed to marshall task %e", err)
			return
		}

		url := "http://localhost:" + strconv.Itoa(process.port) + workertype.RouteProcess
		req, err := http.NewRequestWithContext(process.context, http.MethodPost, url, bytes.NewBuffer(taskBytes))
		if err != nil {
			log.Debug().Msgf("%s %s failed to build process request %e", process.uuid, process.workeruuid, err)
			return
		}

		resp, err := process.client.Do(req)

		if err != nil {
			log.Debug().Msgf("%s %s failed to do process request %e", process.uuid, process.workeruuid, err)
			return
		}

		defer resp.Body.Close()

		var result workertype.ProcessResponse
		json.NewDecoder(resp.Body).Decode(&result) //nolint:all,errcheck

		resultChan <- &result
	}()

	totalTime := time.Duration(0)

	for _, file := range task.Definition.Files {
		totalTime += file.Timeout
	}

	timeout := time.NewTimer(totalTime)
	select {
	case <-process.context.Done():
		log.Debug().Msgf("%s %s doing task closing out", process.workeruuid, process.uuid)
	case result := <-resultChan:
		process.chunkDone <- result
	case <-timeout.C:
		log.Debug().Msgf("%s %s timeout reporting error", process.workeruuid, process.uuid)
		process.processErrored <- &workertype.ProcessResponse{Error: ErrorTimeoutReached}
	}

}

func (process *Process) monitorMemory(pid int) {
	recovery := func() {
		if r := recover(); r != nil {
			log.Debug().Msgf("error recovered %s %s", r, debug.Stack())
		}
	}
	defer recovery()

	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("%s %s memory closing out", process.workeruuid, process.uuid)
			return
		case <-t.C:
			stats, err := pidusage.GetStat(pid)
			if err != nil {
				log.Debug().Msgf("failed to get memory usage %s", err)
				continue
			}

			if stats.Memory > float64(process.config.Worker.MemoryMaximum) {
				log.Debug().Msgf("%s %s memory reporting error", process.workeruuid, process.uuid)
				process.processErrored <- &workertype.ProcessResponse{Error: ErrorOutOfMemory}
				return
			}
		}
	}
}

func (process *Process) monitorRunning(cmd *exec.Cmd) {
	ended := make(chan bool, 1)
	go func() {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("%s %s running closing out", process.workeruuid, process.uuid)
			return
		case <-ended:
			log.Debug().Msgf("%s %s running reporting error", process.workeruuid, process.uuid)
			process.processErrored <- &workertype.ProcessResponse{Error: ErrorProcessCrashed}
			return
		}
	}()

	cmd.Wait() //nolint:all,errcheck
	ended <- true
}
