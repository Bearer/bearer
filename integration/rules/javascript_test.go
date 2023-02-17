package integration_test

import "testing"

const javascriptRulesPath string = "../../pkg/commands/process/settings/rules/javascript/"

func TestJavascriptLangLoggerSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/logger")
}

func TestJavascriptLangSessionSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/session")
}

func TestJavascriptWeakEncryptionSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/weak_encryption")
}

func TestJavascriptJWTSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/jwt")
}

func TestJavascriptHTTPInsecureSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/http_insecure")
}

func TestJavascriptLangExceptionSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/exception")
}

func TestJavascriptLangFileGenerationSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/file_generation")
}

func TestExpressSecureCookieSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"express/insecure_cookie")
}

func TestJavascriptReactGoogleAnalyticsSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"react/google_analytics")
}

func TestJavascriptThirdPartySentrySummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/sentry")
}

func TestJavascriptGTMSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_tag_manager")
}

func TestJavascriptGoogleAnalyticsSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_analytics")
}

func TestJavascriptAlgoliaSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/algolia")
}

func TestJavascriptDataDogSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog")
}

func TestJavascriptDataDogBrowserSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog_browser")
}

func TestJavascriptElasticSearchSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/elasticsearch")
}

func TestJavascriptSegmentDataflow(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/segment")
}

func TestJavascriptNewRelicSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/new_relic")
}

func TestJavascriptRollbarSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/rollbar")
}

func TestJavascriptHoneybadgerSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/honeybadger")
}

func TestJavascriptAirbrakeSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/airbrake")
}

func TestJavascriptOpenTelemetrySummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/javascript_third_parties_open_telemetry")
}

func TestJavascriptBugsnagSummary(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/bugsnag")
}
