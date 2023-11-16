package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

func logAndBuildCommand(ctx context.Context, args ...string) *exec.Cmd {
	log.Debug().Msgf("running command: git %s", strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = environment

	return cmd
}

func basicCommand(ctx context.Context, workingDir string, args ...string) error {
	cmd := logAndBuildCommand(ctx, args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	logWriter := &debugLogWriter{}
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Run(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return newError(err, logWriter.AllOutput())
	}

	return nil
}

func captureCommand(ctx context.Context, workingDir string, args []string, capture func(io.Reader) error) error {
	command := logAndBuildCommand(ctx, args...)
	if workingDir != "" {
		command.Dir = workingDir
	}

	logWriter := &debugLogWriter{}
	command.Stderr = logWriter

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		command.Cancel() //nolint:errcheck
		return err
	}

	if err := capture(stdout); err != nil {
		command.Cancel() //nolint:errcheck
		return err
	}

	stdout.Close()

	if err := command.Wait(); err != nil {
		command.Cancel() //nolint:errcheck
		return newError(err, logWriter.AllOutput())
	}

	return nil
}

func captureCommandBasic(ctx context.Context, workingDir string, args ...string) (output string, err error) {
	err = captureCommand(ctx, workingDir, args, func(r io.Reader) error {
		outputBytes, readErr := io.ReadAll(r)
		output = string(outputBytes)
		return readErr
	})

	return
}

var regexpDefunctProcess = regexp.MustCompile(" git <defunct>")
var regexpPID = regexp.MustCompile("[0-9]+ ")

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
