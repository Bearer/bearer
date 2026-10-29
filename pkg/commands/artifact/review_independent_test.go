package artifact

import (
    "bytes"
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/bearer/bearer/pkg/commands/artifact/scanid"
    "github.com/bearer/bearer/pkg/commands/process/filelist/files"
    "github.com/bearer/bearer/pkg/commands/process/gitrepository"
    "github.com/bearer/bearer/pkg/commands/process/settings"
    "github.com/bearer/bearer/pkg/flag"
    flagtypes "github.com/bearer/bearer/pkg/flag/types"
    "github.com/bearer/bearer/pkg/report/coverage"
    "github.com/bearer/bearer/pkg/scanner/stats"
)

func TestCoverageReviewMalformedSelectedLegacyCache(t *testing.T) {
    for _, tc := range []struct{name, contents string}{
        {"empty", ""},
        {"untyped_object", "{}\n"},
        {"partial_file_list", "{\"type\":\"file_list\"}\n"},
        {"unknown_type", "{\"type\":\"not-a-bearer-record\"}\n"},
    } {
        t.Run(tc.name, func(t *testing.T) {
            target := t.TempDir()
            cfg := settings.Config{
                Scan: flagtypes.ScanOptions{Target: target, Quiet: true},
                Report: flagtypes.ReportOptions{Report: flag.ReportSecurity, Format: flag.FormatJSONV2, Output: filepath.Join(t.TempDir(), "result.json")},
            }
            gitContext := &gitrepository.Context{CurrentCommitHash: "independent-review-cache"}
            id, err := scanid.Build(cfg, gitContext)
            if err != nil { t.Fatal(err) }
            path := filepath.Join(os.TempDir(), "bearer"+id)
            completed := strings.Replace(path, ".jsonl", "-completed.jsonl", 1)
            if err := os.WriteFile(completed, []byte(tc.contents), 0600); err != nil { t.Fatal(err) }
            defer os.Remove(completed)
            defer os.Remove(path)
            r, err := NewRunner(context.Background(), cfg, gitContext, target, nil, stats.New(), nil)
            if err != nil { t.Fatal(err) }
            if !r.CacheUsed() { t.Fatal("cache exception was not selected") }
            failed, err := r.Report(nil, nil)
            if err == nil {
                delivered, readErr := os.ReadFile(cfg.Report.Output)
                t.Errorf("invalid selected legacy cache delivered: failed=%v readError=%v output=%s", failed, readErr, delivered)
            }
        })
    }
}

func TestCoverageReviewDetectorReportRemainsUsable(t *testing.T) {
    discovery := coverage.Discovery{Files: []string{"input.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}
    record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("review-run")), ConfigurationSHA256: coverage.Digest(nil),
        Discovery: discovery, Outcomes: []coverage.Outcome{{Path: "input.go", SHA256: coverage.Digest([]byte("input")), Status: "success"}}, Complete: true}
    var raw bytes.Buffer
    encoder := json.NewEncoder(&raw)
    if err := encoder.Encode(map[string]any{"type":"file_list", "filenames":[]string{"input.go"}}); err != nil { t.Fatal(err) }
    if err := encoder.Encode(map[string]any{"type":"processing_coverage", "coverage":record}); err != nil { t.Fatal(err) }
    path := filepath.Join(t.TempDir(), "report.jsonl")
    if err := os.WriteFile(path, raw.Bytes(), 0600); err != nil { t.Fatal(err) }
    cfg := settings.Config{Scan: flagtypes.ScanOptions{Quiet:true},
        Report: flagtypes.ReportOptions{Report:flag.ReportDetectors, Format:flag.FormatJSON, Output:filepath.Join(t.TempDir(),"detectors.json")}}
    r := &runner{reportPath:path, discovery:discovery, scanSettings:cfg}
    if _, err := r.Report([]files.File{{FilePath:"input.go"}}, nil); err != nil {
        t.Errorf("valid fresh detector report rejected: %v", err)
    }
}

