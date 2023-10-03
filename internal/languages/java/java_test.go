package java_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/internal/languages/java"
	"github.com/bearer/bearer/internal/languages/testhelper"
	patternquerybuilder "github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
	"github.com/bradleyjkemp/cupaloy"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, "Java").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "Java").RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestPattern(t *testing.T) {
	for _, test := range []struct{ name, pattern string }{
		{"method params is a container type", `
				class $<_> {
					void main($<!>$<_>) {}
				}
		`},
		{"catch types is a container type", `
				class $<_> {
					void main() {
						try {} catch ($<!>$<_> e) {}
					}
				}
		`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			result, err := patternquerybuilder.Build(java.Get(), test.pattern, "")
			if err != nil {
				tt.Fatalf("failed to build pattern: %s", err)
			}

			cupaloy.SnapshotT(tt, result)
		})
	}
}
