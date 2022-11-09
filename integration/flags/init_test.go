package integration_test

import (
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func TestInitCommand(t *testing.T) {
	outputPath := "curio.yml"

	testCase := testhelper.NewTestCase("init", []string{"init"})
	testCase.OutputPath = outputPath
	testCase.RunInTempDir = true

	tests := []testhelper.TestCase{*testCase}

	testhelper.RunTests(t, tests)
}
