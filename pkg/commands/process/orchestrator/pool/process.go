package pool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	gopsutilprocess "github.com/shirou/gopsutil/v3/process"

	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/coverage"
)

var (
	ErrorCrashed     = errors.New("exited unexpectedly")
	ErrorNotSpawned  = errors.New("didn't start within expected time")
	ErrorOutOfMemory = errors.New("exceeded memory limit")
)

type Process struct {
	id            string
	command       *exec.Cmd
	context       context.Context
	cancelContext context.CancelFunc
	errorChannel  chan error
	exitChannel   chan struct{}
	client        *http.Client
	baseURL       string
	memoryUsage   uint64
}

type ProcessOptions struct {
	executable    string
	baseArguments []string
	config        *settings.Config
	parentContext context.Context
}

func newProcess(options *ProcessOptions, id string) (*Process, error) {
	parent := options.parentContext
	if parent == nil {
		parent = context.Background()
	}
	if err := parent.Err(); err != nil {
		return nil, err
	}
	port, err := allocatePort()
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("%s spawning on port %d", id, port)

	arguments := append(
		append([]string(nil), options.baseArguments...),
		"--parent-process-id", strconv.Itoa(syscall.Getpid()),
		"--port", strconv.Itoa(port),
		"--worker-id", id,
	)
	command := exec.Command(options.executable, arguments...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	processContext, cancelContext := context.WithCancel(parent)

	process := &Process{
		id:            id,
		command:       command,
		context:       processContext,
		cancelContext: cancelContext,
		errorChannel:  make(chan error, 1),
		exitChannel:   make(chan struct{}),
		client:        &http.Client{Timeout: 0},
		baseURL:       fmt.Sprintf("http://localhost:%d", port),
	}

	if err := process.start(options.config); err != nil {
		process.Close()
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return process, nil
}

func (process *Process) start(config *settings.Config) error {
	if err := process.context.Err(); err != nil {
		close(process.exitChannel)
		return err
	}
	if err := process.command.Start(); err != nil {
		close(process.exitChannel)
		return err
	}

	go process.monitorCommand()
	go process.monitorMemory()

	if err := process.initialize(config); err != nil {
		if process.context.Err() != nil {
			return process.context.Err()
		}
		var result = strings.Split(err.Error(), "failed to create detector customDetector:")
		if len(result) > 1 {
			// custom detector issue ; assume custom rule parse issue
			var ruleName = strings.TrimSpace(strings.Split(result[1], ":")[0])
			log.Debug().Msgf("%s", err.Error())
			return fmt.Errorf("could not parse rule %s. Is this a custom rule? See documentation on rule patterns and format https://docs.bearer.com/guides/custom-rule/: %w", ruleName, err)
		}
		return fmt.Errorf("failed to start bearer, error with your configuration: %w", err)
	}

	return nil
}

func (process *Process) monitorCommand() {
	go func() {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("shutting down %s", process.id)

			if err := process.command.Process.Signal(os.Interrupt); err != nil {
				log.Debug().Msgf("killing %s due to error sending interrupt: %s", process.id, err)
				process.kill()
				return
			}

			timeout := time.NewTimer(settings.TimeoutWorkerShutdown)
			select {
			case <-timeout.C:
				log.Debug().Msgf("killing %s after timeout", process.id)
				process.kill()
			case <-process.exitChannel:
				log.Debug().Msgf("%s stopped", process.id)
			}
		case <-process.exitChannel:
			process.errorChannel <- ErrorCrashed
			return
		}
	}()

	process.command.Wait() //nolint:errcheck
	close(process.exitChannel)
}

func (process *Process) kill() {
	if err := process.command.Process.Kill(); err != nil {
		log.Debug().Msgf("%s failed killing process: %s", process.id, err)
	}
}

func (process *Process) monitorMemory() {
	tick := time.NewTicker(500 * time.Millisecond)
	monitor, err := gopsutilprocess.NewProcessWithContext(process.context, int32(process.command.Process.Pid))
	if err != nil {
		log.Debug().Msgf("failed to start memory monitor: %s", err)
		return
	}

	for {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("%s memory monitor shutting down", process.id)
			return
		case <-tick.C:
			stats, err := monitor.MemoryInfo()
			if err != nil {
				log.Debug().Msgf("failed to get memory usage %s", err)
				continue
			}

			if stats.RSS > settings.MemoryMaximum {
				process.memoryUsage = stats.RSS
				process.errorChannel <- ErrorOutOfMemory
				return
			}

			if stats.RSS > settings.MemorySoftMaximum {
				process.reduceMemoryUsage()
			}
		}
	}
}

func (process *Process) reduceMemoryUsage() {
	request, err := process.buildRequest(work.RouteReduceMemory, nil)
	if err != nil {
		log.Debug().Msgf("%s failed to build memory reduction request %s", process.id, err)
	}

	response, err := process.client.Do(request)
	if err != nil {
		log.Debug().Msgf("%s failed to request memory reduction: %s", process.id, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Debug().Msgf("%s memory reduction not ok response", process.id)
	}
}

func (process *Process) initialize(config *settings.Config) error {
	log.Debug().Msgf("%s initializing", process.id)
	start := time.Now()
	killTime := time.Now().Add(config.Worker.TimeoutWorkerOnline)

	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("couldn't marshal config %w", err)
	}

	request, err := process.buildRequest(work.RouteInitialize, bytes.NewBuffer(marshalledConfig))
	if err != nil {
		return fmt.Errorf("%s failed to build initialization request %s", process.id, err)
	}

	for {
		if process.context.Err() != nil {
			return process.context.Err()
		}
		select {
		case <-process.exitChannel:
			return ErrorCrashed
		default:
		}
		if time.Now().After(killTime) {
			return ErrorNotSpawned
		}

		response, err := process.client.Do(request)
		if err != nil {
			timer := time.NewTimer(100 * time.Millisecond)
			select {
			case <-process.context.Done():
				timer.Stop()
				return process.context.Err()
			case <-process.exitChannel:
				timer.Stop()
				return ErrorCrashed
			case <-timer.C:
			}
			continue
		}

		online, err := process.initializationResponse(response)
		if !online {
			continue
		}

		log.Debug().Msgf("%s spawned after %.2f seconds", process.id, time.Since(start).Seconds())
		if err != nil {
			return err
		}

		log.Debug().Msgf("%s is online", process.id)
		return nil
	}
}

func (process *Process) Scan(scanRequest work.ProcessRequest) (*work.ProcessResponse, error) {
	// Each request gets its own deadline and cancellation. A buffered result lets
	// a late transport finish and close its body after the caller has returned.
	ctx, cancel := context.WithTimeout(process.context, scanRequest.File.Timeout+settings.TimeoutWorkerFileGrace)
	defer cancel()
	taskBytes, err := json.Marshal(scanRequest)
	if err != nil {
		return nil, errors.New("failed to encode worker scan request")
	}
	request, err := process.buildRequest(work.RouteProcess, bytes.NewBuffer(taskBytes))
	if err != nil {
		return nil, errors.New("failed to construct worker scan request")
	}
	request = request.WithContext(ctx)
	type result struct {
		response *work.ProcessResponse
		err      error
	}
	completed := make(chan result, 1)
	go func() {
		response, err := process.client.Do(request)
		if err != nil {
			// Client.Do owns body cleanup on its response-plus-error path. Never echo
			// transport errors, which may include URLs or response text.
			completed <- result{err: errors.New("worker scan transport failed")}
			return
		}
		decoded, err := readScanResponse(response)
		completed <- result{response: decoded, err: err}
	}()
	contextFailure := func() (*work.ProcessResponse, error) {
		cancel()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			process.Close()
			return nil, worker.ErrorTimeoutReached
		}
		return nil, ctx.Err()
	}
	select {
	case result := <-completed:
		// A response arriving with/after cancellation is never accepted as success.
		if ctx.Err() != nil {
			return contextFailure()
		}
		return result.response, result.err
	case err, ok := <-process.errorChannel:
		cancel()
		process.Close()
		if !ok || err == nil {
			return nil, ErrorCrashed
		}
		return nil, err
	case <-ctx.Done():
		return contextFailure()
	}
}

// readScanResponse validates the actual worker envelope before a response can
// become a successful processing outcome. Errors intentionally omit body bytes.
func readScanResponse(response *http.Response) (result *work.ProcessResponse, err error) {
	if response == nil || response.Body == nil {
		return nil, errors.New("worker response body unavailable")
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil && err == nil {
			result = nil
			err = errors.New("failed to close worker response body")
		}
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("worker scan returned HTTP status %d", response.StatusCode)
	}
	data, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, errors.New("failed to read worker response body")
	}
	if coverage.CheckJSON(data) != nil {
		return nil, errors.New("invalid worker response JSON")
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(data, &fields) != nil || len(fields) != 2 {
		return nil, errors.New("incomplete worker response envelope")
	}
	fileStats, hasStats := fields["FileStats"]
	errorText, hasError := fields["Error"]
	if !hasStats || !hasError || bytes.Equal(bytes.TrimSpace(errorText), []byte("null")) {
		return nil, errors.New("incomplete worker response envelope")
	}
	var message string
	if json.Unmarshal(errorText, &message) != nil {
		return nil, errors.New("invalid worker response error field")
	}
	if !bytes.Equal(bytes.TrimSpace(fileStats), []byte("null")) {
		var statsFields map[string]json.RawMessage
		if json.Unmarshal(fileStats, &statsFields) != nil || statsFields == nil {
			return nil, errors.New("invalid worker response stats field")
		}
	}
	var decoded work.ProcessResponse
	if json.Unmarshal(data, &decoded) != nil {
		return nil, errors.New("invalid worker response fields")
	}
	// Worker-reported errors remain in Error for Pool.Scan's existing handling.
	return &decoded, nil
}

func (process *Process) buildRequest(route string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(
		process.context,
		http.MethodPost,
		process.baseURL+route,
		body,
	)
}

func (process *Process) initializationResponse(response *http.Response) (bool, error) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, nil
	}

	var result work.InitializeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return true, fmt.Errorf("error decoding status response: %w", err)
	}

	if result.Error != "" {
		return true, errors.New(result.Error)
	}

	return true, nil
}

func (process *Process) Close() {
	process.cancelContext()
	<-process.exitChannel
}

func allocatePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to resolve localhost %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve address %w", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
