package testhelper

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type TestCase struct {
	name          string
	arguments     []string
	shouldSucceed bool
	options       TestCaseOptions
	displayStdErr bool
	ignoreForce   bool
}

type TestCaseOptions struct {
	DisplayStdErr bool
	IgnoreForce   bool
}

func NewTestCase(name string, arguments []string, options TestCaseOptions) TestCase {
	return TestCase{
		name:          name,
		arguments:     arguments,
		shouldSucceed: true,
		options:       options,
		displayStdErr: options.DisplayStdErr,
		ignoreForce:   options.IgnoreForce,
	}
}

func executeApp(arguments []string) (string, error) {
	cmd := CreateCurioCommand(arguments)

	buffOut := bytes.NewBuffer(nil)
	buffErr := bytes.NewBuffer(nil)
	cmd.Stdout = buffOut
	cmd.Stderr = buffErr

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	combinedOutput := buffOut.String() + "\n--\n" + buffErr.String()

	return combinedOutput, nil
}

func CreateCurioCommand(arguments []string) *exec.Cmd {
	var cmd *exec.Cmd

	if os.Getenv("CURIO_BINARY") != "" {
		cmd = exec.Command("./curio", arguments...)
	} else {
		arguments = append([]string{"run", GetCWD() + "/cmd/curio/main.go"}, arguments...)
		cmd = exec.Command("go", arguments...)
	}

	cmd.Dir = os.Getenv("GITHUB_WORKSPACE")

	return cmd
}

func GetCWD() string {
	return os.Getenv("GITHUB_WORKSPACE")
}

func RunTests(t *testing.T, tests []TestCase) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			arguments := test.arguments

			if !test.displayStdErr {
				arguments = append(arguments, "--quiet")
			}

			if !test.ignoreForce {
				arguments = append(arguments, "--force")
			}

			combinedOutput, err := executeApp(arguments)

			cupaloy.SnapshotT(t, combinedOutput)

			if err != nil {
				if test.shouldSucceed {
					t.Errorf("Expected application to succeed, but it failed: %s", err)
				}
			} else if !test.shouldSucceed {
				t.Error("Expected application to fail, but it did not")
			}
		})
	}
}
