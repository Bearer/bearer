package integration_test

import "testing"

const javascriptRulesPath string = "../../pkg/commands/process/settings/rules/javascript/"

func TestJavascriptLangLogger(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/logger")
}

func TestJavascriptLangSession(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/session")
}

func TestJavascriptWeakEncryption(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/weak_encryption")
}

func TestJavascriptJWT(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/jwt")
}

func TestJavascriptJWTWeakEncryption(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/jwt_weak_encryption")
}

func TestJavascriptJWTHardcodedSecret(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/jwt_hardcoded_secret")
}

func TestJavascriptHTTPInsecure(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/http_insecure")
}

func TestJavascriptLangException(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/exception")
}

func TestJavascriptLangFileGeneration(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"lang/file_generation")
}

func TestJavascriptExpressOpenRedirect(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/open_redirect")
}

func TestJavascriptExpressUnsafeDeserialization(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/unsafe_deserialization")
}

func TestJavascriptExpressInsecureRefResolution(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/insecure_ref_resolution")
}

func TestJavascriptExpressExposedDirListing(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/exposed_dir_listing")
}

func TestJavascriptExpressCrossSiteScripting(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/cross_site_scripting")
}

func TestJavascriptExpressPathTraversal(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/path_traversal")
}

func TestJavascriptExpressServerSideRequestForgery(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/server_side_request_forgery")
}

func TestJavascriptExpressSqlInjection(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/sql_injection")
}

func TestExpressSecureCookie(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/insecure_cookie")
}

func TestExpressXXEVulnerability(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"express/xml_external_entity_vulnerability")
}

func TestJavascriptReactGoogleAnalytics(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"react/google_analytics")
}

func TestJavascriptThirdPartySentry(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/sentry")
}

func TestJavascriptGTM(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_tag_manager")
}

func TestJavascriptGoogleAnalytics(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/google_analytics")
}

func TestJavascriptAlgolia(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/algolia")
}

func TestJavascriptDataDog(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog")
}

func TestJavascriptDataDogBrowser(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/datadog_browser")
}

func TestJavascriptElasticSearch(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/elasticsearch")
}

func TestJavascriptSegmentDataflow(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/segment")
}

func TestJavascriptNewRelic(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/new_relic")
}

func TestJavascriptRollbar(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/rollbar")
}

func TestJavascriptHoneybadger(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/honeybadger")
}

func TestJavascriptAirbrake(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/airbrake")
}

func TestJavascriptOpenTelemetry(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/open_telemetry")
}

func TestJavascriptBugsnag(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, javascriptRulesPath+"third_parties/bugsnag")
}
