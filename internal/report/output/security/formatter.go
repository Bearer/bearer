package security

import (
	"fmt"
	"time"

	"github.com/hhatto/gocloc"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/report/output/gitlab"
	"github.com/bearer/bearer/internal/report/output/html"
	"github.com/bearer/bearer/internal/report/output/reviewdog"
	"github.com/bearer/bearer/internal/report/output/sarif"
	outputtypes "github.com/bearer/bearer/internal/report/output/types"
	outputhandler "github.com/bearer/bearer/internal/util/output"
)

type Formatter struct {
	ReportData   *outputtypes.ReportData
	Config       settings.Config
	GoclocResult *gocloc.Result
	StartTime    time.Time
	EndTime      time.Time
}

func NewFormatter(reportData *outputtypes.ReportData, config settings.Config, goclocResult *gocloc.Result, startTime time.Time, endTime time.Time) *Formatter {
	return &Formatter{
		ReportData:   reportData,
		Config:       config,
		GoclocResult: goclocResult,
		StartTime:    startTime,
		EndTime:      endTime,
	}
}

func (f Formatter) Format(format string) (output string, err error) {
	switch format {
	case flag.FormatEmpty:
		output = BuildReportString(f.ReportData, f.Config, f.GoclocResult).String()
	case flag.FormatSarif:
		sarifContent, sarifErr := sarif.ReportSarif(f.ReportData.FindingsBySeverity, f.Config.Rules)
		if sarifErr != nil {
			return output, fmt.Errorf("error generating sarif report %s", sarifErr)
		}
		return outputhandler.ReportJSON(sarifContent)
	case flag.FormatReviewDog:
		sastContent, reviewdogErr := reviewdog.ReportReviewdog(f.ReportData.FindingsBySeverity)
		if reviewdogErr != nil {
			return output, fmt.Errorf("error generating reviewdog report %s", reviewdogErr)
		}
		return outputhandler.ReportJSON(sastContent)
	case flag.FormatGitLabSast:
		sastContent, sastErr := gitlab.ReportGitLab(f.ReportData.FindingsBySeverity, f.StartTime, f.EndTime)
		if sastErr != nil {
			return output, fmt.Errorf("error generating gitlab-sast report %s", sastErr)
		}
		return outputhandler.ReportJSON(sastContent)
	case flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.FindingsBySeverity)
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.FindingsBySeverity)
	case flag.FormatHTML:
		title := "Security Report"
		body, securityErr := html.ReportSecurityHTML(f.ReportData.FindingsBySeverity)
		if securityErr != nil {
			return output, securityErr
		}

		output, err = html.ReportHTMLWrapper(title, body)
		if err != nil {
			err = fmt.Errorf("could not generate html page %s", err)
		}
	}

	return output, err
}
