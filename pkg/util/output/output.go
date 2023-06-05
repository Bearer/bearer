package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	Debug     bool
	ProcessID string
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

func Setup(cmd *cobra.Command, options SetupRequest) {
	outputWriter = cmd.OutOrStdout()
	ErrorWriter = cmd.ErrOrStderr()

	log.Logger = log.With().Str("process", options.ProcessID).Logger()

	if options.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
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
