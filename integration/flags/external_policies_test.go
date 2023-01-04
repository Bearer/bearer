package flags_test

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func TestExternalPolicies(t *testing.T) {
	basePath := filepath.Join("integration", "flags", "testdata", "external_policies")
	externalDetectorDirFlag := "--external-detector-dir=" + filepath.Join(basePath, "detectors")
	externalPolicyFlag := "--external-policy-dir=" + filepath.Join(basePath, "policies")

	log.Printf("%s", basePath)
	arguments := []string{"scan", filepath.Join(basePath, "scan_data"), externalDetectorDirFlag, externalPolicyFlag, "--report=policies", "--format=yaml"}

	testCase := testhelper.NewTestCase("test external policies", arguments, testhelper.TestCaseOptions{})

	tests := []testhelper.TestCase{testCase}

	testhelper.RunTests(t, tests)
}
