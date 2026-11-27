package artifact

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/structurederrors"
)

func TestEngineErrorItem(t *testing.T) {
	rules := map[string]*settings.Rule{
		"my_rule": {Id: "my_rule", Patterns: []settings.RulePattern{{Pattern: "a()"}, {Pattern: "b("}}},
	}

	t.Run("pattern error", func(t *testing.T) {
		err := fmt.Errorf("error adding my_rule: %w", &ruleset.PatternError{
			RuleID:       "my_rule",
			PatternIndex: 1,
			Err:          errors.New("error parsing pattern at 1:1: b("),
		})

		item := engineErrorItem("javascript", rules, err)

		assert.Equal(t, structurederrors.CodePatternInvalid, item.Code)
		assert.Equal(t, "my_rule", item.RuleID)
		assert.Equal(t, structurederrors.Index(1), item.PatternIndex)
		assert.Equal(t, "b(", item.Pattern)
		assert.Equal(t, "error parsing pattern at 1:1: b(", item.Message)
		assert.Equal(t, "javascript", item.Language)
	})

	t.Run("filter error", func(t *testing.T) {
		err := &ruleset.PatternError{RuleID: "my_rule", IsFilter: true, Err: errors.New("unknown variable 'Y'")}

		item := engineErrorItem("javascript", rules, err)

		assert.Equal(t, structurederrors.CodeFilterInvalid, item.Code)
		assert.Equal(t, "a()", item.Pattern)
	})

	t.Run("other error", func(t *testing.T) {
		item := engineErrorItem("javascript", rules, errors.New("duplicate rule 'string'"))

		assert.Equal(t, structurederrors.CodeRuleSetInvalid, item.Code)
		assert.Empty(t, item.RuleID)
		assert.Equal(t, "javascript", item.Language)
	})
}
