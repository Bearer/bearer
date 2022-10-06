package git

import (
	"bufio"
	"fmt"
	"strings"
)

type RenamedFile struct {
	PreviousFilename string `json:"previous_filename"`
	NewFilename      string `json:"new_filename"`
}

func GetRenames(rootDir, firstCommitSHA, lastCommitSHA string) ([]RenamedFile, error) {
	cmd := logAndBuildCommand(
		"log",
		"--first-parent",
		"--find-renames",
		"--break-rewrites",
		"--name-status",
		"--diff-filter=R",
		"--pretty=tformat:",
		firstCommitSHA+".."+lastCommitSHA,
		"--",
	)
	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		killProcess(cmd)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	renameMap := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, "\t")

		prevFilename, err := unquoteFilename(splitLine[1])
		if err != nil {
			killProcess(cmd)
			return nil, fmt.Errorf("failed to unquote previous filename: %s", err)
		}
		newFilename, err := unquoteFilename(splitLine[2])
		if err != nil {
			killProcess(cmd)
			return nil, fmt.Errorf("failed to unquote new filename: %s", err)
		}

		if latestFilename, alreadyRenamed := renameMap[newFilename]; alreadyRenamed {
			delete(renameMap, newFilename)
			newFilename = latestFilename
		}

		renameMap[prevFilename] = newFilename
	}

	if err := scanner.Err(); err != nil {
		killProcess(cmd)
		return nil, err
	}

	stdout.Close()

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
		return nil, newError(err, logWriter.AllOutput())
	}

	result := []RenamedFile{}
	for prevFilename, newFilename := range renameMap {
		result = append(result, RenamedFile{PreviousFilename: prevFilename, NewFilename: newFilename})
	}

	return result, nil
}
