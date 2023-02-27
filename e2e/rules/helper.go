package rules_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
)

func buildRulesTestCase(testName, path, ruleID string, builtIn bool) testhelper.TestCase {
	arguments := []string{
		"scan",
		path,
		"--only-rule=" + ruleID,
		"--format=yaml",
	}

	if !builtIn {
		arguments = append(arguments,
			"--external-rule-dir="+filepath.Join("e2e", "rules", "testdata", "rules"),
		)
	}

	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(testName, arguments, options)
}

func runRulesTest(folderPath string, ruleID string, builtIn bool, t *testing.T) {
	snapshotDirectory := ".snapshots"

	testDataDir := fmt.Sprintf("testdata/data/%s", folderPath)

	testCases := []testhelper.TestCase{}
	testCases = append(testCases,
		buildRulesTestCase(
			testDataDir,
			filepath.Join("e2e", "rules", testDataDir),
			ruleID,
			builtIn,
		),
	)

	testhelper.RunTestsWithSnapshotSubdirectory(t, testCases, snapshotDirectory)
}
