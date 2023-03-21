package output

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/google/uuid"

	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/types"
	"gopkg.in/yaml.v3"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func ReportCSV(report types.Report, output *zerolog.Event, config settings.Config) error {
	switch config.Report.Report {
	case flag.ReportPrivacy:
		return getPrivacyReportCSVOutput(report, output, config)
	}

	return fmt.Errorf("csv not supported for report type: %s", config.Report.Report)
}

func ReportJSON(outputDetections any, config settings.Config) (*string, error) {
	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal detections: %s", err)
	}

	content := string(jsonBytes)
	return &content, nil
}

func ReportYAML(outputDetections any, config settings.Config) (*string, error) {
	yamlBytes, err := yaml.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to yaml marshal detections: %s", err)
	}

	content := string(yamlBytes)
	return &content, nil
}

func GetOutput(report types.Report, config settings.Config) (any, *gocloc.Result, *dataflow.DataFlow, error) {
	switch config.Report.Report {
	case flag.ReportDetectors:
		return detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		return GetDataflow(report, config, false)
	case flag.ReportSecurity:
		return reportSecurity(report, config)
	case flag.ReportPrivacy:
		return getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		return reportStats(report, config)
	}

	return nil, nil, nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getPrivacyReportCSVOutput(report types.Report, output *zerolog.Event, config settings.Config) error {
	dataflow, _, _, err := GetDataflow(report, config, true)
	if err != nil {
		return err
	}

	csvString, err := privacy.BuildCsvString(dataflow, config)
	if err != nil {
		return err
	}

	output.Msg(csvString.String())
	return nil
}

func getPrivacyReportOutput(report types.Report, config settings.Config) (*privacy.Report, *gocloc.Result, *dataflow.DataFlow, error) {
	dataflow, _, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, nil, err
	}

	return privacy.GetOutput(dataflow, config)
}

func GetDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, *gocloc.Result, *dataflow.DataFlow, error) {
	reportedDetections, _, _, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, detection := range reportedDetections {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func reportStats(report types.Report, config settings.Config) (*stats.Stats, *gocloc.Result, *dataflow.DataFlow, error) {
	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		return nil, nil, nil, err
	}

	dataflowOutput, _, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, nil, err
	}

	return stats.GetOutput(lineOfCodeOutput, dataflowOutput, config)
}

func reportSecurity(
	report types.Report,
	config settings.Config,
) (
	securityResults *security.Results,
	lineOfCodeOutput *gocloc.Result,
	dataflow *dataflow.DataFlow,
	err error,
) {
	lineOfCodeOutput, err = stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		log.Error().Msgf("error in line of code output %s", err)
		return
	}

	dataflow, _, _, err = GetDataflow(report, config, true)
	if err != nil {
		log.Error().Msgf("error in dataflow %s", err)
		return
	}

	securityResults, err = security.GetOutput(dataflow, config)
	if err != nil {
		log.Error().Msgf("error in security %s", err)
		return
	}

	// TODO: Should send report to us
	// content, err := ReportJSON(securityResults, config)
	// log.Error().Msg(*content)

	return
}
