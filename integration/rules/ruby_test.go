package integration_test

import (
	"testing"
)

const rubyRulesPath string = "../../pkg/commands/process/settings/rules/ruby/"

func TestRubyLangCookies(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/cookies")
}

func TestRubyLangDeserializationOfUserInput(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/deserialization_of_user_input")
}

func TestRubyLangEvalUsingUserInput(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/eval_using_user_input")
}

func TestRubyLangFileGeneration(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/file_generation")
}

func TestRubyLangFtpUsingUserInput(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/ftp_using_user_input")
}

func TestRubyLangHttpGetParams(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_get_params")
}

func TestRubyLangHttpInsecure(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_insecure")
}

func TestRubyLangHttpPostInsecureWithData(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_post_insecure_with_data")
}

func TestRubyLangHttpUrlUsingUserInput(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_url_using_user_input")
}

func TestRubyLangInsecureFtp(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/insecure_ftp")
}

func TestRubyLangJwt(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/jwt")
}

func TestRubyLangLogger(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/logger")
}

func TestRubyLangException(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/exception")
}

func TestRubyLangPathUsingUserInput(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/path_using_user_input")
}

func TestRubyLangSslVerification(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/ssl_verification")
}

func TestRubyLangWeakEncryption(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/weak_encryption")
}

func TestRubyLangWeakEncryptionWithData(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"lang/weak_encryption_with_data")
}

func TestRubyRailsDefaultEncryption(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/default_encryption")
}

func TestRubyRailsInsecureCommunication(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/insecure_communication")
}

func TestRubyRailsInsecureHTTPPassowrd(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/insecure_http_password")
}

func TestRubyRailsInsecureSmtp(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/insecure_smtp")
}

func TestRubyRailsLogger(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/logger")
}

func TestRubyRailsPasswordLength(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/password_length")
}

func TestRubyRailsSession(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/session")
}

func TestRubyRailsSessionKeyUsingUserInput(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/session_key_using_user_input")
}

func TestRubyRailsOpenRedirect(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"rails/open_redirect")
}

func TestRubyThirdPartiesAlgolia(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/algolia")
}

func TestRubyThirdPartiesBigQuery(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/bigquery")
}

func TestRubyThirdPartiesDatadog(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/datadog")
}

func TestRubyThirdPartiesElasticsearch(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/elasticsearch")
}

func TestRubyThirdPartiesNewRelic(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/new_relic")
}

func TestRubyThirdPartiesRollbar(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/rollbar")
}

func TestRubyThirdPartiesScoutAPM(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/scout_apm")
}

func TestRubyThirdPartiesSentry(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/sentry")
}

func TestRubyThirdPartiesBugsnag(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/bugsnag")
}

func TestRubyThirdPartiesHoneybadger(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/honeybadger")
}

func TestRubyThirdPartiesAirbrake(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/airbrake")
}

func TestRubyThirdPartiesOpenTelemetry(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/open_telemetry")
}

func TestRubyThirdPartiesSegment(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/segment")
}

func TestRubyThirdPartiesGoogleDataflow(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/google_dataflow")
}

func TestRubyThirdPartiesGoogleAnalytics(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/google_analytics")
}

func TestRubyThirdPartiesClickHouse(t *testing.T) {
	t.Parallel()
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/clickhouse")
}
