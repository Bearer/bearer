package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/bearer/battle_tests/config"
	metricsscan "github.com/bearer/bearer/battle_tests/metrics_scan"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Runtime.ExecutablePath = os.Getenv("BEARER_EXECUTABLE_PATH")
	log.Debug().Msgf("binary path: %s", config.Runtime.ExecutablePath)
	ctx := context.TODO()
	log.Debug().Msg("Starting local test")
	metricsReport := make(chan *metricsscan.MetricsReport, 1)

	metricsscan.ScanRepository("https://github.com/Bearer/bear-publishing", "ruby", metricsReport)

	select {
	case <-ctx.Done():
		log.Error().Msgf("error")
		return
	case metrics := <-metricsReport:
		data, err := json.Marshal(metrics.PolicyBreaches)
		if err != nil {
			log.Error().Msgf("Failed to serialize metrics to JSON: %s", err)
			return
		}
		fmt.Println(string(data[:]))
	}
}
