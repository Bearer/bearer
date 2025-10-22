package output

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hhatto/gocloc"

	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/saas"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/report/output/types"
	globaltypes "github.com/bearer/bearer/pkg/types"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func GetData(
	report globaltypes.Report,
	config settings.Config,
	gitContext *gitrepository.Context,
	baseBranchFindings *basebranchfindings.Findings,
) (*types.ReportData, error) {
	data := &types.ReportData{}

	// add languages
	languages := make(map[string]int32)
	languageFiles := make(map[string]int32)
	uniqueFiles := make(map[string]struct{})
	if report.Inputgocloc != nil {
		for _, language := range report.Inputgocloc.Languages {
			languages[language.Name] = language.Code
			languageFiles[language.Name] = int32(len(language.Files))
			for _, filename := range language.Files {
				uniqueFiles[filename] = struct{}{}
			}
		}
	}
	data.FoundLanguages = languages
	data.LanguageFiles = languageFiles
	data.TotalLanguageFiles = int32(len(uniqueFiles))

	// add detectors
	err := detectors.AddReportData(data, report, config)
	if config.Report.Report == flag.ReportDetectors || err != nil {
		return data, err
	}

	// add dataflow to data
	if err = GetDataflow(data, report, config, true); err != nil {
		return data, err
	}

	// add report-specific items
	switch config.Report.Report {
	case flag.ReportDataFlow:
		return data, err
	case flag.ReportSecurity:
		err = security.AddReportData(data, config, baseBranchFindings, report.HasFiles)
	case flag.ReportSaaS:
		if err = security.AddReportData(data, config, baseBranchFindings, report.HasFiles); err != nil {
			return nil, err
		}
		err = saas.GetReport(data, config, gitContext, false)
	case flag.ReportPrivacy:
		err = privacy.AddReportData(data, config)
	case flag.ReportStats:
		err = stats.AddReportData(data, report.Inputgocloc, config)
	default:
		return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	return data, err
}

func GetDataflow(
	reportData *types.ReportData,
	report globaltypes.Report,
	config settings.Config,
	isInternal bool,
) error {
	if reportData.Detectors == nil {
		if err := detectors.AddReportData(reportData, report, config); err != nil {
			return err
		}
	}
	for _, detection := range reportData.Detectors {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}
	return dataflow.AddReportData(reportData, config, isInternal, report.HasFiles)
}

func FormatOutput(
	reportData *types.ReportData,
	config settings.Config,
	engine engine.Engine,
	goclocResult *gocloc.Result,
	startTime time.Time,
	endTime time.Time,
) (string, error) {
	var formatter types.GenericFormatter
	switch config.Report.Report {
	case flag.ReportDetectors:
		formatter = detectors.NewFormatter(reportData, config)
	case flag.ReportDataFlow:
		formatter = dataflow.NewFormatter(reportData, config)
	case flag.ReportSecurity:
		formatter = security.NewFormatter(reportData, config, engine, goclocResult, startTime, endTime)
	case flag.ReportPrivacy:
		formatter = privacy.NewFormatter(reportData, config)
	case flag.ReportSaaS:
		formatter = saas.NewFormatter(reportData, config)
	case flag.ReportStats:
		formatter = stats.NewFormatter(reportData, config)
	default:
		return "", fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	formatStr, err := formatter.Format(config.Report.Format)
	if err != nil {
		return formatStr, err
	}
	if formatStr == "" {
		return "", fmt.Errorf(`--report flag "%s" does not support --format flag "%s"`, config.Report.Report, config.Report.Format)
	}

	return formatStr, err
}
