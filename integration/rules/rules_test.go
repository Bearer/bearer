package integration_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
)

var rulesFs = &rules.Rules

func buildTestCase(name, reportType, filename string) testhelper.TestCase {
	arguments := []string{
		"scan",
		filepath.Join("pkg", "commands", "process", "settings", "rules", filename),
		"--report=" + reportType,
		"--format=yaml",
	}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestRubyRules(t *testing.T) {
	rubyDirs, err := rulesFs.ReadDir("ruby")
	if err != nil {
		t.Fatalf("failed to read rules/ruby dir %e", err)
	}

	for _, rubyDir := range rubyDirs {
		rubyDirName := rubyDir.Name() // e.g. lang, rails, third_parties
		dirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName)
		if err != nil {
			t.Fatalf("failed to read rules/ruby/%s dir %e", rubyDirName, err)
		}

		for _, dirEntry := range dirEntries {
			dirEntryName := dirEntry.Name() // e.g. cookies, cookies.yml
			ext := filepath.Ext(dirEntryName)

			if ext != "" {
				continue // folder
			}

			rubySubPath := "ruby/" + rubyDirName + "/" + dirEntryName
			snapshotDirectory := "../../pkg/commands/process/settings/rules/" + rubySubPath + "/.snapshots"
			testdataDirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName + "/" + dirEntryName + "/testdata")
			if err != nil {
				t.Fatalf("failed to read rules/ruby/%s/%s dir %e", rubyDirName, dirEntryName, err)
			}

			summaryTests := []testhelper.TestCase{}
			dataflowTests := []testhelper.TestCase{}

			for _, testdataFile := range testdataDirEntries {
				name := testdataFile.Name()

				dataflowTests = append(dataflowTests,
					buildTestCase(
						fmt.Sprintf("dataflow_%s_%s_%s", rubyDirName, dirEntryName, name),
						"dataflow",
						fmt.Sprintf("%s/testdata/%s", rubySubPath, name),
					),
				)
				summaryTests = append(summaryTests,
					buildTestCase(
						fmt.Sprintf("summary_%s_%s_%s", rubyDirName, dirEntryName, name),
						"summary",
						fmt.Sprintf("%s/testdata/%s", rubySubPath, name),
					),
				)
			}

			testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, snapshotDirectory)
			testhelper.RunTestsWithSnapshotSubdirectory(t, summaryTests, snapshotDirectory)
		}
	}
}
