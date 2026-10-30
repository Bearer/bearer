package rules

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
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

func TestLoadCustomDefinitionsLanguageFilter(t *testing.T) {
	dir := fstest.MapFS{
		"rules/ruby_test_rule.yml": &fstest.MapFile{Data: []byte(rubyRuleYAML)},
	}

	t.Run("drops rules for languages that were not found", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		count, err := loadCustomDefinitions(definitions, false, dir, []string{"python"}, nil)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.NotContains(t, definitions, "ruby_test_rule")
	})

	t.Run("keeps rules for languages that were found", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)

		_, err := loadCustomDefinitions(definitions, false, dir, []string{"ruby"}, nil)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "ruby_test_rule")
	})

	t.Run("keeps rules when only-rule is set, regardless of language", func(t *testing.T) {
		definitions := make(map[string]settings.RuleDefinition)
		onlyRule := map[string]bool{"ruby_test_rule": true}

		_, err := loadCustomDefinitions(definitions, false, dir, []string{"python"}, onlyRule)

		assert.NoError(t, err)
		assert.Contains(t, definitions, "ruby_test_rule")
	})
}
