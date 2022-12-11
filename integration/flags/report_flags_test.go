package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/util/tmpfile"
)

func newScanTest(name string, arguments []string, outputPath string) testhelper.TestCase {
	return newScanProject(name, arguments, outputPath, "simple")
}

func newScanProject(name string, arguments []string, outputPath string, projectpath string) testhelper.TestCase {
	arguments = append([]string{"scan", filepath.Join("testdata", projectpath)}, arguments...)
	options := testhelper.TestCaseOptions{
		OutputPath:  outputPath,
		StartWorker: true,
	}
	return testhelper.NewTestCase(name, arguments, options)
}

func TestReportFlags(t *testing.T) {
	outputPath := tmpfile.Create("", "integration_test.jsonl")
	t.Cleanup(func() {
		os.Remove(outputPath)
	})

	tests := []testhelper.TestCase{
		newScanTest("format-json", []string{"--report=detectors", "--format=json"}, ""),
		newScanTest("format-yaml", []string{"--report=detectors", "--format=yaml"}, ""),
		newScanTest("report-detectors", []string{"--report=detectors"}, ""),
		newScanTest("report-dataflow", []string{"--report=dataflow"}, ""),
		newScanProject("report-dataflow-verified-by", []string{"--report=dataflow", "--format=yaml"}, "", "verified_by"),
		newScanProject("report-policies", []string{"--report=policies", "--format=yaml"}, "", "policies"),
		newScanTest("output", []string{"--report=detectors", "--output=" + outputPath}, outputPath),
		newScanTest("health-context", []string{"--report=detectors", "--context=health"}, ""),
		newScanTest("domain-resolution-disabled", []string{"--report=detectors", "--disable-domain-resolution=true"}, ""),
		newScanTest("skipped-paths", []string{"--report=detectors", "--skip-path=\"users/*.go,users/admin.sql\""}, ""),
	}

	testhelper.RunTests(t, tests)
}
