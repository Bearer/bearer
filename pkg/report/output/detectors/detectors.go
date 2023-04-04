package detectors

import (
	"fmt"
	"os"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/jsonlines"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/rs/zerolog/log"
)

func GetOutput(report types.Report, config settings.Config) ([]interface{}, *dataflow.DataFlow, error) {
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Running Detectors")
	}

	var detections []interface{}
	f, err := os.Open(report.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open report: %w", err)
	}

	err = jsonlines.Decode(f, &detections)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode report: %w", err)
	}
	log.Debug().Msgf("got %d detections", len(detections))
	return detections, nil, nil
}
