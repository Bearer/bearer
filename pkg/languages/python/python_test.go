package python_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/testhelper"
)

//go:embed testdata/datatypes_rule.yml
var datatypesRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

//go:embed testdata/flow_rule.yml
var flowRule []byte

//go:embed testdata/import_rule.yml
var importRule []byte

//go:embed testdata/subscript_rule.yml
var subscriptRule []byte

//go:embed testdata/pair_rule.yml
var pairRule []byte

func TestDatatypes(t *testing.T) {
	testhelper.GetRunner(t, datatypesRule, "python").RunTest(t, "./testdata/datatypes", ".snapshots/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "python").RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, flowRule, "python").RunTest(t, "./testdata/flow", ".snapshots/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, "python").RunTest(t, "./testdata/import", ".snapshots/")
}

func TestSubscript(t *testing.T) {
	testhelper.GetRunner(t, subscriptRule, "python").RunTest(t, "./testdata/subscript", ".snapshots/")
}

func TestPair(t *testing.T) {
	testhelper.GetRunner(t, pairRule, "python").RunTest(t, "./testdata/pair", ".snapshots/")
}
