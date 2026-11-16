package rules

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
)

const rubyRuleYAML = `
languages:
  - ruby
patterns:
  - pattern: foo
metadata:
  id: ruby_test_rule
  description: test
`

func TestLoadDefinitionsLanguageFilter(t *testing.T) {
	dir := fstest.MapFS{
		"rules/ruby_test_rule.yml": &fstest.MapFile{Data: []byte(rubyRuleYAML)},
	}

	t.Run("drops rules for languages that were not found", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		count, err := loadDefinitions(definitions, []ruleDir{{path: "rules", fs: dir}}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.NotContains(t, definitions, "ruby_test_rule")
	})

	t.Run("keeps rules for languages that were found", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		_, err := loadDefinitions(definitions, []ruleDir{{path: "rules", fs: dir}}, false, []string{"ruby"}, nil)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "ruby_test_rule")
	})

	t.Run("keeps rules when only-rule is set, regardless of language", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		onlyRule := map[string]bool{"ruby_test_rule": true}

		_, err := loadDefinitions(definitions, []ruleDir{{path: "rules", fs: dir}}, false, []string{"python"}, onlyRule)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "ruby_test_rule")
	})
}

const sharedImportRuleYAML = `
type: shared
languages:
  - python
patterns:
  - import $<!>$<MODULE>
metadata:
  id: python_shared_test_import
  description: test shared rule
`

const customRuleImportingSharedYAML = `
imports:
  - python_shared_test_import
languages:
  - python
patterns:
  - pattern: $<MODULE>.call()
    filters:
      - variable: MODULE
        detection: python_shared_test_import
        scope: cursor
metadata:
  id: custom_python_test_rule
  description: test custom rule
`

const customRuleImportingUnknownYAML = `
imports:
  - python_shared_does_not_exist
languages:
  - python
patterns:
  - pattern: $<MODULE>.call()
    filters:
      - variable: MODULE
        detection: python_shared_does_not_exist
        scope: cursor
metadata:
  id: custom_python_unknown_import
  description: test custom rule
`

// A shared rule may not set a severity, so this one fails validation.
const invalidSharedImportRuleYAML = `
type: shared
severity: high
languages:
  - python
patterns:
  - import $<!>$<MODULE>
metadata:
  id: python_shared_test_import
  description: invalid shared rule
`

const riskRuleYAML = `
languages:
  - python
patterns:
  - import $<!>$<MODULE>
metadata:
  id: python_shared_test_import
  description: test rule that is not shared
`

const sharedImportRuleV2YAML = `
type: shared
languages:
  - python
patterns:
  - from $<!>$<MODULE> import $<_>
metadata:
  id: python_shared_test_import
  description: second version of the shared rule
`

func TestLoadDefinitionsImportsAcrossDirectories(t *testing.T) {
	builtInDir := ruleDir{path: "built-in", fs: fstest.MapFS{
		"python/shared/import.yml": &fstest.MapFile{Data: []byte(sharedImportRuleYAML)},
	}}
	tenantDir := ruleDir{path: "tenant", fs: fstest.MapFS{
		"custom_python_test_rule.yml": &fstest.MapFile{Data: []byte(customRuleImportingSharedYAML)},
	}}

	t.Run("keeps a rule that imports a shared rule from another directory", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		count, err := loadDefinitions(definitions, []ruleDir{builtInDir, tenantDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Contains(t, definitions, "custom_python_test_rule")
		assert.Contains(t, definitions, "python_shared_test_import")
	})

	t.Run("does not depend on the order of the directories", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		_, err := loadDefinitions(definitions, []ruleDir{tenantDir, builtInDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "custom_python_test_rule")
	})

	t.Run("does not let a rule import a non-shared rule from another directory", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		riskDir := ruleDir{path: "built-in", fs: fstest.MapFS{
			"python/import.yml": &fstest.MapFile{Data: []byte(riskRuleYAML)},
		}}

		_, err := loadDefinitions(definitions, []ruleDir{riskDir, tenantDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.NotContains(t, definitions, "custom_python_test_rule")
	})

	t.Run("does not let a rule import a rule that is not in an external directory", func(t *testing.T) {
		definitions := map[string]settings.RuleDefinition{
			"python_shared_test_import": {
				Type:     customdetectors.TypeShared,
				Metadata: &settings.RuleMetadata{ID: "python_shared_test_import"},
			},
		}

		_, err := loadDefinitions(definitions, []ruleDir{tenantDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.NotContains(t, definitions, "custom_python_test_rule")
	})

	t.Run("still drops a rule that imports an unknown rule", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		unknownDir := ruleDir{path: "tenant", fs: fstest.MapFS{
			"custom_python_unknown_import.yml": &fstest.MapFile{Data: []byte(customRuleImportingUnknownYAML)},
		}}

		_, err := loadDefinitions(definitions, []ruleDir{builtInDir, unknownDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.NotContains(t, definitions, "custom_python_unknown_import")
	})

	t.Run("drops a rule that imports an invalid shared rule from another directory", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		invalidSharedDir := ruleDir{path: "built-in", fs: fstest.MapFS{
			"python/shared/import.yml": &fstest.MapFile{Data: []byte(invalidSharedImportRuleYAML)},
		}}

		_, err := loadDefinitions(definitions, []ruleDir{invalidSharedDir, tenantDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.NotContains(t, definitions, "python_shared_test_import")
		assert.NotContains(t, definitions, "custom_python_test_rule")
	})

	t.Run("reports an invalid shared rule once", func(t *testing.T) {
		var logOutput bytes.Buffer
		previousLogger, previousLevel := log.Logger, zerolog.GlobalLevel()
		log.Logger = zerolog.New(&logOutput)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		defer func() {
			log.Logger = previousLogger
			zerolog.SetGlobalLevel(previousLevel)
		}()

		definitions := make(map[string]settings.RuleDefinition)
		invalidSharedDir := ruleDir{path: "built-in", fs: fstest.MapFS{
			"python/shared/import.yml": &fstest.MapFile{Data: []byte(invalidSharedImportRuleYAML)},
		}}

		_, err := loadDefinitions(definitions, []ruleDir{invalidSharedDir, tenantDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Equal(t, 1, strings.Count(logOutput.String(), "python_shared_test_import ignored due to validation errors"))
	})

	t.Run("offers the valid copy of a shared rule defined in two directories", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		invalidCopyDir := ruleDir{path: "other", fs: fstest.MapFS{
			"python/shared/import.yml": &fstest.MapFile{Data: []byte(invalidSharedImportRuleYAML)},
		}}

		_, err := loadDefinitions(
			definitions,
			[]ruleDir{builtInDir, invalidCopyDir, tenantDir},
			false,
			[]string{"python"},
			nil,
		)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "custom_python_test_rule")
		assert.Equal(t, "test shared rule", definitions["python_shared_test_import"].Metadata.Description)
	})

	t.Run("still lets a later directory replace a rule with the same ID", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		laterDir := ruleDir{path: "tenant", fs: fstest.MapFS{
			"copy.yml": &fstest.MapFile{Data: []byte(sharedImportRuleV2YAML)},
		}}

		_, err := loadDefinitions(definitions, []ruleDir{builtInDir, laterDir}, false, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Equal(t, "second version of the shared rule", definitions["python_shared_test_import"].Metadata.Description)
	})
}
