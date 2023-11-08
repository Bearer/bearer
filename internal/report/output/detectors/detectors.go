package detectors

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/report/output/types"
	globaltypes "github.com/bearer/bearer/internal/types"
	"github.com/bearer/bearer/internal/util/jsonlines"
	"github.com/bearer/bearer/internal/util/output"
)

func AddReportData(
	reportData *types.ReportData,
	report globaltypes.Report,
	config settings.Config,
) error {
	if !config.Scan.Quiet && report.HasFiles {
		output.StdErrLog("Running Detectors")
	}

	var detections []interface{}
	f, err := os.Open(report.Path)
	if err != nil {
		return fmt.Errorf("failed to open report: %w", err)
	}

	err = jsonlines.Decode(f, &detections)
	if err != nil {
		return fmt.Errorf("failed to decode report: %w", err)
	}
	log.Debug().Msgf("got %d detections", len(detections))

	reportData.Detectors = detections

	return nil
}
