package pool

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/stats"
)

func TestCoverageS1PreCancelledStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p, err := newProcess(&ProcessOptions{parentContext: ctx, executable: filepath.Join(t.TempDir(), "must-not-start"), config: &settings.Config{}}, "pre-cancelled")
	if p != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("pre-cancelled start = %v %v", p, err)
	}
	// The race between construction and Start also has a closed completion
	// channel even though command.Start never ran.
	exit := make(chan struct{})
	process := &Process{context: ctx, cancelContext: cancel, exitChannel: exit}
	if err := process.start(&settings.Config{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("start=%v", err)
	}
	closed := make(chan struct{})
	go func() { process.Close(); close(closed) }()
	waitClosed(t, closed)
}

func TestCoverageS1InitializationCancellation(t *testing.T) {
	entered, terminated := make(chan struct{}), make(chan struct{})
	p := syntheticProcess(coverageRoundTrip(func(request *http.Request) (*http.Response, error) {
		close(entered)
		<-request.Context().Done()
		close(terminated)
		return nil, request.Context().Err()
	}))
	p.exitChannel = make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- p.initialize(&settings.Config{Worker: settings.WorkerOptions{TimeoutWorkerOnline: time.Minute}})
	}()
	waitClosed(t, entered)
	p.cancelContext()
	waitClosed(t, terminated)
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("initialize error=%v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("initialize ignored cancellation")
	}
	close(p.exitChannel)
	p.Close()
}

type s1HelperRecord struct {
	PID        int    `json:"pid"`
	Parent     int    `json:"parent"`
	Port       string `json:"port"`
	Nonce      string `json:"nonce"`
	Mode       string `json:"mode"`
	Executable string `json:"executable"`
	SHA256     string `json:"sha256"`
}

// TestS1StartupHelper is invoked only as a disposable test subprocess. It is
// not a scanner/engine fixture and never reads source or acquires rules.
func TestS1StartupHelper(t *testing.T) {
	values := make(map[string]string)
	for i, arg := range os.Args {
		if arg == "--" {
			for j := i + 1; j+1 < len(os.Args); j += 2 {
				values[os.Args[j]] = os.Args[j+1]
			}
			break
		}
	}
	if values["--s1-nonce"] == "" {
		return
	}
	dir, mode, port := values["--s1-dir"], values["--s1-mode"], values["--port"]
	if !filepath.IsAbs(dir) || port == "" {
		t.Fatal("invalid owned helper parameters")
	}
	if mode == "blocked-initialize-kill" {
		signal.Ignore(os.Interrupt)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	record := s1HelperRecord{PID: os.Getpid(), Parent: os.Getppid(), Port: port, Nonce: values["--s1-nonce"], Mode: mode, Executable: executable, SHA256: fmt.Sprintf("%x", sha256.Sum256(bytes))}
	writeRecord := func(name string) {
		data, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name+".tmp"), data, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(filepath.Join(dir, name+".tmp"), filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	if mode == "exit-before-ready" {
		writeRecord("identity.json")
		os.Exit(7)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatal(err)
	}
	writeRecord("identity.json")
	mux := http.NewServeMux()
	mux.HandleFunc(work.RouteInitialize, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
		writeRecord("entered.json")
		if mode == "blocked-initialize-kill" {
			<-r.Context().Done()
			writeRecord("cancelled.json")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if mode == "initialize-error" {
			_, _ = io.WriteString(w, `{"Error":"synthetic initialization rejected"}`)
			return
		}
		if mode != "valid-initialize" {
			t.Error("unknown helper mode")
			return
		}
		_, _ = io.WriteString(w, `{"Error":""}`)
	})
	if err := http.Serve(listener, mux); err != nil {
		t.Fatal(err)
	}
}

func s1ReadRecord(t *testing.T, path, nonce string) s1HelperRecord {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil {
			var record s1HelperRecord
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatal(err)
			}
			if record.Nonce != nonce || record.Parent != os.Getpid() {
				t.Fatalf("unowned helper record: %+v", record)
			}
			return record
		}
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("helper observation unavailable: %s", filepath.Base(path))
	return s1HelperRecord{}
}

func TestCoverageS1ActualStartupLifecycle(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(bytes))
	for _, mode := range []string{"blocked-initialize-kill", "valid-initialize", "initialize-error", "exit-before-ready"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			nonce := fmt.Sprintf("s1-%d-%d", os.Getpid(), time.Now().UnixNano())
			cfg := &settings.Config{LogLevel: "error", Worker: settings.WorkerOptions{TimeoutWorkerOnline: 20 * time.Second}}
			cfg.Scan.Scanner = []string{}
			p := &Pool{processOptions: ProcessOptions{executable: executable, baseArguments: []string{"-test.run=^TestS1StartupHelper$", "--", "--s1-dir", dir, "--s1-mode", mode, "--s1-nonce", nonce}, config: cfg}}
			type result struct {
				process *Process
				err     error
			}
			returned := make(chan result, 1)
			go func() { process, err := p.get(); returned <- result{process, err} }()
			record := s1ReadRecord(t, filepath.Join(dir, "identity.json"), nonce)
			if record.Executable != executable || record.SHA256 != digest || record.Mode != mode {
				t.Fatalf("helper identity mismatch: %+v", record)
			}
			t.Logf("owned helper nonce=%s pid=%d parent=%d executable=%q sha256=%s port=%s", nonce, record.PID, record.Parent, executable, digest, record.Port)
			// The outer supervisor also owns every descendant, including failures
			// before this record is published. Local recovery never counts as proof.
			t.Cleanup(func() {
				if syscall.Kill(record.PID, 0) == nil {
					t.Error("supervisor/local recovery required; normal cleanup not proved")
					_ = syscall.Kill(record.PID, syscall.SIGKILL)
				}
			})
			var got result
			if mode == "blocked-initialize-kill" {
				_ = s1ReadRecord(t, filepath.Join(dir, "entered.json"), nonce)
				select {
				case <-returned:
					t.Fatal("factory returned before Close")
				default:
				}
				if err := syscall.Kill(record.PID, 0); err != nil {
					t.Fatalf("helper not live: %v", err)
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				cmd := exec.CommandContext(ctx, "/bin/ps", "-o", "pid=,ppid=,lstart=,stat=,comm=", "-p", strconv.Itoa(record.PID))
				var stdout, stderr strings.Builder
				cmd.Stdout, cmd.Stderr = &stdout, &stderr
				err := cmd.Run()
				cancel()
				code := 0
				if err != nil {
					code = -1
					if ee, ok := err.(*exec.ExitError); ok {
						code = ee.ExitCode()
					}
				}
				t.Logf("identity observation exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
				if err != nil || !strings.Contains(stdout.String(), strconv.Itoa(record.PID)) {
					t.Fatal("owned identity observation failed")
				}
				p.mutex.Lock()
				pending := len(p.starting)
				p.mutex.Unlock()
				if pending != 1 {
					t.Fatalf("pending startup ownership=%d", pending)
				}
				start := time.Now()
				closed := make(chan struct{})
				go func() { p.Close(); close(closed) }()
				select {
				case <-closed:
				case <-time.After(8 * time.Second):
					t.Fatal("startup Close exhausted observation bound; 2s reserved for cleanup")
				}
				elapsed := time.Since(start)
				if elapsed < settings.TimeoutWorkerShutdown {
					t.Fatal("existing five-second fallback not exercised")
				}
				t.Logf("Pool.Close completed after %s without readiness release", elapsed)
				_ = s1ReadRecord(t, filepath.Join(dir, "cancelled.json"), nonce)
				select {
				case got = <-returned:
				case <-time.After(time.Second):
					t.Fatal("factory/get stranded after Close")
				}
				if got.process != nil || !errors.Is(got.err, context.Canceled) {
					t.Fatalf("cancelled startup=%v %v", got.process, got.err)
				}
			} else {
				select {
				case got = <-returned:
				case <-time.After(3 * time.Second):
					t.Fatal("startup did not return")
				}
				if mode == "valid-initialize" {
					if got.process == nil || got.err != nil {
						t.Fatalf("valid startup=%v %v", got.process, got.err)
					}
					_ = s1ReadRecord(t, filepath.Join(dir, "entered.json"), nonce)
				} else if got.process != nil || got.err == nil {
					t.Fatalf("failed startup accepted=%v %v", got.process, got.err)
				}
				closed := make(chan struct{})
				go func() { p.Close(); close(closed) }()
				select {
				case <-closed:
				case <-time.After(8 * time.Second):
					t.Fatal("Close exhausted observation bound; 2s reserved for cleanup")
				}
				if got.process != nil {
					select {
					case <-got.process.exitChannel:
						t.Log("Wait/exitChannel complete")
					default:
						t.Fatal("exitChannel open after Close")
					}
				}
			}
			t.Logf("factory/get completed with %v", got.err)
			if err := syscall.Kill(record.PID, 0); err != syscall.ESRCH {
				t.Fatalf("owned PID still present: %v", err)
			}
			conn, err := net.DialTimeout("tcp", "127.0.0.1:"+record.Port, 100*time.Millisecond)
			if err == nil {
				conn.Close()
				t.Fatal("owned listener survived")
			}
			p.mutex.Lock()
			pending, owned, available := len(p.starting), len(p.owned), len(p.available)
			p.mutex.Unlock()
			if pending != 0 || owned != 0 || available != 0 {
				t.Fatalf("ownership remains: %d/%d/%d", pending, owned, available)
			}
			t.Logf("PID %d absent; listener %s absent; startup/owned/available empty", record.PID, record.Port)
		})
	}
}

type coverageBody struct {
	reader *strings.Reader
	closed chan struct{}
}

func (b *coverageBody) Read(p []byte) (int, error) { return b.reader.Read(p) }
func (b *coverageBody) Close() error {
	select {
	case <-b.closed:
	default:
		close(b.closed)
	}
	return nil
}

type coverageRoundTrip func(*http.Request) (*http.Response, error)

func (f coverageRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func syntheticProcess(rt http.RoundTripper) *Process {
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan struct{})
	close(exit)
	return &Process{
		context:       ctx,
		cancelContext: cancel,
		errorChannel:  make(chan error, 1),
		exitChannel:   exit,
		client:        &http.Client{Transport: rt},
		baseURL:       "http://synthetic.invalid",
	}
}

func processRequest() work.ProcessRequest {
	return work.ProcessRequest{File: files.File{FilePath: "input.go", Timeout: time.Second}}
}

func responseBody(body string) (*coverageBody, *http.Response) {
	b := &coverageBody{reader: strings.NewReader(body), closed: make(chan struct{})}
	return b, &http.Response{StatusCode: http.StatusOK, Body: b, Header: make(http.Header)}
}

func waitClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("response body was not closed")
	}
}

type failingBody struct {
	reader   *strings.Reader
	readErr  error
	closeErr error
	closed   chan struct{}
}

func (b *failingBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	if b.reader != nil {
		return b.reader.Read(p)
	}
	return 0, io.EOF
}
func (b *failingBody) Close() error {
	select {
	case <-b.closed:
	default:
		close(b.closed)
	}
	return b.closeErr
}

func TestCoverageProcessScanAcceptsValidNullStatsAndWorkerError(t *testing.T) {
	for _, tc := range []struct {
		name      string
		body      string
		workerErr string
	}{
		{name: "null stats", body: `{"FileStats":null,"Error":""}`},
		{name: "populated stats", body: `{"FileStats":{"Rules":{}},"Error":""}`},
		{name: "worker error", body: `{"FileStats":null,"Error":"detector failed"}`, workerErr: "detector failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, response := responseBody(tc.body)
			p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return response, nil }))
			got, err := p.Scan(processRequest())
			if err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if got == nil || got.Error != tc.workerErr {
				t.Fatalf("response = %#v, want worker error %q", got, tc.workerErr)
			}
			waitClosed(t, body.closed)
		})
	}
}

func TestCoverageProcessScanRejectsMalformedWorkerResponses(t *testing.T) {
	cases := []string{"", "{", "null", `{"FileStats":null}`, `{"Error":""}`, `{"FileStats":null,"Error":"","Extra":1}`, `{"FileStats":null,"Error":""}{}`, `{"FileStats":null,"Error":"","Error":"second"}`}
	for _, raw := range cases {
		t.Run(fmt.Sprintf("%q", raw), func(t *testing.T) {
			body, response := responseBody(raw)
			p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return response, nil }))
			if got, err := p.Scan(processRequest()); err == nil || got != nil {
				t.Fatalf("Scan() = %#v, %v; want strict decode failure", got, err)
			}
			waitClosed(t, body.closed)
		})
	}
}

func TestCoverageProcessScanRejectsTransportAndStatusFailures(t *testing.T) {
	p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return nil, errors.New("PRIVATE_RESPONSE_BODY") }))
	if got, err := p.Scan(processRequest()); got != nil || err == nil || strings.Contains(err.Error(), "PRIVATE_RESPONSE_BODY") {
		t.Fatalf("transport returned %v %v", got, err)
	}
	for _, status := range []int{http.StatusBadRequest, http.StatusBadGateway, http.StatusNoContent} {
		body, response := responseBody("PRIVATE_RESPONSE_BODY")
		response.StatusCode = status
		p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return response, nil }))
		if got, err := p.Scan(processRequest()); got != nil || err == nil || strings.Contains(err.Error(), "PRIVATE_RESPONSE_BODY") {
			t.Fatalf("status returned %v %v", got, err)
		}
		waitClosed(t, body.closed)
	}
}

func TestCoverageProcessScanRejectsBodyReadAndCloseFailures(t *testing.T) {
	cases := []struct {
		name string
		body *failingBody
	}{
		{name: "read", body: &failingBody{readErr: errors.New("synthetic body read failure"), closed: make(chan struct{})}},
		{name: "close", body: &failingBody{reader: strings.NewReader(`{"FileStats":null,"Error":""}`), closeErr: errors.New("synthetic body close failure"), closed: make(chan struct{})}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := tc.body
			p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
			}))
			if got, err := p.Scan(processRequest()); err == nil || got != nil {
				t.Fatalf("Scan() = %#v, %v; want body failure", got, err)
			}
			waitClosed(t, body.closed)
		})
	}
}

func TestCoverageProcessScanRejectsInvalidRequestURL(t *testing.T) {
	called := make(chan struct{}, 1)
	p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		called <- struct{}{}
		return nil, errors.New("unexpected transport invocation")
	}))
	p.baseURL = "://invalid"
	if got, err := p.Scan(processRequest()); err == nil || got != nil {
		t.Fatalf("Scan() = %#v, %v; want request construction failure", got, err)
	}
	select {
	case <-called:
		t.Fatal("transport invoked for invalid request")
	default:
	}
}

func TestCoveragePoolScanPropagatesActualProcessResponseError(t *testing.T) {
	for _, raw := range []string{`{"FileStats":{"Rules":{}},"Error":"detector failed"}`, `{"FileStats":null,"Error":`, "PRIVATE_RESPONSE_BODY"} {
		body, response := responseBody(raw)
		process := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) { return response, nil }))
		p := &Pool{stats: stats.New(), available: []*Process{process}}
		if err := p.Scan(processRequest()); err == nil || strings.Contains(err.Error(), "PRIVATE_RESPONSE_BODY") {
			t.Fatalf("Pool.Scan() error = %v; expected redacted failure", err)
		}
		waitClosed(t, body.closed)
	}
}

func TestCoverageProcessScanPreservesWorkerError(t *testing.T) {
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)
	p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		<-release
		return nil, errors.New("late transport")
	}))
	want := errors.New("synthetic worker failure")
	p.errorChannel <- want
	done := make(chan error, 1)
	go func() { _, err := p.Scan(processRequest()); done <- err }()
	select {
	case err := <-done:
		if !errors.Is(err, want) {
			t.Fatalf("error = %v, want worker sentinel", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Scan did not return worker error")
	}
	unblock()
}

func TestCoverageProcessScanClosesBodyAfterLateTimeout(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)
	var body *coverageBody
	p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		body, _ = responseBody(`{"FileStats":null,"Error":""}`)
		close(entered)
		<-release
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
	}))
	request := processRequest()
	request.File.Timeout = -settings.TimeoutWorkerFileGrace + 100*time.Millisecond
	done := make(chan error, 1)
	go func() { _, err := p.Scan(request); done <- err }()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("RoundTripper was not entered")
	}
	select {
	case err := <-done:
		if !errors.Is(err, worker.ErrorTimeoutReached) {
			t.Fatalf("error = %v, want timeout", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Scan did not time out")
	}
	unblock()
	waitClosed(t, body.closed)
}

func TestCoverageProcessScanCancellationClosesLateBody(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)
	var body *coverageBody
	p := syntheticProcess(coverageRoundTrip(func(*http.Request) (*http.Response, error) {
		body, _ = responseBody(`{"FileStats":null,"Error":""}`)
		close(entered)
		<-release
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
	}))
	done := make(chan error, 1)
	go func() { _, err := p.Scan(processRequest()); done <- err }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("RoundTripper was not entered")
	}
	p.cancelContext()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancellation returned nil error")
		}
	case <-time.After(time.Second):
		t.Fatal("Scan did not return after cancellation")
	}
	unblock()
	waitClosed(t, body.closed)
}

func TestCoverageActualWorkerShutdown(t *testing.T) {
	binary := os.Getenv("BEARER_COVERAGE_EXECUTABLE")
	if !filepath.IsAbs(binary) {
		t.Fatal("required owned candidate executable absent")
	}
	binaryBytes, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("actual candidate sha256=%x", sha256.Sum256(binaryBytes))
	for _, stopped := range []bool{false, true} {
		name := "cooperative"
		if stopped {
			name = "stopped-kill-fallback"
		}
		t.Run(name, func(t *testing.T) {
			cfg := &settings.Config{LogLevel: "error", Worker: settings.WorkerOptions{TimeoutWorkerOnline: 5 * time.Second}}
			cfg.Scan.Scanner = []string{}
			proc, err := newProcess(&ProcessOptions{executable: binary, baseArguments: []string{"processing-worker", "--log-level", "error"}, config: cfg}, "coverage-owned")
			if err != nil {
				t.Fatal(err)
			}
			pid := proc.command.Process.Pid
			t.Logf("owned worker pid=%d url=%s", pid, proc.baseURL)
			t.Cleanup(func() {
				select {
				case <-proc.exitChannel:
					return
				default:
				}
				proc.cancelContext()
				_ = proc.command.Process.Kill()
				select {
				case <-proc.exitChannel:
				case <-time.After(2 * time.Second):
					t.Error("owned cleanup failed; outer supervisor must recover")
				}
			})
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
			req := work.ProcessRequest{Repository: work.Repository{Dir: dir}, File: files.File{FilePath: "input.txt", Timeout: 2 * time.Second}, ReportPath: filepath.Join(dir, "report.jsonl"), InputExpectation: bound.Expectation()}
			if !stopped {
				response, err := proc.Scan(req)
				if err != nil || response == nil || response.Error != "" {
					t.Fatalf("actual scan response=%+v error=%v", response, err)
				}
				t.Log("actual Scan returned a valid empty-scanner response")
			} else {
				t.Logf("owned identity pid=%d parent=%d executable=%q", pid, os.Getpid(), binary)
				if err := syscall.Kill(pid, 0); err != nil {
					t.Fatalf("owned worker not live before SIGSTOP: %v", err)
				}
				if err := proc.command.Process.Signal(syscall.SIGSTOP); err != nil {
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
					t.Fatal("owned worker stop not observed")
				}
			}
			var scanDone chan error
			if stopped {
				sent := make(chan struct{}, 1)
				proc.client.Transport = coverageRoundTrip(func(r *http.Request) (*http.Response, error) {
					ctx := httptrace.WithClientTrace(r.Context(), &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) {
						select {
						case sent <- struct{}{}:
						default:
						}
					}})
					return http.DefaultTransport.RoundTrip(r.WithContext(ctx))
				})
				scanDone = make(chan error, 1)
				go func() { _, err := proc.Scan(req); scanDone <- err }()
				select {
				case <-sent:
				case <-time.After(time.Second):
					t.Fatal("synthetic request dispatch not observed")
				}
			}
			start := time.Now()
			closed := make(chan struct{})
			go func() { proc.Close(); close(closed) }()
			select {
			case <-closed:
			case <-time.After(10 * time.Second):
				t.Fatal("normal Close exceeded shutdown bound")
			}
			elapsed := time.Since(start)
			if stopped && elapsed < settings.TimeoutWorkerShutdown {
				t.Fatal("kill fallback interval not observed")
			}
			if stopped {
				select {
				case err := <-scanDone:
					if err == nil {
						t.Fatal("cancelled Scan returned success")
					}
					t.Logf("Scan terminated: %v", err)
				case <-time.After(time.Second):
					t.Fatal("Scan stranded after Close")
				}
			}
			select {
			case <-proc.exitChannel:
				t.Log("exitChannel observed")
			default:
				t.Fatal("Wait/exitChannel not complete")
			}
			if proc.command.ProcessState == nil || !proc.command.ProcessState.Exited() && !stopped {
				t.Fatal("process Wait state missing")
			}
			if err := syscall.Kill(pid, 0); err != syscall.ESRCH {
				t.Fatalf("owned PID still present: %v", err)
			}
			t.Logf("PID %d absent after Close", pid)
			conn, err := net.DialTimeout("tcp", strings.TrimPrefix(proc.baseURL, "http://"), 200*time.Millisecond)
			if err == nil {
				conn.Close()
				t.Fatal("owned listener still accepts connections")
			}
			t.Logf("listener %s absent after Close: %v", proc.baseURL, err)
			t.Logf("normal Close completed in %s; Wait complete; PID and listener absent; no recovery cleanup used", elapsed)
		})
	}
}
