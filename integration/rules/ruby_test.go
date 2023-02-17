package integration_test

import (
	"testing"
)

const rubyRulesPath string = "../../pkg/commands/process/settings/rules/ruby/"

func TestRubyLangCookiesSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/cookies")
}

func TestRubyLangDeserializationOfUserInputSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/deserialization_of_user_input")
}

func TestRubyLangEvalUsingUserInput(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/eval_using_user_input")
}

func TestRubyLangFileGenerationSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/file_generation")
}

func TestRubyLangHttpGetParamsSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_get_params")
}

func TestRubyLangHttpInsecureSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_insecure")
}

func TestRubyLangHttpPostInsecureWithDataSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/http_post_insecure_with_data")
}

func TestRubyLangInsecureFtpSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/insecure_ftp")
}

func TestRubyLangJwtSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/jwt")
}

func TestRubyLangLoggerSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/logger")
}

func TestRubyLangExceptionSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/exception")
}

func TestRubyLangSslVerificationSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/ssl_verification")
}

func TestRubyLangWeakEncryptionSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/weak_encryption")
}

func TestRubyLangWeakEncryptionWithDataSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"lang/weak_encryption_with_data")
}

func TestRubyRailsDefaultEncryptionSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/default_encryption")
}

func TestRubyRailsInsecureCommunicationSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/insecure_communication")
}

func TestRubyRailsInsecureSmtpSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/insecure_smtp")
}

func TestRubyRailsLoggerSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/logger")
}

func TestRubyRailsPasswordLengthSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/password_length")
}

func TestRubyRailsSessionSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"rails/session")
}

func TestRubyThirdPartiesAlgoliaSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/algolia")
}

func TestRubyThirdPartiesBigQuerySummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/bigquery")
}

func TestRubyThirdPartiesDatadogSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/datadog")
}

func TestRubyThirdPartiesElasticsearchSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/elasticsearch")
}

func TestRubyThirdPartiesNewRelicSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/new_relic")
}

func TestRubyThirdPartiesRollbarSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/rollbar")
}

func TestRubyThirdPartiesScoutAPMSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/scout_apm")
}

func TestRubyThirdPartiesSentrySummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/sentry")
}

func TestRubyThirdPartiesBugsnagSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/bugsnag")
}

func TestRubyThirdPartiesHoneybadgerSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/honeybadger")
}

func TestRubyThirdPartiesAirbrakeSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/airbrake")
}

func TestRubyThirdPartiesOpenTelemetrySummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/open_telemetry")
}

func TestRubyThirdPartiesSegmentSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/segment")
}

func TestRubyThirdPartiesGoogleDataflowSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/google_dataflow")
}

func TestRubyThirdPartiesGoogleAnalyticsSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/google_analytics")
}

func TestRubyThirdPartiesClickHouseSummary(t *testing.T) {
	getRunner(t).runTest(t, rubyRulesPath+"third_parties/clickhouse")
}
