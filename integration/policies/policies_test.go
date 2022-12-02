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
		newPolicyTest("http_get_parameters", []string{"ruby/http_get_parameters.rb"}),
		newPolicyTest("insecure_smtp_with_sensitive_data", []string{"ruby/insecure_smtp/with_sensitive_data.rb"}),
		newPolicyTest("insecure_smtp_without_sensitive_data", []string{"ruby/insecure_smtp/without_sensitive_data.rb"}),
		newPolicyTest("insecure_communication_with_sensitive_data", []string{"ruby/insecure_communication/with_sensitive_data.rb"}),
		newPolicyTest("insecure_communication_without_sensitive_data", []string{"ruby/insecure_communication/without_sensitive_data.rb"}),
		newPolicyTest("insecure_ftp_with_sensitive_data", []string{"ruby/insecure_ftp/with_sensitive_data.rb"}),
		newPolicyTest("insecure_ftp_without_sensitive_data", []string{"ruby/insecure_ftp/without_sensitive_data.rb"}),
	}

	testhelper.RunTests(t, tests)
}
