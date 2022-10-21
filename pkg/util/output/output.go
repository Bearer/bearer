package output

import (
	"io"
	"os"

	"github.com/bearer/curio/pkg/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultLogger returns default output logger
func StdOutLogger() *zerolog.Event {
	return PlainLogger(os.Stdout)
}

func StdErrLogger() *zerolog.Event {
	return PlainLogger(os.Stderr)
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

func Setup(options flag.Options) {
	if options.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
