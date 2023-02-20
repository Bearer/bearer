package rules_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/curio/e2e/internal/testhelper"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
)

var rulesFs = &rules.Rules

func buildRulesTestCase(testName, fileName, ruleID string) testhelper.TestCase {
	arguments := []string{
		"scan",
		fileName,
		"--external-rule-dir=" + filepath.Join("testdata", "rules"),
		"--only-rule=" + ruleID,
		"--format=yaml",
	}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(testName, arguments, options)
}

func runRulesTest(folderPath string, ruleID string, t *testing.T) {
	snapshotDirectory := "/.snapshots"

	testDataDir := fmt.Sprintf("testdata/data/%s", folderPath)

	testdataDirEntries, err := rulesFs.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read rules/%s dir %e", folderPath, err)
	}

	testCases := []testhelper.TestCase{}
	for _, testdataFile := range testdataDirEntries {
		filePath := testdataFile.Name()
		ext := filepath.Ext(filePath)
		testName := strings.TrimSuffix(filePath, ext)
		testName = strings.TrimPrefix(testName, testDataDir)

		testCases = append(testCases,
			buildRulesTestCase(
				testName,
				filepath.Join(testDataDir, filePath),
				ruleID,
			),
		)
	}

	testhelper.RunTestsWithSnapshotSubdirectory(t, testCases, snapshotDirectory)
}
