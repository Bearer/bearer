package detectors

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/types"
	globaltypes "github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/jsonlines"
	"github.com/bearer/bearer/pkg/util/output"
)

func AddReportData(reportData *types.ReportData, report globaltypes.Report, config settings.Config) error {
	if !config.Scan.Quiet {
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
