package rules_test

import (
	"testing"
)

func TestAuxilary(t *testing.T) {
	t.Parallel()
	runRulesTest("auxilary", "javascript_third_parties_datadog_test", false, t)
}

func TestSimpleRuby(t *testing.T) {
	t.Parallel()
	runRulesTest("simple_ruby", "ruby_rails_insecure_communication_test", false, t)
}

func TestRubyRailsDefaultEncryptionStructure(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby_rails_default_encryption_structure_sql", "ruby_rails_default_encryption", true, t)
}

func TestRubyRailsDefaultEncryptionSchema(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby_rails_default_encryption_schema_rb", "ruby_rails_default_encryption", true, t)
}
