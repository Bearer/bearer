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
	scanner.TempDir, err = os.MkdirTemp(config.Runtime.EFSLocation, "broker_battle_test*")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp directory repository %w", err)
	}

	// clone to temp directory
	cloneDir, err := git.CloneLatest(scanner.TempDir, scanner.repositoryUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clone repository %w", err)
	}

	// calculate repository size
	scanner.FSSize, err = fs.DirSize(scanner.TempDir)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't determine size of repository %e", err)
	}

	// create file for report
	scanner.ReportFilePath = tmpfile.Create(config.Runtime.EFSLocation, ".json")

	processStartedAt := time.Now()
	// process scan
	output, err := exec.Command(
		"./curio",
		"scan",
		"--report=stats",
		"--quiet",
		"--format=json",
		cloneDir,
	).Output()

	return output, &processStartedAt, err
}

func (scanner *Scanner) Cleanup() {
	if _, err := os.Stat(scanner.ReportFilePath); err == nil {
		log.Printf("removing report %s", scanner.ReportFilePath)
		err := os.RemoveAll(scanner.ReportFilePath)
		if err != nil {
			log.Printf("failed to remove report %s %s", scanner.ReportFilePath, err)
		}
	}

	if _, err := os.Stat(scanner.TempDir); err == nil {
		log.Printf("removing clone directory %s", scanner.TempDir)
		err = os.RemoveAll(scanner.TempDir)
		if err != nil {
			log.Printf("failed to remove directory %s", err)
		}
	}
}
