package integration_test

import (
	"os"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/util/tmpfile"
	"github.com/bradleyjkemp/cupaloy"
)

func TestOuputFlag(t *testing.T) {
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
