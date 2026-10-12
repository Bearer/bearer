package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bearer/bearer/pkg/commands/artifact/scanid"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	detectorapi "github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/coverage"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/report/writer"
	scannerstats "github.com/bearer/bearer/pkg/scanner/stats"
	globaltypes "github.com/bearer/bearer/pkg/types"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/hhatto/gocloc"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func coverageCachedRunner(t *testing.T, contents string, force bool) Runner {
	return coverageCachedRunnerConfigured(t, contents, force, nil)
}
func coverageCachedRunnerConfigured(t *testing.T, contents string, force bool, configure func(*settings.Config)) Runner {
	t.Helper()
	target := t.TempDir()
	config := settings.Config{Scan: flagtypes.ScanOptions{Target: target, Force: force, Quiet: true}, Report: flagtypes.ReportOptions{Report: flag.ReportSecurity, Format: flag.FormatJSONV2, Output: filepath.Join(t.TempDir(), "output.json")}}
	if configure != nil {
		configure(&config)
	}
	gitContext := &gitrepository.Context{CurrentCommitHash: "coverage-cache-head"}
	scanID, err := scanid.Build(config, gitContext)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(os.TempDir(), "bearer"+scanID)
	completed := strings.Replace(path, ".jsonl", "-completed.jsonl", 1)
	if err := os.WriteFile(completed, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path); _ = os.Remove(completed) })
	r, err := NewRunner(context.Background(), config, gitContext, target, nil, scannerstats.New(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestCoverageCachedRunnerRendersLegacyWithoutCoverage(t *testing.T) {
	r := coverageCachedRunner(t, `{"type":"file_list","filenames":["input.go"]}`+"\n", false)
	if !r.CacheUsed() {
		t.Fatal("completed report was not selected")
	}
	if _, _, err := r.Scan(context.Background(), flagtypes.Options{}); err != nil {
		t.Fatalf("cached Scan() error = %v", err)
	}
	before, err := os.ReadFile(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	if failed, err := r.Report(nil, nil); err != nil || failed {
		t.Fatalf("legacy cached Report() = failed=%v err=%v, want successful null-coverage output", failed, err)
	}
	out, err := os.ReadFile(r.(*runner).scanSettings.Report.Output)
	if err != nil || !bytes.Contains(out, []byte(`"coverage":null`)) {
		t.Fatalf("cached output = %s, err=%v; want null coverage", out, err)
	}
	after, err := os.ReadFile(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) || !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatal("cached report changed during Report")
	}
}

func TestCoverageCachedRunnerRejectsMalformedAndErrorRecords(t *testing.T) {
	for _, contents := range []string{
		`{"type":"file_list"` + "\n",
		`{"type":"error","file":"input.go","message":"synthetic"}` + "\n",
		`{"type":"file_error","file":"input.go","error":"synthetic timeout"}` + "\n",
		`null` + "\n",
		`{"type":"file_list","type":"error"}` + "\n",
		`{"type":"processing_coverage","coverage":{"complete":true}}` + "\n",
	} {
		t.Run(contents, func(t *testing.T) {
			r := coverageCachedRunner(t, contents, false)
			if _, err := r.Report(nil, nil); err == nil {
				t.Fatal("malformed/error cache accepted")
			}
			if _, err := os.Stat(r.(*runner).scanSettings.Report.Output); !os.IsNotExist(err) {
				t.Fatalf("failure published output: %v", err)
			}
		})
	}
}

func TestCoverageCachedRunnerDoesNotReusePriorCoverage(t *testing.T) {
	record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("prior")), ConfigurationSHA256: coverage.Digest(nil), Discovery: coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{}, Complete: true}
	b, err := json.Marshal(map[string]any{"type": "processing_coverage", "coverage": record})
	if err != nil {
		t.Fatal(err)
	}
	contents := `{"type":"file_list","filenames":[]}` + "\n" + string(b) + "\n"
	r := coverageCachedRunner(t, contents, false)
	beforeInfo, err := os.Stat(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	if failed, err := r.Report(nil, nil); err != nil || failed {
		t.Fatalf("prior-coverage Report() = failed=%v err=%v, want successful null-coverage output", failed, err)
	}
	out, err := os.ReadFile(r.(*runner).scanSettings.Report.Output)
	if err != nil || !bytes.Contains(out, []byte(`"coverage":null`)) {
		t.Fatalf("prior-coverage output = %s, err=%v; want null coverage", out, err)
	}
	after, err := os.ReadFile(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(r.ReportPath())
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != contents || !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatal("prior coverage source or timestamp mutated")
	}
}

func TestCoverageForceBypassesCompletedCache(t *testing.T) {
	r := coverageCachedRunner(t, `{"type":"file_list","filenames":["input.go"]}`+"\n", true)
	if r.CacheUsed() {
		t.Fatal("force run reused completed cache")
	}
	if _, err := r.Report(nil, nil); err == nil {
		t.Fatal("force fresh report without coverage was accepted")
	}
}

type faultFile struct {
	bytes.Buffer
	fault  string
	closed bool
}

func (f *faultFile) Name() string { return "synthetic-temp" }
func (f *faultFile) Write(p []byte) (int, error) {
	if f.fault == "write" {
		return 0, errors.New("write error")
	}
	if f.fault == "short" {
		return len(p) - 1, nil
	}
	return f.Buffer.Write(p)
}
func (f *faultFile) Sync() error {
	if f.fault == "sync" {
		return errors.New("sync error")
	}
	return nil
}
func (f *faultFile) Close() error {
	f.closed = true
	if f.fault == "close" {
		return errors.New("close error")
	}
	return nil
}
func TestCoverageDeliveryFaults(t *testing.T) {
	for _, fault := range []string{"create", "write", "short", "sync", "close", "rename", "none"} {
		t.Run(fault, func(t *testing.T) {
			f := &faultFile{fault: fault}
			published, removed := false, false
			err := publishReport("synthetic-output", []byte("report"), deliveryOps{
				create: func(string) (deliveryFile, error) {
					if fault == "create" {
						return nil, errors.New("create error")
					}
					return f, nil
				},
				rename: func(string, string) error {
					if fault == "rename" {
						return errors.New("rename error")
					}
					published = true
					return nil
				},
				remove: func(string) error { removed = true; return nil },
			})
			if (err == nil) != (fault == "none") {
				t.Fatalf("error=%v", err)
			}
			if fault != "none" && published {
				t.Fatal("incomplete report published")
			}
			if fault != "create" && (!removed || !f.closed) {
				t.Fatal("temporary report not closed/cleaned")
			}
		})
	}
}
func TestCoverageAtomicDelivery(t *testing.T) {
	target := filepath.Join(t.TempDir(), "report.json")
	if err := deliverReport(target, []byte("synthetic report")); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(target)
	if err != nil || string(b) != "synthetic report" {
		t.Fatalf("delivery %q %v", b, err)
	}
}
func TestCoverageBinding(t *testing.T) {
	d := coverage.Discovery{Files: []string{"input.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
	makeData := func() *outputtypes.ReportData {
		return &outputtypes.ReportData{Dataflow: &outputtypes.DataFlow{}, Coverage: &coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("run")), ConfigurationSHA256: coverage.Digest(nil), Discovery: d, Outcomes: []coverage.Outcome{{Path: "input.go", SHA256: coverage.Digest(nil), Status: "success"}}, Complete: true}}
	}
	if err := bindCoverage(makeData(), d, false); err != nil {
		t.Fatal(err)
	}
	if bindCoverage(makeData(), d, true) == nil {
		t.Fatal("cache accepted as fresh execution")
	}
	missing := d
	missing.Files = []string{"omitted.go"}
	if bindCoverage(makeData(), missing, false) == nil {
		t.Fatal("omission accepted")
	}
	incomplete := d
	incomplete.Complete = false
	if bindCoverage(makeData(), incomplete, false) == nil {
		t.Fatal("incomplete discovery accepted")
	}
	empty := d
	empty.Files = []string{}
	if bindCoverage(makeData(), empty, false) == nil {
		t.Fatal("unknown work accepted")
	}
	data := makeData()
	data.Coverage = nil
	if bindCoverage(data, d, false) == nil {
		t.Fatal("legacy coverage accepted")
	}
}

func TestCoverageDiffCompatibilityBinding(t *testing.T) {
	unsupported := coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{"git diff discovery does not expose a supported same-run terminal file inventory"}}
	r := &runner{gitContext: &gitrepository.Context{BaseCommitHash: "local-fixture-base"}, discovery: unsupported}
	r.scanSettings.Scan.Diff = true
	makeData := func() *outputtypes.ReportData {
		d := coverage.Discovery{Files: []string{"changed.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
		return &outputtypes.ReportData{Dataflow: &outputtypes.DataFlow{}, Coverage: &coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("run")), ConfigurationSHA256: coverage.Digest(nil), Discovery: d, Outcomes: []coverage.Outcome{{Path: "changed.go", SHA256: coverage.Digest(nil), Status: "success"}}, Complete: true}}
	}
	if err := r.validateDiscovery(true); err != nil {
		t.Fatal(err)
	}
	data := makeData()
	if err := r.bindReportCoverage(data, false); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(struct {
		Coverage *coverage.Record `json:"coverage"`
	}{data.Coverage})
	if err != nil || string(b) != `{"coverage":null}` {
		t.Fatalf("unsupported coverage=%s %v", b, err)
	}
	for _, name := range []string{"ordinary", "mode-mismatch", "extra-error", "wrong-error", "missing-base", "complete-false", "processing-error", "missing-outcome", "diagnostic", "cached"} {
		t.Run(name, func(t *testing.T) {
			rr := *r
			rr.discovery = unsupported
			dd := makeData()
			cached := false
			discoveryFailure := true
			diff := true
			switch name {
			case "ordinary":
				rr.scanSettings.Scan.Diff = false
				diff = false
			case "mode-mismatch":
				diff = false
			case "extra-error":
				rr.discovery.Errors = append([]string{}, unsupported.Errors...)
				rr.discovery.Errors = append(rr.discovery.Errors, "acquisition failed")
			case "wrong-error":
				rr.discovery.Errors = []string{"acquisition failed"}
			case "missing-base":
				rr.gitContext = nil
			case "complete-false":
				rr.discovery.Errors = nil
			case "processing-error":
				dd.Coverage.Outcomes[0].Status = "error"
				dd.Coverage.Outcomes[0].Error = "synthetic"
				discoveryFailure = false
			case "missing-outcome":
				dd.Coverage.Outcomes = nil
				discoveryFailure = false
			case "diagnostic":
				dd.Dataflow.Errors = []dataflowtypes.Error{{}}
				discoveryFailure = false
			case "cached":
				cached = true
				discoveryFailure = false
			}
			if discoveryFailure {
				if rr.validateDiscovery(diff) == nil {
					t.Fatal("incomplete discovery admitted")
				}
			} else if rr.bindReportCoverage(dd, cached) == nil {
				t.Fatal("failure became legacy success")
			}
		})
	}
}

type coverageDiffEngine struct {
	engine.Engine
	calls  []string
	failAt int
}

func (e *coverageDiffEngine) Scan(_ *settings.Config, _ *scannerstats.Stats, reportPath, target string, input []files.File) error {
	if len(input) != 1 || input[0].FilePath != "input.go" {
		return errors.New("unexpected synthetic inputs")
	}
	b, err := os.ReadFile(filepath.Join(target, "input.go"))
	if err != nil {
		return err
	}
	e.calls = append(e.calls, string(b))
	if len(e.calls) == e.failAt {
		return errors.New("synthetic engine failure")
	}
	return os.WriteFile(reportPath, nil, 0600)
}
func TestCoverageDiffRunnerPreservesBaseAndHead(t *testing.T) {
	for _, failAt := range []int{0, 1, 2} {
		t.Run(fmt.Sprintf("failure-at-%d", failAt), func(t *testing.T) {
			dir := t.TempDir()
			git := func(input string, args ...string) string {
				t.Helper()
				c := exec.Command("git", args...)
				c.Dir = dir
				c.Stdin = strings.NewReader(input)
				c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=coverage", "GIT_AUTHOR_EMAIL=coverage@example.invalid", "GIT_COMMITTER_NAME=coverage", "GIT_COMMITTER_EMAIL=coverage@example.invalid")
				b, err := c.CombinedOutput()
				if err != nil {
					t.Fatalf("fixture git %v: %v %s", args, err, b)
				}
				return strings.TrimSpace(string(b))
			}
			git("", "init", "--quiet")
			makeCommit := func(content, parent string) string {
				blob := git(content, "hash-object", "-w", "--stdin")
				tree := git("100644 blob "+blob+"\tinput.go\n", "mktree")
				args := []string{"commit-tree", tree, "-m", "synthetic"}
				if parent != "" {
					args = append(args, "-p", parent)
				}
				return git("", args...)
			}
			base := makeCommit("package fixture\nvar value = 1\n", "")
			head := makeCommit("package fixture\nvar value = 2\n", base)
			git("", "update-ref", "refs/heads/main", head)
			git("", "symbolic-ref", "HEAD", "refs/heads/main")
			git("", "read-tree", "--reset", "-u", head)
			e := &coverageDiffEngine{failAt: failAt}
			cfg := settings.Config{Worker: settings.WorkerOptions{TimeoutFileMinimum: time.Second, TimeoutFileMaximum: 5 * time.Second, TimeoutFileBytesPerSecond: 100000}}
			cfg.Scan.Diff = true
			cfg.Scan.Quiet = true
			cfg.Scan.SkipGitIgnore = true
			cfg.Report.Report = flag.ReportSecurity
			r := &runner{targetPath: dir, reportPath: filepath.Join(t.TempDir(), "report.jsonl"), scanSettings: cfg, goclocResult: &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, gitContext: &gitrepository.Context{RootDir: dir, CurrentBranch: "main", CurrentCommitHash: head, BaseCommitHash: base}, engine: e}
			_, _, err := r.Scan(context.Background(), flagtypes.Options{ScanOptions: cfg.Scan})
			if (err != nil) != (failAt != 0) {
				t.Fatalf("Scan error=%v calls=%v", err, e.calls)
			}
			if failAt == 0 && (len(e.calls) != 2 || !strings.Contains(e.calls[0], "value = 1") || !strings.Contains(e.calls[1], "value = 2")) {
				t.Fatalf("base/head sequence=%v", e.calls)
			}
			if failAt != 0 && (err == nil || !strings.Contains(err.Error(), "synthetic engine failure")) {
				t.Fatalf("processing error lost: %v", err)
			}
			if got := git("", "rev-parse", "HEAD"); got != head {
				t.Fatalf("HEAD not restored: %s", got)
			}
		})
	}
}

func TestCoverageCachedLargeValidRecord(t *testing.T) {
	body, err := json.Marshal(map[string]any{"type": "file_list", "filenames": []string{strings.Repeat("x", 6*1024*1024)}})
	if err != nil {
		t.Fatal(err)
	}
	r := coverageCachedRunner(t, string(body)+"\n", false)
	if _, err := r.Report(nil, nil); err != nil {
		t.Fatalf("valid large cached record rejected: %v", err)
	}
	out, err := os.ReadFile(r.(*runner).scanSettings.Report.Output)
	if err != nil || !bytes.Contains(out, []byte(`"coverage":null`)) {
		t.Fatalf("large cached output lacks null coverage: %v", err)
	}
}
func TestCoverageCachedDeliveryFailure(t *testing.T) {
	r := coverageCachedRunner(t, `{"type":"file_list","filenames":["input.go"]}`+"\n", false).(*runner)
	destination := t.TempDir()
	r.scanSettings.Report.Output = destination
	if _, err := r.Report(nil, nil); err == nil {
		t.Fatal("cache report delivery failure accepted")
	}
	if info, err := os.Stat(destination); err != nil || !info.IsDir() {
		t.Fatalf("destination changed: %v", err)
	}
}
func TestCoverageCachedFindingExitSemantics(t *testing.T) {
	r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":["input.go"]}`+"\n", false, func(cfg *settings.Config) {
		cfg.Report.NoExtract = true
		cfg.Report.Severity = set.Set[string]{"high": {}}
		cfg.Report.FailOnSeverity = set.Set[string]{"high": {}}
		cfg.Rules = map[string]*settings.Rule{"synthetic": {Id: "synthetic", Type: "risk", Severity: "high", Description: "Synthetic cache finding"}}
		cfg.Policies = map[string]*settings.Policy{"risk": {Query: "policy_failure = data.synthetic.policy_failure", Modules: settings.Modules{&settings.PolicyModule{Name: "synthetic.rego", Content: `package synthetic
policy_failure := [{"filename":"input.go", "full_filename":"input.go", "line_number":1, "sink":{"start":1,"end":1}}]
`}}}}
	}).(*runner)
	if !r.CacheUsed() {
		t.Fatal("cache not selected")
	}
	failed, err := r.Report(nil, nil)
	if err != nil || !failed {
		t.Fatalf("finding exit semantics lost: failed=%v error=%v", failed, err)
	}
	b, err := os.ReadFile(r.scanSettings.Report.Output)
	if err != nil {
		t.Fatal(err)
	}
	var output struct {
		Coverage *coverage.Record `json:"coverage"`
		Findings []struct {
			Severity string `json:"severity"`
		}
	}
	if err := json.Unmarshal(b, &output); err != nil {
		t.Fatal(err)
	}
	if output.Coverage != nil || len(output.Findings) != 1 || output.Findings[0].Severity != "high" {
		t.Fatalf("cached finding/coverage output=%s", b)
	}
}

func TestCoverageCacheStructuralControls(t *testing.T) {
	for _, body := range []string{
		" \n\t", `{"type":"file_list","filenames":null}`, `{"type":"file_list","filenames":[null]}`,
		`{"type":"file_list","filenames":["a.go","a.go"]}`,
		`{"type":"file_list","filenames":[]}` + "\n" + `{"type":"file_list","filenames":[]}`,
		`{"type":"file_list","filenames":[]}` + "\n" + `{"type":"schema"}`,
		`{"type":"file_list","filenames":[],"TYPE":"error"}`,
	} {
		r := coverageCachedRunner(t, body, false)
		if _, err := r.Report(nil, nil); err == nil {
			t.Fatalf("invalid legacy stream accepted %q", body)
		}
	}
	r := coverageCachedRunner(t, `{"type":"file_list","filenames":[]}`+"\n", false)
	if _, err := r.Report(nil, nil); err != nil {
		t.Fatalf("valid empty legacy report rejected: %v", err)
	}
}
func TestCoverageDetectorOnlyDiagnosticAndDeliveryControls(t *testing.T) {
	for _, body := range []string{`{"type":"error","message":"synthetic","file":"input.go"}`, `{"type":"file_error","error":"synthetic","file":"input.go"}`} {
		r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":["input.go"]}`+"\n"+body, false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = flag.FormatJSON })
		if _, err := r.Report(nil, nil); err == nil {
			t.Fatal("detector-only diagnostics accepted")
		}
	}
	r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":[]}`+"\n", false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = flag.FormatJSON }).(*runner)
	r.scanSettings.Report.Output = t.TempDir()
	if _, err := r.Report(nil, nil); err == nil {
		t.Fatal("detector-only delivery error lost")
	}
}

func TestCoverageLegacyEnvelopeCompatibility(t *testing.T) {
	valid := []string{
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go"},"value":{}}`,
		`{"type":"framework","detector_type":"simple","detector_language":"","commit_sha":"","source":{"filename":"input.go","full_filename":"","language":"","language_type":"","start_line_number":null,"start_column_number":0,"end_line_number":2,"end_column_number":null,"text":null,"legacy_extension":{"a":1}},"value":null,"legacy_extension":true}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","text":"text"},"value":[1,true,"text",null]}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go"},"value":42}`,
	}
	for _, mode := range []string{flag.ReportSecurity, flag.ReportDetectors} {
		for i, record := range valid {
			t.Run(fmt.Sprintf("%s/valid-%d", mode, i), func(t *testing.T) {
				r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":["input.go"]}`+"\n"+record+"\n", false, func(c *settings.Config) {
					c.Report.Report = mode
					if mode == flag.ReportDetectors {
						c.Report.Format = flag.FormatJSON
					}
				}).(*runner)
				if !r.CacheUsed() {
					t.Fatal("cache not selected")
				}
				if _, err := r.Report(nil, nil); err != nil {
					t.Fatal(err)
				}
				if _, err := os.ReadFile(r.scanSettings.Report.Output); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
	invalid := []string{
		`{"type":"schema","detector_type":null,"source":{"filename":"input.go"},"value":{}}`,
		`{"type":"schema","detector_type":"simple","detector_language":null,"source":{"filename":"input.go"},"value":{}}`,
		`{"type":"framework","detector_type":"simple","framework_detection_type":[],"source":{"filename":"input.go"},"value":{}}`,
		`{"type":"schema","detector_type":"simple","commit_sha":false,"source":{"filename":"input.go"},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","full_filename":null},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","language_type":[]},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","start_line_number":1.5},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","end_column_number":9223372036854775808},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","text":false},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","Filename":"hidden.go"},"value":{}}`,
		`{"type":"schema","detector_type":"simple","ſource":{"filename":"input.go"},"source":{"filename":"input.go"},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go"},"VALUE":{},"value":{}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go"},"value":{"a":1,"\u0061":2}}`,
		`{"type":"file_list","filenameſ":[]}`,
		`{"type":"file_list","filenames":[],"\u0066ilenames":["hidden.go"]}`,
	}
	for _, mode := range []string{flag.ReportSecurity, flag.ReportDetectors} {
		for i, record := range invalid {
			t.Run(fmt.Sprintf("%s/invalid-%d", mode, i), func(t *testing.T) {
				body := record + "\n"
				if !strings.Contains(record, `"type":"file_list"`) {
					body = `{"type":"file_list","filenames":["input.go"]}` + "\n" + body
				}
				r := coverageCachedRunnerConfigured(t, body, false, func(c *settings.Config) {
					c.Report.Report = mode
					if mode == flag.ReportDetectors {
						c.Report.Format = flag.FormatJSON
					}
				}).(*runner)
				if _, err := r.Report(nil, nil); err == nil {
					t.Fatalf("invalid record accepted: %s", record)
				}
				if _, err := os.Stat(r.scanSettings.Report.Output); !os.IsNotExist(err) {
					t.Fatalf("invalid record delivered: %v", err)
				}
			})
		}
	}
}

func TestCoverageUnavailableRepresentations(t *testing.T) {
	d := coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
	record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("prior")), ConfigurationSHA256: coverage.Digest(nil), Discovery: d, Outcomes: []coverage.Outcome{}, Complete: true}
	raw, err := json.Marshal(map[string]any{"type": "processing_coverage", "coverage": record})
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{flag.FormatEmpty, flag.FormatJSON, flag.FormatYAML} {
		for _, mode := range []string{"cached", "diff", "fresh"} {
			t.Run(mode+"/"+format, func(t *testing.T) {
				body := `{"type":"file_list","filenames":[]}` + "\n" + string(raw) + "\n"
				r := coverageCachedRunnerConfigured(t, body, false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = format }).(*runner)
				if mode != "cached" {
					r.reuseDetection = false
					r.discovery = d
				}
				if mode == "diff" {
					r.scanSettings.Scan.Diff = true
					r.gitContext.BaseCommitHash = "fixture-base"
					r.discovery = coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{"git diff discovery does not expose a supported same-run terminal file inventory"}}
				}
				before, err := os.ReadFile(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = r.Report(nil, nil); err != nil {
					t.Fatal(err)
				}
				output, err := os.ReadFile(r.scanSettings.Report.Output)
				if err != nil {
					t.Fatal(err)
				}
				has := bytes.Contains(output, []byte("processing_coverage"))
				if has != (mode == "fresh") {
					t.Fatalf("coverage representation %s", output)
				}
				if !bytes.Contains(output, []byte("file_list")) {
					t.Fatalf("legacy records lost: %s", output)
				}
				after, err := os.ReadFile(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				afterInfo, err := os.Stat(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(before, after) || !info.ModTime().Equal(afterInfo.ModTime()) {
					t.Fatal("cached bytes/timestamp changed")
				}
			})
		}
	}
	t.Run("malformed-before-removal", func(t *testing.T) {
		bad := strings.Replace(string(raw), `"complete":true`, `"complete":null`, 1)
		r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":[]}`+"\n"+bad+"\n", false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = flag.FormatJSON })
		if _, err := r.Report(nil, nil); err == nil {
			t.Fatal("malformed prior evidence stripped before validation")
		}
	})
	t.Run("later-dataflow-does-not-reload", func(t *testing.T) {
		data := &outputtypes.ReportData{Detectors: []any{map[string]any{"type": "processing_coverage", "coverage": record}}, Coverage: &record}
		clearUnavailableCoverage(data)
		if data.Detectors == nil || data.Coverage != nil {
			t.Fatal("unavailable representation not cleared")
		}
		if err := reportoutput.GetDataflow(data, globaltypes.Report{Path: filepath.Join(t.TempDir(), "absent")}, settings.Config{}, true); err != nil {
			t.Fatal(err)
		}
		if data.Coverage != nil || len(data.Detectors) != 0 {
			t.Fatal("prior coverage reloaded")
		}
	})
}

type coverageEligibilityEngine struct {
	engine.Engine
	calls [][]string
}

func (e *coverageEligibilityEngine) Scan(_ *settings.Config, _ *scannerstats.Stats, reportPath, target string, input []files.File) error {
	paths := []string{}
	for _, f := range input {
		b, err := os.ReadFile(filepath.Join(target, f.FilePath))
		if err != nil {
			return err
		}
		paths = append(paths, f.FilePath+":"+string(b))
	}
	e.calls = append(e.calls, paths)
	return os.WriteFile(reportPath, nil, 0600)
}
func TestCoverageDiffEligibilitySeparateTrees(t *testing.T) {
	for _, baseBinary := range []bool{false, true} {
		t.Run(fmt.Sprintf("base-binary-%v", baseBinary), func(t *testing.T) {
			dir := t.TempDir()
			git := func(input string, args ...string) string {
				t.Helper()
				c := exec.Command("git", args...)
				c.Dir = dir
				c.Stdin = strings.NewReader(input)
				c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=coverage", "GIT_AUTHOR_EMAIL=coverage@example.invalid", "GIT_COMMITTER_NAME=coverage", "GIT_COMMITTER_EMAIL=coverage@example.invalid")
				b, err := c.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v %s", args, err, b)
				}
				return strings.TrimSpace(string(b))
			}
			git("", "init", "--quiet")
			commit := func(content, parent string) string {
				blob := git(content, "hash-object", "-w", "--stdin")
				tree := git("100644 blob "+blob+"\tinput.go\n", "mktree")
				args := []string{"commit-tree", tree, "-m", "synthetic"}
				if parent != "" {
					args = append(args, "-p", parent)
				}
				return git("", args...)
			}
			baseContent, headContent := "package fixture\n", "binary\x00fixture"
			if baseBinary {
				baseContent, headContent = headContent, baseContent
			}
			base := commit(baseContent, "")
			head := commit(headContent, base)
			git("", "update-ref", "refs/heads/main", head)
			git("", "symbolic-ref", "HEAD", "refs/heads/main")
			git("", "read-tree", "--reset", "-u", head)
			cfg := settings.Config{Worker: settings.WorkerOptions{TimeoutFileMinimum: time.Second, TimeoutFileMaximum: 5 * time.Second, TimeoutFileBytesPerSecond: 100000}}
			cfg.Scan.Diff = true
			cfg.Scan.Quiet = true
			cfg.Scan.SkipGitIgnore = true
			cfg.Report.Report = flag.ReportSecurity
			e := &coverageEligibilityEngine{}
			r := &runner{targetPath: dir, reportPath: filepath.Join(t.TempDir(), "report.jsonl"), scanSettings: cfg, goclocResult: &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, gitContext: &gitrepository.Context{RootDir: dir, CurrentBranch: "main", CurrentCommitHash: head, BaseCommitHash: base}, engine: e}
			inputs, _, err := r.Scan(context.Background(), flagtypes.Options{ScanOptions: cfg.Scan})
			if err != nil {
				t.Fatal(err)
			}
			if baseBinary {
				if len(e.calls) != 1 || len(e.calls[0]) != 1 || len(inputs) != 1 {
					t.Fatalf("base exclusion suppressed eligible head: calls=%v inputs=%v", e.calls, inputs)
				}
			} else {
				if len(e.calls) != 2 || len(e.calls[0]) != 1 || len(e.calls[1]) != 0 || len(inputs) != 0 {
					t.Fatalf("head exclusion suppressed eligible base or was processed: calls=%v inputs=%v", e.calls, inputs)
				}
			}
			if git("", "rev-parse", "HEAD") != head {
				t.Fatal("head not restored")
			}
			data := &outputtypes.ReportData{Dataflow: &outputtypes.DataFlow{}, Coverage: &coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest(nil), ConfigurationSHA256: coverage.Digest(nil), Discovery: coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{}, Complete: true}}
			if err := r.bindReportCoverage(data, false); err != nil {
				t.Fatal(err)
			}
			if data.Coverage != nil {
				t.Fatal("diff invented coverage")
			}
		})
	}
}

func checkpointYAMLKey(t *testing.T, node *yaml.Node, key string) *yaml.Node {
	t.Helper()
	if node.Kind != yaml.MappingNode {
		t.Fatalf("not a YAML map for %s: %#v", key, node)
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	t.Fatalf("missing YAML key %s", key)
	return nil
}

func TestCoverageCheckpointOpaqueNumbers(t *testing.T) {
	value := `{"large":9007199254740993,"huge":1e400,"small":1.2300e-400,"nested":[-0,{"zero":0}]}`
	record := `{"type":"schema","detector_type":"simple","source":{"filename":"input.go","extension":1e400},"value":` + value + `}`
	for _, format := range []string{flag.FormatJSON, flag.FormatEmpty, flag.FormatYAML} {
		t.Run(format, func(t *testing.T) {
			body := `{"type":"file_list","filenames":["input.go"]}` + "\n" + record + "\n"
			r := coverageCachedRunnerConfigured(t, body, false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = format }).(*runner)
			before, err := os.Stat(r.reportPath)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = r.Report(nil, nil); err != nil {
				t.Fatal(err)
			}
			out, err := os.ReadFile(r.scanSettings.Report.Output)
			if err != nil {
				t.Fatal(err)
			}
			if format == flag.FormatYAML {
				var doc yaml.Node
				if err = yaml.Unmarshal(out, &doc); err != nil {
					t.Fatal(err)
				}
				rows := doc.Content[0]
				if rows.Kind != yaml.SequenceNode || len(rows.Content) != 2 {
					t.Fatalf("bad YAML records: %s", out)
				}
				got := checkpointYAMLKey(t, rows.Content[1], "value")
				for key, want := range map[string][2]string{"large": {"9007199254740993", "!!int"}, "huge": {"1e400", "!!float"}, "small": {"1.2300e-400", "!!float"}} {
					n := checkpointYAMLKey(t, got, key)
					if n.Value != want[0] || n.Tag != want[1] {
						t.Fatalf("%s token/type changed: %#v", key, n)
					}
				}
				nested := checkpointYAMLKey(t, got, "nested")
				if nested.Kind != yaml.SequenceNode || nested.Content[0].Value != "-0" || nested.Content[0].Tag != "!!int" {
					t.Fatalf("nested number changed: %#v", nested)
				}
			} else {
				d := json.NewDecoder(bytes.NewReader(out))
				d.UseNumber()
				var rows []map[string]any
				if err = d.Decode(&rows); err != nil {
					t.Fatal(err)
				}
				d = json.NewDecoder(strings.NewReader(record))
				d.UseNumber()
				var expected map[string]any
				if err = d.Decode(&expected); err != nil {
					t.Fatal(err)
				}
				if len(rows) != 2 || !reflect.DeepEqual(rows[1], expected) {
					t.Fatalf("opaque record changed: %s", out)
				}
			}
			after, err := os.ReadFile(r.reportPath)
			if err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(r.reportPath)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != body || !info.ModTime().Equal(before.ModTime()) {
				t.Fatal("cache bytes/timestamp changed")
			}
		})
	}
}

func TestCoverageCheckpointNumericValidationRemainsStrict(t *testing.T) {
	for _, record := range []string{
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","start_line_number":1e400},"value":1e400}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go"},"value":{"n":1e400,"n":0}}`,
		`{"type":"schema","detector_type":"simple","source":{"filename":"input.go","ſtart_line_number":0},"value":1e400}`,
	} {
		r := coverageCachedRunnerConfigured(t, `{"type":"file_list","filenames":["input.go"]}`+"\n"+record+"\n", false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = flag.FormatJSON }).(*runner)
		if _, err := r.Report(nil, nil); err == nil {
			t.Fatal("malformed producer field/duplicate/alias accepted")
		}
		if _, err := os.Stat(r.scanSettings.Report.Output); !os.IsNotExist(err) {
			t.Fatal("rejected evidence delivered")
		}
	}
}

func TestCoverageCheckpointBoundRepresentations(t *testing.T) {
	d := coverage.Discovery{Files: []string{"input.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
	record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("checkpoint")), ConfigurationSHA256: coverage.Digest(nil), Discovery: d, Outcomes: []coverage.Outcome{{Path: "input.go", SHA256: coverage.Digest(nil), Status: "success"}}, Complete: true}
	raw, err := json.Marshal(map[string]any{"type": "processing_coverage", "coverage": record})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ mode, format string }{{flag.ReportDetectors, flag.FormatJSON}, {flag.ReportDetectors, flag.FormatEmpty}, {flag.ReportDetectors, flag.FormatYAML}, {flag.ReportSecurity, flag.FormatJSONV2}} {
		for _, state := range []string{"fresh", "cached", "diff"} {
			t.Run(tc.mode+"/"+tc.format+"/"+state, func(t *testing.T) {
				body := `{"type":"file_list","filenames":["input.go"]}` + "\n" + string(raw) + "\n"
				r := coverageCachedRunnerConfigured(t, body, false, func(c *settings.Config) { c.Report.Report = tc.mode; c.Report.Format = tc.format }).(*runner)
				if state != "cached" {
					r.reuseDetection = false
					r.discovery = d
					r.discovery.Excluded = []coverage.Exclusion{{Path: "excluded.log", Reason: "built-in-path"}}
				}
				if state == "diff" {
					r.scanSettings.Scan.Diff = true
					r.gitContext.BaseCommitHash = "synthetic-base"
					r.discovery = coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{"git diff discovery does not expose a supported same-run terminal file inventory"}}
				}
				before, err := os.Stat(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = r.Report([]files.File{{FilePath: "input.go"}}, nil); err != nil {
					t.Fatal(err)
				}
				out, err := os.ReadFile(r.scanSettings.Report.Output)
				if err != nil {
					t.Fatal(err)
				}
				var payload any
				if tc.mode == flag.ReportSecurity {
					var v map[string]any
					if err = json.Unmarshal(out, &v); err != nil {
						t.Fatal(err)
					}
					payload = v["coverage"]
				} else {
					var rows []map[string]any
					if tc.format == flag.FormatYAML {
						err = yaml.Unmarshal(out, &rows)
					} else {
						err = json.Unmarshal(out, &rows)
					}
					if err != nil {
						t.Fatal(err)
					}
					for _, row := range rows {
						if row["type"] == "processing_coverage" {
							payload = row["coverage"]
						}
					}
				}
				if state != "fresh" {
					if payload != nil {
						t.Fatalf("unavailable coverage reappeared: %s", out)
					}
				} else {
					b, err := json.Marshal(payload)
					if err != nil {
						t.Fatal(err)
					}
					got, err := coverage.Decode(b)
					if err != nil {
						t.Fatalf("invalid emitted coverage: %v %s", err, out)
					}
					if !reflect.DeepEqual(got.Discovery.Excluded, r.discovery.Excluded) {
						t.Fatalf("discovery lost: %s", out)
					}
					if err = got.ValidateInputs(map[string]string{"input.go": coverage.Digest(nil)}); err != nil {
						t.Fatal(err)
					}
					if err = got.ValidateInputs(map[string]string{"input.go": coverage.Digest(nil), "excluded.log": coverage.Digest(nil)}); err == nil {
						t.Fatal("excluded input satisfied unchanged promise")
					}
				}
				after, err := os.ReadFile(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				if string(after) != body || !info.ModTime().Equal(before.ModTime()) {
					t.Fatal("source report bytes/timestamp changed")
				}
			})
		}
	}
}

func TestCoverageReportLifetimeFindingExtraction(t *testing.T) {
	for _, mode := range []string{"valid", "single", "replacement", "read-failure", "missing-capability", "closed-capability", "delivery-failure"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			name := "input.go"
			target := root
			full := filepath.Join(root, name)
			if err := os.WriteFile(full, []byte("acquired synthetic source\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if mode == "single" {
				target = full
				name = "."
			}
			scope, err := inputfile.OpenInputScope(target)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			bound, err := scope.Acquire(name, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err = bound.Close(); err != nil {
				t.Fatal(err)
			}
			cfg := settings.Config{ReportInputs: scope, Scan: flagtypes.ScanOptions{Target: target, Quiet: true}, Report: flagtypes.ReportOptions{Report: flag.ReportSecurity, Format: flag.FormatJSONV2, Output: filepath.Join(t.TempDir(), "out.json"), Severity: set.Set[string]{"high": {}}, FailOnSeverity: set.Set[string]{"high": {}}}}
			cfg.Rules = map[string]*settings.Rule{"synthetic": {Id: "synthetic", Type: "risk", Severity: "high"}}
			cfg.Policies = map[string]*settings.Policy{"risk": {Query: "policy_failure = data.synthetic.policy_failure", Modules: settings.Modules{&settings.PolicyModule{Name: "synthetic.rego", Content: fmt.Sprintf(`package synthetic
policy_failure := [{"filename":%q,"full_filename":"/unsafe/ignored","line_number":1,"source":{"start":1,"end":1},"sink":{"start":1,"end":1}}]`, name)}}}}
			record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest(nil), ConfigurationSHA256: coverage.Digest(nil), Discovery: coverage.Discovery{Files: []string{name}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{{Path: name, SHA256: bound.Digest(), Status: "success"}}, Complete: true}
			b, _ := json.Marshal(map[string]any{"type": "processing_coverage", "coverage": record})
			names, _ := json.Marshal(map[string]any{"type": "file_list", "filenames": []string{name}})
			reportPath := filepath.Join(t.TempDir(), "report.jsonl")
			if err := os.WriteFile(reportPath, append(append(names, '\n'), append(b, '\n')...), 0600); err != nil {
				t.Fatal(err)
			}
			r := &runner{targetPath: target, reportPath: reportPath, scanSettings: cfg, discovery: record.Discovery}
			switch mode {
			case "replacement":
				if err = os.Rename(full, full+".old"); err != nil {
					t.Fatal(err)
				}
				if err = os.Symlink(full+".old", full); err != nil {
					t.Fatal(err)
				}
			case "read-failure":
				if err = os.Remove(full); err != nil {
					t.Fatal(err)
				}
			case "missing-capability":
				r.scanSettings.ReportInputs = nil
			case "closed-capability":
				if err = scope.Close(); err != nil {
					t.Fatal(err)
				}
			case "delivery-failure":
				r.scanSettings.Report.Output = filepath.Join(t.TempDir(), "absent", "out.json")
			}
			failed, err := r.Report([]files.File{{FilePath: name}}, nil)
			valid := mode == "valid" || mode == "single"
			if valid {
				if err != nil || !failed {
					t.Fatalf("finding report failed=%v err=%v", failed, err)
				}
				out, e := os.ReadFile(cfg.Report.Output)
				if e != nil || !bytes.Contains(out, []byte("acquired synthetic source")) {
					t.Fatalf("actual source not emitted: %s %v", out, e)
				}
			} else if err == nil {
				t.Fatal("unsafe/incomplete report accepted")
			}
			if r.scanSettings.ReportInputs != nil {
				t.Fatal("scope ownership leaked")
			}
		})
	}
}
func TestCoverageReportLifetimeFinalAcceptance(t *testing.T) {
	for _, mode := range []string{"valid", "post-delivery-replacement", "delivery-error", "latched-close-error"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "input.go")
			if err := os.WriteFile(path, []byte("same"), 0600); err != nil {
				t.Fatal(err)
			}
			scope, err := inputfile.OpenInputScope(root)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			b, err := scope.Acquire("input.go", nil)
			if err != nil {
				t.Fatal(err)
			}
			if err = b.Close(); err != nil {
				t.Fatal(err)
			}
			r := &runner{scanSettings: settings.Config{ReportInputs: scope}}
			err = r.acceptReport(func() error {
				if mode == "post-delivery-replacement" {
					if e := os.Rename(path, path+".old"); e != nil {
						return e
					}
					return os.Symlink(path+".old", path)
				}
				if mode == "delivery-error" {
					return errors.New("synthetic delivery failure")
				}
				if mode == "latched-close-error" {
					_, e := scope.ReportReader("not-acquired")
					if e == nil {
						t.Fatal("fault injection did not fail")
					}
				}
				return nil
			})
			if (err == nil) != (mode == "valid") {
				t.Fatalf("acceptance %s=%v", mode, err)
			}
			if r.scanSettings.ReportInputs != nil {
				t.Fatal("scope not released")
			}
			if _, err = scope.Acquire("input.go", nil); err == nil {
				t.Fatal("scope still open")
			}
		})
	}
}

// This author control uses the real statistics formatter and final scope close.
func TestCoverageStatisticsReporting(t *testing.T) {
	for _, mode := range []string{"selected", "context", "single", "context-final-change"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			body := []byte("#!/bin/sh\necho synthetic\n")
			full := filepath.Join(root, "script")
			if err := os.WriteFile(full, body, 0600); err != nil {
				t.Fatal(err)
			}
			target, logical := root, "script"
			if mode == "single" {
				target = full
				logical = "."
			}
			scope, err := inputfile.OpenInputScope(target)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			if mode == "selected" || mode == "single" {
				b, e := scope.Acquire(logical, nil)
				if e != nil {
					t.Fatal(e)
				}
				if e = b.Close(); e != nil {
					t.Fatal(e)
				}
			}
			reportPath := filepath.Join(t.TempDir(), "report.jsonl")
			if err = os.WriteFile(reportPath, []byte("{\"type\":\"file_list\",\"filenames\":[]}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			cfg := settings.Config{ReportInputs: scope, Scan: flagtypes.ScanOptions{Target: target, Quiet: true}, Report: flagtypes.ReportOptions{Report: flag.ReportDetectors, Format: flag.FormatJSON, IncludeStats: true}}
			loc := &gocloc.Result{Files: map[string]*gocloc.ClocFile{full: {Name: full, Lang: "Shell", Code: 2}}, Languages: map[string]*gocloc.Language{}}
			data, err := reportoutput.GetData(globaltypes.Report{Path: reportPath, Inputgocloc: loc, HasFiles: true}, cfg, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(data.LanguageStats) != 1 || data.LanguageStats[0].Language != "Shell" || data.LanguageStats[0].Bytes != int64(len(body)) {
				t.Fatalf("stats: %+v", data.LanguageStats)
			}
			if data.Coverage != nil {
				t.Fatal("statistics context manufactured processing coverage")
			}
			if mode == "context-final-change" {
				if err = os.WriteFile(full, []byte("#!/bin/sh\necho different\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			r := &runner{scanSettings: cfg}
			err = r.acceptReport(func() error { return nil })
			if (err != nil) != (mode == "context-final-change") {
				t.Fatalf("final acceptance: %v", err)
			}
		})
	}
}

func TestCoverageStatisticsCachedCompatibility(t *testing.T) {
	rr := coverageCachedRunnerConfigured(t, "{\"type\":\"file_list\",\"filenames\":[]}\n", false, func(c *settings.Config) { c.Report.IncludeStats = true })
	r := rr.(*runner)
	path := filepath.Join(r.targetPath, "script")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho synthetic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	r.goclocResult = &gocloc.Result{Files: map[string]*gocloc.ClocFile{path: {Name: path, Lang: "Shell", Code: 2}}, Languages: map[string]*gocloc.Language{}}
	if !r.CacheUsed() {
		t.Fatal("fixture did not select cache")
	}
	_, err := r.Report(nil, nil)
	if err != nil {
		t.Fatalf("valid cached statistics rejected: %v", err)
	}
}

func TestCoverageCachedStatisticsOwnership(t *testing.T) {
	for _, mode := range []string{"valid", "missing", "nonregular", "outside-link", "delivery-error", "invalid-cache"} {
		t.Run(mode, func(t *testing.T) {
			contents := "{\"type\":\"file_list\",\"filenames\":[]}\n"
			if mode == "invalid-cache" {
				contents = "{truncated"
			}
			r := coverageCachedRunnerConfigured(t, contents, false, func(c *settings.Config) { c.Report.IncludeStats = true }).(*runner)
			name := filepath.Join(r.targetPath, "script")
			body := []byte("#!/bin/sh\necho synthetic\n")
			if err := os.WriteFile(name, body, 0600); err != nil {
				t.Fatal(err)
			}
			r.goclocResult = &gocloc.Result{Files: map[string]*gocloc.ClocFile{name: {Name: name, Lang: "Shell", Code: 2}}, Languages: map[string]*gocloc.Language{}}
			cacheBefore, err := os.ReadFile(r.reportPath)
			if err != nil {
				t.Fatal(err)
			}
			cacheInfo, err := os.Stat(r.reportPath)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "missing" || mode == "nonregular" || mode == "outside-link" {
				if err = os.Remove(name); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "nonregular" {
				if err = os.Mkdir(name, 0700); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "outside-link" {
				outside := filepath.Join(t.TempDir(), "script")
				if err = os.WriteFile(outside, body, 0600); err != nil {
					t.Fatal(err)
				}
				if err = os.Symlink(outside, name); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "delivery-error" {
				r.scanSettings.Report.Output = t.TempDir()
			}
			_, err = r.Report(nil, nil)
			t.Logf("mode=%s Report=%v scopeReleased=%v", mode, err, r.scanSettings.ReportInputs == nil)
			if (err == nil) != (mode == "valid") {
				t.Fatalf("report outcome: %v", err)
			}
			if r.scanSettings.ReportInputs != nil {
				t.Fatal("cached display scope leaked")
			}
			after, e := os.ReadFile(r.reportPath)
			if e != nil {
				t.Fatal(e)
			}
			afterInfo, e := os.Stat(r.reportPath)
			if e != nil {
				t.Fatal(e)
			}
			if !bytes.Equal(cacheBefore, after) || !cacheInfo.ModTime().Equal(afterInfo.ModTime()) {
				t.Fatal("cache bytes/timestamps changed")
			}
			if mode == "valid" {
				out, e := os.ReadFile(r.scanSettings.Report.Output)
				if e != nil {
					t.Fatal(e)
				}
				if !bytes.Contains(out, []byte(`"coverage":null`)) {
					t.Fatalf("cached processing coverage fabricated: %s", out)
				}
			}
		})
	}
}

func TestCoverageCachedStatisticsFinalReconciliation(t *testing.T) {
	for _, mode := range []string{"valid", "changed-during-delivery", "closed-before-delivery", "delivery-and-final-error"} {
		t.Run(mode, func(t *testing.T) {
			r := coverageCachedRunnerConfigured(t, "{\"type\":\"file_list\",\"filenames\":[]}\n", false, func(c *settings.Config) { c.Report.IncludeStats = true }).(*runner)
			name := filepath.Join(r.targetPath, "script")
			body := []byte("#!/bin/sh\necho synthetic\n")
			if err := os.WriteFile(name, body, 0600); err != nil {
				t.Fatal(err)
			}
			if err := r.prepareCachedStatisticsInputs(r.CacheUsed()); err != nil {
				t.Fatal(err)
			}
			scope := r.scanSettings.ReportInputs
			if scope == nil {
				t.Fatal("no owned display scope")
			}
			if _, _, err := scope.ReportSample("script", 16*1024); err != nil {
				t.Fatal(err)
			}
			if mode == "closed-before-delivery" {
				if err := scope.Close(); err != nil {
					t.Fatal(err)
				}
			}
			marker := errors.New("synthetic delivery failure")
			entered := false
			err := r.acceptReport(func() error {
				entered = true
				if mode == "changed-during-delivery" || mode == "delivery-and-final-error" {
					if e := os.WriteFile(name, []byte("changed"), 0600); e != nil {
						return e
					}
				}
				if mode == "delivery-and-final-error" {
					return marker
				}
				return nil
			})
			t.Logf("mode=%s deliveryEntered=%v accepted=%v", mode, entered, err)
			if (err == nil) != (mode == "valid") {
				t.Fatal("incorrect final close result")
			}
			if (mode == "closed-before-delivery") == entered {
				t.Fatal("invalid delivery precondition")
			}
			if mode == "delivery-and-final-error" && !errors.Is(err, marker) {
				t.Fatal("delivery failure lost")
			}
			if r.scanSettings.ReportInputs != nil {
				t.Fatal("scope not released")
			}
			if _, _, e := scope.ReportSample("script", 16); e == nil {
				t.Fatal("scope still readable after acceptance")
			}
		})
	}
}
func TestCoverageCachedStatisticsNeverMasksFreshOwnership(t *testing.T) {
	r := &runner{targetPath: t.TempDir(), scanSettings: settings.Config{Report: flagtypes.ReportOptions{IncludeStats: true}}}
	if err := r.prepareCachedStatisticsInputs(false); err != nil {
		t.Fatal(err)
	}
	if r.scanSettings.ReportInputs != nil {
		t.Fatal("display scope fabricated for fresh scan")
	}
	path := filepath.Join(r.targetPath, "script")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0600); err != nil {
		t.Fatal(err)
	}
	r.scanSettings.Scan.Target = r.targetPath
	_, err := reportoutput.GetData(globaltypes.Report{Inputgocloc: &gocloc.Result{Files: map[string]*gocloc.ClocFile{path: {Name: path, Lang: "Shell"}}}}, r.scanSettings, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "scope unavailable") {
		t.Fatalf("missing fresh ownership hidden: %v", err)
	}
}

func TestCoverageCacheFinalOwnershipErrors(t *testing.T) {
	for _, stats := range []bool{false, true} {
		for _, mode := range []string{"valid", "read-failure", "delivery-failure", "closed-capability"} {
			t.Run(fmt.Sprintf("stats=%v/%s", stats, mode), func(t *testing.T) {
				r := coverageCachedRunnerConfigured(t, "{\"type\":\"file_list\",\"filenames\":[\"input.go\"]}\n", false, func(c *settings.Config) {
					c.Report.IncludeStats = stats
					c.Report.Severity = set.Set[string]{"high": {}}
					c.Report.FailOnSeverity = set.Set[string]{"high": {}}
					c.Rules = map[string]*settings.Rule{"synthetic": {Id: "synthetic", Type: "risk", Severity: "high"}}
					c.Policies = map[string]*settings.Policy{"risk": {Query: "policy_failure = data.synthetic.policy_failure", Modules: settings.Modules{&settings.PolicyModule{Name: "synthetic.rego", Content: `package synthetic
policy_failure := [{"filename":"input.go","full_filename":"/unsafe/ignored","line_number":1,"source":{"start":1,"end":1},"sink":{"start":1,"end":1}}]`}}}}
				}).(*runner)
				name := filepath.Join(r.targetPath, "input.go")
				if err := os.WriteFile(name, []byte("acquired display source\n"), 0600); err != nil {
					t.Fatal(err)
				}
				before, err := os.ReadFile(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(r.reportPath)
				if err != nil {
					t.Fatal(err)
				}
				if mode == "read-failure" {
					if err = os.Remove(name); err != nil {
						t.Fatal(err)
					}
				}
				if mode == "delivery-failure" {
					r.scanSettings.Report.Output = t.TempDir()
				}
				if mode == "closed-capability" {
					if err = r.prepareCachedStatisticsInputs(r.CacheUsed()); err != nil {
						t.Fatal(err)
					}
					if err = r.scanSettings.ReportInputs.Close(); err != nil {
						t.Fatal(err)
					}
				}
				failed, err := r.Report(nil, nil)
				t.Logf("stats=%v mode=%s failed=%v error=%v", stats, mode, failed, err)
				if mode == "valid" {
					if err != nil || !failed {
						t.Fatalf("finding behavior lost: %v", err)
					}
					out, e := os.ReadFile(r.scanSettings.Report.Output)
					if e != nil || !bytes.Contains(out, []byte("acquired display source")) || !bytes.Contains(out, []byte(`"coverage":null`)) {
						t.Fatalf("display or coverage changed: %s %v", out, e)
					}
				} else if err == nil {
					t.Fatal("incomplete report accepted")
				}
				if r.scanSettings.ReportInputs != nil {
					t.Fatal("display scope leaked")
				}
				after, e := os.ReadFile(r.reportPath)
				if e != nil {
					t.Fatal(e)
				}
				afterInfo, e := os.Stat(r.reportPath)
				if e != nil {
					t.Fatal(e)
				}
				if !bytes.Equal(before, after) || !info.ModTime().Equal(afterInfo.ModTime()) {
					t.Fatal("cached report changed")
				}
			})
		}
	}
}

type contextHandoffEngine struct {
	engine.Engine
	scan func(*settings.Config, string, string, []files.File) error
}

func (e *contextHandoffEngine) Scan(c *settings.Config, _ *scannerstats.Stats, out, root string, inputs []files.File) error {
	return e.scan(c, out, root, inputs)
}
func TestCoverageContextHandoffReport(t *testing.T) {
	for _, mode := range []string{"valid", "missing", "duplicate", "wrong-root", "malformed", "cached-stale", "cached-malformed", "diff", "delivery-failure", "during-delivery", "context-during-delivery", "fresh-legacy"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "input.go"), []byte("package synthetic\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if mode == "context-during-delivery" {
				if e := os.WriteFile(filepath.Join(root, "manage.py"), []byte("from django.core.management import execute_from_command_line\n"), 0600); e != nil {
					t.Fatal(e)
				}
			}
			discovery := coverage.Discovery{Files: []string{"input.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
			cfg := settings.Config{Scan: flagtypes.ScanOptions{Target: root, Quiet: true, SkipGitIgnore: true}, Report: flagtypes.ReportOptions{Report: flag.ReportDetectors, Format: flag.FormatJSON, Output: filepath.Join(t.TempDir(), "out.json")}}
			eng := &contextHandoffEngine{}
			eng.scan = func(c *settings.Config, out, root string, inputs []files.File) error {
				b, e := c.ReportInputs.Acquire("input.go", nil)
				if e != nil {
					return e
				}
				expected := b.Expectation()
				var raw bytes.Buffer
				e = detectorapi.ExtractWithInput(context.Background(), root, "input.go", &writer.Detectors{File: &raw}, nil, []string{"sast"}, nil, false, true, expected)
				if e = errors.Join(e, b.Close()); e != nil {
					return e
				}
				var result bytes.Buffer
				decoder := json.NewDecoder(&raw)
				for {
					var record json.RawMessage
					e := decoder.Decode(&record)
					if e == io.EOF {
						break
					}
					if e != nil {
						return e
					}
					var head struct{ Type string }
					if e = json.Unmarshal(record, &head); e != nil {
						return e
					}
					if head.Type == inputfile.ContextReceiptType {
						if mode == "missing" {
							continue
						}
						if mode == "wrong-root" || mode == "cached-stale" {
							r, e := inputfile.DecodeContextReceipt(record)
							if e != nil {
								return e
							}
							r.Input.Root = root + "/different"
							record, e = json.Marshal(r)
							if e != nil {
								return e
							}
						}
						if mode == "malformed" || mode == "cached-malformed" {
							record = bytes.Replace(record, []byte(`"observations":`), []byte(`"Observations":`), 1)
						}
						if mode == "fresh-legacy" {
							var fields map[string]json.RawMessage
							if e = json.Unmarshal(record, &fields); e != nil {
								return e
							}
							fields["schema"] = json.RawMessage(`"bearer-context-v1"`)
							delete(fields, "reads")
							record, e = json.Marshal(fields)
							if e != nil {
								return e
							}
						}
						if mode == "duplicate" {
							result.Write(record)
							result.WriteByte('\n')
						}
					}
					result.Write(record)
					result.WriteByte('\n')
				}
				enc := json.NewEncoder(&result)
				if e = enc.Encode(map[string]any{"type": "file_list", "filenames": []string{"input.go"}}); e != nil {
					return e
				}
				record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("synthetic")), ConfigurationSHA256: coverage.Digest(nil), Discovery: discovery, Outcomes: []coverage.Outcome{{Path: "input.go", SHA256: expected.SHA256, Status: "success"}}, Complete: true}
				if e = enc.Encode(map[string]any{"type": "processing_coverage", "coverage": record}); e != nil {
					return e
				}
				return os.WriteFile(out, result.Bytes(), 0600)
			}
			r := &runner{targetPath: root, reportPath: filepath.Join(t.TempDir(), "raw.jsonl"), scanSettings: cfg, discovery: discovery, engine: eng}
			defer r.closeReportInputs()
			if err := r.scanWithInputs(r.reportPath, []files.File{{FilePath: "input.go"}}); err != nil {
				t.Fatal(err)
			}
			if mode == "cached-stale" || mode == "cached-malformed" {
				raw, e := os.ReadFile(r.reportPath)
				if e != nil {
					t.Fatal(e)
				}
				if e = r.closeReportInputs(); e != nil {
					t.Fatal(e)
				}
				r = coverageCachedRunnerConfigured(t, string(raw), false, func(c *settings.Config) { c.Report.Report = flag.ReportDetectors; c.Report.Format = flag.FormatJSON }).(*runner)
				defer r.closeReportInputs()
			}
			if mode == "diff" {
				r.scanSettings.Scan.Diff = true
				r.gitContext = &gitrepository.Context{BaseCommitHash: "synthetic-base"}
				r.discovery = coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{"git diff discovery does not expose a supported same-run terminal file inventory"}, Complete: false}
			}
			if mode == "delivery-failure" {
				r.scanSettings.Report.Output = t.TempDir()
			}
			if mode == "during-delivery" || mode == "context-during-delivery" {
				raw, e := os.ReadFile(r.reportPath)
				if e != nil {
					t.Fatal(e)
				}
				decoder := json.NewDecoder(bytes.NewReader(raw))
				for {
					var record json.RawMessage
					e = decoder.Decode(&record)
					if e == io.EOF {
						break
					}
					if e != nil {
						t.Fatal(e)
					}
					var header struct{ Type string }
					if e = json.Unmarshal(record, &header); e != nil {
						t.Fatal(e)
					}
					if header.Type == inputfile.ContextReceiptType {
						receipt, e := inputfile.DecodeContextReceipt(record)
						if e != nil {
							t.Fatal(e)
						}
						if e = r.scanSettings.ReportInputs.ImportContextReceipt(receipt, nil); e != nil {
							t.Fatal(e)
						}
					}
				}
				data, e := reportoutput.GetData(globaltypes.Report{Path: r.reportPath, HasFiles: true}, r.scanSettings, nil, nil)
				if e != nil {
					t.Fatal(e)
				}
				if e = r.bindReportCoverage(data, false); e != nil {
					t.Fatal(e)
				}
				rendered, e := reportoutput.FormatOutput(data, r.scanSettings, r.engine, nil, time.Now(), time.Now())
				if e != nil {
					t.Fatal(e)
				}
				delivered := false
				err := r.acceptReport(func() error {
					if e := deliverReport(r.scanSettings.Report.Output, []byte(rendered)); e != nil {
						return e
					}
					delivered = true
					if mode == "context-during-delivery" {
						return os.WriteFile(filepath.Join(root, "manage.py"), []byte("from foobar.core.management import execute_from_command_line\n"), 0600)
					}
					return os.WriteFile(filepath.Join(root, "go.mod"), []byte("module created\n"), 0600)
				})
				if !delivered {
					t.Fatal("actual delivery not reached")
				}
				t.Logf("mode=%s delivered=%v terminal=%v", mode, delivered, err)
				if err == nil || r.scanSettings.ReportInputs != nil {
					t.Fatalf("delivery-time mutation accepted/scope retained: %v", err)
				}
				return
			}
			_, err := r.Report([]files.File{{FilePath: "input.go"}}, nil)
			good := mode == "valid" || mode == "cached-stale" || mode == "diff"
			if (err == nil) != good {
				t.Fatalf("mode=%s err=%v", mode, err)
			}
			if r.scanSettings.ReportInputs != nil {
				t.Fatal("scope retained")
			}
			if good {
				out, e := os.ReadFile(r.scanSettings.Report.Output)
				if e != nil {
					t.Fatal(e)
				}
				if bytes.Contains(out, []byte(inputfile.ContextReceiptType)) {
					t.Fatal("internal receipt emitted")
				}
				if mode != "valid" && bytes.Contains(out, []byte("processing_coverage")) {
					t.Fatal("unavailable coverage emitted")
				}
			}
		})
	}
}
