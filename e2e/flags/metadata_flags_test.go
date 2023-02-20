package flags_test

import (
	"testing"

	"github.com/bearer/curio/e2e/internal/testhelper"
)

func newMetadataTest(name string, arguments []string) testhelper.TestCase {
	return testhelper.NewTestCase(name, arguments, testhelper.TestCaseOptions{DisplayStdErr: true, IgnoreForce: true})
}

func TestMetadataFlags(t *testing.T) {
	t.Parallel()
	tests := []testhelper.TestCase{
		newMetadataTest("help", []string{"help"}),
		newMetadataTest("version", []string{"version"}),
		newMetadataTest("scan-help", []string{"scan", "--help"}),
		newMetadataTest("help-scan", []string{"help", "scan"}),
	}

	testhelper.RunTests(t, tests)
}
