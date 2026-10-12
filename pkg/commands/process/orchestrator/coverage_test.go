package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bearer/bearer/pkg/commands/process/filelist"
	detectorapi "github.com/bearer/bearer/pkg/detectors"
	reportapi "github.com/bearer/bearer/pkg/report"
	reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/writer"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/hhatto/gocloc"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/coverage"
)

type fakeWorker struct {
	scan func(work.ProcessRequest) error
}

func (f *fakeWorker) Scan(r work.ProcessRequest) error { return f.scan(r) }
func (f *fakeWorker) Close()                           {}
func syntheticOrchestrator(t *testing.T, fn func(work.ProcessRequest) error) *Orchestrator {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("synthetic"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg := &settings.Config{}
	cfg.Scan.Quiet = true
	return &Orchestrator{repository: work.Repository{Dir: dir}, config: cfg, maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{}), pool: &fakeWorker{fn}}
}
func TestCoverageWorkerTerminal(t *testing.T) {
	for _, c := range []struct {
		name    string
		body    string
		failure error
		success bool
	}{
		{"clean", "", nil, true}, {"malformed", "{", nil, false}, {"null", "null", nil, false}, {"error-diagnostic", `{"type":"error","message":"synthetic"}`, nil, false}, {"timeout", "", errors.New("timeout"), false}, {"processing-fault", "", errors.New("synthetic processing fault"), false},
	} {
		t.Run(c.name, func(t *testing.T) {
			o := syntheticOrchestrator(t, func(r work.ProcessRequest) error {
				if c.failure != nil {
					return c.failure
				}
				return os.WriteFile(r.ReportPath, []byte(c.body), 0600)
			})
			var report bytes.Buffer
			out := o.scanFile(&report, files.File{FilePath: "input.go"})
			if (out.Status == "success") != c.success {
				t.Fatalf("outcome=%+v", out)
			}
		})
	}
}
func TestCoverageInputMutationAndCancellation(t *testing.T) {
	o := syntheticOrchestrator(t, func(r work.ProcessRequest) error {
		return os.WriteFile(filepath.Join(r.Repository.Dir, r.File.FilePath), []byte("changed"), 0600)
	})
	if out := o.scanFile(io.Discard, files.File{FilePath: "input.go"}); out.Status == "success" {
		t.Fatal("mutation accepted")
	}
	o.Close()
	if out := o.scanFile(io.Discard, files.File{FilePath: "input.go"}); out.Status == "success" {
		t.Fatal("cancelled processing accepted")
	}
	if _, err := o.waitForScan(make(chan coverage.Outcome), 1); err == nil {
		t.Fatal("cancelled wait accepted")
	}
	ch := make(chan coverage.Outcome)
	close(ch)
	o = syntheticOrchestrator(t, func(work.ProcessRequest) error { return nil })
	if _, err := o.waitForScan(ch, 1); err == nil {
		t.Fatal("missing terminal accepted")
	}
}
func TestCoverageScanCompletion(t *testing.T) {
	o := syntheticOrchestrator(t, func(r work.ProcessRequest) error { return os.WriteFile(r.ReportPath, nil, 0600) })
	destination := filepath.Join(t.TempDir(), "report.jsonl")
	if err := o.Scan(destination, []files.File{{FilePath: "input.go"}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"type":"processing_coverage"`) {
		t.Fatal("terminal record absent")
	}
	if err := o.Scan(destination, []files.File{{FilePath: "input.go"}, {FilePath: "input.go"}}); err == nil {
		t.Fatal("duplicate input accepted")
	}
	o.Close()
	if err := o.Scan(destination, []files.File{{FilePath: "input.go"}}); err == nil {
		t.Fatal("cancelled dispatch accepted")
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("synthetic write failure") }
func TestCoverageWorkerReportFailure(t *testing.T) {
	o := syntheticOrchestrator(t, func(r work.ProcessRequest) error {
		return os.WriteFile(r.ReportPath, []byte(`{"type":"custom"}`), 0600)
	})
	if out := o.scanFile(failedWriter{}, files.File{FilePath: "input.go"}); out.Status == "success" {
		t.Fatal("result write error accepted")
	}
	o.pool = &fakeWorker{func(work.ProcessRequest) error { return errors.New("timeout") }}
	if out := o.scanFile(failedWriter{}, files.File{FilePath: "input.go"}); out.Status == "success" || !strings.Contains(out.Error, "write failure") {
		t.Fatalf("error write lost: %+v", out)
	}
}

func TestCoverageCancellationDuringWorker(t *testing.T) {
	o := syntheticOrchestrator(t, func(work.ProcessRequest) error { return nil })
	o.pool = &fakeWorker{func(r work.ProcessRequest) error { o.Close(); return os.WriteFile(r.ReportPath, nil, 0600) }}
	if err := o.Scan(filepath.Join(t.TempDir(), "report.jsonl"), []files.File{{FilePath: "input.go"}}); err == nil {
		t.Fatal("cancelled worker accepted")
	}
}

type coverageErrorDetector struct{}

func (coverageErrorDetector) AcceptDir(*file.Path) (bool, error) { return true, nil }
func (coverageErrorDetector) ProcessFile(*file.FileInfo, *file.Path, reportapi.Report) (bool, error) {
	return false, errors.New("synthetic terminal failure")
}
func TestCoverageActualDetectorDiagnosticRejected(t *testing.T) {
	o := syntheticOrchestrator(t, func(r work.ProcessRequest) error {
		f, err := os.Create(r.ReportPath)
		if err != nil {
			return err
		}
		report := &writer.Detectors{File: f}
		err = detectorapi.ExtractWithDetectors(context.Background(), r.Dir, r.File.FilePath, report, nil, []detectorapi.InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: coverageErrorDetector{}}}, nil, false, true)
		closeErr := f.Close()
		if err != nil {
			return err
		}
		return closeErr
	})
	out := o.scanFile(io.Discard, files.File{FilePath: "input.go"})
	if out.Status == "success" || !strings.Contains(out.Error, "error or unexpected terminal diagnostic") {
		t.Fatalf("actual worker report rejection boundary not reached: %+v", out)
	}
}

type coverageClosingWorker struct {
	entered chan struct{}
	cancel  chan struct{}
}

func (w *coverageClosingWorker) Scan(work.ProcessRequest) error {
	close(w.entered)
	<-w.cancel
	return errors.New("synthetic cancelled worker")
}
func (w *coverageClosingWorker) Close() { close(w.cancel) }
func TestCoverageCloseCancelsBeforeJoining(t *testing.T) {
	o := syntheticOrchestrator(t, nil)
	w := &coverageClosingWorker{make(chan struct{}), make(chan struct{})}
	o.pool = w
	done := make(chan error, 1)
	go func() { done <- o.Scan(filepath.Join(t.TempDir(), "report"), []files.File{{FilePath: "input.go"}}) }()
	select {
	case <-w.entered:
	case <-time.After(time.Second):
		t.Fatal("worker not entered")
	}
	o.Close()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled scan accepted")
		}
	case <-time.After(time.Second):
		t.Fatal("Scan joined without cancelling worker")
	}
}
func TestCoverageSingleFileAndDirectoryTargets(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "input.go")
	body := []byte("package synthetic\n")
	if err := os.WriteFile(filename, body, 0600); err != nil {
		t.Fatal(err)
	}
	cfg := settings.Config{IgnoreGit: true, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 100000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: time.Second}}
	cfg.Scan.Quiet = true
	cfg.Scan.SkipGitIgnore = true
	list, err := filelist.Discover(nil, filename, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Files) != 1 || list.Files[0].FilePath != "." {
		t.Fatalf("single-file discovery=%+v", list)
	}
	worker := &fakeWorker{scan: func(request work.ProcessRequest) error {
		f, err := os.Create(request.ReportPath)
		if err != nil {
			return err
		}
		return errors.Join(detectorapi.ExtractWithDetectors(context.Background(), request.Dir, request.File.FilePath, &writer.Detectors{File: f}, nil, nil, nil, false, true), f.Close())
	}}
	o := &Orchestrator{repository: work.Repository{Dir: filename}, config: &cfg, maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{}), pool: worker}
	report := filepath.Join(t.TempDir(), "single.jsonl")
	if err := o.Scan(report, list.Files); err != nil {
		t.Fatal(err)
	}
	if err := reviewTerminal(t, report).ValidateInputs(map[string]string{".": coverage.Digest(body)}); err != nil {
		t.Fatal(err)
	}
	o.repository.Dir = dir
	if err := o.Scan(filepath.Join(t.TempDir(), "directory.jsonl"), []files.File{{FilePath: "."}}); err == nil {
		t.Fatal("directory represented as processed file")
	}
}

// The immutable independent probes retain their original schedule-size
// preconditions. This author control exercises the authorized earlier exclusion
// decision without changing their assertions or the caller's fixed promise.
func TestCoverageResidualEligibilityPromise(t *testing.T) {
	for _, tc := range []struct {
		name, path, body    string
		skipTest, gitIgnore bool
	}{
		{"test", "input_test.go", "package fixture\n", true, false},
		{"builtin", "excluded.log", "synthetic text\n", false, false},
		{"binary", "binary.go", "binary\x00payload", false, false},
		{"gitignore", "ignored.go", "package fixture\n", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			keep := []byte("package fixture\n")
			fixed := map[string]string{"input.go": coverage.Digest(keep), tc.path: coverage.Digest([]byte(tc.body))}
			for name, b := range map[string][]byte{"input.go": keep, tc.path: []byte(tc.body)} {
				if err := os.WriteFile(filepath.Join(root, name), b, 0600); err != nil {
					t.Fatal(err)
				}
			}
			cfg := settings.Config{IgnoreGit: true, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 10000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: time.Second}}
			cfg.Scan.Quiet = true
			cfg.Scan.SkipTest = tc.skipTest
			cfg.Scan.SkipGitIgnore = !tc.gitIgnore
			if tc.gitIgnore {
				if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(tc.path+"\n"), 0600); err != nil {
					t.Fatal(err)
				}
				cfg.Scan.SkipPath = []string{".gitignore"}
			}
			list, err := filelist.Discover(nil, root, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, cfg)
			if err != nil {
				t.Fatal(err)
			}
			if len(list.Files) != 1 || list.Files[0].FilePath != "input.go" {
				t.Fatalf("compatible scheduling failed: %+v", list)
			}
			excluded := false
			for _, x := range list.Coverage.Excluded {
				if x.Path == tc.path && x.Reason != "" {
					excluded = true
				}
			}
			if !excluded {
				t.Fatal("exclusion lost")
			}
			calls := 0
			newOrchestrator := func() *Orchestrator {
				worker := &fakeWorker{scan: func(req work.ProcessRequest) error {
					report, err := os.Create(req.ReportPath)
					if err != nil {
						return err
					}
					scanErr := detectorapi.ExtractWithDetectors(context.Background(), req.Dir, req.File.FilePath, &writer.Detectors{File: report}, nil, []detectorapi.InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: &reviewDetector{calls: &calls}}}, nil, tc.skipTest, !tc.gitIgnore)
					return errors.Join(scanErr, report.Close())
				}}
				return &Orchestrator{repository: work.Repository{Dir: root}, config: &cfg, maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{}), pool: worker}
			}
			output := filepath.Join(t.TempDir(), "eligible.jsonl")
			if err := newOrchestrator().Scan(output, list.Files); err != nil {
				t.Fatalf("ordinary findings path failed: %v", err)
			}
			terminal := reviewTerminal(t, output)
			if calls != 1 || len(terminal.Outcomes) != 1 || terminal.Outcomes[0].Path != "input.go" || terminal.Outcomes[0].Status != "success" {
				t.Fatalf("visitation/outcomes: %d %+v", calls, terminal)
			}
			if err := terminal.ValidateInputs(fixed); err == nil {
				t.Fatal("excluded promise was incorrectly certified")
			}
			// An input already scheduled independently may not become success if later
			// processing excludes it. Keep the fixed promise and inject that schedule.
			output = filepath.Join(t.TempDir(), "fixed.jsonl")
			calls = 0
			err = newOrchestrator().Scan(output, []files.File{{FilePath: "input.go", Timeout: time.Second}, {FilePath: tc.path, Timeout: time.Second}})
			if err == nil || calls != 1 {
				t.Fatalf("unvisited scheduled input accepted: calls=%d err=%v", calls, err)
			}
			raw, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(raw, []byte(`"type":"processing_coverage"`)) {
				t.Fatal("failed fixed schedule emitted complete terminal")
			}
			if len(fixed) != 2 || fixed[tc.path] != coverage.Digest([]byte(tc.body)) {
				t.Fatal("external promise mutated")
			}
		})
	}
}

func TestCoverageCheckpointSingleFileGitIgnoreEligibility(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, "input.go")
	if err := os.WriteFile(filename, []byte("package fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("input.go\n"), 0600); err != nil {
		t.Fatal(err)
	}

	visited := []string{}
	err := file.IterateFilesListWithExclusions(
		filename,
		[]string{"."},
		false,
		false,
		func(path *file.Path) (bool, error) {
			if !path.Exists() {
				t.Fatalf("single-file parent callback path does not exist: %+v", path)
			}
			info, statErr := os.Stat(path.AbsolutePath)
			if statErr != nil {
				t.Fatal(statErr)
			}
			if !info.IsDir() {
				t.Fatalf("single-file parent callback received a file: %+v", path)
			}
			return true, nil
		},
		func(info *file.FileInfo) error {
			visited = append(visited, info.RelativePath)
			return nil
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(visited) != 0 {
		t.Fatalf("parent .gitignore exclusion failed: %v", visited)
	}

	visited = nil
	err = file.IterateFilesListWithExclusions(
		filename,
		[]string{"."},
		false,
		true,
		func(path *file.Path) (bool, error) { return true, nil },
		func(info *file.FileInfo) error {
			visited = append(visited, info.RelativePath, info.Base, info.Extension)
			return nil
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(visited, "|"), ".|input.go|.go"; got != want {
		t.Fatalf("disabled-gitignore single-file visitation=%q, want %q", got, want)
	}
}

func TestCoverageCheckpointGitIgnoreAcquisitionFailure(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "input.go"), []byte("package fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".gitignore"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := file.IterateFilesListWithExclusions(root, []string{"input.go"}, false, false, func(*file.Path) (bool, error) { return true, nil }, func(*file.FileInfo) error { return nil }, nil); err == nil {
		t.Fatal("directory at .gitignore was accepted")
	}
}

func TestCoverageCheckpointOutsideTargetFilenameRejected(t *testing.T) {
	owned := t.TempDir()
	root := filepath.Join(owned, "target")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(owned, "outside.go"), []byte("package fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	visited := false
	err := file.IterateFilesListWithExclusions(root, []string{"../outside.go"}, false, true, func(*file.Path) (bool, error) { return true, nil }, func(*file.FileInfo) error { visited = true; return nil }, nil)
	if err == nil || !strings.Contains(err.Error(), "input outside target") || visited {
		t.Fatalf("outside-target control: error=%v visited=%v", err, visited)
	}
}

func TestCoverageF4TerminalReconciliation(t *testing.T) {
	for _, mode := range []string{"valid", "leaf-identical", "component-identical", "deleted"} {
		t.Run(mode, func(t *testing.T) {
			var o *Orchestrator
			o = syntheticOrchestrator(t, func(req work.ProcessRequest) error {
				if req.InputExpectation == nil {
					return errors.New("missing binding")
				}
				target := filepath.Join(req.Dir, req.File.FilePath)
				if mode != "valid" {
					if mode == "deleted" {
						os.Remove(target)
					} else {
						outside := t.TempDir()
						os.WriteFile(filepath.Join(outside, "input.go"), []byte("synthetic"), 0600)
						dest := filepath.Join(outside, "input.go")
						if mode == "component-identical" {
							target = req.Dir
							dest = outside
						}
						os.Rename(target, target+".held")
						os.Symlink(dest, target)
					}
				}
				return os.WriteFile(req.ReportPath, nil, 0600)
			})
			var output bytes.Buffer
			out := o.scanFile(&output, files.File{FilePath: "input.go"})
			if (out.Status == "success") != (mode == "valid") {
				t.Fatalf("%s: %+v", mode, out)
			}
		})
	}
}
func TestCoverageF4PromiseDigestBinding(t *testing.T) {
	var seen *file.InputExpectation
	o := syntheticOrchestrator(t, func(req work.ProcessRequest) error {
		seen = req.InputExpectation
		return os.WriteFile(req.ReportPath, nil, 0600)
	})
	var output bytes.Buffer
	out := o.scanFile(&output, files.File{FilePath: "input.go"})
	if out.Status != "success" || seen == nil || out.SHA256 != seen.SHA256 || out.SHA256 != coverage.Digest([]byte("synthetic")) || seen.Relative != "input.go" {
		t.Fatalf("unbound outcome %+v expectation %+v", out, seen)
	}
}

func TestCoverageReportLifetimeWorkerOwnership(t *testing.T) {
	o := syntheticOrchestrator(t, func(r work.ProcessRequest) error { return os.WriteFile(r.ReportPath, nil, 0600) })
	scope, err := file.OpenInputScope(o.repository.Dir)
	if err != nil {
		t.Fatal(err)
	}
	defer scope.Close()
	o.config.ReportInputs = scope
	if err = o.Scan(filepath.Join(t.TempDir(), "report.jsonl"), []files.File{{FilePath: "input.go"}}); err != nil {
		t.Fatal(err)
	}
	reader, err := scope.ReportReader("input.go")
	if err != nil {
		t.Fatal("processing closed report capability:", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil || string(data) != "synthetic" {
		t.Fatal(string(data), err)
	}
	if err = reader.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(o.repository.Dir, "input.go")
	if err = os.Rename(path, path+".old"); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(path+".old", path); err != nil {
		t.Fatal(err)
	}
	if err = scope.Close(); err == nil {
		t.Fatal("post-worker replacement accepted at final close")
	}
}

func TestCoverageContextHandoffOrchestrator(t *testing.T) {
	for _, mode := range []string{"valid", "missing", "duplicate", "wrong-input", "short-write", "changed-after"} {
		t.Run(mode, func(t *testing.T) {
			var cfg settings.Config
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "input.go"), []byte("package synthetic\n"), 0600); err != nil {
				t.Fatal(err)
			}
			scope, err := file.OpenInputScope(root)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			cfg.ReportInputs = scope
			o := &Orchestrator{repository: work.Repository{Dir: root}, config: &cfg, maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{})}
			o.pool = &fakeWorker{scan: func(req work.ProcessRequest) error {
				var raw bytes.Buffer
				if err := detectorapi.ExtractWithInput(context.Background(), root, "input.go", &writer.Detectors{File: &raw}, nil, []string{"sast"}, nil, false, true, req.InputExpectation); err != nil {
					return err
				}
				body := raw.Bytes()
				if mode == "missing" {
					body = nil
				}
				if mode == "duplicate" {
					body = append(append([]byte{}, body...), body...)
				}
				if mode == "wrong-input" {
					body = bytes.ReplaceAll(body, []byte(`"relative":"input.go"`), []byte(`"relative":"other.go"`))
				}
				return os.WriteFile(req.ReportPath, body, 0600)
			}}
			dest := filepath.Join(t.TempDir(), "raw.jsonl")
			if mode == "short-write" {
				held, e := scope.Acquire("input.go", nil)
				if e != nil {
					t.Fatal(e)
				}
				o.inputScope = scope
				child, e := file.OpenInputScope(root)
				if e != nil {
					t.Fatal(e)
				}
				defer child.Close()
				if _, e = child.Acquire("input.go", held.Expectation()); e != nil {
					t.Fatal(e)
				}
				receipt, e := child.ContextReceipt(held.Expectation())
				if e != nil {
					t.Fatal(e)
				}
				raw, e := json.Marshal(receipt)
				if e != nil {
					t.Fatal(e)
				}
				if e = os.WriteFile(dest, raw, 0600); e != nil {
					t.Fatal(e)
				}
				if e = o.writeFileResult(coverageHandoffShortWriter{}, dest, held.Expectation()); e == nil {
					t.Fatal("parent diagnostic delivery accepted")
				}
				return
			}
			err = o.Scan(dest, []files.File{{FilePath: "input.go"}})
			if mode == "missing" || mode == "duplicate" || mode == "wrong-input" {
				if err == nil {
					t.Fatal("invalid handoff accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if err = scope.RequireContextReceipts([]string{"input.go"}); err != nil {
				t.Fatal(err)
			}
			if mode == "changed-after" {
				if err = os.WriteFile(filepath.Join(root, "go.mod"), []byte("module synthetic"), 0600); err != nil {
					t.Fatal(err)
				}
				if scope.Close() == nil {
					t.Fatal("parent lost absent marker")
				}
			}
		})
	}
}

type coverageHandoffShortWriter struct{}

func (coverageHandoffShortWriter) Write(b []byte) (int, error) { return len(b) - 1, nil }
