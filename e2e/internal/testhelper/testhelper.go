package testhelper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
)

var TestTimeout = 1 * time.Minute

type TestCase struct {
	name               string
	arguments          []string
	ShouldSucceed      bool
	options            TestCaseOptions
	displayStdErr      bool
	displayProgressBar bool
	ignoreForce        bool
}

type TestCaseOptions struct {
	DisplayStdErr      bool
	DisplayProgressBar bool
	IgnoreForce        bool
}

func NewTestCase(name string, arguments []string, options TestCaseOptions) TestCase {
	return TestCase{
		name:               name,
		arguments:          arguments,
		ShouldSucceed:      true,
		options:            options,
		displayStdErr:      options.DisplayStdErr,
		displayProgressBar: options.DisplayProgressBar,
		ignoreForce:        options.IgnoreForce,
	}
}

func executeApp(t *testing.T, arguments []string) (string, string, error) {
	cmd, cancel := CreateCommand(arguments)

	buffOut := bytes.NewBuffer(nil)
	buffErr := bytes.NewBuffer(nil)
	cmd.Stdout = buffOut
	cmd.Stderr = buffErr

	var err error

	timer := time.NewTimer(TestTimeout)
	commandFinished := make(chan struct{}, 1)

	go func() {
		err = cmd.Start()

		if err != nil {
			commandFinished <- struct{}{}
			return
		}

		err = cmd.Wait()
		commandFinished <- struct{}{}
	}()

	select {
	case <-timer.C:
		cancel()
		t.Fatalf(
			"command failed to complete on time 'bearer %s':\n%s\n--\n%s",
			strings.Join(arguments, " "),
			buffOut,
			buffErr,
		)
	case <-commandFinished:
		cancel()
	}

	errStr := buffErr.String()
	// make output from `go run` match a compiled executable
	errStr = strings.TrimSuffix(errStr, "exit status 1\n")

	return buffOut.String(), errStr, err
}

func CreateCommand(arguments []string) (*exec.Cmd, context.CancelFunc) {
	var cmd *exec.Cmd

	ctx, cancel := context.WithCancel(context.Background())

	if os.Getenv("USE_BINARY") != "" {
		cmd = exec.CommandContext(ctx, executablePath(), arguments...)
	} else {
		arguments = append([]string{"run", GetCWD() + "/cmd/bearer/bearer.go"}, arguments...)
		cmd = exec.CommandContext(ctx, "go", arguments...)
	}

	cmd.Dir = GetCWD()
	cmd.Env = getEnvironment()

	return cmd, cancel
}

func getEnvironment() []string {
	var result []string

	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "BEARER_") {
			result = append(result, variable)
		}
	}

	return result
}

func executablePath() string {
	if value, ok := os.LookupEnv("BEARER_EXECUTABLE_PATH"); ok {
		return value
	}

	return "./bearer"
}

func GetCWD() string {
	return os.Getenv("GITHUB_WORKSPACE")
}

func RunTestsWithSnapshotSubdirectory(t *testing.T, tests []TestCase, snapshotSubdirectory string) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, stdErr := ExecuteTest(test, t)
			cupaloyCopy := cupaloy.NewDefaultConfig().WithOptions(cupaloy.SnapshotSubdirectory(snapshotSubdirectory))
			cupaloyCopy.SnapshotT(t, combineOutput(stdOut, stdErr))
		})
	}
}

func RunTests(t *testing.T, tests []TestCase) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdOut, stdErr := ExecuteTest(test, t)
			cupaloy.SnapshotT(t, combineOutput(stdOut, stdErr))
		})
	}
}

func ExecuteTest(test TestCase, t *testing.T) (string, string) {
	arguments := test.arguments

	if !test.displayProgressBar {
		arguments = append(arguments, "--hide-progress-bar")
	}

	if !test.displayStdErr {
		arguments = append(arguments, "--quiet")
	}

	if !test.ignoreForce {
		arguments = append(arguments, "--force")
	}

	stdOut, stdErr, err := executeApp(t, arguments)
	if test.ShouldSucceed && err != nil {
		t.Fatalf("command completed with error %s %s", err, combineOutput(stdOut, stdErr))
	}

	if !test.ShouldSucceed && err == nil {
		t.Fatal("expected command to fail but it succeeded instead")
	}

	return stdOut, stdErr
}

func combineOutput(stdOut, stdErr string) string {
	return fmt.Sprintf("%s\n--\n%s", stdOut, stdErr)
}
