package integration_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
)

var rulesFs = &rules.Rules

func buildRulesTestCase(name, reportType, filename string) testhelper.TestCase {
	arguments := []string{
		"scan",
		filepath.Join("pkg", "commands", "process", "settings", "rules", filename),
		"--report=" + reportType,
		"--format=yaml",
	}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestRubyLangCookiesSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/cookies", "summary", t)
}

func TestRubyLangCookiesDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/cookies", "dataflow", t)
}

func TestRubyLangFileGenerationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/file_generation", "summary", t)
}

func TestRubyLangFileGenerationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/file_generation", "dataflow", t)
}

func TestRubyLangHttpGetParamsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_get_params", "summary", t)
}

func TestRubyLangHttpGetParamsDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_get_params", "dataflow", t)
}

func TestRubyLangHttpInsecureSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_insecure", "summary", t)
}

func TestRubyLangHttpInsecureDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_insecure", "dataflow", t)
}

func TestRubyLangHttpPostInsecureWithDataSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_post_insecure_with_data", "summary", t)
}

func TestRubyLangHttpPostInsecureWithDataDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_post_insecure_with_data", "dataflow", t)
}

func TestRubyLangJwtSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/jwt", "summary", t)
}

func TestRubyLangJwtDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/jwt", "dataflow", t)
}

func TestRubyLangLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/logger", "summary", t)
}

func TestRubyLangLoggerDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/logger", "dataflow", t)
}

func TestRubyLangSslVerificationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/ssl_verification", "summary", t)
}

func TestRubyLangSslVerificationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/ssl_verification", "dataflow", t)
}

func TestRubyLangWeakEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption", "summary", t)
}

func TestRubyLangWeakEncryptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption", "dataflow", t)
}

func TestRubyLangWeakEncryptionWithDataSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption_with_data", "summary", t)
}

func TestRubyLangWeakEncryptionWithDataDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption_with_data", "dataflow", t)
}

func TestRubyRailsDefaultEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/default_encryption", "summary", t)
}

func TestRubyRailsDefaultEncryptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/default_encryption", "dataflow", t)
}

func TestRubyRailsInsecureCommunicationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_communication", "summary", t)
}

func TestRubyRailsInsecureCommunicationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_communication", "dataflow", t)
}

func TestRubyRailsInsecureFtpSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_ftp", "summary", t)
}

func TestRubyRailsInsecureFtpDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_ftp", "dataflow", t)
}

func TestRubyRailsInsecureSmtpSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_smtp", "summary", t)
}

func TestRubyRailsInsecureSmtpDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_smtp", "dataflow", t)
}

func TestRubyRailsLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/logger", "summary", t)
}

func TestRubyRailsLoggerDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/logger", "dataflow", t)
}

func TestRubyRailsPasswordLengthSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/password_length", "summary", t)
}

func TestRubyRailsPasswordLengthDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/password_length", "dataflow", t)
}

func TestRubyRailsSessionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/session", "summary", t)
}

func TestRubyRailsSessionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/session", "dataflow", t)
}

func TestRubyThirdPartiesSentrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/sentry", "summary", t)
}

func TestRubyThirdPartiesSentryDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/sentry", "dataflow", t)
}

func runRulesTest(folderPath string, format string, t *testing.T) {
	snapshotDirectory := "../../pkg/commands/process/settings/rules/" + folderPath + "/.snapshots"
	testdataDirEntries, err := rulesFs.ReadDir(fmt.Sprintf("%s/testdata", folderPath))
	if err != nil {
		t.Fatalf("failed to read rules/%s dir %e", folderPath, err)
	}

	dataflowTests := []testhelper.TestCase{}
	for _, testdataFile := range testdataDirEntries {
		name := testdataFile.Name()

		testName := strings.Replace(fmt.Sprintf("%s_%s_%s", format, folderPath, name), "/", "_", -1)
		dataflowTests = append(dataflowTests,
			buildRulesTestCase(
				testName,
				format,
				fmt.Sprintf("%s/testdata/%s", folderPath, name),
			),
		)
	}

	testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, snapshotDirectory)
}
