package output

import (
	"os"

	"github.com/bearer/curio/pkg/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DefaultLogger() *zerolog.Event {
	logger := log.Output(zerolog.ConsoleWriter{
		Out:     os.Stdout,
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
