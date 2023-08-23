package output

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	globaltypes "github.com/bearer/bearer/pkg/types"

	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/saas"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/report/output/types"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func GetOutput(
	report globaltypes.Report,
	config settings.Config,
	baseBranchFindings *basebranchfindings.Findings,
) (*types.ReportData, error) {
	sendToCloud := false

	data := &types.ReportData{}
	// add detectors
	err := detectors.AddReportData(data, report, config)
	if config.Report.Report == flag.ReportDetectors {
		return data, err
	}

	if config.Report.Report == flag.ReportDataFlow {
		err = GetDataflow(data, report, config, false)
		return data, err
	}

	// add dataflow to data for internal use
	if err = GetDataflow(data, report, config, true); err != nil {
		return data, err
	}

	// add report-specific items
	switch config.Report.Report {
	case flag.ReportSecurity:
		sendToCloud = true
		err = security.AddReportData(data, config, baseBranchFindings)
	case flag.ReportSaaS:
		if err = security.AddReportData(data, config, baseBranchFindings); err != nil {
			return nil, err
		}

		sendToCloud = true
		err = saas.GetReport(data, config)
	case flag.ReportPrivacy:
		err = privacy.AddReportData(data, config)
	case flag.ReportStats:
		err = stats.AddReportData(data, report.Inputgocloc, config)
	default:
		return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	if sendToCloud && config.Client != nil && config.Client.Error == nil {
		// send SaaS report to Cloud
		saas.SendReport(config, data)
	}

	return data, err
}

func GetPrivacyReportCSVOutput(reportOutput *types.ReportData, report globaltypes.Report, config settings.Config) (*string, error) {
	csvString, err := privacy.BuildCsvString(reportOutput, config)
	if err != nil {
		return nil, err
	}

	content := csvString.String()

	return &content, nil
}

func GetDataflow(reportData *types.ReportData, report globaltypes.Report, config settings.Config, isInternal bool) error {
	if reportData.Detectors == nil {
		if err := detectors.AddReportData(reportData, report, config); err != nil {
			return err
		}
	}
	for _, detection := range reportData.Detectors {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}
	return dataflow.AddReportData(reportData, config, isInternal)
}
