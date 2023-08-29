package reviewdog_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/report/output/reviewdog"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	"github.com/bearer/bearer/pkg/util/output"
)

func TestRailsGoatReviewdog(t *testing.T) {
	securityOutput, err := os.ReadFile("testdata/rails-goat-security-report.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var securityFindings map[string][]securitytypes.Finding
	err = json.Unmarshal(securityOutput, &securityFindings)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}

	res, err := reviewdog.ReportReviewdog(securityFindings)
	if err != nil {
		t.Fatalf("failed to generate security output, err: %s", err)
	}

	sarifOutput, err := output.ReportJSON(res)
	if err != nil {
		t.Fatalf("failed to generate JSON output, err: %s", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(*sarifOutput), "", "\t")
	if err != nil {
		t.Fatalf("error indenting output, err: %s", err)
	}
	cupaloy.SnapshotT(t, prettyJSON.String())
}
