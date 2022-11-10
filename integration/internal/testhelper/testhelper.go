package testhelper

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bearer/curio/cmd/curio/build"
	"github.com/bearer/curio/pkg/commands"
	"github.com/bearer/curio/pkg/commands/process/balancer"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog/log"
)

type TestCase struct {
	Name          string
	Arguments     []string
	ShouldSucceed bool
	RunInTempDir  bool
	OutputPath    string
	StartWorker   bool
}

func NewTestCase(name string, arguments []string) *TestCase {
	return &TestCase{
		Name:          name,
		Arguments:     arguments,
		ShouldSucceed: true,
	}
}

func executeApp(arguments []string, port int) (string, error) {
	app := commands.NewApp(build.Version, build.CommitSHA)

	if arguments[0] == "scan" {
		arguments = append(arguments, "--existing-worker=http://localhost:"+strconv.Itoa(port), "--quiet")
	}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	app.SetOut(stdoutWriter)
	app.SetErr(stderrWriter)
	app.SetArgs(arguments)

	if err := app.Execute(); err != nil {
		return "", err
	}

	stdoutWriter.Close()
	stderrWriter.Close()

	var stdoutBuf bytes.Buffer
	_, err = io.Copy(&stdoutBuf, stdoutReader)
	if err != nil {
		panic(err)
	}

	var stderrBuf bytes.Buffer
	_, err = io.Copy(&stderrBuf, stderrReader)
	if err != nil {
		panic(err)
	}

	combinedOutput := stdoutBuf.String() + "\n--\n" + stderrBuf.String()

	return combinedOutput, nil
}

func startWorker(port int) error {
	app := commands.NewApp(build.Version, build.CommitSHA)

	arguments := []string{"processing-worker", "--port=" + strconv.Itoa(port), "--debug"}

	app.SetArgs(arguments)

	log.Debug().Msgf("Starting worker on port: %d...", port)
	if err := app.Execute(); err != nil {
		return err
	}

	return nil
}

func RunTests(t *testing.T, tests []TestCase) {
	var port int

	shouldStartWorker := false
	for _, test := range tests {
		if test.StartWorker {
			shouldStartWorker = true
			break
		}
	}

	if shouldStartWorker {
		port = balancer.GetFreePort()
		go func() {
			err := startWorker(port)
			if err != nil {
				log.Fatal().Msgf("failed to start worker: %s", err)
			}
		}()

		// this needs to be here since otherwise viper is getting written twice concurrently from 2 gorutines
		// we need to find a way to let main program know viper has finished loading config
		time.Sleep(1 * time.Second)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var originalDir string
			var err error

			if test.RunInTempDir {
				originalDir, err = os.Getwd()
				if err != nil {
					t.Fatal(err)
				}

				tempDir, err := os.MkdirTemp("", "curio-integration-test")
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() {
					if err := os.Chdir(originalDir); err != nil {
						t.Fatal(err)
					}

					os.RemoveAll(tempDir)
				})

				if err := os.Chdir(tempDir); err != nil {
					t.Fatal(err)
				}
			}

			combinedOutput, err := executeApp(test.Arguments, port)

			if test.OutputPath != "" {
				fileContent, err := os.ReadFile(test.OutputPath)
				if err != nil {
					t.Fatalf("Failed to read file %s: %s", test.OutputPath, err)
				}
				combinedOutput = string(fileContent)
			}

			if test.RunInTempDir {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatal(err)
				}
			}

			cupaloy.SnapshotT(t, combinedOutput)

			if err != nil {
				if test.ShouldSucceed {
					t.Errorf("Expected application to succeed, but it failed: %s", err)
				}
			} else if !test.ShouldSucceed {
				t.Error("Expected application to fail, but it did not")
			}
		})
	}
}
