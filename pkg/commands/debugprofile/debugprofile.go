package debugprofile

import (
	"os"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/pkg/flag"
)

var file *os.File

func Start() {
	log.Debug().Msgf("starting cpu profiling")

	processID := viper.GetString(flag.WorkerIDFlag.ConfigName)
	if processID == "" {
		processID = "main"
	}

	var err error
	file, err = os.Create(processID + ".pprof")
	if err != nil {
		log.Err(err).Msg("failed to create profiling file")
		return
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		log.Err(err).Msg("failed to start cpu profile")
		file.Close()
		file = nil
		return
	}
}

func Stop() {
	if file == nil {
		return
	}

	log.Debug().Msg("stopping cpu profiling")
	pprof.StopCPUProfile()
	file.Close()

	file = nil
}
