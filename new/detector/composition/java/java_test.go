package java_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/logger.yml
var loggerRule []byte

func TestFlow(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, loggerRule, "Javascript").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}
