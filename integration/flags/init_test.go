package integration_test

import (
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func TestInitCommand(t *testing.T) {
	options := testhelper.TestCaseOptions{
		RunInTempDir: true,
		OutputPath:   "curio.yml",
	}
	testCase := testhelper.NewTestCase("init", []string{"init"}, options)

	tests := []testhelper.TestCase{testCase}

	testhelper.RunTests(t, tests)
}
