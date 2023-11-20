package ignore_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
)

func newIgnoreTest(name string, arguments []string) testhelper.TestCase {
	arguments = append([]string{
		"ignore"},
		arguments...,
	)
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{
		DisplayProgressBar: true,
		DisplayStdErr:      true,
		IgnoreForce:        true,
	})
}

func TestShowAll(t *testing.T) {
	tests := []testhelper.TestCase{
		newIgnoreTest("show-all", []string{
			"show",
			"--all",
			"--ignore-file",
			filepath.Join("e2e", "ignore", "testdata/test.ignore"),
		}),
	}

	testhelper.RunTests(t, tests)
}

func TestShowIndividual(t *testing.T) {
	tests := []testhelper.TestCase{
		newIgnoreTest("show-individual", []string{
			"show",
			"68a86d90f28db878612eb5f699c06543_0",
			"--ignore-file",
			filepath.Join("e2e", "ignore", "testdata/test.ignore"),
		}),
	}

	testhelper.RunTests(t, tests)
}
