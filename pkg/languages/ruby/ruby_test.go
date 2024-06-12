package ruby_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/pkg/languages/ruby"
	"github.com/bearer/bearer/pkg/languages/testhelper"
)

//go:embed testdata/rule.yml
var loggerRule []byte

//go:embed testdata/pattern_variables_rule.yml
var patternVariablesRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestRuby(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, ruby.Get()).RunTest(t, "./testdata/testcases", ".snapshots/")
}

func TestPatternVariables(t *testing.T) {
	testhelper.GetRunner(t, patternVariablesRule, ruby.Get()).RunTest(t, "./testdata/pattern_variables", ".snapshots/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, ruby.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}
