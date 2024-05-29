package privacy_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	settingsloader "github.com/bearer/bearer/pkg/commands/process/settings/loader"
	"github.com/bearer/bearer/pkg/engine"
	engineimpl "github.com/bearer/bearer/pkg/engine/implementation"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/languages"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/testhelper"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/version_check"
)

func TestBuildCsvString(t *testing.T) {
	engine := engineimpl.New(languages.Default())
	config, err := generateConfig(engine, flagtypes.ReportOptions{Report: "privacy"})
	config.Rules = map[string]*settings.Rule{
		"ruby_third_parties_sentry": testhelper.RubyThirdPartiesSentryRule(),
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	output := &outputtypes.ReportData{
		Dataflow: dummyDataflow(),
	}
	err = privacy.AddReportData(output, config)
	if err != nil {
		t.Fatalf("failed to add privacy report:%s", err)
	}
	stringBuilder, _ := privacy.BuildCsvString(output, config)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestAddReportData(t *testing.T) {
	engine := engineimpl.New(languages.Default())
	config, err := generateConfig(engine, flagtypes.ReportOptions{Report: "privacy"})
	config.Rules = map[string]*settings.Rule{
		"ruby_third_parties_sentry": testhelper.RubyThirdPartiesSentryRule(),
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	output := &outputtypes.ReportData{
		Dataflow: dummyDataflow(),
	}
	if err = privacy.AddReportData(output, config); err != nil {
		t.Fatalf("failed to generate privacy output err:%s", err)
	}

	cupaloy.SnapshotT(t, output.PrivacyReport)
}

func generateConfig(engine engine.Engine, reportOptions flagtypes.ReportOptions) (settings.Config, error) {
	opts := flagtypes.Options{
		ScanOptions: flagtypes.ScanOptions{
			Scanner: []string{"sast"},
		},
		RuleOptions:    flagtypes.RuleOptions{},
		ReportOptions:  reportOptions,
		GeneralOptions: flagtypes.GeneralOptions{},
	}

	meta := &version_check.VersionMeta{
		Rules: version_check.RuleVersionMeta{
			Packages: make(map[string]string),
		},
		Binary: version_check.BinaryVersionMeta{
			Latest:  true,
			Message: "",
		},
	}

	return settingsloader.FromOptions(opts, meta, engine)
}

func dummyDataflow() *outputtypes.DataFlow {
	subject := "User"
	thirdPartyRisk := types.RiskDetector{
		DetectorID: "ruby_third_parties_sentry",
		Locations: []types.RiskLocation{
			{
				Filename:        "/app/controllers/application_controller.rb",
				StartLineNumber: 39,
				Source: &schema.Source{
					StartLineNumber:   38,
					StartColumnNumber: 10,
					EndLineNumber:     38,
					EndColumnNumber:   28,
				},
				DataTypes: []types.RiskDatatype{
					{
						Name:         "Email Address",
						CategoryUUID: "cef587dd-76db-430b-9e18-7b031e1a193b",
						Schemas: []types.RiskSchema{
							{
								FieldName:   "email",
								ObjectName:  "current_user",
								SubjectName: &subject,
							},
						},
					},
				},
			},
		},
	}

	// build risk []interface
	risks := make([]types.RiskDetector, 1)
	risks[0] = thirdPartyRisk

	return &outputtypes.DataFlow{
		Datatypes: []types.Datatype{
			{
				Name:         "Email Address",
				CategoryName: "Contact",
				Detectors: []types.DatatypeDetector{
					{
						Name: "ruby",
						Locations: []types.DatatypeLocation{
							{
								Filename:          "/app/controllers/application_controller.rb",
								StartLineNumber:   39,
								StartColumnNumber: 10,
								EndColumnNumber:   17,
								FieldName:         "email",
								ObjectName:        "current_user",
								SubjectName:       &subject,
							},
						},
					},
				},
			},
			{
				Name:         "Country",
				CategoryName: "Location",
				Detectors: []types.DatatypeDetector{
					{
						Name: "ruby",
						Locations: []types.DatatypeLocation{
							{
								Filename:          "/app/models/location.rb",
								StartLineNumber:   112,
								StartColumnNumber: 10,
								EndColumnNumber:   17,
								FieldName:         "country",
								ObjectName:        "Address",
							},
						},
					},
				},
			},
		},
		Risks: risks,
		Components: []types.Component{
			{
				Name:    "Sentry",
				Type:    "external_service",
				SubType: "third_party",
				Locations: []types.ComponentLocation{
					{
						Detector:   "gemfile-lock",
						Filename:   "/Users/elsapet/Bearer/bear-publishing/Gemfile.lock",
						LineNumber: 204,
					},
				},
			},
		},
	}
}
