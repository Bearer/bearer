package git

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BlameResult []BlameRange

type BlameRange struct {
	SHA       string
	StartLine int
	LineCount int
}

func Blame(rootDir string, commitSHAsFilename string, filename string) (BlameResult, error) {
	result := []BlameRange{}

	cmd := logAndBuildCommand(
		"blame",
		"--first-parent",
		"--incremental",
		"-M",
		"-S",
		commitSHAsFilename,
		filename,
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
	newRange := true
	blameRange := BlameRange{}

	for scanner.Scan() {
		line := scanner.Text()

		if newRange {
			splitLine := strings.Split(line, " ")
			blameRange.SHA = splitLine[0]

			startLine, err := strconv.Atoi(splitLine[2])
			if err != nil {
				killProcess(cmd)
				return nil, fmt.Errorf("failed to decode start line: %s", err)
			}
			blameRange.StartLine = startLine

			lineCount, err := strconv.Atoi(splitLine[3])
			if err != nil {
				killProcess(cmd)
				return nil, fmt.Errorf("failed to decode line count: %s", err)
			}
			blameRange.LineCount = lineCount

			newRange = false
			continue
		}

		splitLine := strings.SplitN(line, " ", 2)
		if splitLine[0] != "filename" {
			continue
		}

		result = append(result, blameRange)
		newRange = true
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

	return result, nil
}

func WriteCommitsForBlame(file *os.File, commitList []CommitInfo) error {
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, commitInfo := range commitList {
		if _, err := fmt.Fprintf(w, "%s ", commitInfo.SHA); err != nil {
			return err
		}
	}

	return w.Flush()
}

func (blame BlameResult) SHAForLine(lineNumber int) string {
	for _, blameRange := range blame {
		if lineNumber >= blameRange.StartLine && lineNumber < blameRange.StartLine+blameRange.LineCount {
			return blameRange.SHA
		}
	}

	return ""
}
