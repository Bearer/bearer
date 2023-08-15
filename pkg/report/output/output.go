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
) (*types.Output[any], error) {
	var err error
	var output types.GenericOutput

	switch config.Report.Report {
	case flag.ReportDetectors:
		output, err = detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		output, err = GetDataflow(report, config, false)
	case flag.ReportSecurity:
		output, err = reportSecurity(report, config, baseBranchFindings)
	case flag.ReportSaaS:
		securityOutput, secErr := reportSecurity(report, config, baseBranchFindings)
		if secErr != nil {
			return nil, secErr
		}

		output, err = saas.GetReport(config, securityOutput)
	case flag.ReportPrivacy:
		output, err = getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		output, err = reportStats(report, config)
	default:
		return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}

	return output.ToGeneric(), err
}

func GetPrivacyReportCSVOutput(report globaltypes.Report, dataflow *types.DataFlow, config settings.Config) (*string, error) {
	csvString, err := privacy.BuildCsvString(dataflow, config)
	if err != nil {
		return nil, err
	}

	content := csvString.String()

	return &content, nil
}

func getPrivacyReportOutput(report globaltypes.Report, config settings.Config) (*types.Output[*privacy.Report], error) {
	dataflowOutput, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, err
	}

	return privacy.GetOutput(dataflowOutput.Dataflow, config)
}

func GetDataflow(report globaltypes.Report, config settings.Config, isInternal bool) (*types.Output[*types.DataFlow], error) {
	detectorsOutput, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, err
	}

	for _, detection := range detectorsOutput.Data {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(detectorsOutput.Data, config, isInternal)
}

func reportStats(report globaltypes.Report, config settings.Config) (*types.Output[stats.Stats], error) {
	dataflowOutput, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, err
	}

	return stats.GetOutput(report.Inputgocloc, dataflowOutput.Dataflow, config)
}

func reportSecurity(
	report globaltypes.Report,
	config settings.Config,
	baseBranchFindings *basebranchfindings.Findings,
) (*types.Output[security.Results], error) {
	dataflowOutput, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, fmt.Errorf("error in dataflow %w", err)
	}

	output, err := security.GetOutput(dataflowOutput, config, baseBranchFindings)
	if err != nil {
		return nil, fmt.Errorf("error in security %w", err)
	}

	if config.Client != nil && config.Client.Error == nil {
		saas.SendReport(config, output)
	}

	return output, nil
}
