package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

func ReportSecurity(report types.Report, output *zerolog.Event, config settings.Config) (reportPassed bool, err error) {
	reportSupported := false
	reportPassed = false
	err = nil

	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		return
	}

	reportSupported, err = anySupportedLanguagesPresent(lineOfCodeOutput, config)
	if err != nil {
		return
	}

	if !reportSupported {
		var placeholderStr *strings.Builder
		placeholderStr, err = getPlaceholderOutput(report, config, lineOfCodeOutput)
		if err != nil {
			return
		}

		reportPassed = true
		output.Msg(placeholderStr.String())
		return
	}

	dataflow, err := getDataflow(report, config, true)
	if err != nil {
		return
	}

	securityResults, err := security.GetOutput(dataflow, config)
	if err != nil {
		return
	}

	reportStr, reportPassed := security.BuildReportString(config, securityResults, lineOfCodeOutput, dataflow)

	output.Msg(reportStr.String())

	return
}

func ReportCSV(report types.Report, output *zerolog.Event, config settings.Config) error {
	switch config.Report.Report {
	case flag.ReportPrivacy:
		return getPrivacyReportCsvOutput(report, output, config)
	}

	return fmt.Errorf("csv not supported for report type: %s", config.Report.Report)
}

func ReportJSON(report types.Report, output *zerolog.Event, config settings.Config) error {
	outputDetections, err := getReportOutput(report, config)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
}

func ReportYAML(report types.Report, output *zerolog.Event, config settings.Config) error {
	ouputDetections, err := getReportOutput(report, config)
	if err != nil {
		return err
	}

	jsonBytes, err := yaml.Marshal(&ouputDetections)
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
}

func getReportOutput(report types.Report, config settings.Config) (any, error) {

	switch config.Report.Report {
	case flag.ReportDetectors:
		return detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		return getDataflow(report, config, false)
	case flag.ReportSecurity:
		dataflow, err := getDataflow(report, config, true)
		if err != nil {
			return nil, err
		}
		return security.GetOutput(dataflow, config)
	case flag.ReportPrivacy:
		return getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		return reportStats(report, config)
	}

	return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getPrivacyReportCsvOutput(report types.Report, output *zerolog.Event, config settings.Config) error {
	dataflow, err := getDataflow(report, config, true)
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

func getPrivacyReportOutput(report types.Report, config settings.Config) (*privacy.Report, error) {
	dataflow, err := getDataflow(report, config, true)
	if err != nil {
		return nil, err
	}

	return privacy.GetOutput(dataflow, config)
}

func getDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, error) {
	reportedDetections, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, err
	}

	for _, detection := range reportedDetections {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func reportStats(report types.Report, config settings.Config) (*stats.Stats, error) {
	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	dataflowOutput, err := getDataflow(report, config, true)
	if err != nil {
		return nil, err
	}

	return stats.GetOutput(lineOfCodeOutput, dataflowOutput, config)
}

func anySupportedLanguagesPresent(inputgocloc *gocloc.Result, config settings.Config) (bool, error) {
	ruleLanguages := make(map[string]bool)
	for _, rule := range config.Rules {
		for _, language := range rule.Languages {
			ruleLanguages[language] = true
		}
	}

	foundLanguages := make(map[string]bool)
	for _, language := range inputgocloc.Languages {
		foundLanguages[strings.ToLower(language.Name)] = true
	}

	// typescript includes tsx, javascript includes jsx
	supportedLanguages := []string{"ruby", "javascript", "typescript"}
	for _, supportedLanguage := range supportedLanguages {
		_, found := foundLanguages[supportedLanguage]
		if found {
			return true, nil
		}
	}

	log.Debug().Msgf("No language found for which rules are applicable, found languages %#v", foundLanguages)
	return false, nil
}

func getPlaceholderOutput(report types.Report, config settings.Config, inputgocloc *gocloc.Result) (outputStr *strings.Builder, err error) {
	dataflowOutput, err := getDataflow(report, config, true)
	if err != nil {
		return
	}

	return stats.GetPlaceholderOutput(inputgocloc, dataflowOutput, config)
}
