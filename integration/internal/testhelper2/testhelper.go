package testhelper2

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
)

var rulesFs = &rules.Rules

func buildRulesTestCase(name, reportType, ruleID, filename string) testhelper.TestCase {
	arguments := []string{
		"scan",
		filepath.Join("pkg", "commands", "process", "settings", "rules", filename),
		"--report=" + reportType,
		"--format=yaml",
		"--only-rule=" + ruleID,
	}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func runRulesTest(folderPath, format, ruleID string, t *testing.T) {
	snapshotDirectory := "../../pkg/commands/process/settings/rules/" + folderPath + "/.snapshots"

	testDataDir := fmt.Sprintf("%s/testdata", folderPath)

	testdataDirEntries, err := rulesFs.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read rules/%s dir %e", folderPath, err)
	}

	dataflowTests := []testhelper.TestCase{}
	for _, testdataFile := range testdataDirEntries {
		name := testdataFile.Name()

		testName := strings.Replace(fmt.Sprintf("%s_%s_%s", format, folderPath, name), "/", "_", -1)
		dataflowTests = append(dataflowTests,
			buildRulesTestCase(
				testName,
				format,
				ruleID,
				fmt.Sprintf("%s/testdata/%s", folderPath, name),
			),
		)
	}

	testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, snapshotDirectory)
}
