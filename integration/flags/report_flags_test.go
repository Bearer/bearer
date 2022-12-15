package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newScanTest(name string, arguments []string) testhelper.TestCase {
	return newScanProject(name, arguments, "simple")
}

func newScanProject(name string, arguments []string, projectpath string) testhelper.TestCase {
	arguments = append([]string{"scan", filepath.Join("integration", "flags", "testdata", projectpath)}, arguments...)
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{})
}

func TestReportFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		newScanTest("format-json", []string{"--report=detectors", "--format=json"}),
		newScanTest("format-yaml", []string{"--report=detectors", "--format=yaml"}),
		newScanTest("report-detectors", []string{"--report=detectors"}),
		newScanTest("report-dataflow", []string{"--report=dataflow"}),
		newScanProject("report-dataflow-verified-by", []string{"--report=dataflow", "--format=yaml"}, "verified_by"),
		newScanProject("report-policies", []string{"--report=policies", "--format=yaml"}, "policies"),
		newScanTest("health-context", []string{"--report=detectors", "--context=health"}),
		newScanTest("domain-resolution-disabled", []string{"--report=detectors", "--disable-domain-resolution=true"}),
		newScanTest("skipped-paths", []string{"--report=detectors", "--skip-path=\"users/*.go,users/admin.sql\""}),
	}

	testhelper.RunTests(t, tests)
}
