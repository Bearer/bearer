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

func TestJavascriptHTTPInsecureDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/http_insecure", "dataflow", "javascript_http_insecure", t)
}

func TestJavascriptHTTPInsecureSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/http_insecure", "summary", "javascript_http_insecure", t)
}

func TestJavascriptThirdPartySentrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/sentry", "summary", "javascript_third_parties_sentry", t)
}

func TestJavascriptThirdPartySentryDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/sentry", "dataflow", "javascript_third_parties_sentry", t)
}
