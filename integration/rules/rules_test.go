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

func rulesTest(name, filename string) testhelper.TestCase {
	arguments := []string{"scan", filepath.Join("pkg", "commands", "process", "settings", "rules", filename), "--report=dataflow", "--format=yaml"}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func TestRubyRules(t *testing.T) {
	tests := []testhelper.TestCase{}

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

			if ext == "" || ext == ".yaml" || ext == ".yml" {
				// it's a YAML-format rule, or something else we don't care about
				continue
			}

			tests = append(tests, rulesTest("ruby_"+subDir+"_"+name, "ruby/"+subDir+"/"+filename))
		}
	}

	testhelper.RunTests(t, tests)
}
