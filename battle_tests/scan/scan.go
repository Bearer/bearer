package scan

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/battle_tests/config"
	"github.com/bearer/bearer/battle_tests/fs"
	"github.com/bearer/bearer/battle_tests/git"
	"github.com/bearer/bearer/pkg/util/tmpfile"
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

func (scanner *Scanner) Start(reportType string) (outputBytes []byte, startTime *time.Time, err error) {
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
		return nil, nil, fmt.Errorf("couldn't determine size of repository %s", err.Error())
	}
	log.Debug().Msgf("FSSize is %d", scanner.FSSize)

	// create file for report
	scanner.ReportFilePath = tmpfile.Create(config.Runtime.EFSLocation, ".json")
	log.Debug().Msgf("file for report path %s", scanner.ReportFilePath)

	processStartedAt := time.Now()
	log.Debug().Msgf("scan process (%s) on %s start at %s", reportType, scanner.repositoryUrl, processStartedAt)
	// process scan
	output, err := exec.Command(
		config.Runtime.ExecutablePath,
		"scan",
		fmt.Sprintf("--report=%s", reportType),
		"--quiet",
		"--format=json",
		cloneDir,
	).Output()

	if err != nil {
		log.Error().Msgf("Error with cmd %s", err.Error())
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
