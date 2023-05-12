package privacy_test

import (
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bradleyjkemp/cupaloy"
)

func TestBuildCsvString(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "privacy"})
	config.Rules = map[string]*settings.Rule{
		"ruby_third_parties_sentry": config.Rules["ruby_third_parties_sentry"],
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflow := dummyDataflow()

	stringBuilder, _ := privacy.BuildCsvString(&dataflow, config)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestGetOutput(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "privacy"})
	config.Rules = map[string]*settings.Rule{
		"ruby_third_parties_sentry": config.Rules["ruby_third_parties_sentry"],
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflow := dummyDataflow()
	results, _, err := privacy.GetOutput(&dataflow, config)
	if err != nil {
		t.Fatalf("failed to generate privacy output err:%s", err)
	}

	cupaloy.SnapshotT(t, results)
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

func dummyDataflow() dataflow.DataFlow {
	subject := "User"
	thirdPartyRisk := types.RiskDetector{
		DetectorID: "ruby_third_parties_sentry",
		Locations: []types.RiskLocation{
			{
				Filename:        "/app/controllers/application_controller.rb",
				StartLineNumber: 39,
				Parent: &schema.Parent{
					StartLineNumber:   38,
					StartColumnNumber: 10,
					EndLineNumber:     38,
					EndColumnNumber:   28,
					Content:           "Sentry.set_user(email: current_user.email)",
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

	return dataflow.DataFlow{
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
