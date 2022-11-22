package git

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	mailmapFilename = ".mailmap"
)

var (
	environment []string = append(
		os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GCM_INTERACTIVE=never",
		"GIT_LFS_SKIP_SMUDGE=1",
	)

	specialFiles []string = []string{
		mailmapFilename,
		".gitignore",
		".gitattributes",
		".gitmodules",
	}
)

type CommitIdentifier struct {
	SHA       string    `json:"sha" yaml:"sha"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

func urlWithCredentials(originalURL *url.URL, token string) *url.URL {
	result := *originalURL
	result.User = url.UserPassword("x", token)
	return &result
}

func logAndBuildCommand(args ...string) *exec.Cmd {
	log.Debug().Msgf("running command: git %s", strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	cmd.Env = environment

	return cmd
}

func basicCommand(workingDir string, args ...string) error {
	cmd := logAndBuildCommand(args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	logWriter := &debugLogWriter{}
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Run(); err != nil {
		killProcess(cmd)
		return newError(err, logWriter.AllOutput())
	}

	return nil
}

var regexpDefunctProcess = regexp.MustCompile(" git <defunct>")
var regexpPID = regexp.MustCompile("[0-9]+ ")

func killProcess(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill() //nolint:all,errcheck
	}
}

func unquoteFilename(quoted string) (string, error) {
	if len(quoted) == 0 || quoted[0] != '"' {
		return quoted, nil
	}

	return strconv.Unquote(quoted)
}

var logSplitPattern = regexp.MustCompile(`[\r\n]+`)

type debugLogWriter struct {
	allOutput bytes.Buffer
}

func (writer *debugLogWriter) Write(data []byte) (int, error) {
	n := len(data)

	for _, line := range logSplitPattern.Split(string(data), -1) {
		line = strings.TrimRight(line, " ")
		if len(line) != 0 {
			log.Debug().Msgf("[git] %s", line)
		}
	}

	writer.allOutput.Write(data)

	return n, nil
}

func (writer *debugLogWriter) AllOutput() string {
	return writer.allOutput.String()
}

type GitError struct {
	Err    error
	Output string
}

func (err *GitError) Error() string {
	return fmt.Sprintf("%s\n-- git output --\n%s-- end git output --", err.Err, err.Output)
}

func newError(err error, output string) *GitError {
	return &GitError{
		Err:    err,
		Output: output,
	}
}
