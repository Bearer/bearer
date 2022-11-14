package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/policies"
	"github.com/bearer/curio/pkg/types"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

func ReportJSON(report types.Report, output *zerolog.Event, config settings.Config) error {
	ouputDetections, err := getReportOutput(report, config)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(&ouputDetections)
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
	var ouputDetections any
	var err error

	if config.Report.Report == flag.ReportDetectors {
		ouputDetections, err = GetDetectorsOutput(report)
		if err != nil {
			return nil, err
		}
	} else if config.Report.Report == flag.ReportDataFlow {
		detections, err := GetDetectorsOutput(report)
		if err != nil {
			return nil, err
		}

		ouputDetections, err = dataflow.GetOuput(detections, config)
		// output, err := policies.GetDataflow(detections)
		if err != nil {
			return nil, err
		}

		return ouputDetections, nil

	} else if config.Report.Report == flag.ReportPolicies {
		detections, err := GetDetectorsOutput(report)
		if err != nil {
			return nil, err
		}

		policiesData, err := dataflow.GetOuput(detections, config)
		if err != nil {
			return nil, err
		}

		data, err := policies.GetPolicies(policiesData, config)
		if err != nil {
			return nil, err
		}

		return data, err
		// log.Debug().Msgf("%s", data)
	}

	return nil, nil
}

func GetDetectorsOutput(report types.Report) ([]interface{}, error) {
	var detections []interface{}
	f, err := os.Open(report.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open report: %w", err)
	}

	err = jsonlines.Decode(f, &detections)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report: %w", err)
	}
	log.Debug().Msgf("got %d detections", len(detections))

	return detections, nil
}
