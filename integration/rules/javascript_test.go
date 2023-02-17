package integration_test

import "testing"

func TestJavascriptLangLoggerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/logger", "summary", "javascript_lang_logger", t)
}

func TestJavascriptLangSessionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/session", "summary", "javascript_session", t)
}

func TestJavascriptWeakEncryptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/weak_encryption", "summary", "javascript_weak_encryption", t)
}

func TestExpressSecureCookieSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/express/insecure_cookie", "summary", "express_insecure_cookie", t)
}

func TestJavascriptJWTSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/jwt", "summary", "javascript_jwt", t)
}

func TestJavascriptHTTPInsecureSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/http_insecure", "summary", "javascript_http_insecure", t)
}

func TestJavascriptThirdPartySentrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/sentry", "summary", "javascript_third_parties_sentry", t)
}

func TestJavascriptLangExceptionSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/exception", "summary", "javascript_lang_exception", t)
}

func TestJavascriptLangFileGenerationSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/lang/file_generation", "summary", "javascript_lang_file_generation", t)
}

func TestJavascriptGTMSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_tag_manager", "summary", "javascript_google_tag_manager", t)
}

func TestJavascriptGoogleAnalyticsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_analytics", "summary", "javascript_google_analytics", t)
}

func TestJavascriptReactGoogleAnalyticsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/react/google_analytics", "summary", "javascript_react_google_analytics", t)
}

func TestJavascriptAlgoliaSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/algolia", "summary", "javascript_third_parties_algolia", t)
}

func TestJavascriptDataDogSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog", "summary", "javascript_third_parties_datadog", t)
}

func TestJavascriptDataDogBrowserSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog_browser", "summary", "javascript_third_parties_datadog_browser", t)
}

func TestJavascriptElasticSearchSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/elasticsearch", "summary", "javascript_elasticsearch", t)
}

func TestJavascriptSegmentDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/segment", "dataflow", "javascript_third_parties_segment", t)
}

func TestJavascriptSegmentSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/segment", "summary", "javascript_third_parties_segment", t)
}

func TestJavascriptNewRelicSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/new_relic", "summary", "javascript_third_parties_new_relic", t)
}

func TestJavascriptRollbarSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/rollbar", "summary", "javascript_rollbar", t)
}

func TestJavascriptHoneybadgerSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/honeybadger", "summary", "javascript_honeybadger", t)
}

func TestJavascriptAirbrakeSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/airbrake", "summary", "javascript_third_parties_airbrake", t)
}

func TestJavascriptOpenTelemetrySummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/open_telemetry", "summary", "javascript_third_parties_open_telemetry", t)
}

func TestJavascriptBugsnagSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/bugsnag", "summary", "javascript_third_parties_bugsnag", t)
}
