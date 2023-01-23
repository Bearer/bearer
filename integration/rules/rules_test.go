package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
	"github.com/rs/zerolog/log"
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
	rubyDirs, err := rulesFs.ReadDir("ruby")
	if err != nil {
		log.Error().Msgf("failed to read rules/ruby dir %e", err)
	}

	for _, rubyDir := range rubyDirs {
		rubyDirName := rubyDir.Name() // e.g. lang, rails, third_parties
		dirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName)
		if err != nil {
			log.Error().Msgf("failed to read rules/ruby/%s dir %e", rubyDirName, err)
		}

		for _, dirEntry := range dirEntries {
			dirEntryName := dirEntry.Name() // e.g. cookies, cookies.yml
			ext := filepath.Ext(dirEntryName)

			if ext != "" {
				continue // not a folder
			}
			policyTests := []testhelper.TestCase{}
			dataflowTests := []testhelper.TestCase{}

			testdataDirEntries, err := rulesFs.ReadDir("ruby/" + rubyDirName + "/" + dirEntryName + "/testdata")

			if err != nil {
				// FIXME: do we want to fail silently here?
				log.Error().Msgf("failed to read rules/ruby/%s/%s dir %e", rubyDirName, dirEntryName, err)
			}

			for _, testdataFile := range testdataDirEntries {
				name := testdataFile.Name()

				dataflowTests = append(dataflowTests, buildTestCase("dataflow_"+dirEntryName+"_"+name, "dataflow", "ruby/"+rubyDirName+"/"+dirEntryName+"/testdata/"+name))
				policyTests = append(policyTests, buildTestCase("policy_"+dirEntryName+"_"+name, "policies", "ruby/"+rubyDirName+"/"+dirEntryName+"/testdata/"+name))
			}
			testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, "ruby/"+rubyDirName+"/"+dirEntryName+"/.snapshots")
			testhelper.RunTestsWithSnapshotSubdirectory(t, policyTests, "ruby/"+rubyDirName+"/"+dirEntryName+"/.snapshots")
		}
	}
}
