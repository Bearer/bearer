package dataflow

import (
	"encoding/json"
	"github.com/bearer/bearer/pkg/report/coverage"
	"github.com/bearer/bearer/pkg/report/output/types"
	"testing"
)

func syntheticCoverage(t *testing.T) map[string]any {
	t.Helper()
	r := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest(nil), ConfigurationSHA256: coverage.Digest(nil), Discovery: coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{}, Complete: true}
	b, err := json.Marshal(map[string]any{"type": "processing_coverage", "coverage": r})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}
func TestCoverageReportDecoding(t *testing.T) {
	m := syntheticCoverage(t)
	d := &types.ReportData{Detectors: []any{m}}
	if err := addCoverage(d); err != nil {
		t.Fatal(err)
	}
	if d.Coverage == nil {
		t.Fatal("missing decoded evidence")
	}
	d.Detectors = append(d.Detectors, m)
	if addCoverage(d) == nil {
		t.Fatal("duplicate terminal accepted")
	}
	for _, bad := range []any{nil, "broken", map[string]any{"type": "processing_coverage", "coverage": map[string]any{"schema": coverage.Schema}}, map[string]any{"type": "processing_coverage", "coverage": nil}} {
		d.Detectors = []any{bad}
		if addCoverage(d) == nil {
			t.Fatalf("malformed accepted: %#v", bad)
		}
	}
	d.Detectors = []any{}
	if err := addCoverage(d); err != nil {
		t.Fatal(err)
	}
	if d.Coverage != nil {
		t.Fatal("invented legacy coverage")
	}
}
