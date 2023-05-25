package ruby_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/rule.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestRuby(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, loggerRule, "Ruby").RunTest(t, "./testdata/testcases", ".snapshots/")
}

func TestScope(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, scopeRule, "Ruby").RunTest(t, "./testdata/scope", ".snapshots/")
}
