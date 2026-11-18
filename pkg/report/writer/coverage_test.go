package writer

import (
	"bytes"
	"errors"
	inputfile "github.com/bearer/bearer/pkg/util/file"
	"io"
	"log"
	"strings"
	"testing"
)

type coverageFailWriter struct{}

func (coverageFailWriter) Write([]byte) (int, error) {
	return 0, errors.New("synthetic writer failure")
}

func TestCoverageAddErrorWritesErrorDetection(t *testing.T) {
	var output bytes.Buffer
	report := &Detectors{File: &output}
	report.AddError("input.go", errors.New("synthetic detector failure"))
	if !strings.Contains(output.String(), `"type":"error"`) || !strings.Contains(output.String(), "synthetic detector failure") {
		t.Fatalf("output = %q, want error detection", output.String())
	}
}

func TestCoverageAddErrorLogsWriterFailure(t *testing.T) {
	previous := log.Writer()
	defer log.SetOutput(previous)
	var logs bytes.Buffer
	log.SetOutput(&logs)

	report := &Detectors{File: coverageFailWriter{}}
	report.AddError("input.go", errors.New("synthetic detector failure"))
	if !strings.Contains(logs.String(), "failed to encode data line") {
		t.Fatalf("logs = %q, want writer failure log", logs.String())
	}
}

type coverageShortWriter struct{}

func (coverageShortWriter) Write(b []byte) (int, error) { return len(b) - 1, nil }
func TestCoverageWriterRetainsDeliveryFaults(t *testing.T) {
	for _, sink := range []io.Writer{coverageFailWriter{}, coverageShortWriter{}} {
		report := &Detectors{File: sink}
		report.AddError("input.go", errors.New("synthetic"))
		if report.Err() == nil {
			t.Fatal("write failure not retained")
		}
		var valid bytes.Buffer
		report.File = &valid
		report.AddError("input.go", errors.New("later diagnostic"))
		if report.Err() == nil {
			t.Fatal("later success cleared prior failure")
		}
	}
	var valid bytes.Buffer
	report := &Detectors{File: &valid}
	report.AddError("input.go", errors.New("recorded diagnostic"))
	if report.Err() != nil {
		t.Fatal("successful error record mistaken for delivery failure")
	}
}

func TestCoverageContextHandoffWriterFaults(t *testing.T) {
	for _, sink := range []io.Writer{coverageFailWriter{}, coverageShortWriter{}} {
		r := &Detectors{File: sink}
		if r.AddContextReceipt(&inputfile.ContextReceipt{}) == nil || r.Err() == nil {
			t.Fatal("context write error lost")
		}
		var good bytes.Buffer
		r.File = &good
		if r.AddContextReceipt(&inputfile.ContextReceipt{}) == nil {
			t.Fatal("later write cleared fault")
		}
	}
}
