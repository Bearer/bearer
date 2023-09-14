package git

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// MonitorDefunct periodically cleans up defunct zombie git proccess that happen on git error and pile up over time
func MonitorDefunct(ctx context.Context) {
	timer := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer.C:
			cleanupDefunct()
		case <-ctx.Done():
			return
		}
	}
}

func cleanupDefunct() {
	log.Debug().Msgf("[git] defunct checking for abandoned process")

	cmdProcess := exec.Command("ps", "-A", "-H")
	stdout, err := cmdProcess.CombinedOutput()
	if err != nil {
		log.Debug().Msgf("failed to get output of command %s", err)
	}

	defunctPIDs := make([]string, 0)

	lines := strings.Split(string(stdout), "\n")

	for _, line := range lines {
		if !regexpDefunctProcess.Match([]byte(line)) {
			continue
		}

		pid := regexpPID.FindString(line)
		if pid == "" {
			continue
		}

		defunctPIDs = append(defunctPIDs, strings.Trim(pid, " "))
	}

	for _, pid := range defunctPIDs {
		cmdKill := exec.Command("kill", pid)
		if err := cmdKill.Run(); err != nil {
			log.Debug().Msgf("failed to kill git process %s", err)
		}
	}

	log.Debug().Msgf("[git] defunct process cleanup found %d processes", len(defunctPIDs))
}
