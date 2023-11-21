package rules_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
)

func TestSecrets(t *testing.T) {
	testCases := []testhelper.TestCase{
		testhelper.NewTestCase(
			"secrets",
			[]string{
				"scan",
				filepath.Join("e2e", "rules", "testdata/data/secrets"),
				"--scanner=secrets",
				"--only-rule=gitleaks",
				"--format=yaml",
				"--disable-version-check",
				"--disable-default-rules",
				"--exit-code=0",
			},
			testhelper.TestCaseOptions{},
		),
	}

	testhelper.RunTestsWithSnapshotSubdirectory(t, testCases, ".snapshots")
}

func TestAuxilary(t *testing.T) {
	runRulesTest("auxilary", "javascript_third_parties_datadog_test", t)
}

func TestDisabledRules(t *testing.T) {
	runRulesTest("disabled_rules", "match_sink", t)
}

func TestReferenceFilters(t *testing.T) {
	runRulesTest("reference_filters", "reference_filters_test", t)
}

func TestSanitizer(t *testing.T) {
	runRulesTest("sanitizer", "sanitizer_test", t)
}

func TestSimpleRuby(t *testing.T) {
	runRulesTest("simple_ruby", "ruby_rails_insecure_communication_test", t)
}

func TestRubyRailsDefaultEncryptionStructure(t *testing.T) {
	runRulesTest("ruby_rails_default_encryption_structure_sql", "ruby_rails_default_encryption", t)
}

func TestRubyRailsDefaultEncryptionSchema(t *testing.T) {
	runRulesTest("ruby_rails_default_encryption_schema_rb", "ruby_rails_default_encryption", t)
}

func TestExpectedRule(t *testing.T) {
	testDataDir := "testdata/data/expected_rule"

	testCases := []testhelper.TestCase{}
	testCases = append(testCases,
		buildRulesTestCaseJsonV2(
			testDataDir,
			filepath.Join("e2e", "rules", testDataDir),
			"expected_rule",
		),
	)

	testhelper.RunTestsWithSnapshotSubdirectory(t, testCases, ".snapshots")
}
