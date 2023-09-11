package sarif_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/sarif"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	util "github.com/bearer/bearer/pkg/util/output"
)

func TestJuiceShopSarif(t *testing.T) {
	securityOutput, err := os.ReadFile("testdata/juice-shop-security-report.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var securityResults map[string][]securitytypes.Finding
	err = json.Unmarshal(securityOutput, &securityResults)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}
	var rules = make(map[string]*settings.Rule)
	rules["rule_1"] = &settings.Rule{
		Id:                 "rule_1",
		AssociatedRecipe:   "",
		Type:               "risk",
		Trigger:            settings.RuleTrigger{},
		IsLocal:            false,
		Detectors:          []string{},
		Processors:         []string{},
		Stored:             false,
		AutoEncrytPrefix:   "",
		HasDetailedContext: false,
		SkipDataTypes:      []string{},
		OnlyDataTypes:      []string{},
		Severity:           "high",
		Description:        "rule 1",
		RemediationMessage: "## Rule 1\nremediation message",
		CWEIDs:             []string{"cwe-10"},
		Languages:          []string{"ruby"},
		Patterns:           []settings.RulePattern{},
		SanitizerRuleID:    "",
		DocumentationUrl:   "",
		IsAuxilary:         false,
		Metavars:           map[string]settings.MetaVar{},
		ParamParenting:     false,
		DetectPresence:     false,
		OmitParent:         false,
	}

	res, err := sarif.ReportSarif(securityResults, rules)
	if err != nil {
		t.Fatalf("failed to generate security output, err: %s", err)
	}

	sarifOutput, err := util.ReportJSON(res)
	if err != nil {
		t.Fatalf("failed to generate JSON output, err: %s", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(sarifOutput), "", "\t")
	if err != nil {
		t.Fatalf("error indenting output, err: %s", err)
	}
	cupaloy.SnapshotT(t, prettyJSON.String())
}
