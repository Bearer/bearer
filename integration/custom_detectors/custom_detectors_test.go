package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func newScanTest(language, name, filename string) testhelper.TestCase {
	arguments := []string{"scan", filepath.Join("integration", "custom_detectors", "testdata", language, filename), "--report=dataflow", "--format=yaml"}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestCustomDetectors(t *testing.T) {
	tests := []testhelper.TestCase{
		newScanTest("ruby", "detect_rails_jwt", "detect_rails_jwt.rb"),
		newScanTest("ruby", "detect_rails_session", "detect_rails_session.rb"),
		newScanTest("ruby", "detect_rails_cookies", "detect_rails_cookies.rb"),
		newScanTest("ruby", "detect_ruby_logger", "detect_ruby_logger.rb"),
		newScanTest("ruby", "ftp", "ftp.rb"),
		newScanTest("ruby", "ruby_file_detection", "ruby_file_detection.rb"),
		newScanTest("ruby", "ssl_certificate_verification_disabled", "ssl_certificate_verification_disabled.rb"),
		newScanTest("ruby", "ruby_http_detection", "ruby_http_detection.rb"),
		newScanTest("ruby", "ruby_weak_password_encryption", "weak_password_encryption.rb"),
		newScanTest("ruby", "detect_password_length", "detect_password_length.rb"),
		newScanTest("ruby", "ruby_weak_encryption_library", "weak_encryption_library.rb"),
	}

	testhelper.RunTests(t, tests)
}
