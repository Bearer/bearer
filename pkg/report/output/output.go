package output

import (
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/output/dataflow"

	dataflowtypes "github.com/bearer/curio/pkg/report/output/dataflow/types"
	"github.com/bearer/curio/pkg/report/output/detectors"
	"github.com/bearer/curio/pkg/report/output/policies"
	"github.com/bearer/curio/pkg/report/output/stats"
	"github.com/bearer/curio/pkg/types"

	"github.com/bearer/curio/pkg/util/rego"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
)

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
		return detectors.GetOutput(report)
	case flag.ReportDataFlow:
		return getDataflow(report, config)
	case flag.ReportPolicies:
		dataflow, err := getDataflow(report, config)
		if err != nil {
			return nil, err
		}

		return policies.GetOutput(dataflow, config)
	case flag.ReportStats:
		lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
		if err != nil {
			return nil, err
		}

		detectorsOutput, err := detectors.GetOutput(report)
		if err != nil {
			return nil, err
		}

		dataflowOutput, err := dataflow.GetOutput(detectorsOutput, config)
		if err != nil {
			return nil, err
		}

		return stats.GetOutput(lineOfCodeOutput, dataflowOutput, config)
	}

	return nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getDataflow(report types.Report, config settings.Config) (*dataflow.DataFlow, error) {
	reportedDetections, err := detectors.GetOutput(report)
	if err != nil {
		return nil, err
	}

	for _, processor := range config.Processors {
		result, err := rego.RunQuery(processor.Query, reportedDetections, processor.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		resultRisks, ok := result["risks"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected to get slice from preprocessor but got %#v instead", result["detections"])
		}

		for _, risk := range resultRisks {
			value := make(map[string]interface{})
			value["type"] = string(detections.TypeComputedDataflowRisk)
			value["value"] = risk

			reportedDetections = append(reportedDetections, value)
		}

	}

	return dataflow.GetOutput(reportedDetections, config)
}

type RiskDetection struct {
	Type  string                     `json:"type"`
	Value dataflowtypes.RiskDetector `json:"value"`
}
