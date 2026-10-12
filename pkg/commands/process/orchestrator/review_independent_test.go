package orchestrator

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "io"
    stdlog "log"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/bearer/bearer/pkg/commands/process/filelist"
    "github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
    "github.com/bearer/bearer/pkg/commands/process/settings"
    detectorapi "github.com/bearer/bearer/pkg/detectors"
    reportapi "github.com/bearer/bearer/pkg/report"
    "github.com/bearer/bearer/pkg/report/coverage"
    reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
    "github.com/bearer/bearer/pkg/report/writer"
    "github.com/bearer/bearer/pkg/util/file"
    "github.com/hhatto/gocloc"
)

type reviewPool struct { scan func(work.ProcessRequest) error }
func (p *reviewPool) Scan(r work.ProcessRequest) error { return p.scan(r) }
func (*reviewPool) Close() {}

type reviewDetector struct { calls *int; fault bool }
func (*reviewDetector) AcceptDir(*file.Path) (bool, error) { return true, nil }
func (d *reviewDetector) ProcessFile(_ *file.FileInfo, _ *file.Path, _ reportapi.Report) (bool, error) {
    *d.calls = *d.calls + 1
    if d.fault { return false, errors.New("independent detector failure") }
    return true, nil
}
type reviewWriteFailure struct { target io.Writer; failures *int }
func (w reviewWriteFailure) Write([]byte) (int, error) {
    *w.failures = *w.failures + 1
    return 0, errors.New("independent diagnostic capture failure")
}

func reviewTerminal(t *testing.T, reportPath string) *coverage.Record {
    t.Helper()
    raw, err := os.ReadFile(reportPath)
    if err != nil { t.Fatal(err) }
    decoder := json.NewDecoder(bytes.NewReader(raw))
    for {
        var item struct { Type string `json:"type"`; Coverage *coverage.Record `json:"coverage"` }
        err := decoder.Decode(&item)
        if err == io.EOF { break }
        if err != nil { t.Fatal(err) }
        if item.Type == "processing_coverage" { return item.Coverage }
    }
    t.Fatal("terminal record missing")
    return nil
}

func TestCoverageReviewExtractorSkipsAreNotSuccessfulProcessing(t *testing.T) {
    for _, tc := range []struct { name, path, contents string; skipTest, gitIgnore, shouldProcess bool }{
        {"control", "input.go", "package synthetic\n", false, false, true},
        {"test_filter", "input_test.go", "package synthetic\n", true, false, false},
        {"built_in_path_filter", "input.log", "synthetic text\n", false, false, false},
        {"binary_filter", "input.go", "\x00synthetic\x00", false, false, false},
        {"gitignore_filter", "input.go", "package synthetic\n", false, true, false},
    } {
        t.Run(tc.name, func(t *testing.T) {
            dir := t.TempDir()
            if err := os.WriteFile(filepath.Join(dir, tc.path), []byte(tc.contents), 0600); err != nil { t.Fatal(err) }
            // Fix this independently before discovery.
            promise := map[string]string{tc.path: coverage.Digest([]byte(tc.contents))}
            cfg := settings.Config{IgnoreGit: true, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 100000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: time.Second}}
            cfg.Scan.Quiet = true
            cfg.Scan.SkipTest = tc.skipTest
            cfg.Scan.SkipGitIgnore = !tc.gitIgnore
            if tc.gitIgnore {
                if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(tc.path+"\n"), 0600); err != nil { t.Fatal(err) }
                cfg.Scan.SkipPath = []string{".gitignore"}
            }
            loc := &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}
            list, err := filelist.Discover(nil, dir, loc, cfg)
            if err != nil || !list.Coverage.Complete || len(list.Files) != 1 || list.Files[0].FilePath != tc.path {
                t.Fatalf("unexpected independent discovery: %#v, %v", list, err)
            }
            calls := 0
            detector := &reviewDetector{calls: &calls}
            p := &reviewPool{scan: func(r work.ProcessRequest) error {
                f, err := os.Create(r.ReportPath)
                if err != nil { return err }
                extractErr := detectorapi.ExtractWithDetectors(context.Background(), r.Dir, r.File.FilePath,
                    &writer.Detectors{File: f}, nil,
                    []detectorapi.InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: detector}},
                    nil, tc.skipTest, !tc.gitIgnore)
                return errors.Join(extractErr, f.Close())
            }}
            o := &Orchestrator{repository: work.Repository{Dir: dir}, config: &cfg,
                maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{}), pool: p}
            output := filepath.Join(t.TempDir(), "report.jsonl")
            err = o.Scan(output, list.Files)
            if tc.shouldProcess {
                if err != nil || calls != 1 { t.Fatalf("control failed: scan=%v calls=%d", err, calls) }
                if err := reviewTerminal(t, output).ValidateInputs(promise); err != nil { t.Fatal(err) }
                return
            }
            if calls != 0 { t.Fatalf("probe did not reach a skip: calls=%d", calls) }
            if err == nil {
                record := reviewTerminal(t, output)
                t.Errorf("extractor calls=0 but Scan succeeded; terminal=%+v; fixed-promise reconciliation=%v", record, record.ValidateInputs(promise))
            }
        })
    }
}

func TestCoverageReviewDiagnosticCaptureFailureReachesTerminal(t *testing.T) {
    dir := t.TempDir()
    input := []byte("package synthetic\n")
    if err := os.WriteFile(filepath.Join(dir, "input.go"), input, 0600); err != nil { t.Fatal(err) }
    promise := map[string]string{"input.go": coverage.Digest(input)}
    cfg := settings.Config{IgnoreGit: true, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 100000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: time.Second}}
    cfg.Scan.Quiet = true
    loc := &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}
    list, err := filelist.Discover(nil, dir, loc, cfg)
    if err != nil { t.Fatal(err) }
    calls, failures := 0, 0
    var logs bytes.Buffer
    previous := stdlog.Writer()
    stdlog.SetOutput(&logs)
    defer stdlog.SetOutput(previous)
    var extractErr error
    p := &reviewPool{scan: func(r work.ProcessRequest) error {
        f, err := os.Create(r.ReportPath)
        if err != nil { return err }
        sink := reviewWriteFailure{target: f, failures: &failures}
        extractErr = detectorapi.ExtractWithDetectors(context.Background(), r.Dir, r.File.FilePath,
            &writer.Detectors{File: sink}, nil,
            []detectorapi.InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: &reviewDetector{calls: &calls, fault: true}}},
            nil, false, true)
        return errors.Join(extractErr, f.Close())
    }}
    o := &Orchestrator{repository: work.Repository{Dir: dir}, config: &cfg,
        maxWorkersSemaphore: make(chan struct{}, 1), done: make(chan struct{}), pool: p}
    output := filepath.Join(t.TempDir(), "report.jsonl")
    err = o.Scan(output, list.Files)
    if calls != 1 || failures != 1 || !strings.Contains(logs.String(), "failed to encode data line") {
        t.Fatalf("diagnostic failure injection not reached: calls=%d failures=%d logs=%q", calls, failures, logs.String())
    }
    if err == nil {
        record := reviewTerminal(t, output)
        t.Errorf("real AddError writer failure logged but terminal succeeded; Extract error=%v; fixed-promise reconciliation=%v; logs=%q", extractErr, record.ValidateInputs(promise), logs.String())
    }
}

