package ruby_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/rule.yml
var loggerRule []byte

func TestRuby(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, loggerRule, "Ruby").RunTest(t, "./testdata/testcases", ".snapshots/")
}
