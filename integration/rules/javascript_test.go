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

func TestJavascriptGoogleAnalyticsDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_analytics", "dataflow", "javascript_google_analytics", t)
}

func TestJavascriptGoogleAnalyticsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/google_analytics", "summary", "javascript_google_analytics", t)
}

func TestJavascriptReactGoogleAnalyticsDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/react/google_analytics", "dataflow", "javascript_react_google_analytics", t)
}

func TestJavascriptReactGoogleAnalyticsSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/react/google_analytics", "summary", "javascript_react_google_analytics", t)
}

func TestJavascriptAlgoliaDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/algolia", "dataflow", "javascript_third_parties_algolia", t)
}

func TestJavascriptAlgoliaSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/algolia", "summary", "javascript_third_parties_algolia", t)
}

func TestJavascriptDataDogDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog", "dataflow", "javascript_third_parties_datadog", t)
}

func TestJavascriptDataDogSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog", "summary", "javascript_third_parties_datadog", t)
}

func TestJavascriptDataDogBrowserDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog_browser", "dataflow", "javascript_third_parties_datadog_browser", t)
}

func TestJavascriptDataDogBrowserSummary(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/datadog_browser", "summary", "javascript_third_parties_datadog_browser", t)
}
func TestJavascriptElasticSearchDataflow(t *testing.T) {
	t.Parallel()
	runRulesTest("javascript/third_parties/elasticsearch", "dataflow", "javascript_elasticsearch", t)
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
