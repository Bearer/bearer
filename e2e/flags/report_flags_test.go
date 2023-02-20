package flags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/curio/e2e/internal/testhelper"
	"github.com/bearer/curio/pkg/util/tmpfile"
	"github.com/bradleyjkemp/cupaloy"
)

func newScanTest(name string, arguments []string) testhelper.TestCase {
	arguments = append([]string{"scan", filepath.Join("integration", "flags", "testdata", "simple")}, arguments...)
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{})
}

func TestReportFlags(t *testing.T) {
	t.Parallel()
	tests := []testhelper.TestCase{
		newScanTest("report-dataflow", []string{"--report=dataflow"}),
	}

	testhelper.RunTests(t, tests)
}

func TestReportFlagsShouldFail(t *testing.T) {
	t.Parallel()
	tests := []testhelper.TestCase{
		newScanTest("invalid-report-flag", []string{"--report=testing"}),
		newScanTest("invalid-format-flag", []string{"--format=testing"}),
		newScanTest("invalid-context-flag", []string{"--context=testing"}),
	}

	for i := range tests {
		tests[i].ShouldSucceed = false
	}

	testhelper.RunTests(t, tests)

}

func TestOuputFlag(t *testing.T) {
	t.Parallel()
	outputPath := tmpfile.Create("", "test_output.jsonl")
	defer func() {
		err := os.Remove(outputPath)
		if err != nil {
			t.Fatalf("failed to clean up created output file %s", err)
		}
	}()

	tests := []testhelper.TestCase{
		newScanTest("output", []string{"--report=detectors", "--output=" + outputPath}),
	}

	testhelper.RunTests(t, tests)

	fileContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read created output file, err: %s", err)
	}

	cupaloy.SnapshotT(t, string(fileContent))
}
