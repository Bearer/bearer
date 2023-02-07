package integration_test

import (
	"testing"
)

func TestRubyLangCookiesSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/cookies", "summary", "ruby_lang_cookies", t)
}

func TestRubyLangCookiesDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/cookies", "dataflow", "ruby_lang_cookies", t)
}

func TestRubyLangFileGenerationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/file_generation", "summary", "ruby_lang_file_generation", t)
}

func TestRubyLangFileGenerationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/file_generation", "dataflow", "ruby_lang_file_generation", t)
}

func TestRubyLangHttpGetParamsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_get_params", "summary", "ruby_lang_http_get_params", t)
}

func TestRubyLangHttpGetParamsDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_get_params", "dataflow", "ruby_lang_http_get_params", t)
}

func TestRubyLangHttpInsecureSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_insecure", "summary", "ruby_lang_http_insecure", t)
}

func TestRubyLangHttpInsecureDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_insecure", "dataflow", "ruby_lang_http_insecure", t)
}

func TestRubyLangHttpPostInsecureWithDataSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_post_insecure_with_data", "summary", "ruby_lang_http_post_insecure_with_data", t)
}

func TestRubyLangHttpPostInsecureWithDataDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/http_post_insecure_with_data", "dataflow", "ruby_lang_http_post_insecure_with_data", t)
}

func TestRubyLangInsecureFtpSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/insecure_ftp", "summary", "ruby_lang_insecure_ftp", t)
}

func TestRubyLangInsecureFtpDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/insecure_ftp", "dataflow", "ruby_lang_insecure_ftp", t)
}

func TestRubyLangJwtSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/jwt", "summary", "ruby_lang_jwt", t)
}

func TestRubyLangJwtDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/jwt", "dataflow", "ruby_lang_jwt", t)
}

func TestRubyLangLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/logger", "summary", "ruby_lang_logger", t)
}

func TestRubyLangLoggerDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/logger", "dataflow", "ruby_lang_logger", t)
}

func TestRubyLangExceptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/exception", "summary", "ruby_lang_exception", t)
}

func TestRubyLangExceptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/exception", "dataflow", "ruby_lang_exception", t)
}

func TestRubyLangSslVerificationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/ssl_verification", "summary", "ruby_lang_ssl_verification", t)
}

func TestRubyLangSslVerificationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/ssl_verification", "dataflow", "ruby_lang_ssl_verification", t)
}

func TestRubyLangWeakEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption", "summary", "ruby_lang_weak_encryption", t)
}

func TestRubyLangWeakEncryptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption", "dataflow", "ruby_lang_weak_encryption", t)
}

func TestRubyLangWeakEncryptionWithDataSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption_with_data", "summary", "ruby_lang_weak_encryption_with_data", t)
}

func TestRubyLangWeakEncryptionWithDataDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/lang/weak_encryption_with_data", "dataflow", "ruby_lang_weak_encryption_with_data", t)
}

func TestRubyRailsDefaultEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/default_encryption", "summary", "ruby_rails_default_encryption", t)
}

func TestRubyRailsDefaultEncryptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/default_encryption", "dataflow", "ruby_rails_default_encryption", t)
}

func TestRubyRailsInsecureCommunicationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_communication", "summary", "ruby_rails_insecure_communication", t)
}

func TestRubyRailsInsecureCommunicationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_communication", "dataflow", "ruby_rails_insecure_communication", t)
}

func TestRubyRailsInsecureSmtpSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_smtp", "summary", "ruby_rails_insecure_smtp", t)
}

func TestRubyRailsInsecureSmtpDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/insecure_smtp", "dataflow", "ruby_rails_insecure_smtp", t)
}

func TestRubyRailsLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/logger", "summary", "ruby_rails_logger", t)
}

func TestRubyRailsLoggerDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/logger", "dataflow", "ruby_rails_logger", t)
}

func TestRubyRailsPasswordLengthSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/password_length", "summary", "ruby_rails_password_length", t)
}

func TestRubyRailsPasswordLengthDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/password_length", "dataflow", "ruby_rails_password_length", t)
}

func TestRubyRailsSessionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/session", "summary", "ruby_rails_session", t)
}

func TestRubyRailsSessionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/rails/session", "dataflow", "ruby_rails_session", t)
}

func TestRubyThirdPartiesNewRelicSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/new_relic", "summary", "ruby_third_parties_new_relic", t)
}

func TestRubyThirdPartiesNewRelicDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/new_relic", "dataflow", "ruby_third_parties_new_relic", t)
}

func TestRubyThirdPartiesSentrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/sentry", "summary", "ruby_third_parties_sentry", t)
}

func TestRubyThirdPartiesSentryDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("ruby/third_parties/sentry", "dataflow", "ruby_third_parties_sentry", t)
}
