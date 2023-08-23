package debugprofile

import (
	"os"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/internal/flag"
)

var cpuFile *os.File

func Start() {
	log.Debug().Msgf("starting cpu profiling")

	var err error
	cpuFile, err = os.Create(getProcessID() + "-cpu.pprof")
	if err != nil {
		log.Err(err).Msg("failed to create cpu profiling file")
		return
	}

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Err(err).Msg("failed to start cpu profile")
		cpuFile.Close()
		cpuFile = nil
		return
	}
}

func Stop() {
	if cpuFile == nil {
		return
	}

	log.Debug().Msg("stopping cpu profiling")
	pprof.StopCPUProfile()
	cpuFile.Close()
	cpuFile = nil

	log.Debug().Msgf("writing memory profile")
	memFile, err := os.Create(getProcessID() + "-mem.pprof")
	if err != nil {
		log.Err(err).Msg("failed to create memory profiling file")
		return
	}

	if err = pprof.Lookup("allocs").WriteTo(memFile, 0); err != nil {
		log.Err(err).Msg("failed to write memory profile")
	}

	memFile.Close()
}

func getProcessID() string {
	processID := viper.GetString(flag.WorkerIDFlag.ConfigName)
	if processID != "" {
		return processID
	}

	return "main"
}
