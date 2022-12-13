package detectors

import (
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/types"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

func GetOutput(report types.Report, config settings.Config) ([]interface{}, error) {
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Running Detectors")
	}

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
