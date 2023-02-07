package summary_test

import (
	"testing"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	"github.com/bearer/curio/pkg/report/output/summary"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bradleyjkemp/cupaloy"
)

func TestSummary(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "summary",
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

	dataflow := dummyDataflow()

	res, err := summary.GetOutput(&dataflow, config)
	if err != nil {
		t.Fatalf("failed to generate summary output err:%s", err)
	}

	cupaloy.SnapshotT(t, res)
}

func TestSummaryWithSeverity(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{
		Report: "summary",
		Severity: map[string]bool{
			"critical": true,
			"high":     true,
			"medium":   false,
			"low":      false,
			"warning":  false,
		},
	})

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	dataflow := dummyDataflow()

	res, err := summary.GetOutput(&dataflow, config)
	if err != nil {
		t.Fatalf("failed to generate summary output err:%s", err)
	}

	cupaloy.SnapshotT(t, res)
}

func generateConfig(reportOptions flag.ReportOptions) (settings.Config, error) {
	opts := flag.Options{
		ScanOptions:    flag.ScanOptions{},
		RuleOptions:    flag.RuleOptions{},
		RepoOptions:    flag.RepoOptions{},
		ReportOptions:  reportOptions,
		GeneralOptions: flag.GeneralOptions{},
	}

	return settings.FromOptions(opts)
}

func dummyDataflow() dataflow.DataFlow {
	subject := "User"
	loggerRisk := types.RiskDetector{
		DetectorID: "ruby_rails_logger",
		DataTypes: []types.RiskDatatype{
			{
				Name:         "Email Address",
				Stored:       false,
				UUID:         "22e24c62-82d3-4b72-827c-e261533331bd",
				CategoryUUID: "cef587dd-76db-430b-9e18-7b031e1a193b",
				Locations: []types.RiskLocation{
					{
						Filename:    "pkg/datatype_leak.rb",
						LineNumber:  1,
						FieldName:   "email",
						ObjectName:  "user",
						SubjectName: &subject,
						Parent: &schema.Parent{
							LineNumber: 1,
							Content:    "Rails.logger.info(user.email)",
						},
					},
				},
			},
			{
				Name:         "Browsing Behavior",
				Stored:       false,
				UUID:         "c73ae276-b1b1-4b70-b6d5-ed73a83e87ed",
				CategoryUUID: "8099225c-7e49-414f-aac2-e7045379bb40",
				Locations: []types.RiskLocation{
					{
						Filename:    "pkg/datatype_leak.rb",
						LineNumber:  2,
						FieldName:   "browsing_behavior",
						ObjectName:  "user",
						SubjectName: &subject,
						Parent: &schema.Parent{
							LineNumber: 2,
							Content:    "Rails.logger.info(user.browsing_behavior)",
						},
					},
				},
			},
		},
	}

	risks := make([]interface{}, 1)
	risks[0] = loggerRisk

	return dataflow.DataFlow{
		Datatypes: []types.Datatype{
			{
				Name:         "Email Address",
				UUID:         "22e24c62-82d3-4b72-827c-e261533331bd",
				CategoryUUID: "cef587dd-76db-430b-9e18-7b031e1a193b",
				Detectors: []types.DatatypeDetector{
					{
						Name: "ruby",
						Locations: []types.DatatypeLocation{
							{
								Filename:    "pkg/datatype_leak.rb",
								LineNumber:  1,
								FieldName:   "email",
								ObjectName:  "user",
								SubjectName: &subject,
							},
						},
					},
				},
			},
			{
				Name:         "Browsing Behavior",
				UUID:         "c73ae276-b1b1-4b70-b6d5-ed73a83e87ed",
				CategoryUUID: "8099225c-7e49-414f-aac2-e7045379bb40",
				Detectors: []types.DatatypeDetector{
					{
						Name: "ruby",
						Locations: []types.DatatypeLocation{
							{
								Filename:    "pkg/datatype_leak.rb",
								LineNumber:  2,
								FieldName:   "browsing_behavior",
								ObjectName:  "user",
								SubjectName: &subject,
							},
						},
					},
				},
			},
		},
		Risks:      risks,
		Components: []types.Component{},
	}
}
