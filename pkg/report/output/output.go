package output

import (
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"

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
	if config.Report.Report == flag.ReportDetectors {
		return detectors.GetOutput(report)
	} else if config.Report.Report == flag.ReportDataFlow {
		return getDataflow(report, config)
	} else if config.Report.Report == flag.ReportPolicies {
		dataflow, err := getDataflow(report, config)
		if err != nil {
			return nil, err
		}

		return policies.GetOutput(dataflow, config)
	} else if config.Report.Report == flag.ReportStats {
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

		processedDetections, ok := result["detections"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected to get slice from preprocessor but got %#v instead", result["detections"])
		}

		reportedDetections = append(reportedDetections, processedDetections...)
	}

	return dataflow.GetOutput(reportedDetections, config)
}
