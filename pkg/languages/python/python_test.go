package python_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/python"
	"github.com/bearer/bearer/pkg/languages/testhelper"
	patternquerybuilder "github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery/builder"
	"github.com/bradleyjkemp/cupaloy"
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

//go:embed testdata/decorator_rule.yml
var decoratorRule []byte

func TestDatatypes(t *testing.T) {
	testhelper.GetRunner(t, datatypesRule, python.Get()).RunTest(t, "./testdata/datatypes", ".snapshots/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, python.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, flowRule, python.Get()).RunTest(t, "./testdata/flow", ".snapshots/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, python.Get()).RunTest(t, "./testdata/import", ".snapshots/")
}

func TestSubscript(t *testing.T) {
	testhelper.GetRunner(t, subscriptRule, python.Get()).RunTest(t, "./testdata/subscript", ".snapshots/")
}

func TestPair(t *testing.T) {
	testhelper.GetRunner(t, pairRule, python.Get()).RunTest(t, "./testdata/pair", ".snapshots/")
}

func TestDecorator(t *testing.T) {
	testhelper.GetRunner(t, decoratorRule, python.Get()).RunTest(t, "./testdata/decorator", ".snapshots/")
}

func TestPattern(t *testing.T) {
	for _, test := range []struct{ name, pattern string }{
		{"catch function decorator", `
				$<!>@$<_>.route()
				def $<_>():
		`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			result, err := patternquerybuilder.Build(python.Get(), test.pattern, "")
			if err != nil {
				tt.Fatalf("failed to build pattern: %s", err)
			}

			cupaloy.SnapshotT(tt, result)
		})
	}
}
