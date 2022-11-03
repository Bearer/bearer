package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/types"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

func ReportJSON(report types.Report, output *zerolog.Event) error {
	var ouputDetections any
	var err error

	if report.Type == flag.ReportDetectors {
		ouputDetections, err = GetDetectorsOutput(report)
		if err != nil {
			return err
		}
	} else if report.Type == flag.ReportDataFlow {
		ouputDetections, err = GetDataFlowOutput(report)
	}

	jsonBytes, err := json.MarshalIndent(&ouputDetections, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.Msg(string(jsonBytes))

	return nil
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
