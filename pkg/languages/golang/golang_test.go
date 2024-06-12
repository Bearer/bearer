package golang_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/golang"
	"github.com/bearer/bearer/pkg/languages/testhelper"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

//go:embed testdata/import_rule.yml
var importRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, golang.Get()).RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, golang.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, golang.Get()).RunTest(t, "./testdata/import", ".snapshots/")
}
