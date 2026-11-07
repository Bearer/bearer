package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"os"
	"path/filepath"
	"testing"
)

type coverageSink struct {
	bytes.Buffer
	closed bool
	fault  error
}

func (s *coverageSink) Close() error { s.closed = true; return s.fault }
func TestCoverageWorkerReportCloseFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("package synthetic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	scope, err := inputfile.OpenInputScope(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer scope.Close()
	bound, err := scope.Acquire("input.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, fault := range []error{nil, errors.New("synthetic close failure")} {
		sink := &coverageSink{fault: fault}
		w := &Worker{skipGitIgnore: true}
		_, err := w.scanReport(context.Background(), work.ProcessRequest{Repository: work.Repository{Dir: dir}, File: files.File{FilePath: "input.go"}, InputExpectation: bound.Expectation()}, sink)
		if !sink.closed || (err == nil) != (fault == nil) {
			t.Fatalf("close=%v err=%v expected fault=%v", sink.closed, err, fault)
		}
	}
}

func TestCoverageF4WorkerExpectation(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "input.go"), []byte("package synthetic\n"), 0600)
	scope, err := inputfile.OpenInputScope(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer scope.Close()
	bound, err := scope.Acquire("input.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"valid", "missing", "digest", "selector", "identity"} {
		t.Run(kind, func(t *testing.T) {
			req := work.ProcessRequest{Repository: work.Repository{Dir: dir}, File: files.File{FilePath: "input.go"}, InputExpectation: bound.Expectation()}
			switch kind {
			case "missing":
				req.InputExpectation = nil
			case "digest":
				req.InputExpectation.SHA256 = "bad"
			case "selector":
				req.InputExpectation.Relative = "elsewhere"
			case "identity":
				req.InputExpectation.RootIdentity.Inode++
			}
			data, err := json.Marshal(req)
			if err != nil {
				t.Fatal(err)
			}
			var decoded work.ProcessRequest
			if err = json.Unmarshal(data, &decoded); err != nil {
				t.Fatal(err)
			}
			sink := &coverageSink{}
			w := &Worker{skipGitIgnore: true}
			_, err = w.scanReport(context.Background(), decoded, sink)
			if (err == nil) != (kind == "valid") || !sink.closed {
				t.Fatalf("%s result %v closed=%v", kind, err, sink.closed)
			}
		})
	}
}
