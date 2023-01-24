package integration_test

import (
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

// test dataflow and policies
func TestRubyRules(t *testing.T) {
	// rubyDirs, err := rulesFs.ReadDir("ruby")
	// if err != nil {
	// 	t.Fatalf("failed to read rules/ruby dir %e", err)
	// }

	// for _, rubyDir := range rubyDirs {
	// 	rubyDirName := rubyDir.Name() // e.g. lang, rails, third_parties
	// dirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName)
	// if err != nil {
	// 	t.Fatalf("failed to read rules/ruby/%s dir %e", rubyDirName, err)
	// }

	rubyDirName := "lang"
	dirEntryName := "http_insecure"

	// for _, dirEntry := range dirEntries {
	// 	dirEntryName := dirEntry.Name() // e.g. cookies, cookies.yml
	// 	ext := filepath.Ext(dirEntryName)

	// 	if ext != "" {
	// 		continue // not a folder
	// 	}
	policyTests := []testhelper.TestCase{}
	dataflowTests := []testhelper.TestCase{}

	testdataDirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName + "/" + dirEntryName + "/testdata")

	if err != nil {
		t.Fatalf("failed to read rules/ruby/%s/%s dir %e", rubyDirName, dirEntryName, err)
	}

	rubySubPath := "ruby/" + rubyDirName + "/" + dirEntryName

	for _, testdataFile := range testdataDirEntries {
		name := testdataFile.Name()

		dataflowTests = append(dataflowTests, buildTestCase("dataflow_"+rubyDirName+"_"+dirEntryName+"_"+name, "dataflow", rubySubPath+"/testdata/"+name))
		policyTests = append(policyTests, buildTestCase("policy_"+rubyDirName+"_"+dirEntryName+"_"+name, "policies", rubySubPath+"/testdata/"+name))
	}

	snapshotDirectory := "../../pkg/commands/process/settings/rules/" + rubySubPath + "/.snapshots"
	testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, snapshotDirectory)
	testhelper.RunTestsWithSnapshotSubdirectory(t, policyTests, snapshotDirectory)
}

// }

// }
