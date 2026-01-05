package rust_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/rust"
	"github.com/bearer/bearer/pkg/languages/testhelper"
	patternquerybuilder "github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery/builder"
	"github.com/bradleyjkemp/cupaloy"
)

//go:embed testdata/scope_rule.yml
var scopeRule []byte

//go:embed testdata/flow_rule.yml
var flowRule []byte

//go:embed testdata/import_rule.yml
var importRule []byte

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, rust.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, flowRule, rust.Get()).RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestImport(t *testing.T) {
	testhelper.GetRunner(t, importRule, rust.Get()).RunTest(t, "./testdata/import", ".snapshots/")
}

func TestPattern(t *testing.T) {
	for _, test := range []struct{ name, pattern string }{
		// Using single named variable or only anonymous variables
		// to avoid non-deterministic ordering issues in snapshots
		{"function call", `$<FUNC>($<_>)`},
		{"method call", `$<_>.$<_>()`},
		{"struct initialization", `$<_> { $<...> }`},
		{"let declaration", `let $<_> = $<_>`},
		{"unsafe block", `unsafe { $<...> }`},
		{"field access", `$<_>.$<_>`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			result, err := patternquerybuilder.Build(rust.Get(), test.pattern, "")
			if err != nil {
				tt.Fatalf("failed to build pattern: %s", err)
			}

			cupaloy.SnapshotT(tt, result)
		})
	}
}

