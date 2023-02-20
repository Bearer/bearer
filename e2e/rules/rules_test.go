package rules_test

import (
	"testing"
)

func TestAuxilary(t *testing.T) {
	t.Parallel()
	runRulesTest("auxilary", "javascript_third_parties_datadog_test", t)
}

func TestSimpleRuby(t *testing.T) {
	t.Parallel()
	runRulesTest("simple_ruby", "ruby_rails_insecure_communication_test", t)
}
