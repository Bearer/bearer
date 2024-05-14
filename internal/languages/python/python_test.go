package python_test

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
	testhelper.GetRunner(t, loggerRule, "python").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "python").RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, "python").RunTest(t, "./testdata/import", ".snapshots/")
}
