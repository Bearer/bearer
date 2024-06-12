package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/javascript"
	"github.com/bearer/bearer/pkg/languages/testhelper"
)

//go:embed testdata/import_rule.yml
var importRule []byte

//go:embed testdata/insecureURL.yml
var insecureURLRule []byte

//go:embed testdata/datatype.yml
var datatypeRule []byte

//go:embed testdata/deconstructing.yml
var deconstructingRule []byte

//go:embed testdata/pattern_variables_rule.yml
var patternVariablesRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, datatypeRule, javascript.Get()).RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestObjectDeconstructing(t *testing.T) {
	testhelper.GetRunner(t, deconstructingRule, javascript.Get()).RunTest(t, "./testdata/testcases/object-deconstructing", ".snapshots/object-deconstructing/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, javascript.Get()).RunTest(t, "./testdata/import", ".snapshots/import/")
}

func TestString(t *testing.T) {
	testhelper.GetRunner(t, insecureURLRule, javascript.Get()).RunTest(t, "./testdata/testcases/string", ".snapshots/string/")
}

func TestPatternVariables(t *testing.T) {
	testhelper.GetRunner(t, patternVariablesRule, javascript.Get()).RunTest(t, "./testdata/pattern_variables", ".snapshots/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, javascript.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}
