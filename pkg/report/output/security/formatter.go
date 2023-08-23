package security

import (
	"fmt"
	"time"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/gitlab"
	"github.com/bearer/bearer/pkg/report/output/html"
	"github.com/bearer/bearer/pkg/report/output/reviewdog"
	"github.com/bearer/bearer/pkg/report/output/sarif"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	outputhandler "github.com/bearer/bearer/pkg/util/output"
	"github.com/hhatto/gocloc"
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

func (f Formatter) Format(format string) (output *string, err error) {
	switch format {
	case flag.FormatEmpty:
		reportStr := BuildReportString(f.ReportData, f.Config, f.GoclocResult).String()
		output = &reportStr
	case flag.FormatSarif:
		sarifContent, sarifErr := sarif.ReportSarif(f.ReportData.FindingsBySeverity, f.Config.Rules)
		if sarifErr != nil {
			return output, fmt.Errorf("error generating sarif report %s", sarifErr)
		}
		output, err = outputhandler.ReportJSON(sarifContent)
	case flag.FormatReviewDog:
		sastContent, reviewdogErr := reviewdog.ReportReviewdog(f.ReportData.FindingsBySeverity)
		if reviewdogErr != nil {
			return output, fmt.Errorf("error generating reviewdog report %s", reviewdogErr)
		}
		output, err = outputhandler.ReportJSON(sastContent)
	case flag.FormatGitLabSast:
		sastContent, sastErr := gitlab.ReportGitLab(f.ReportData.FindingsBySeverity, f.StartTime, f.EndTime)
		if sastErr != nil {
			return output, fmt.Errorf("error generating gitlab-sast report %s", sastErr)
		}
		output, err = outputhandler.ReportJSON(sastContent)
	case flag.FormatJSON:
		output, err = outputhandler.ReportJSON(f.ReportData.FindingsBySeverity)
	case flag.FormatYAML:
		output, err = outputhandler.ReportYAML(f.ReportData.FindingsBySeverity)
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
	default:
		err = fmt.Errorf(`--report flag "%s" is not supported`, f.Config.Report.Report)
	}

	return output, err
}
