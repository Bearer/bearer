package golang_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/internal/languages/testhelper"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, "<language>").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "<language>").RunTest(t, "./testdata/scope", ".snapshots/")
}
