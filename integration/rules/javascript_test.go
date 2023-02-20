package integration_test

import "testing"

const javascriptRulesPath string = "../../pkg/commands/process/settings/rules/javascript/"

func TestJavascriptLangLogger(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/logger")
}

func TestJavascriptLangSession(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/session")
}

func TestJavascriptWeakEncryption(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/weak_encryption")
}

func TestJavascriptJWT(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/jwt")
}

func TestJavascriptHTTPInsecure(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/http_insecure")
}

func TestJavascriptLangException(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/exception")
}

func TestJavascriptLangFileGeneration(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"lang/file_generation")
}

func TestJavascriptExpressExposedDirListing(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"express/exposed_dir_listing")
}

func TestExpressSecureCookie(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"express/insecure_cookie")
}

func TestJavascriptReactGoogleAnalytics(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"react/google_analytics")
}

func TestJavascriptThirdPartySentry(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/sentry")
}

func TestJavascriptGTM(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_tag_manager")
}

func TestJavascriptGoogleAnalytics(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_analytics")
}

func TestJavascriptAlgolia(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/algolia")
}

func TestJavascriptDataDog(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog")
}

func TestJavascriptDataDogBrowser(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog_browser")
}

func TestJavascriptElasticSearch(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/elasticsearch")
}

func TestJavascriptSegmentDataflow(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/segment")
}

func TestJavascriptNewRelic(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/new_relic")
}

func TestJavascriptRollbar(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/rollbar")
}

func TestJavascriptHoneybadger(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/honeybadger")
}

func TestJavascriptAirbrake(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/airbrake")
}

func TestJavascriptOpenTelemetry(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/open_telemetry")
}

func TestJavascriptBugsnag(t *testing.T) {
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/bugsnag")
}
