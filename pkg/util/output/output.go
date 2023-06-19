package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bearer/bearer/pkg/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	outputWriter io.Writer = os.Stdout
	ErrorWriter  io.Writer = os.Stderr
)

type SetupRequest struct {
	Quiet     bool
	ProcessID string
	LogLevel  string
}

// DefaultLogger returns default output logger
func StdOutLogger() *zerolog.Event {
	return PlainLogger(outputWriter)
}

func StdErrLogger() *zerolog.Event {
	return PlainLogger(ErrorWriter)
}

func PlainLogger(out io.Writer) *zerolog.Event {
	logger := log.Output(zerolog.ConsoleWriter{
		Out:     out,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			return ""
		},
		FormatLevel: func(i interface{}) string {
			return ""
		},
		FieldsExclude: []string{"process"},
	})

	return logger.Info()
}

func debugLogger(out io.Writer, processID string) zerolog.Logger {
	baseLogger := log.Output(zerolog.ConsoleWriter{
		Out:     out,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			timestamp, _ := time.Parse(time.RFC3339, i.(string))
			return timestamp.Format("2006-01-02 15:04:05")
		},
	})

	return baseLogger.With().Str("process", processID).Logger()
}

func Setup(cmd *cobra.Command, options SetupRequest) {
	outputWriter = cmd.OutOrStdout()
	ErrorWriter = cmd.ErrOrStderr()

	log.Logger = debugLogger(ErrorWriter, options.ProcessID)

	switch options.LogLevel {
	case flag.ErrorLogLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case flag.InfoLogLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case flag.DebugLogLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case flag.TraceLogLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

func ReportJSON(outputDetections any) (*string, error) {
	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal detections: %s", err)
	}

	content := string(jsonBytes)
	return &content, nil
}

func ReportYAML(outputDetections any) (*string, error) {
	yamlBytes, err := yaml.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to yaml marshal detections: %s", err)
	}

	content := string(yamlBytes)
	return &content, nil
}
