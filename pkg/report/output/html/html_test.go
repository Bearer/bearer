package html

import (
	"encoding/json"
	"os"
	"testing"

	privacy "github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bradleyjkemp/cupaloy"
)

func TestJuiceShopSecurityHtml(t *testing.T) {
	securityOutput, err := os.ReadFile("testdata/juice-shop-security-report.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var securityResults map[string][]security.Result
	err = json.Unmarshal(securityOutput, &securityResults)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}

	output, err := ReportSecurityHTML(&securityResults)
	if err != nil {
		t.Fatalf("failed to generate security output, err: %s", err)
	}

	snapshotter := cupaloy.New(cupaloy.SnapshotFileExtension(".html"))
	snapshotter.SnapshotT(t, []byte(*output))
}
func TestBearPublishingPrivacyHtml(t *testing.T) {
	privacyOutput, err := os.ReadFile("testdata/bear-publishing-privacy-report.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var privacyResults privacy.Report
	err = json.Unmarshal(privacyOutput, &privacyResults)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}

	output, err := ReportPrivacyHTML(&privacyResults)
	if err != nil {
		t.Fatalf("failed to generate security output, err: %s", err)
	}
	snapshotter := cupaloy.New(cupaloy.SnapshotFileExtension(".html"))
	snapshotter.SnapshotT(t, []byte(*output))
}
