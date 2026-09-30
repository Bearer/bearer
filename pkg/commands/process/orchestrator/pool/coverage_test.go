package pool

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/stats"
)

// Exercise the real Scan -> remove -> Process.Close handoff. The held channel
// represents process-exit completion, not startup readiness or a real worker.
func TestCoverageS1FC1FailureCleanupJoin(t *testing.T) {
	for _, mode := range []string{"transport", "http", "decode", "read", "body-close"} {
		t.Run(mode, func(t *testing.T) {
			created := make(chan *Process, 1)
			var cancelCalls atomic.Int32
			exit := make(chan struct{})
			var finish sync.Once
			complete := func() { finish.Do(func() { close(exit) }) }
			t.Cleanup(complete)
			p := &Pool{stats: stats.New()}
			p.processFactory = func(options *ProcessOptions, _ string) (*Process, error) {
				ctx, cancel := context.WithCancel(options.parentContext)
				process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
					if mode == "transport" {
						return nil, errors.New("synthetic transport failure")
					}
					status, raw := 200, `{"FileStats":{},"Error":""}`
					if mode == "http" {
						status = 503
					}
					if mode == "decode" {
						raw = "{"
					}
					body := &failingBody{reader: strings.NewReader(raw), closed: make(chan struct{})}
					if mode == "read" {
						body.readErr = errors.New("synthetic body read failure")
					}
					if mode == "body-close" {
						body.closeErr = errors.New("synthetic body close failure")
					}
					return &http.Response{StatusCode: status, Body: body}, nil
				}))
				process.context, process.exitChannel = ctx, exit
				process.cancelContext = func() { cancelCalls.Add(1); cancel() }
				created <- process
				return process, nil
			}
			scanDone := make(chan error, 1)
			go func() { scanDone <- p.Scan(processRequest()) }()
			var process *Process
			select {
			case process = <-created:
			case <-time.After(time.Second):
				t.Fatal("process not promoted")
			}
			waitClosed(t, process.context.Done())
			select {
			case <-scanDone:
				t.Fatal("failed Scan returned before held cleanup")
			default:
			}
			p.mutex.Lock()
			owned, cleaning, starting := len(p.owned), len(p.cleaning), len(p.starting)
			p.mutex.Unlock()
			if owned != 0 || cleaning != 1 || starting != 0 {
				t.Fatalf("handoff owned=%d cleaning=%d starting=%d", owned, cleaning, starting)
			}

			closed, entered := make(chan struct{}, 4), make(chan struct{}, 4)
			for i := 0; i < 4; i++ {
				go func() { entered <- struct{}{}; p.Close(); closed <- struct{}{} }()
			}
			for i := 0; i < 4; i++ {
				<-entered
			}
			waitClosed(t, p.context.Done())
			getDone := make(chan error, 1)
			go func() { _, err := p.get(); getDone <- err }()
			select {
			case err := <-getDone:
				if err == nil {
					t.Fatal("closed pool accepted get")
				}
			case <-time.After(time.Second):
				t.Fatal("pool mutex held during cleanup")
			}
			select {
			case <-closed:
				t.Fatal("Pool.Close returned before failed process exit")
			case <-time.After(20 * time.Millisecond):
			}
			if cancelCalls.Load() != 1 {
				t.Fatalf("cleanup owners=%d", cancelCalls.Load())
			}
			complete()
			select {
			case err := <-scanDone:
				if err == nil {
					t.Fatal("failed scan accepted")
				}
				t.Logf("terminal Scan failure: %v", err)
			case <-time.After(time.Second):
				t.Fatal("Scan cleanup stranded")
			}
			for i := 0; i < 4; i++ {
				select {
				case <-closed:
				case <-time.After(time.Second):
					t.Fatal("Close did not join completed cleanup")
				}
			}
			if cancelCalls.Load() != 1 || len(p.cleaning) != 0 || len(p.owned) != 0 {
				t.Fatal("cleanup ownership retained or repeated")
			}
			t.Log("actual failed Scan cleanup remained joinable; four Close callers waited; one cleanup owner")
		})
	}
}

func TestCoverageS1FC1FailureCleanupWithoutClose(t *testing.T) {
	failed := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return nil, errors.New("synthetic failure") }))
	healthy := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"FileStats":{},"Error":""}`))}, nil
	}))
	p := &Pool{stats: stats.New(), available: []*Process{failed}, processFactory: func(*ProcessOptions, string) (*Process, error) { return healthy, nil }}
	if err := p.Scan(processRequest()); err == nil {
		t.Fatal("failure lost")
	}
	if failed.context.Err() == nil || len(p.cleaning) != 0 || len(p.owned) != 0 || p.closed {
		t.Fatal("failed process not discharged in an open pool")
	}
	if err := p.Scan(processRequest()); err != nil {
		t.Fatalf("new healthy scan: %v", err)
	}
	if err := p.Scan(processRequest()); err != nil {
		t.Fatalf("healthy reuse: %v", err)
	}
	p.Close()
	if healthy.context.Err() == nil {
		t.Fatal("healthy replacement not closed")
	}
}

func TestCoverageS1FC1ResponseControls(t *testing.T) {
	for _, message := range []string{"", "synthetic worker failure"} {
		t.Run(map[bool]string{true: "valid", false: "worker-error"}[message == ""], func(t *testing.T) {
			process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"FileStats":{},"Error":"` + message + `"}`))}, nil
			}))
			p := &Pool{stats: stats.New(), available: []*Process{process}}
			err := p.Scan(processRequest())
			if (err == nil) != (message == "") {
				t.Fatalf("Scan=%v", err)
			}
			if err != nil && err.Error() != message {
				t.Fatalf("worker error changed: %v", err)
			}
			if process.context.Err() != nil || len(p.available) != 1 || len(p.cleaning) != 0 {
				t.Fatal("valid worker envelope entered failure cleanup")
			}
			p.Close()
			if process.context.Err() == nil {
				t.Fatal("Close did not cancel control")
			}
		})
	}
}

func TestCoveragePoolCloseClosesIdleProcess(t *testing.T) {
	process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("idle process must not scan")
	}))
	p := &Pool{stats: stats.New(), available: []*Process{process}}

	p.Close()

	select {
	case <-process.context.Done():
	default:
		t.Fatal("idle process context was not cancelled")
	}
	if !p.closed {
		t.Fatal("pool was not closed")
	}
	if len(p.available) != 0 {
		t.Fatalf("available processes = %d, want 0", len(p.available))
	}
}

func TestCoveragePoolCloseCancelsActiveProcess(t *testing.T) {
	entered := make(chan struct{})
	process := syntheticProcess(coverageRoundTrip(func(request *http.Request) (*http.Response, error) {
		close(entered)
		<-request.Context().Done()
		return nil, request.Context().Err()
	}))
	p := &Pool{stats: stats.New(), available: []*Process{process}}
	scanDone := make(chan error, 1)
	go func() { scanDone <- p.Scan(processRequest()) }()

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("active process did not receive a scan")
	}

	closed := make(chan struct{})
	go func() { p.Close(); close(closed) }()
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("Pool.Close stranded on active process")
	}

	select {
	case err := <-scanDone:
		if err == nil {
			t.Fatal("cancelled Pool.Scan returned nil")
		}
	case <-time.After(time.Second):
		t.Fatal("active Pool.Scan did not terminate after Pool.Close")
	}
	if len(p.available) != 0 {
		t.Fatalf("available processes = %d, want 0 after close", len(p.available))
	}
}

func TestCoveragePoolLateStartReturnIsClosed(t *testing.T) {
	started := make(chan struct{})
	process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("late process must not scan")
	}))
	p := &Pool{
		stats: stats.New(),
		processFactory: func(options *ProcessOptions, _ string) (*Process, error) {
			close(started)
			<-options.parentContext.Done()
			return process, nil
		},
	}
	startedGet := make(chan error, 1)
	go func() {
		_, err := p.get()
		startedGet <- err
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("process factory was not entered")
	}

	closed := make(chan struct{})
	go func() { p.Close(); close(closed) }()
	waitClosed(t, closed)
	select {
	case err := <-startedGet:
		if err == nil || err.Error() != "pool closed" {
			t.Fatalf("late get error = %v, want pool closed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("late process startup did not return")
	}
	select {
	case <-process.context.Done():
	default:
		t.Fatal("late process was not closed")
	}
}

func TestCoverageS1PendingStartupCancellation(t *testing.T) {
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	p := &Pool{processFactory: func(options *ProcessOptions, _ string) (*Process, error) {
		if options.parentContext == nil {
			return nil, errors.New("startup cancellation unavailable")
		}
		close(entered)
		<-options.parentContext.Done()
		close(cancelled)
		return nil, options.parentContext.Err()
	}}
	getDone := make(chan error, 1)
	go func() { _, err := p.get(); getDone <- err }()
	waitClosed(t, entered)
	p.mutex.Lock()
	pending, original := len(p.starting), p.processOptions.parentContext
	p.mutex.Unlock()
	if pending != 1 || original != nil {
		t.Fatalf("pending=%d, shared options mutated=%v", pending, original != nil)
	}
	select {
	case <-getDone:
		t.Fatal("factory returned before cancellation")
	default:
	}
	closed := make(chan struct{})
	go func() { p.Close(); close(closed) }()
	waitClosed(t, cancelled)
	waitClosed(t, closed)
	select {
	case err := <-getDone:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("startup error=%v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("get did not finish cancellation")
	}
	if len(p.starting) != 0 || len(p.owned) != 0 {
		t.Fatal("startup ownership remains")
	}
	if _, err := p.get(); err == nil {
		t.Fatal("closed pool accepted startup")
	}
}

func TestCoverageS1StartupFailureCleanup(t *testing.T) {
	for _, mode := range []string{"factory-error", "process-and-error", "spawn-error"} {
		t.Run(mode, func(t *testing.T) {
			var process *Process
			failure := errors.New("synthetic startup failure")
			p := &Pool{}
			if mode == "spawn-error" {
				p.processOptions = ProcessOptions{executable: filepath.Join(t.TempDir(), "missing-owned-binary"), config: &settings.Config{}}
			} else {
				p.processFactory = func(options *ProcessOptions, _ string) (*Process, error) {
					if mode == "process-and-error" {
						process = syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return nil, failure }))
					}
					return process, failure
				}
			}
			if got, err := p.get(); got != nil || err == nil {
				t.Fatalf("get=%v %v", got, err)
			}
			if len(p.starting) != 0 || len(p.owned) != 0 {
				t.Fatal("failed startup retained ownership")
			}
			if process != nil && process.context.Err() == nil {
				t.Fatal("process accompanying error was leaked")
			}
			closed := make(chan struct{})
			go func() { p.Close(); close(closed) }()
			waitClosed(t, closed)
		})
	}
}

func TestCoverageS1ConcurrentCloseJoinsStartup(t *testing.T) {
	entered, cleanupCompleted := make(chan struct{}), make(chan struct{})
	p := &Pool{processFactory: func(options *ProcessOptions, _ string) (*Process, error) {
		close(entered)
		<-options.parentContext.Done()
		// Completion happens because cancellation arrived, not because the test
		// releases a blocked factory after Close returns.
		close(cleanupCompleted)
		return nil, options.parentContext.Err()
	}}
	getDone := make(chan struct{})
	go func() { _, _ = p.get(); close(getDone) }()
	waitClosed(t, entered)
	results := make(chan bool, 2)
	for i := 0; i < 2; i++ {
		go func() {
			p.Close()
			select {
			case <-cleanupCompleted:
				results <- true
			default:
				results <- false
			}
		}()
	}
	for i := 0; i < 2; i++ {
		select {
		case complete := <-results:
			if !complete {
				t.Fatal("Close returned before startup cleanup")
			}
		case <-time.After(time.Second):
			t.Fatal("concurrent Close stranded")
		}
	}
	waitClosed(t, getDone)
}

func TestCoveragePoolActiveCancellation(t *testing.T) {
	binary := os.Getenv("BEARER_COVERAGE_EXECUTABLE")
	if !filepath.IsAbs(binary) {
		t.Fatal("required owned candidate executable absent")
	}
	binaryBytes, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("actual candidate sha256=%x", sha256.Sum256(binaryBytes))

	cfg := &settings.Config{LogLevel: "error", Worker: settings.WorkerOptions{TimeoutWorkerOnline: 5 * time.Second}}
	cfg.Scan.Scanner = []string{}
	p := New(cfg, stats.New())
	p.processOptions.executable = binary
	p.processOptions.baseArguments = []string{"processing-worker", "--log-level", "error"}
	process, err := p.get()
	if err != nil {
		t.Fatal(err)
	}
	pid := process.command.Process.Pid
	t.Logf("active pool worker pid=%d url=%s", pid, process.baseURL)
	closed := false
	t.Cleanup(func() {
		if !closed {
			select {
			case <-process.exitChannel:
				return
			default:
			}
			t.Log("recovery cleanup required; normal shutdown not proved")
			process.cancelContext()
			_ = process.command.Process.Kill()
			select {
			case <-process.exitChannel:
			case <-time.After(2 * time.Second):
				t.Error("owned cleanup failed; outer supervisor must recover")
			}
		}
	})

	t.Logf("owned identity pid=%d parent=%d executable=%q", pid, os.Getpid(), binary)
	if err := syscall.Kill(pid, 0); err != nil {
		t.Fatalf("owned worker not live before SIGSTOP: %v", err)
	}
	if err := process.command.Process.Signal(syscall.SIGSTOP); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	observed := false
	for time.Now().Before(deadline) {
		var stdout, stderr strings.Builder
		pollContext, cancelPoll := context.WithDeadline(context.Background(), deadline)
		cmd := exec.CommandContext(pollContext, "/bin/ps", "-o", "state=", "-p", strconv.Itoa(pid))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		cancelPoll()
		exitCode := 0
		if err != nil {
			exitCode = -1
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
		}
		t.Logf("stop poll pid=%d exit=%d err=%v stdout=%q stderr=%q", pid, exitCode, err, stdout.String(), stderr.String())
		if err == nil && strings.Contains(stdout.String(), "T") {
			observed = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !observed {
		t.Fatal("active pool worker stop not observed")
	}

	requestSent := make(chan struct{}, 1)
	process.client.Transport = coverageRoundTrip(func(request *http.Request) (*http.Response, error) {
		ctx := httptrace.WithClientTrace(request.Context(), &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) {
			select {
			case requestSent <- struct{}{}:
			default:
			}
		}})
		return http.DefaultTransport.RoundTrip(request.WithContext(ctx))
	})
	p.release(process)
	scanDone := make(chan error, 1)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.txt"), []byte("disposable synthetic text\n"), 0600); err != nil {
		t.Fatal(err)
	}
	scope, err := inputfile.OpenInputScope(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer scope.Close()
	bound, err := scope.Acquire("input.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		scanDone <- p.Scan(work.ProcessRequest{Repository: work.Repository{Dir: dir}, File: files.File{FilePath: "input.txt", Timeout: 2 * time.Second}, InputExpectation: bound.Expectation()})
	}()
	select {
	case <-requestSent:
	case <-time.After(time.Second):
		t.Fatal("active pool scan dispatch not observed")
	}

	start := time.Now()
	poolClosed := make(chan struct{})
	go func() { p.Close(); close(poolClosed) }()
	select {
	case <-poolClosed:
	case <-time.After(10 * time.Second):
		t.Fatal("active pool Close exceeded shutdown bound")
	}
	closed = true
	elapsed := time.Since(start)
	if elapsed < settings.TimeoutWorkerShutdown {
		t.Fatal("active pool kill fallback interval not observed")
	}
	t.Logf("active pool Close completed in %s", elapsed)

	select {
	case err := <-scanDone:
		if err == nil {
			t.Fatal("cancelled active pool Scan returned success")
		}
		t.Logf("Scan termination observed: %v", err)
	case <-time.After(time.Second):
		t.Fatal("active pool Scan stranded after Close")
	}
	select {
	case <-process.exitChannel:
		t.Log("exitChannel observed")
	default:
		t.Fatal("exitChannel did not close")
	}
	if err := syscall.Kill(pid, 0); err != syscall.ESRCH {
		t.Fatalf("active pool worker PID still present: %v", err)
	}
	t.Logf("PID %d absent after Close", pid)
	conn, err := net.DialTimeout("tcp", strings.TrimPrefix(process.baseURL, "http://"), 200*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Fatal("active pool worker listener still accepts connections")
	}
	t.Logf("listener %s absent after Close: %v", process.baseURL, err)
}

func TestCoveragePoolLateReleaseAndConcurrentClose(t *testing.T) {
	process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return nil, errors.New("unused") }))
	p := &Pool{stats: stats.New(), available: []*Process{process}}
	p.Close()
	if p.release(process) {
		t.Fatal("late successful response accepted after closure")
	}
	if len(p.available) != 0 {
		t.Fatal("late response reentered closed pool")
	}
	// An ongoing Close remains a join point for another caller.
	first := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return nil, errors.New("unused") }))
	first.exitChannel = make(chan struct{})
	p = &Pool{stats: stats.New(), available: []*Process{first}}
	doneA, doneB := make(chan struct{}), make(chan struct{})
	go func() { p.Close(); close(doneA) }()
	select {
	case <-first.context.Done():
	case <-time.After(time.Second):
		t.Fatal("Close not entered")
	}
	go func() { p.Close(); close(doneB) }()
	select {
	case <-doneB:
		t.Fatal("concurrent Close returned before process exit")
	case <-time.After(10 * time.Millisecond):
	}
	close(first.exitChannel)
	for _, done := range []chan struct{}{doneA, doneB} {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Close did not join")
		}
	}
}

func TestCoverageTimeoutClassificationExactJoinedMarker(t *testing.T) {
	for _, tc := range []struct {
		name, message string
		killed        bool
		want          stats.FailedReason
	}{
		{"direct", worker.ErrorTimeoutReached.Error(), false, stats.TimeoutReason},
		{"joined", errors.Join(errors.New("processing interrupted"), worker.ErrorTimeoutReached).Error(), false, stats.TimeoutReason},
		{"joined-killed", errors.Join(worker.ErrorTimeoutReached, errors.New("processing interrupted")).Error(), true, stats.KilledAfterTimeoutReason},
		{"substring", "prefix " + worker.ErrorTimeoutReached.Error(), false, stats.UnexpectedReason},
		{"suffix", worker.ErrorTimeoutReached.Error() + " suffix", true, stats.UnexpectedReason},
		{"oom", ErrorOutOfMemory.Error(), false, stats.MemoryLimitReason},
		{"other", "synthetic failure", false, stats.UnexpectedReason},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := translateErrorForStats(tc.message, tc.killed); got != tc.want {
				t.Fatalf("classification %v want %v", got, tc.want)
			}
		})
	}
}
