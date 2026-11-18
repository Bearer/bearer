package security

import (
	"encoding/json"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/coverage"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	"testing"
)

func TestCoverageJSONV2Envelope(t *testing.T) {
	r := &coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest(nil), ConfigurationSHA256: coverage.Digest(nil), Discovery: coverage.Discovery{Files: []string{}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true}, Outcomes: []coverage.Outcome{}, Complete: true}
	f := Formatter{ReportData: &outputtypes.ReportData{Dataflow: &outputtypes.DataFlow{}, Coverage: r}}
	for _, failed := range []bool{false, true} {
		f.ReportData.ReportFailed = failed
		text, err := f.Format(flag.FormatJSONV2)
		if err != nil {
			t.Fatal(err)
		}
		var out JsonV2Output
		if err = json.Unmarshal([]byte(text), &out); err != nil {
			t.Fatal(err)
		}
		if out.Coverage == nil || out.Coverage.Validate([]string{}) != nil {
			t.Fatal("same-run coverage lost")
		}
		if out.Findings != nil || out.Errors != nil {
			t.Fatal("clean null output changed")
		}
	}
	f.ReportData.Coverage = nil
	text, err := f.Format(flag.FormatJSONV2)
	if err != nil {
		t.Fatal(err)
	}
	var out JsonV2Output
	if err = json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatal(err)
	}
	if out.Coverage != nil {
		t.Fatal("legacy output invented coverage")
	}
}
