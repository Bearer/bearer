package integration_test

import (
	"os"
	"testing"
)

func TestRuby(t *testing.T) {
	rulesPath, _ := os.LookupEnv("RULES_PATH") // defaults to "" if not present
	var rubyRulesPath string = rulesPath + "/ruby/"

	tests := []RuleTestCase{}
	entries, err := os.ReadDir(rubyRulesPath)
	if err != nil {
		t.Fatalf("failed to read /ruby folder: %s", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		ruleDirs, err := os.ReadDir(rubyRulesPath + entry.Name())
		if err != nil {
			t.Fatalf("failed to read /ruby/%s folder: %s", rubyRulesPath+entry.Name(), err)
		}
		for _, ruleDir := range ruleDirs {
			if !ruleDir.IsDir() {
				continue
			}
			tests = append(tests, RuleTestCase{
				ProjectPath: rubyRulesPath + entry.Name() + "/" + ruleDir.Name(),
			})
		}
	}

	t.Parallel()
	runner := getRunner(t)
	for _, testCase := range tests {
		runner.runTest(t, testCase.ProjectPath)
	}
}
