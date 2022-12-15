package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/google/uuid"

	"github.com/bearer/curio/pkg/report/output/detectors"
	"github.com/bearer/curio/pkg/report/output/policies"
	"github.com/bearer/curio/pkg/report/output/stats"
	"github.com/bearer/curio/pkg/types"
	"gopkg.in/yaml.v3"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func ReportPolicies(report types.Report, output *zerolog.Event, config settings.Config) (reportPassed bool, err error) {
	reportSupported := false
	reportPassed = false
	err = nil

	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		return
	}

	reportSupported, err = anySupportedLanguagesPresent(lineOfCodeOutput , config)
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

	policyResults, err := getPolicyReportOutput(report, config)
	if err != nil {
		return
	}

	outputToFile := config.Report.Output != ""
	reportStr := policies.BuildReportString(policyResults, config.Policies, outputToFile)

	output.Msg(reportStr.String())

	reportPassed = len(policyResults) == 0
	return
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
	case flag.ReportPolicies:
		return getPolicyReportOutput(report, config)
	case flag.ReportStats:
		return reportStats(report, config)
	}

	return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getPolicyReportOutput(report types.Report, config settings.Config) (map[string][]policies.PolicyResult, error) {
	dataflow, err := getDataflow(report, config, true)
	if err != nil {
		return nil, err
	}

	return policies.GetOutput(dataflow, config)
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
	for _, rule := range config.CustomDetector {
		for _, language := range rule.Languages {
			ruleLanguages[language] = true
		}
	}

	foundLanguages := make(map[string]bool)
	for _, language := range inputgocloc.Languages {
		foundLanguages[strings.ToLower(language.Name)] = true
	}

	_, rubyPresent := foundLanguages["ruby"]
	if rubyPresent {
		return true, nil
	}

	// for foundLanguage := range foundLanguages {
	//	for ruleLanguage := range ruleLanguages {
	//		if foundLanguage == ruleLanguage {
	//			log.Debug().Msgf("Found a language for which policies are applicable: %s", foundLanguage)
	//			return true, nil
	//		}
	//	}
	// }

	log.Debug().Msg("No language found for which policies are applicable")
	return false, nil
}

func getPlaceholderOutput(report types.Report, config settings.Config, inputgocloc *gocloc.Result) (outputStr *strings.Builder, err error) {
	dataflowOutput, err := getDataflow(report, config, true)
	if err != nil {
		return
	}

	return stats.GetPlaceholderOutput(inputgocloc, dataflowOutput, config)
}
