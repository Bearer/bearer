package flags_test

import (
	"os"
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
	"github.com/bradleyjkemp/cupaloy"
)

func TestInitCommand(t *testing.T) {
	t.Parallel()
	filePath := testhelper.GetCWD() + "/bearer.yml"
	arguments := []string{"init"}

	cmd, _ := testhelper.CreateCommand(arguments)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run init command err: %s output: %s", err, string(output))
	}
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			t.Fatalf("failed to cleanup file err: %s", err)
		}
	}()

	_, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to find created bearer.yaml file err:%s", err)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read created config file err: %s", err)
	}

	cupaloy.SnapshotT(t, string(fileContent))
}
