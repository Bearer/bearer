package policies_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newPolicyTest(name string, testFiles []string) testhelper.TestCase {
	filenames := []string{}
	for _, testFile := range testFiles {
		filenames = append(filenames, filepath.Join("testdata", testFile))
	}

	arguments := append(
		append(
			[]string{"scan"},
			filenames...,
		),
		"--report=policies",
		"--format=yaml",
	)

	options := testhelper.TestCaseOptions{StartWorker: true}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestPolicies(t *testing.T) {
	tests := []testhelper.TestCase{
		newPolicyTest("http_get_parameters", []string{"ruby/http_get_parameters.rb"}),
	}

	testhelper.RunTests(t, tests)
}
