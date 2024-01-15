package saas

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/bearer/bearer/internal/commands/process/gitrepository"
	"github.com/bearer/bearer/internal/commands/process/settings"
	flagtypes "github.com/bearer/bearer/internal/flag/types"
	securitytypes "github.com/bearer/bearer/internal/report/output/security/types"
	"github.com/bearer/bearer/internal/report/output/types"
	util "github.com/bearer/bearer/internal/util/output"
	"github.com/bradleyjkemp/cupaloy"
)

type reportFixture struct {
	Findings map[string][]securitytypes.Finding `json:"findings"`
	Dataflow types.DataFlow                     `json:"dataflow"`
	Files    []string                           `json:"files"`
}

func TestBearerPublishingSaas(t *testing.T) {
	reportData := reportDataFixture(t)
	err := GetReport(
		reportData,
		configFixture(),
		gitContextFixture(),
		true,
	)
	if err != nil {
		t.Fatalf("failed to update struct with saas data, err: %s", err)
	}

	saasOutput, err := util.ReportJSON(reportData.SaasReport)
	if err != nil {
		t.Fatalf("failed to generate JSON output, err: %s", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(saasOutput), "", "\t")
	if err != nil {
		t.Fatalf("error indenting output, err: %s", err)
	}
	cupaloy.SnapshotT(t, prettyJSON.String())
}

func TestBearerPublishingGitlabMetaSaas(t *testing.T) {
	os.Setenv("CI_PIPELINE_ID", "123")
	os.Setenv("CI_JOB_ID", "456")
	defer os.Unsetenv("CI_PIPELINE_ID")
	defer os.Unsetenv("CI_JOB_ID")

	reportData := reportDataFixture(t)
	err := GetReport(
		reportData,
		configFixture(),
		gitContextFixture(),
		true,
	)
	if err != nil {
		t.Fatalf("failed to update struct with saas data, err: %s", err)
	}

	saasOutput, err := util.ReportJSON(reportData.SaasReport.Meta)
	if err != nil {
		t.Fatalf("failed to generate JSON output, err: %s", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(saasOutput), "", "\t")
	if err != nil {
		t.Fatalf("error indenting output, err: %s", err)
	}
	cupaloy.SnapshotT(t, prettyJSON.String())
}

func gitContextFixture() *gitrepository.Context {
	return &gitrepository.Context{
		ID:            "github.com/Bearer/bear-publishing",
		Host:          "github.com",
		Owner:         "Bearer",
		OriginURL:     "git@github.com:Bearer/bear-publishing.git",
		Name:          "bear-publishing",
		FullName:      "Bearer/bear-publishing",
		CommitHash:    "9e54ffa8633898ab65bc4b4e804f7ef24cc068c4",
		Branch:        "main",
		DefaultBranch: "main",
		BaseBranch:    ""}
}

func reportDataFixture(t *testing.T) *types.ReportData {
	reportFixtureOutput, err := os.ReadFile("testdata/report_fixture.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var reportFixture reportFixture
	err = json.Unmarshal(reportFixtureOutput, &reportFixture)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}

	return &types.ReportData{
		FoundLanguages: map[string]int32{
			"CSS":        26,
			"HTML":       413,
			"JavaScript": 16,
			"Markdown":   11,
			"Plain Text": 2,
			"Ruby":       1198,
			"YAML":       59,
		},
		FindingsBySeverity:        reportFixture.Findings,
		IgnoredFindingsBySeverity: make(map[string][]securitytypes.IgnoredFinding, 0),
		Dataflow:                  &reportFixture.Dataflow,
		Files:                     reportFixture.Files,
	}
}

func configFixture() settings.Config {
	return settings.Config{
		Scan: flagtypes.ScanOptions{
			Target: ".",
		},
		BearerRulesVersion: "v0.0.0",
	}
}
