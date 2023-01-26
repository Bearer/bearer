package flags_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newScanTest(name string, arguments []string) testhelper.TestCase {
	arguments = append([]string{"scan", filepath.Join("integration", "flags", "testdata", "simple")}, arguments...)
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{})
}

func TestReportFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		newScanTest("report-dataflow", []string{"--report=dataflow"}),
	}

	testhelper.RunTests(t, tests)
}

func TestReportFlagsShouldFail(t *testing.T) {
	tests := []testhelper.TestCase{
		newScanTest("invalid-report-flag", []string{"--report=testing"}),
		newScanTest("invalid-format-flag", []string{"--format=testing"}),
		newScanTest("invalid-context-flag", []string{"--format=testing"}),
	}

	for i := range tests {
		tests[i].ShouldSucceed = false
	}

	testhelper.RunTests(t, tests)

}
