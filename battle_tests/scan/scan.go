package scan

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/curio/battle_tests/config"
	"github.com/bearer/curio/battle_tests/fs"
	"github.com/bearer/curio/battle_tests/git"
	"github.com/bearer/curio/pkg/util/tmpfile"
)

type Scanner struct {
	repositoryUrl  string
	TempDir        string
	ReportFilePath string
	FSSize         int64
}

func NewScan(repositoryUrl string) *Scanner {
	return &Scanner{
		repositoryUrl: repositoryUrl,
	}
}

func (scanner *Scanner) Start() (outputBytes []byte, startTime *time.Time, err error) {
	// create temp directory
	log.Debug().Msg("creating temporary folder")
	scanner.TempDir, err = os.MkdirTemp(config.Runtime.EFSLocation, "broker_battle_test*")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp directory repository %w", err)
	}
	log.Debug().Msgf("temporary folder created at %s", scanner.TempDir)

	// clone to temp directory
	cloneDir, err := git.CloneLatest(scanner.TempDir, scanner.repositoryUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clone repository %w", err)
	}
	log.Debug().Msgf("cloneDir is %s", cloneDir)

	// calculate repository size
	scanner.FSSize, err = fs.DirSize(scanner.TempDir)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't determine size of repository %e", err)
	}
	log.Debug().Msgf("FS Size is %s", scanner.FSSize)

	// create file for report
	scanner.ReportFilePath = tmpfile.Create(config.Runtime.EFSLocation, ".json")
	log.Debug().Msgf("file for report path %s", scanner.ReportFilePath)

	processStartedAt := time.Now()
	log.Debug().Msgf("scan process start at %s", processStartedAt)
	// process scan
	output, err := exec.Command(
		config.Runtime.CurioExecutablePath,
		"scan",
		"--report=stats",
		"--debug",
		"--format=json",
		cloneDir,
	).Output()

	if err != nil {
		log.Error().Msgf("Error with cmd %e", err.Error())
	}

	return output, &processStartedAt, err
}

func (scanner *Scanner) Cleanup() {
	if _, err := os.Stat(scanner.ReportFilePath); err == nil {
		log.Printf("removing report %s", scanner.ReportFilePath)
		err := os.RemoveAll(scanner.ReportFilePath)
		if err != nil {
			log.Error().Msgf("failed to remove report %s %s", scanner.ReportFilePath, err.Error())
		}
	}

	if _, err := os.Stat(scanner.TempDir); err == nil {
		log.Printf("removing clone directory %s", scanner.TempDir)
		err = os.RemoveAll(scanner.TempDir)
		if err != nil {
			log.Error().Msgf("failed to remove directory %s", err.Error())
		}
	}
}
