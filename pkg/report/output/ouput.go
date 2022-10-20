package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/types"
	"github.com/bearer/curio/pkg/util/output"

	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

func ReportJSON(report types.Report) error {
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

	jsonBytes, err := json.MarshalIndent(&detections, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	output.DefaultLogger().Msg(string(jsonBytes))

	return nil
}
