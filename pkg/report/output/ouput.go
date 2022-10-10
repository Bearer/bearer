package output

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bearer/curio/pkg/types"

	// "github.com/rs/zerolog/log"
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

	jsonBytes, err := json.MarshalIndent(&detections, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to json marshal detections: %w", err)
	}

	log.Print(string(jsonBytes))
	return nil
}
