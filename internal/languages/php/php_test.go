package php_test

import (
	_ "embed"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/languages/php"
	"github.com/bearer/bearer/internal/languages/testhelper"
	patternquerybuilder "github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, "PHP").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "PHP").RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestPattern(t *testing.T) {
	for _, test := range []struct{ name, pattern string }{
		{"class name in object creation is unanchored", `
				new $<!>Foo;
		`},
		{"named arguments are unanchored", `
				foo(x: $<!>$<_>)
		`},
		{"property names are unanchored", `
				class $<_> {
					public $<!>$<_>;
				}
		`},
		{"parameter names are unanchored", `
				class $<_> {
					public function $<_>($<_> $<!>$<_>) {}
				}
		`},
		{"catch clauses and types are unanchored", `
				try {} catch ($<_> $<!>$$<_>) {}
		`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			result, err := patternquerybuilder.Build(php.Get(), test.pattern, "")
			if err != nil {
				tt.Fatalf("failed to build pattern: %s", err)
			}

			cupaloy.SnapshotT(tt, result)
		})
	}
}
