package policies_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newPolicyTest(name string, testFiles []string) testhelper.TestCase {
	filenames := []string{}
	for _, testFile := range testFiles {
		filenames = append(filenames, filepath.Join("testdata", testFile))
	}

	arguments := append(
		append(
			[]string{"scan"},
			filenames...,
		),
		"--report=policies",
		"--format=yaml",
	)

	options := testhelper.TestCaseOptions{StartWorker: true}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestPolicies(t *testing.T) {
	tests := []testhelper.TestCase{
		newPolicyTest("logger_leaking", []string{"ruby/logger_leaking.rb"}),
		newPolicyTest("http", []string{"ruby/http.rb"}),
		newPolicyTest("insecure_smtp", []string{"ruby/insecure_smtp.rb"}),
		newPolicyTest("insecure_communication", []string{"ruby/insecure_communication.rb"}),
		newPolicyTest("insecure_ftp", []string{"ruby/insecure_ftp.rb"}),
		newPolicyTest("sending_data_in_category_to_third_party", []string{"ruby/sending_data_in_category_to_third_party.rb"}),
		newPolicyTest("application_level_encryption_missing_structure_sql", []string{"ruby/application_level_encryption_missing/structure_sql"}),
		newPolicyTest("application_level_encryption_missing_schema_rb", []string{"ruby/application_level_encryption_missing/schema_rb"}),
	}

	testhelper.RunTests(t, tests)
}
