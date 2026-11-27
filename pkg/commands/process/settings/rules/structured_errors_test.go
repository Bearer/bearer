package rules

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/structurederrors"
)

func loadWithStructuredErrors(
	t *testing.T,
	files map[string]string,
) (map[string]settings.RuleDefinition, []structurederrors.Item, error) {
	t.Helper()

	dir := fstest.MapFS{}
	for name, content := range files {
		dir[name] = &fstest.MapFile{Data: []byte(content)}
	}

	errorCollector := structurederrors.New()
	definitions := make(map[string]settings.RuleDefinition)

	_, err := loadDefinitions(
		definitions,
		[]ruleDir{{path: "rules", fs: dir}},
		false,
		[]string{"python"},
		nil,
		errorCollector,
	)

	return definitions, errorCollector.Report().Items, err
}

func codes(items []structurederrors.Item) []structurederrors.Code {
	result := make([]structurederrors.Code, len(items))
	for i, item := range items {
		result[i] = item.Code
	}

	return result
}

func TestLoadDefinitionsStructuredErrors(t *testing.T) {
	t.Run("reports invalid YAML and still fails", func(t *testing.T) {
		_, items, err := loadWithStructuredErrors(t, map[string]string{
			"broken.yml": "languages: [python\npatterns:\n  - foo\n",
		})

		require.Error(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, structurederrors.CodeRuleYAMLInvalid, items[0].Code)
		assert.Equal(t, "rules/broken.yml", items[0].RuleFile)
	})

	t.Run("reports languages missing from the target", func(t *testing.T) {
		_, items, err := loadWithStructuredErrors(t, map[string]string{
			"lang.yml": "languages: [pyhton]\npatterns: [foo]\nmetadata:\n  id: lang\n",
		})

		require.NoError(t, err)
		assert.Equal(t, []structurederrors.Code{structurederrors.CodeRuleLanguageNotFound}, codes(items))
	})

	t.Run("reports missing metadata and id", func(t *testing.T) {
		_, items, err := loadWithStructuredErrors(t, map[string]string{
			"a.yml": "languages: [python]\npatterns: [foo]\n",
			"b.yml": "languages: [python]\npatterns: [foo]\nmetadata:\n  description: x\n",
		})

		require.NoError(t, err)
		assert.Equal(t, []structurederrors.Code{structurederrors.CodeRuleMissingMetadata, structurederrors.CodeRuleMissingID}, codes(items))
	})

	t.Run("reports duplicate rule IDs and still fails", func(t *testing.T) {
		_, items, err := loadWithStructuredErrors(t, map[string]string{
			"a.yml": "languages: [python]\npatterns: [foo]\nmetadata:\n  id: dup\n",
			"b.yml": "languages: [python]\npatterns: [bar]\nmetadata:\n  id: dup\n",
		})

		require.Error(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, structurederrors.CodeRuleDuplicateID, items[0].Code)
		assert.Equal(t, "rules/b.yml", items[0].RuleFile)
	})

	t.Run("reports validation issues with the related rule", func(t *testing.T) {
		definitions, items, err := loadWithStructuredErrors(t, map[string]string{
			"importer.yml": `
languages: [python]
imports: [missing_shared, not_shared]
patterns:
  - pattern: foo($<X>)
    filters:
      - variable: X
        detection: other
metadata:
  id: importer
`,
			"not_shared.yml": "languages: [python]\npatterns: [foo]\nmetadata:\n  id: not_shared\n",
		})

		require.NoError(t, err)
		assert.NotContains(t, definitions, "importer")
		assert.ElementsMatch(t, []structurederrors.Code{
			structurederrors.CodeImportUnknownRule,
			structurederrors.CodeImportNotShared,
			structurederrors.CodeFilterUnknownRule,
		}, codes(items))

		for _, item := range items {
			assert.Equal(t, "importer", item.RuleID)
			assert.Equal(t, "rules/importer.yml", item.RuleFile)
			assert.NotEmpty(t, item.RelatedRuleID)
		}
	})
}
