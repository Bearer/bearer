package sarif

import (
	"github.com/bearer/bearer/pkg/commands/process/settings"
	sarif "github.com/bearer/bearer/pkg/report/output/sarif/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
)

func ReportSarif(outputDetections map[string][]securitytypes.Finding, rules map[string]*settings.Rule) (sarif.SarifOutput, error) {
	var sarifRules []sarif.Rule

	for _, rule := range rules {
		if !rule.PolicyType() {
			continue
		}

		sarifRules = append(sarifRules, sarif.Rule{
			Id:   rule.Id,
			Name: rule.Id,
			ShortDescription: sarif.Description{
				Text: rule.Description,
			},
			FullDescription: sarif.Description{
				Text: rule.Description,
			},
			Help: sarif.Help{
				Text:     rule.RemediationMessage,
				Markdown: rule.RemediationMessage,
			},
			DefaultConfiguration: sarif.Configuration{
				Level: "error", // rule.Severity, accepted values are ("none", "note", "warning", "error")
			},
			// Properties: sarif.Properties{
			// 		Tags:      []string{"maintainability"},
			// 		Precision: "very-high",
			// },
		})
	}

	var results []sarif.Result

	for _, level := range []string{"critical", "high", "medium", "low", "warning"} {
		if findings, ok := outputDetections[level]; ok {
			for _, finding := range findings {
				results = append(results, sarif.Result{
					RuleId: finding.Id,
					Message: sarif.Message{
						Text: finding.Title,
					},
					Locations: []sarif.Location{
						{
							PhysicalLocation: sarif.PhysicalLocation{
								ArtifactLocation: sarif.ArtifactLocation{
									URI: finding.Filename,
								},
								Region: sarif.Region{
									StartLine:   finding.Sink.Start,
									EndLine:     finding.Sink.End,
									StartColumn: finding.Sink.Column.Start,
									EndColumn:   finding.Sink.Column.End,
								},
							},
						},
					},
					PartialFingerprints: &sarif.PartialFingerprints{
						PrimaryLocationLineHash: finding.Fingerprint,
					},
				})
			}
		}
	}

	output := sarif.SarifOutput{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarif.Run{
			{
				Tool: sarif.Tool{
					Driver: sarif.Driver{
						Name:  "Bearer",
						Rules: sarifRules,
					},
				},
				Results: results,
			},
		},
	}

	return output, nil
}
