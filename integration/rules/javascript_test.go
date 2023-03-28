package integration_test

import (
	"os"
	"testing"
)

func TestJavascript(t *testing.T) {
	rulesPath, _ := os.LookupEnv("RULES_PATH") // defaults to "" if not present
	var javascriptRulesPath string = rulesPath + "/javascript/"

	tests := []RuleTestCase{}
	entries, err := os.ReadDir(javascriptRulesPath)
	if err != nil {
		t.Fatalf("failed to read /javascript folder: %s", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		ruleDirs, err := os.ReadDir(javascriptRulesPath + entry.Name())
		if err != nil {
			t.Fatalf("failed to read /javascript/%s folder: %s", javascriptRulesPath+entry.Name(), err)
		}
		for _, ruleDir := range ruleDirs {
			if !ruleDir.IsDir() {
				continue
			}

			tests = append(tests, RuleTestCase{
				ProjectPath: javascriptRulesPath + entry.Name() + "/" + ruleDir.Name(),
			})
		}
	}

	t.Parallel()
	runner := getRunner(t)
	for _, testCase := range tests {
		runner.runTest(t, testCase.ProjectPath)
	}
}
