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
