package filters_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/languages/ruby"
	"github.com/bearer/bearer/pkg/scanner/ast"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/filters"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/scanner/variableshape"
	"github.com/bearer/bearer/pkg/util/pointers"
)

type MockFilter struct {
	result *filters.Result
	err    error
}

type MockDetectorContext struct {
	t        *testing.T
	filename string
	scan     func(
		rootNode *tree.Node,
		rule *ruleset.Rule,
		traversalStrategy traversalstrategy.Strategy,
	) ([]*detectortypes.Detection, error)
}

func (ctx *MockDetectorContext) Filename() string {
	return ctx.filename
}

func (ctx *MockDetectorContext) Scan(
	rootNode *tree.Node,
	rule *ruleset.Rule,
	traversalStrategy traversalstrategy.Strategy,
) ([]*detectortypes.Detection, error) {
	if ctx.scan != nil {
		return ctx.scan(rootNode, rule, traversalStrategy)
	}

	ctx.t.Fatal("MockDetectorContext.scan called but no scan function was set")
	panic("unreachable")
}

func (filter *MockFilter) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*filters.Result, error) {
	return filter.result, filter.err
}

func newDefaultDetectorContext(t *testing.T) *MockDetectorContext {
	return &MockDetectorContext{t: t, filename: "src/foo.go"}
}

func TestNot(t *testing.T) {
	patternVariables := []*tree.Node{{ID: 42}}
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("child filter has a match returns NO matches", func(t *testing.T) {
		filter := &filters.Not{
			Child: &MockFilter{result: filters.NewResult(filters.NewMatch(nil, "", nil))},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("child filter has NO matches returns match with pattern variables", func(t *testing.T) {
		filter := &filters.Not{
			Child: &MockFilter{result: filters.NewResult()},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("child filter result is unknown returns unknown", func(t *testing.T) {
		filter := &filters.Not{
			Child: &MockFilter{result: nil},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestEither(t *testing.T) {
	patternVariables := []*tree.Node{{ID: 42}}
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("child filter matches are combined", func(t *testing.T) {
		match1 := filters.NewMatch([]*tree.Node{{ID: 1}}, "", nil)
		match2 := filters.NewMatch([]*tree.Node{{ID: 2}}, "", nil)
		match3 := filters.NewMatch([]*tree.Node{{ID: 3}}, "", nil)

		filter := &filters.Either{
			Children: []filters.Filter{
				&MockFilter{result: filters.NewResult(match1, match2)},
				&MockFilter{result: filters.NewResult(match3)},
				&MockFilter{result: nil},
				&MockFilter{result: filters.NewResult()},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, result.Matches(), match1)
		assert.Contains(t, result.Matches(), match2)
		assert.Contains(t, result.Matches(), match3)
	})

	t.Run("no child filter matches returns NO matches", func(t *testing.T) {
		filter := &filters.Either{
			Children: []filters.Filter{
				&MockFilter{result: nil},
				&MockFilter{result: filters.NewResult()},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("all child filter results unknown returns unknown", func(t *testing.T) {
		filter := &filters.Either{
			Children: []filters.Filter{
				&MockFilter{result: nil},
				&MockFilter{result: nil},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("no child filters returns unknown", func(t *testing.T) {
		filter := &filters.Either{Children: nil}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestAll(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	datatype1 := &detectortypes.Detection{RuleID: "dt1"}
	datatype2 := &detectortypes.Detection{RuleID: "dt2"}
	datatype3 := &detectortypes.Detection{RuleID: "dt3"}
	datatype4 := &detectortypes.Detection{RuleID: "dt4"}
	datatype5 := &detectortypes.Detection{RuleID: "dt5"}
	discordantDatatype := &detectortypes.Detection{RuleID: "dtd"}

	setup := func(t *testing.T) ([]*tree.Node, variableshape.Values, filters.Match, filters.Match, filters.Match, filters.Match, filters.Match, filters.Match) {
		t.Helper()

		nodes := parseNodes(t, []string{"n1", "n2", "n3", "n4", "n5", "n6", "n7", "n8"})
		patternVariables := []*tree.Node{nodes[0], nil, nil, nil}

		match1 := filters.NewMatch([]*tree.Node{nodes[0], nil, nil, nil}, "", []*detectortypes.Detection{datatype1})
		match2 := filters.NewMatch([]*tree.Node{nil, nodes[2], nodes[4], nil}, "", []*detectortypes.Detection{datatype2})
		match3 := filters.NewMatch([]*tree.Node{nil, nodes[3], nodes[5], nil}, "", []*detectortypes.Detection{datatype3})
		match4 := filters.NewMatch([]*tree.Node{nodes[0], nodes[3], nil, nodes[6]}, "", []*detectortypes.Detection{datatype4})
		match5 := filters.NewMatch([]*tree.Node{nodes[0], nodes[3], nil, nodes[7]}, "", []*detectortypes.Detection{datatype5})
		discordantMatch := filters.NewMatch([]*tree.Node{nodes[1], nil, nil, nil}, "", []*detectortypes.Detection{discordantDatatype})

		return nodes, patternVariables, match1, match2, match3, match4, match5, discordantMatch
	}

	t.Run("single child filter with matches returns child matches", func(t *testing.T) {
		_, patternVariables, match1, match2, _, _, _, _ := setup(t)
		filter := &filters.All{
			Children: []filters.Filter{
				&MockFilter{result: filters.NewResult(match1, match2)},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, result.Matches(), match1)
		assert.Contains(t, result.Matches(), match2)
	})

	t.Run("multiple child filters all match returns joined matches", func(t *testing.T) {
		nodes, patternVariables, match1, _, match3, match4, match5, _ := setup(t)
		filter := &filters.All{
			Children: []filters.Filter{
				&MockFilter{result: filters.NewResult(match1)},
				&MockFilter{result: filters.NewResult(
					filters.NewMatch([]*tree.Node{nil, nodes[2], nodes[4], nil}, "", []*detectortypes.Detection{datatype2}),
					match3,
				)},
				&MockFilter{result: filters.NewResult(match4, match5)},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, result.Matches(), filters.NewMatch(
			[]*tree.Node{nodes[0], nodes[3], nodes[5], nodes[6]},
			"",
			[]*detectortypes.Detection{datatype1, datatype3, datatype4},
		))
		assert.Contains(t, result.Matches(), filters.NewMatch(
			[]*tree.Node{nodes[0], nodes[3], nodes[5], nodes[7]},
			"",
			[]*detectortypes.Detection{datatype1, datatype3, datatype5},
		))
	})

	t.Run("one child has NO joining matches returns NO matches", func(t *testing.T) {
		_, patternVariables, match1, match2, match3, match4, match5, discordantMatch := setup(t)
		filter := &filters.All{
			Children: []filters.Filter{
				&MockFilter{result: filters.NewResult(match1)},
				&MockFilter{result: filters.NewResult(match2, match3)},
				&MockFilter{result: filters.NewResult(match4, match5)},
				&MockFilter{result: filters.NewResult(discordantMatch)},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("at least one child unknown returns unknown", func(t *testing.T) {
		_, patternVariables, match1, _, _, _, _, _ := setup(t)
		filter := &filters.All{
			Children: []filters.Filter{
				&MockFilter{result: filters.NewResult(match1)},
				&MockFilter{result: nil},
			},
		}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("no child filters returns match with pattern variables", func(t *testing.T) {
		_, patternVariables, _, _, _, _, _, _ := setup(t)
		filter := &filters.All{Children: nil}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})
}

func TestFilenameRegex(t *testing.T) {
	patternVariables := []*tree.Node{{ID: 42}}
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("filename matches regex returns match", func(t *testing.T) {
		filter := &filters.FilenameRegex{Regex: regexp.MustCompile(`foo`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("filename does NOT match regex returns NO matches", func(t *testing.T) {
		filter := &filters.FilenameRegex{Regex: regexp.MustCompile(`bar`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})
}

func TestValues(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	variable, patternVariables := setupContentTest(t, "hello", "other")

	t.Run("content equals one of the values returns match", func(t *testing.T) {
		filter := &filters.Values{Variable: variable, Values: []string{"hello", "other"}}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("content does NOT match any value returns NO matches", func(t *testing.T) {
		filter := &filters.Values{Variable: variable, Values: []string{"other"}}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})
}

func TestRegex(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	variable, patternVariables := setupContentTest(t, "hello", "other")

	t.Run("content matches regex returns match", func(t *testing.T) {
		filter := &filters.Regex{Variable: variable, Regex: regexp.MustCompile(`l{2}`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("content does NOT match regex returns NO matches", func(t *testing.T) {
		filter := &filters.Regex{Variable: variable, Regex: regexp.MustCompile(`other`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})
}

func TestStringLengthLessThan(t *testing.T) {
	variable, patternVariables := setupContentTest(t, "hello", "other")

	t.Run("string length less than filter value returns match", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("foo"))
		filter := &filters.StringLengthLessThan{Variable: variable, Value: 4}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("string length greater or equal returns NO matches", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("foo"))
		filter := &filters.StringLengthLessThan{Variable: variable, Value: 3}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("no string detector value returns unknown", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), nil)
		filter := &filters.StringLengthLessThan{Variable: variable, Value: 4}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestStringRegex(t *testing.T) {
	variable, patternVariables := setupContentTest(t, "hello", "other")

	t.Run("string value matches regex returns match", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("foo"))
		filter := &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`o{2}`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "foo", nil)), result)
	})

	t.Run("string value does NOT match regex returns NO matches", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("foo"))
		filter := &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`bar`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("no string detector value returns unknown", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), nil)
		filter := &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`o{2}`)}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestEntropyGreaterThan(t *testing.T) {
	variable, patternVariables := setupContentTest(t, "hello", "other")

	t.Run("entropy greater than filter value returns match", func(t *testing.T) {
		// entropy("Au+u1hvsvJeEXxky") == 3.75
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("Au+u1hvsvJeEXxky"))
		filter := &filters.EntropyGreaterThan{Variable: variable, Value: 3.7}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("entropy less than or equal returns NO matches", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), pointers.String("Au+u1hvsvJeEXxky"))
		filter := &filters.EntropyGreaterThan{Variable: variable, Value: 3.8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("not a string value returns unknown", func(t *testing.T) {
		detectorCtx := setupStringTest(t, patternVariables.Node(variable), nil)
		filter := &filters.EntropyGreaterThan{Variable: variable, Value: 3.7}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestIntegerLessThan(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("integer less than filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerLessThan{Variable: variable, Value: 10}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer greater or equal returns NO matches", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerLessThan{Variable: variable, Value: 9}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("not an integer value returns unknown", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "hello", "other")
		filter := &filters.IntegerLessThan{Variable: variable, Value: 10}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestIntegerLessThanOrEqual(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("integer less than filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerLessThanOrEqual{Variable: variable, Value: 10}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer equal to filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerLessThanOrEqual{Variable: variable, Value: 9}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer greater than filter value returns NO matches", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerLessThanOrEqual{Variable: variable, Value: 8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("not an integer value returns unknown", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "hello", "other")
		filter := &filters.IntegerLessThanOrEqual{Variable: variable, Value: 10}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestIntegerGreaterThan(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("integer greater than filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerGreaterThan{Variable: variable, Value: 8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer less than or equal returns NO matches", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerGreaterThan{Variable: variable, Value: 9}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("not an integer value returns unknown", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "hello", "other")
		filter := &filters.IntegerGreaterThan{Variable: variable, Value: 8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestIntegerGreaterThanOrEqual(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)

	t.Run("integer greater than filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer equal to filter value returns match", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 9}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(filters.NewMatch(patternVariables, "", nil)), result)
	})

	t.Run("integer less than filter value returns NO matches", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "9", "other")
		filter := &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 10}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Equal(t, filters.NewResult(), result)
	})

	t.Run("not an integer value returns unknown", func(t *testing.T) {
		variable, patternVariables := setupContentTest(t, "hello", "other")
		filter := &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 8}
		result, err := filter.Evaluate(detectorCtx, patternVariables)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestUnknown(t *testing.T) {
	detectorCtx := newDefaultDetectorContext(t)
	filter := &filters.Unknown{}
	result, err := filter.Evaluate(detectorCtx, nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func parseNodes(t *testing.T, content []string) []*tree.Node {
	t.Helper()

	tree, err := ast.Parse(context.Background(), ruby.Get(), []byte(strings.Join(content, "\n")))
	require.NoError(t, err)
	return tree.RootNode().NamedChildren()
}

func setupContentTest(t *testing.T, content, otherContent string) (*variableshape.Variable, variableshape.Values) {
	t.Helper()

	variableShape := variableshape.NewBuilder().Add("one").Add("two").Build()

	otherVariable, err := variableShape.Variable("one")
	require.NoError(t, err)
	variable, err := variableShape.Variable("two")
	require.NoError(t, err)

	nodes := parseNodes(t, []string{otherContent, content})

	patternVariables := variableShape.NewValues()
	patternVariables.Set(otherVariable, nodes[0])
	patternVariables.Set(variable, nodes[1])

	return variable, patternVariables
}

func setupStringTest(t *testing.T, node *tree.Node, value *string) detectortypes.Context {
	t.Helper()

	return &MockDetectorContext{
		t:        t,
		filename: "src/foo.go",
		scan: func(
			rootNode *tree.Node,
			rule *ruleset.Rule,
			traversalStrategy traversalstrategy.Strategy,
		) ([]*detectortypes.Detection, error) {
			if rootNode == node &&
				rule == ruleset.BuiltinStringRule &&
				traversalStrategy == traversalstrategy.Cursor {

				if value == nil {
					return nil, nil
				}

				return []*detectortypes.Detection{{
					RuleID:    rule.ID(),
					MatchNode: rootNode,
					Data:      common.String{Value: *value, IsLiteral: true},
				}}, nil
			}

			t.Fatal("unexpected call to MockDetectorContext.scan")
			panic("unreachable")
		},
	}
}
