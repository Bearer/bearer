package gitlab

import (
	"fmt"
	"strings"
	"time"

	"github.com/bearer/bearer/cmd/bearer/build"
	gitlab "github.com/bearer/bearer/pkg/report/output/gitlab/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
)

func ReportGitLab(
	outputDetections map[string][]securitytypes.Finding,
	startTime time.Time,
	endTime time.Time,
) (gitlab.GitLabOutput, error) {
	var vulnerabilities []gitlab.Vulnerability
	for _, level := range []string{"critical", "high", "medium", "low", "warning"} {
		if findings, ok := outputDetections[level]; ok {
			for _, finding := range findings {
				identifiers := []gitlab.Identifier{
					{
						Type:  "bearer",
						Name:  finding.Id,
						Value: finding.Id,
						Url:   finding.DocumentationUrl,
					},
				}
				for _, cwe := range finding.CWEIDs {
					identifiers = append(identifiers, gitlab.Identifier{
						Type:  "cwe",
						Name:  "CWE-" + cwe,
						Value: cwe,
						Url:   fmt.Sprintf("https://cwe.mitre.org/data/definitions/%s.html", cwe),
					})
				}

				vulnerabilities = append(vulnerabilities, gitlab.Vulnerability{
					Id:                   finding.Fingerprint,
					Category:             "sast",
					Name:                 finding.Title,
					Description:          extractDescription(finding.Description),
					Solution:             extractSolution(finding.Description),
					Severity:             formatSeverity(level), // level,
					Confidence:           "Unknown",
					RawSourceCodeExtract: finding.Sink.Content,
					Scanner: gitlab.VulnerabilityScanner{
						Id:   "bearer",
						Name: "Bearer",
					},
					Location: gitlab.Location{
						File:      finding.Filename,
						Startline: finding.Sink.Start,
						Endline:   finding.Sink.End,
					},
					Identifiers: identifiers,
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
		return "Info"
	default:
		return strings.ToUpper(level[:1]) + level[1:]
	}
}

func extractDescription(body string) string {
	split := strings.Split(body, "## Remediations")
	if len(split) < 2 {
		return strings.ReplaceAll(body, "## ", "")
	}
	return strings.Replace(split[0], "## Description\n", "", 1)
}

func extractSolution(body string) string {
	split := strings.Split(body, "## Remediations")
	if len(split) < 2 {
		return ""
	}
	return strings.Replace(split[len(split)-1], "## Resources\n", "Resources:\n", 1)
}
