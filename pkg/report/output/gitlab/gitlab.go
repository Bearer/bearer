package gitlab

import (
	"strings"
	"time"

	"github.com/bearer/bearer/cmd/bearer/build"
	gitlab "github.com/bearer/bearer/pkg/report/output/gitlab/types"
	"github.com/bearer/bearer/pkg/report/output/security"
)

func ReportGitLab(
	outputDetections *map[string][]security.Result,
	startTime time.Time,
	endTime time.Time,
) (gitlab.GitLabOutput, error) {
	var vulnerabilities []gitlab.Vulnerability
	for _, level := range []string{"critical", "high", "medium", "low", "warning"} {
		if findings, ok := (*outputDetections)[level]; ok {
			for _, finding := range findings {
				vulnerabilities = append(vulnerabilities, gitlab.Vulnerability{
					Id:          finding.Fingerprint,
					Category:    "sast",
					Name:        finding.Rule.Title,
					Message:     finding.Description,
					Description: finding.Description,
					// CVE:                  "",
					Severity:             formatSeverity(level), // level,
					Confidence:           "Unknown",
					RawSourceCodeExtract: finding.Source.Content,
					Scanner: gitlab.VulnerabilityScanner{
						Id:   "bearer",
						Name: "Bearer",
					},
					Location: gitlab.Location{
						File:      finding.Filename,
						Startline: finding.Source.Start,
					},
					Identifiers: []gitlab.Identifier{
						{
							Type:  finding.Rule.Id,
							Name:  finding.Rule.Title,
							Value: finding.Rule.Title,
						},
					},
				})
			}
		}
	}

	output := gitlab.GitLabOutput{
		Schema:          "https://gitlab.com/gitlab-org/security-products/security-report-schemas/-/raw/master/dist/sast-report-format.json",
		Version:         "15.0.4",
		Vulnerabilities: vulnerabilities,
		Scan: gitlab.Scan{
			Analyzer: gitlab.Analyzer{
				Id:   "bearer-sast",
				Name: "Bearer SAST",
				URL:  "https://github.com/bearer/bearer",
				Vendor: gitlab.Vendor{
					Name: "Bearer",
				},
				Version: build.Version,
			},
			Scanner: gitlab.Scanner{
				Id:   "bearer",
				Name: "Bearer",
				URL:  "https://github.com/bearer/bearer",
				Vendor: gitlab.Vendor{
					Name: "Bearer",
				},
				Version: build.Version,
			},
			Type:      "sast",
			StartTime: startTime.Format("2006-01-02T15:04:05"),
			EndTime:   endTime.Format("2006-01-02T15:04:05"),
			Status:    calculateStatus(vulnerabilities),
		},
	}

	return output, nil
}

func calculateStatus(vulnerabilities []gitlab.Vulnerability) string {
	if len(vulnerabilities) > 0 {
		return "failure"
	} else {
		return "success"
	}
}

func formatSeverity(level string) string {
	switch level {
	case "warning":
		return "Unknown"
	default:
		return strings.ToUpper(level[:1]) + level[1:]
	}
}
