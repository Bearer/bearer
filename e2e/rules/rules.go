package rules_test

import (
	"testing"
)

func TestAuxilary(t *testing.T) {
	t.Parallel()
	runRulesTest("auxilary", "javascript_third_parties_datadog_test", t)
}
