package integration_test

import "testing"

func TestJavascriptLangLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/logger", "summary", "javascript_lang_logger", t)
}

func TestJavascriptLangLoggerDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/logger", "dataflow", "javascript_lang_logger", t)
}

func TestJavascriptLangSessionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/session", "summary", "javascript_session", t)
}

func TestJavascriptLangSessionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/session", "dataflow", "javascript_session", t)
}

func TestJavascriptWeakEncryptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/weak_encryption", "dataflow", "javascript_weak_encryption", t)
}

func TestJavascriptWeakEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/weak_encryption", "summary", "javascript_weak_encryption", t)
}

func TestExpressSecureCookieDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/express/insecure_cookie", "dataflow", "express_insecure_cookie", t)
}

func TestExpressSecureCookieSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/express/insecure_cookie", "summary", "express_insecure_cookie", t)
}

func TestJavascriptJWTDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/jwt", "dataflow", "javascript_jwt", t)
}

func TestJavascriptJWTSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/jwt", "summary", "javascript_jwt", t)
}

func TestJavascriptThirdPartySentrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/sentry", "summary", "javascript_third_parties_sentry", t)
}

func TestJavascriptThirdPartySentryDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/sentry", "dataflow", "javascript_third_parties_sentry", t)
}

func TestJavascriptLangExceptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/exception", "summary", "javascript_lang_exception", t)
}

func TestJavascriptLangExceptionDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/exception", "dataflow", "javascript_lang_exception", t)
}

func TestJavascriptLangFileGenerationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/file_generation", "summary", "javascript_lang_file_generation", t)
}

func TestJavascriptLangFileGenerationDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/file_generation", "dataflow", "javascript_lang_file_generation", t)
}

func TestJavascriptGTMDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_tag_manager", "dataflow", "javascript_google_tag_manager", t)
}

func TestJavascriptGTMSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_tag_manager", "summary", "javascript_google_tag_manager", t)
}
