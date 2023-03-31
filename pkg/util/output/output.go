package output

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	outputWriter io.Writer = os.Stdout
	ErrorWriter  io.Writer = os.Stderr
)

type SetupRequest struct {
	Quiet bool
	Debug bool
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
	})

	return logger.Info()
}

func Setup(cmd *cobra.Command, options SetupRequest) {
	outputWriter = cmd.OutOrStdout()
	ErrorWriter = cmd.ErrOrStderr()

	if options.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
