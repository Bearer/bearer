package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bearer/bearer/internal/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	outputWriter io.Writer = os.Stdout
	errorWriter  io.Writer = os.Stderr

	stdOutLogger func(message string)
	stdErrLogger func(message string)
)

type SetupRequest struct {
	Quiet     bool
	ProcessID string
	LogLevel  string
}

func ErrorWriter() io.Writer {
	return errorWriter
}

func PlainLogger(out io.Writer) func(message string) {
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

	return func(message string) { logger.Info().Msg(message) }
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
	errorWriter = cmd.ErrOrStderr()
	stdOutLogger = PlainLogger(outputWriter)
	stdErrLogger = PlainLogger(errorWriter)
	log.Logger = debugLogger(errorWriter, options.ProcessID)

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

func StdErrLog(message string) {
	if stdErrLogger != nil {
		stdErrLogger(message)
	}
}

func StdOutLog(message string) {
	if stdOutLogger != nil {
		stdOutLogger(message)
	}
}

func Fatal(message string) {
	StdErrLog(message)
	os.Exit(1)
}

func ReportJSON(outputDetections any) (string, error) {
	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return "", fmt.Errorf("failed to json marshal detections: %s", err)
	}

	return string(jsonBytes), nil
}

func ReportYAML(outputDetections any) (string, error) {
	yamlBytes, err := yaml.Marshal(&outputDetections)
	if err != nil {
		return "", fmt.Errorf("failed to yaml marshal detections: %s", err)
	}

	return string(yamlBytes), nil
}
