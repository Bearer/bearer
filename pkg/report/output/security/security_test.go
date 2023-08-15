package security_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hhatto/gocloc"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/schema"

	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/report/output/types"
)

func TestBuildReportString(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "security",
		Severity: map[string]bool{
			"critical": true,
			"high":     true,
			"medium":   true,
			"low":      true,
			"warning":  true,
		},
	})
	// set rule version
	config.BearerRulesVersion = "TEST"

	// new rules are added
	customRule := &settings.Rule{
		Id:          "custom_test_rule",
		Description: "Its a test!",
		CWEIDs:      []string{},
		Type:        "risk",
		Languages:   []string{"ruby"},
		Severity:    "low",
		IsLocal:     false,
	}

	// limit rules so that test doesn't fail just because
	config.Rules = map[string]*settings.Rule{
		"ruby_lang_ssl_verification": config.Rules["ruby_lang_ssl_verification"],
		"ruby_rails_logger":          config.Rules["ruby_rails_logger"],
		"custom_test_rule":           customRule,
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflowOutput := dummyDataflowOutput()
	output, err := security.GetOutput(dataflowOutput, config, nil)
	if err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{},
		Languages: map[string]*gocloc.Language{
			"Ruby": {},
		},
		MaxPathLength: 0,
	}

	stringBuilder := security.BuildReportString(config, output, &dummyGoclocResult)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestNoRulesBuildReportString(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "security",
		Severity: map[string]bool{
			"critical": true,
			"high":     true,
			"medium":   true,
			"low":      true,
			"warning":  true,
		},
	})
	// set rule version
	config.BearerRulesVersion = "TEST"
	config.Rules = map[string]*settings.Rule{}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflowOutput := dummyDataflowOutput()
	output, err := security.GetOutput(dataflowOutput, config, nil)
	if err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{},
		Languages: map[string]*gocloc.Language{
			"Ruby": {},
		},
		MaxPathLength: 0,
	}

	stringBuilder := security.BuildReportString(config, output, &dummyGoclocResult)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestGetOutput(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "security",
		Severity: map[string]bool{
			"critical": true,
			"high":     true,
			"medium":   true,
			"low":      true,
			"warning":  true,
		},
	})

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflowOutput := dummyDataflowOutput()
	output, err := security.GetOutput(dataflowOutput, config, nil)
	if err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	assert.Equal(t, dataflowOutput.Dataflow, dataflowOutput.Dataflow)
	assert.Equal(t, dataflowOutput.Files, dataflowOutput.Files)
	cupaloy.SnapshotT(t, output.Data)
}

func TestTestGetOutputWithSeverity(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "security",
		Severity: map[string]bool{
			"critical": true,
			"high":     false,
			"medium":   false,
			"low":      false,
			"warning":  false,
		},
	})

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflowOutput := dummyDataflowOutput()
	output, err := security.GetOutput(dataflowOutput, config, nil)
	if err != nil {
		t.Fatalf("failed to generate security output err: %s", err)
	}

	cupaloy.SnapshotT(t, output.Data)
}

func TestCalculateSeverity(t *testing.T) {
	res := []string{
		security.CalculateSeverity([]string{"PHI", "Personal Data"}, "low", true),
		security.CalculateSeverity([]string{"Personal Data (Sensitive)"}, "low", false),
		security.CalculateSeverity([]string{"Personal Data"}, "low", false),
		security.CalculateSeverity([]string{"Personal Data"}, "warning", false),
		security.CalculateSeverity([]string{}, "warning", false),
	}

	cupaloy.SnapshotT(t, res)
}

func generateConfig(reportOptions flag.ReportOptions) (settings.Config, error) {
	opts := flag.Options{
		ScanOptions: flag.ScanOptions{
			Scanner: []string{"sast"},
		},
		RuleOptions:    flag.RuleOptions{},
		RepoOptions:    flag.RepoOptions{},
		ReportOptions:  reportOptions,
		GeneralOptions: flag.GeneralOptions{},
	}

	return settings.FromOptions(opts, []string{"ruby"})
}

func dummyDataflowOutput() *types.Output[*types.DataFlow] {
	subject := "User"
	lowRisk := dataflowtypes.RiskDetector{
		DetectorID: "ruby_lang_ssl_verification",
		Locations: []dataflowtypes.RiskLocation{
			{
				Filename:        "config/application.rb",
				StartLineNumber: 2,
				Source: &schema.Source{
					StartLineNumber:   2,
					StartColumnNumber: 10,
					EndLineNumber:     2,
					EndColumnNumber:   28,
					Content:           "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
				},
				PresenceMatches: []dataflowtypes.RiskPresence{
					{
						Name: "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
					},
				},
			},
		},
	}

	criticalRisk := dataflowtypes.RiskDetector{
		DetectorID: "ruby_rails_logger",
		Locations: []dataflowtypes.RiskLocation{
			{
				Filename:        "pkg/datatype_leak.rb",
				StartLineNumber: 1,
				Source: &schema.Source{
					StartLineNumber:   1,
					StartColumnNumber: 10,
					EndLineNumber:     2,
					EndColumnNumber:   28,
					Content:           "Rails.logger.info(user.biometric_data)",
				},
				DataTypes: []dataflowtypes.RiskDatatype{
					{
						Name:         "Biometric Data",
						CategoryUUID: "35b94efa-9b67-49b2-abb9-29b6a759a030",
						Schemas: []dataflowtypes.RiskSchema{
							{
								FieldName:   "",
								ObjectName:  "",
								SubjectName: &subject,
							},
						},
					},
				},
			},
		},
	}

	// build risk []interface
	risks := make([]dataflowtypes.RiskDetector, 2)
	risks[0] = criticalRisk
	risks[1] = lowRisk

	dataflow := &types.DataFlow{
		Datatypes: []dataflowtypes.Datatype{
			{
				Name:         "Email Address",
				UUID:         "02bb0d3a-2c8c-4842-be1c-c057f0dccd63",
				CategoryUUID: "dd88aee5-9d40-4ad2-8983-0c791ddec47c",
				Detectors: []dataflowtypes.DatatypeDetector{
					{
						Name: "ruby",
						Locations: []dataflowtypes.DatatypeLocation{
							{
								Filename:          "app/model/user.rb",
								StartLineNumber:   1,
								StartColumnNumber: 10,
								EndColumnNumber:   17,
								FieldName:         "email",
								ObjectName:        "user",
								SubjectName:       &subject,
							},
						},
					},
				},
			},
		},
		Risks:      risks,
		Components: []dataflowtypes.Component{},
	}

	return &types.Output[*types.DataFlow]{
		Data:     dataflow,
		Dataflow: dataflow,
		Files:    []string{"config/application.rb", "pkg/datatype_leak.rb", "app/model/user.rb"},
	}
}
