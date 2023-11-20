package security_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hhatto/gocloc"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/report/basebranchfindings"
	basebranchfindingstypes "github.com/bearer/bearer/internal/report/basebranchfindings/types"
	"github.com/bearer/bearer/internal/report/schema"
	globaltypes "github.com/bearer/bearer/internal/types"
	"github.com/bearer/bearer/internal/util/set"
	"github.com/bearer/bearer/internal/version_check"

	dataflowtypes "github.com/bearer/bearer/internal/report/output/dataflow/types"
	"github.com/bearer/bearer/internal/report/output/security"
	securitytypes "github.com/bearer/bearer/internal/report/output/security/types"
	"github.com/bearer/bearer/internal/report/output/testhelper"
	outputtypes "github.com/bearer/bearer/internal/report/output/types"
)

func TestBuildReportString(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "security"})
	// set rule version
	config.BearerRulesVersion = "TEST"

	config.Rules = map[string]*settings.Rule{
		"ruby_lang_ssl_verification": testhelper.RubyLangSSLVerificationRule(),
		"ruby_rails_logger":          testhelper.RubyRailsLoggerRule(),
		"custom_test_rule":           testhelper.CustomRule(),
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	data := dummyDataflowData()
	if err := security.AddReportData(data, config, nil, true); err != nil {
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

	stringBuilder := security.BuildReportString(data, config, &dummyGoclocResult)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestNoRulesBuildReportString(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "security"})
	// set rule version
	config.BearerRulesVersion = "TEST"
	config.Rules = map[string]*settings.Rule{}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	output := dummyDataflowData()
	if err := security.AddReportData(output, config, nil, true); err != nil {
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

	stringBuilder := security.BuildReportString(output, config, &dummyGoclocResult)
	cupaloy.SnapshotT(t, stringBuilder.String())
}

func TestAddReportData(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "security"})

	config.Rules = map[string]*settings.Rule{
		"ruby_lang_ssl_verification": testhelper.RubyLangSSLVerificationRule(),
		"ruby_rails_logger":          testhelper.RubyRailsLoggerRule(),
	}

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	output := dummyDataflowData()
	if err = security.AddReportData(output, config, nil, true); err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	assert.Equal(t, output.Dataflow, output.Dataflow)
	assert.Equal(t, output.Files, output.Files)
	cupaloy.SnapshotT(t, output.FindingsBySeverity)
}

func TestAddReportDataWithSeverity(t *testing.T) {
	severity := set.New[string]()
	severity.Add(globaltypes.LevelCritical)

	config, err := generateConfig(flag.ReportOptions{
		Report:   "security",
		Severity: severity,
	})

	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	config.Rules = map[string]*settings.Rule{
		"ruby_rails_logger": testhelper.RubyRailsLoggerRule(),
	}

	data := dummyDataflowData()
	if err = security.AddReportData(data, config, nil, true); err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	cupaloy.SnapshotT(t, data.FindingsBySeverity)
}

func TestAddReportDataWithFailOnSeverity(t *testing.T) {
	for _, test := range []struct {
		FailOnSeverity,
		Severity string
		Expected bool
	}{
		{FailOnSeverity: globaltypes.LevelCritical, Expected: true},
		{FailOnSeverity: globaltypes.LevelHigh, Expected: true},
		{FailOnSeverity: globaltypes.LevelHigh, Severity: globaltypes.LevelCritical, Expected: false},
		{FailOnSeverity: globaltypes.LevelMedium, Expected: false},
		{FailOnSeverity: globaltypes.LevelLow, Expected: false},
		{FailOnSeverity: globaltypes.LevelWarning, Expected: false},
	} {
		t.Run(test.FailOnSeverity, func(tt *testing.T) {
			failOnSeverity := set.New[string]()
			failOnSeverity.Add(test.FailOnSeverity)

			var severity set.Set[string]
			if test.Severity != "" {
				severity = set.New[string]()
				severity.Add(test.Severity)
			}

			config, err := generateConfig(flag.ReportOptions{
				Report:         "security",
				Severity:       severity,
				FailOnSeverity: failOnSeverity,
			})

			if err != nil {
				tt.Fatalf("failed to generate config:%s", err)
			}

			config.Rules = map[string]*settings.Rule{
				"ruby_rails_logger":          testhelper.RubyRailsLoggerRule(),
				"ruby_lang_ssl_verification": testhelper.RubyLangSSLVerificationRule(),
			}

			data := dummyDataflowData()
			if err = security.AddReportData(data, config, nil, true); err != nil {
				tt.Fatalf("failed to generate security output err:%s", err)
			}

			assert.Equal(tt, test.Expected, data.ReportFailed)
		})
	}
}

func TestCalculateSeverity(t *testing.T) {
	res := []securitytypes.SeverityMeta{
		security.CalculateSeverity([]string{"PHI", "Personal Data"}, "low", true),
		security.CalculateSeverity([]string{"Personal Data (Sensitive)"}, "low", false),
		security.CalculateSeverity([]string{"Personal Data"}, "low", false),
		security.CalculateSeverity([]string{"Personal Data"}, "warning", false),
		security.CalculateSeverity([]string{}, "warning", false),
	}

	cupaloy.SnapshotT(t, res)
}

func TestFingerprintIsStableWithBaseBranchFindings(t *testing.T) {
	config, err := generateConfig(flag.ReportOptions{Report: "security"})
	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	config.Rules = map[string]*settings.Rule{
		"ruby_lang_ssl_verification": testhelper.RubyLangSSLVerificationRule(),
	}

	filename := "config/application.rb"

	data := &outputtypes.ReportData{
		Dataflow: &outputtypes.DataFlow{
			Risks: []dataflowtypes.RiskDetector{
				{
					DetectorID: "ruby_lang_ssl_verification",
					Locations: []dataflowtypes.RiskLocation{
						{
							Filename:        filename,
							StartLineNumber: 1,
							Source: &schema.Source{
								StartLineNumber:   1,
								StartColumnNumber: 1,
								EndLineNumber:     1,
								EndColumnNumber:   44,
								Content:           "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
							},
							PresenceMatches: []dataflowtypes.RiskPresence{
								{
									Name: "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
								},
							},
						},
					},
				},
				{
					DetectorID: "ruby_lang_ssl_verification",
					Locations: []dataflowtypes.RiskLocation{
						{
							Filename:        filename,
							StartLineNumber: 2,
							Source: &schema.Source{
								StartLineNumber:   2,
								StartColumnNumber: 1,
								EndLineNumber:     2,
								EndColumnNumber:   44,
								Content:           "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
							},
							PresenceMatches: []dataflowtypes.RiskPresence{
								{
									Name: "http.verify_mode = OpenSSL::SSL::VERIFY_NONE",
								},
							},
						},
					},
				},
			},
		},
		Files: []string{filename},
	}

	if err = security.AddReportData(data, config, nil, true); err != nil {
		t.Fatalf("failed to generate security output err:%s", err)
	}

	fullScanFinding := data.FindingsBySeverity[globaltypes.LevelMedium][1]

	chunks := basebranchfindings.NewChunks()
	chunks.Add(basebranchfindingstypes.ChunkEqual, 1)
	chunks.Add(basebranchfindingstypes.ChunkAdd, 1)

	file := files.File{FilePath: filename}
	fileList := &files.List{
		Files:     []files.File{file},
		BaseFiles: []files.File{file},
		Chunks: map[string]basebranchfindingstypes.Chunks{
			filename: chunks,
		},
	}

	baseBranchFindings := basebranchfindings.New(fileList)
	baseBranchFindings.Add("ruby_lang_ssl_verification", filename, 1, 1)

	if err = security.AddReportData(data, config, baseBranchFindings, true); err != nil {
		t.Fatalf("failed to generate security output with base branch findings err:%s", err)
	}

	diffFinding := data.FindingsBySeverity[globaltypes.LevelMedium][0]

	assert.Equal(t, fullScanFinding.LineNumber, diffFinding.LineNumber)
	assert.Equal(t, fullScanFinding.Fingerprint, diffFinding.Fingerprint)
}

func generateConfig(reportOptions flag.ReportOptions) (settings.Config, error) {
	if reportOptions.Severity == nil {
		reportOptions.Severity = set.New[string]()
		reportOptions.Severity.AddAll(globaltypes.Severities)
	}

	if reportOptions.FailOnSeverity == nil {
		reportOptions.FailOnSeverity = set.New[string]()
		reportOptions.FailOnSeverity.Add(globaltypes.LevelCritical)
		reportOptions.FailOnSeverity.Add(globaltypes.LevelHigh)
		reportOptions.FailOnSeverity.Add(globaltypes.LevelMedium)
		reportOptions.FailOnSeverity.Add(globaltypes.LevelLow)
	}

	opts := flag.Options{
		ScanOptions: flag.ScanOptions{
			Scanner: []string{"sast"},
		},
		RuleOptions:    flag.RuleOptions{},
		ReportOptions:  reportOptions,
		GeneralOptions: flag.GeneralOptions{},
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
	return settings.FromOptions(opts, meta)
}

func dummyDataflowData() *outputtypes.ReportData {
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

	dataflow := &outputtypes.DataFlow{
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

	return &outputtypes.ReportData{
		Dataflow: dataflow,
		Files:    []string{"config/application.rb", "pkg/datatype_leak.rb", "app/model/user.rb"},
	}
}
