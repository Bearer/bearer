package reviewdog

import (
	reviewdog "github.com/bearer/bearer/pkg/report/output/reviewdog/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
)

func ReportReviewdog(outputDetections map[string][]securitytypes.Finding) (reviewdog.ReviewdogOutput, error) {
	var reviewdogDiagnostics []reviewdog.Diagnostic

	for _, level := range []string{"critical", "high", "medium", "low", "warning"} {
		if findings, ok := outputDetections[level]; ok {
			for _, finding := range findings {
				var severity string
				if level == "warning" {
					severity = "WARNING"
				} else {
					severity = "ERROR"
				}

				message := "\n# " + finding.Title + "\n" + finding.Description

				reviewdogDiagnostics = append(reviewdogDiagnostics, reviewdog.Diagnostic{
					Message:  message,
					Severity: severity,
					Location: reviewdog.Location{
						Path: finding.Filename,
						Range: reviewdog.LocationRange{
							Start: reviewdog.LocationPosition{
								Line:   finding.Sink.Start,
								Column: finding.Sink.Column.Start,
							},
							End: reviewdog.LocationPosition{
								Line:   finding.Sink.End,
								Column: finding.Sink.Column.End,
							},
						},
					},
					Code: reviewdog.Code{
						RuleId:           finding.Id,
						DocumentationUrl: finding.DocumentationUrl,
					},
					Suggestions: []reviewdog.Suggestion{},
				})
			}
		}
	}

	output := reviewdog.ReviewdogOutput{
		Source: reviewdog.Source{
			Name: "Bearer",
			Url:  "https://docs.bearer.com/",
		},
		Diagnostics: reviewdogDiagnostics,
	}

	return output, nil
}
