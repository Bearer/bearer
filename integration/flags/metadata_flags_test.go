package integration_test

import (
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newTest(name string) *testhelper.TestCase {
	return testhelper.NewTestCase(name, []string{name}, testhelper.TestCaseOptions{})
}

func TestMetadataFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		*newTest("help"),
		*newTest("version"),
	}

	testhelper.RunTests(t, tests)
}
