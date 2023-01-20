package integration_test

import (
	"path/filepath"
	"strings"
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
	policyTests := []testhelper.TestCase{}
	dataflowTests := []testhelper.TestCase{}

	rubyDirs, err := rulesFs.ReadDir("ruby")
	if err != nil {
		log.Error().Msgf("failed to read rules/ruby dir %e", err)
	}

	for _, rubyDir := range rubyDirs {
		subDir := rubyDir.Name()
		dirEntries, err := rulesFs.ReadDir("ruby/" + subDir)
		if err != nil {
			log.Error().Msgf("failed to read rules/ruby/%s dir %e", subDir, err)
		}

		for _, dirEntry := range dirEntries {
			filename := dirEntry.Name()
			ext := filepath.Ext(filename)
			name := strings.TrimSuffix(filename, ext)

			if ext == ".yaml" || ext == ".yml" {
				// it's a YAML-format rule
				continue
			}

			if ext == "" && !strings.HasSuffix(name, "testdata") {
				// it's a folder we don't care about
				continue
			}

			dataflowTests = append(dataflowTests, buildTestCase("dataflow_"+subDir+"_"+name, "dataflow", "ruby/"+subDir+"/"+filename))
			policyTests = append(policyTests, buildTestCase("policy_"+subDir+"_"+name, "policies", "ruby/"+subDir+"/"+filename))
		}
	}

	testhelper.RunTests(t, dataflowTests)
	testhelper.RunTests(t, policyTests)
}
