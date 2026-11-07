package detectors

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/bearer/pkg/report"
	reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/writer"
	"github.com/bearer/bearer/pkg/scanner/stats"
	"github.com/bearer/bearer/pkg/util/file"
)

type coverageDetector struct {
	acceptErr      error
	processErr     error
	panicOnProcess bool
}

func (d coverageDetector) AcceptDir(*file.Path) (bool, error) {
	if d.acceptErr != nil {
		return false, d.acceptErr
	}
	return true, nil
}

func (d coverageDetector) ProcessFile(*file.FileInfo, *file.Path, report.Report) (bool, error) {
	if d.panicOnProcess {
		panic("synthetic detector panic")
	}
	return false, d.processErr
}

func coverageExtract(t *testing.T, detector coverageDetector) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("package input\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	report := &writer.Detectors{File: &output}
	err := ExtractWithDetectors(
		context.Background(), dir, "input.go", report, stats.NewFileStats(),
		[]InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: detector}}, nil, false, true,
	)
	if err != nil {
		t.Fatalf("ExtractWithDetectors() error = %v", err)
	}
	return output.String()
}

func TestCoverageExtractAcceptDirErrorWritesDiagnostic(t *testing.T) {
	output := coverageExtract(t, coverageDetector{acceptErr: errors.New("synthetic accept-dir failure")})
	if !strings.Contains(output, `"type":"error"`) || !strings.Contains(output, "synthetic accept-dir failure") {
		t.Fatalf("output = %q, want error diagnostic", output)
	}
}

func TestCoverageExtractProcessFileErrorWritesDiagnostic(t *testing.T) {
	output := coverageExtract(t, coverageDetector{processErr: errors.New("synthetic process-file failure")})
	if !strings.Contains(output, `"type":"error"`) || !strings.Contains(output, "synthetic process-file failure") {
		t.Fatalf("output = %q, want error diagnostic", output)
	}
}

func TestCoverageExtractRecoveredPanicWritesDiagnostic(t *testing.T) {
	output := coverageExtract(t, coverageDetector{panicOnProcess: true})
	if !strings.Contains(output, `"type":"error"`) || !strings.Contains(output, "synthetic detector panic") {
		t.Fatalf("output = %q, want recovered panic diagnostic", output)
	}
}

func TestCoverageMissingExcludedAndEmptyScannerInputs(t *testing.T) {
	dir := t.TempDir()
	for _, test := range []struct {
		name     string
		contents []byte
		create   bool
		good     bool
	}{
		{"input.go", []byte("package synthetic\n"), true, true},
		{"missing.go", nil, false, false},
		{"input.log", []byte("synthetic"), true, false},
		{"binary.go", []byte{0, 1, 0}, true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.create {
				if err := os.WriteFile(filepath.Join(dir, test.name), test.contents, 0600); err != nil {
					t.Fatal(err)
				}
			}
			var b bytes.Buffer
			report := &writer.Detectors{File: &b}
			err := ExtractWithDetectors(context.Background(), dir, test.name, report, nil, nil, nil, false, true)
			if (err == nil) != test.good {
				t.Fatalf("empty-scanner visitation error=%v want success=%v", err, test.good)
			}
		})
	}
}

func TestCoverageUnreadableInput(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "unreadable.go")
	if err := os.WriteFile(name, []byte("package synthetic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(name, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(name, 0600)
	if _, err := os.ReadFile(name); err == nil {
		t.Fatal("unreadable fixture precondition unavailable")
	}
	var b bytes.Buffer
	if err := ExtractWithDetectors(context.Background(), dir, "unreadable.go", &writer.Detectors{File: &b}, nil, nil, nil, false, true); err == nil {
		t.Fatal("unreadable promised input accepted")
	}
}

type coverageReadFaultDetector struct{}

func (coverageReadFaultDetector) AcceptDir(*file.Path) (bool, error) { return true, nil }
func (coverageReadFaultDetector) ProcessFile(f *file.FileInfo, _ *file.Path, _ report.Report) (bool, error) {
	if err := os.Remove(f.AbsolutePath); err != nil {
		return false, err
	}
	// Deliberately swallow the read error. Scope-level accounting must still fail.
	_, _ = f.ReadBytes()
	return true, nil
}
func TestCoverageF4DiagnosticReadFailure(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "input.go"), []byte("package synthetic\n"), 0600)
	var output bytes.Buffer
	err := ExtractWithDetectors(context.Background(), dir, "input.go", &writer.Detectors{File: &output}, nil, []InitializedDetector{{Type: reportdetectors.DetectorSimple, Detector: coverageReadFaultDetector{}}}, nil, false, true)
	if err == nil {
		t.Fatal("swallowed input fault became successful processing")
	}
}
