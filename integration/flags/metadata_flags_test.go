package integration_test

import (
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newMetadataTest(name string, arguments []string) *testhelper.TestCase {
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{})
}

func TestMetadataFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		*newMetadataTest("help", []string{"help"}),
		*newMetadataTest("version", []string{"version"}),
		*newMetadataTest("scan-help", []string{"scan", "--help"}),
	}

	testhelper.RunTests(t, tests)
}
