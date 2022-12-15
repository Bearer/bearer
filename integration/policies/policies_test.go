package policies_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newPolicyTest(name string, testFiles []string, healthContext bool) testhelper.TestCase {
	filenames := []string{}
	for _, testFile := range testFiles {
		filenames = append(filenames, filepath.Join("integration", "policies", "testdata", testFile))
	}

	arguments := append(
		append(
			[]string{"scan"},
			filenames...,
		),
		"--report=policies",
		"--format=yaml",
	)

	if healthContext {
		arguments = append(arguments, "--context=health")
	}

	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestPolicies(t *testing.T) {
	tests := []testhelper.TestCase{
		newPolicyTest("logger_leaking", []string{"ruby/logger_leaking.rb"}, false),
		newPolicyTest("http", []string{"ruby/http.rb"}, false),
		newPolicyTest("insecure_smtp", []string{"ruby/insecure_smtp.rb"}, false),
		newPolicyTest("insecure_communication", []string{"ruby/insecure_communication.rb"}, false),
		newPolicyTest("insecure_ftp", []string{"ruby/insecure_ftp.rb"}, false),
		newPolicyTest("sending_data_in_category_to_third_party", []string{"ruby/sending_data_in_category_to_third_party.rb"}, false),
		newPolicyTest("application_level_encryption_missing_structure_sql", []string{"ruby/application_level_encryption_missing/structure_sql"}, false),
		newPolicyTest("application_level_encryption_missing_schema_rb", []string{"ruby/application_level_encryption_missing/schema_rb"}, false),
		newPolicyTest("ruby_weak_password_encryption", []string{"ruby/weak_password_encryption.rb"}, false),
	}

	testhelper.RunTests(t, tests)
}

func TestPolicesWithHealthContext(t *testing.T) {
	tests := []testhelper.TestCase{
		newPolicyTest("logger_leaking", []string{"ruby/logger_leaking.rb"}, true),
		newPolicyTest("sending_data_in_category_to_third_party", []string{"ruby/sending_data_in_category_to_third_party.rb"}, true),
	}

	testhelper.RunTests(t, tests)
}

func TestPolicyFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		testhelper.NewTestCase(
			"skip_policy",
			[]string{
				"scan",
				filepath.Join("integration", "policies", "testdata", "ruby/logger_leaking.rb"),
				"--report=policies",
				"--format=yaml",
				"--skip-policy=CR-001",
			},
			testhelper.TestCaseOptions{},
		),
		testhelper.NewTestCase(
			"only_policy",
			[]string{
				"scan",
				filepath.Join("integration", "policies", "testdata", "ruby/logger_leaking.rb"),
				"--report=policies",
				"--format=yaml",
				"--only-policy=CR-001",
			},
			testhelper.TestCaseOptions{},
		),
	}

	testhelper.RunTests(t, tests)
}
