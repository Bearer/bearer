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

//go:embed testdata/import_rule.yml
var importRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, "Go").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "Go").RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, "Go").RunTest(t, "./testdata/import", ".snapshots/")
}
