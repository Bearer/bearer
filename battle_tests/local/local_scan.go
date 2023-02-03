package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bearer/curio/battle_tests/config"
	metricsscan "github.com/bearer/curio/battle_tests/metrics_scan"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Runtime.CurioExecutablePath = os.Getenv("CURIO_EXECUTABLE_PATH")
	log.Debug().Msgf("binary path: %s", config.Runtime.CurioExecutablePath)
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
		log.Debug().Msgf("result %s", string(data[:]))
	}
}
