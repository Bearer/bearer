package filters_test

import (
	"context"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
	filename string
	scan     func(
		rootNode *tree.Node,
		rule *ruleset.Rule,
		traversalStrategy traversalstrategy.Strategy,
	) ([]*detectortypes.Detection, error)
}

func (context *MockDetectorContext) Filename() string {
	return context.filename
}

func (context *MockDetectorContext) Scan(
	rootNode *tree.Node,
	rule *ruleset.Rule,
	traversalStrategy traversalstrategy.Strategy,
) ([]*detectortypes.Detection, error) {
	if context.scan != nil {
		return context.scan(rootNode, rule, traversalStrategy)
	}

	Fail("MockDetectorContext.scan called but no scan function was set")
	panic("unreachable")
}

func (filter *MockFilter) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*filters.Result, error) {
	return filter.result, filter.err
}

var defaultDetectorContext = &MockDetectorContext{
	filename: "src/foo.go",
}

var _ = Describe("Not", func() {
	var filter *filters.Not
	var patternVariables = []*tree.Node{{ID: 42}}

	When("the child filter has a match", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Not{
				Child: &MockFilter{result: filters.NewResult(filters.NewMatch(nil, "", nil))},
			}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the child filter has NO matches", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Not{
				Child: &MockFilter{result: filters.NewResult()},
			}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the child filter result is unknown", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Not{
				Child: &MockFilter{result: nil},
			}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("Either", func() {
	var filter *filters.Either
	patternVariables := []*tree.Node{{ID: 42}}

	When("there are child filter matches", func() {
		match1 := filters.NewMatch([]*tree.Node{{ID: 1}}, "", nil)
		match2 := filters.NewMatch([]*tree.Node{{ID: 2}}, "", nil)
		match3 := filters.NewMatch([]*tree.Node{{ID: 3}}, "", nil)

		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Either{
				Children: []filters.Filter{
					&MockFilter{result: filters.NewResult(match1, match2)},
					&MockFilter{result: filters.NewResult(match3)},
					&MockFilter{result: nil},
					&MockFilter{result: filters.NewResult()},
				},
			}
		})

		It("returns a result with all matches combined", func(ctx SpecContext) {
			result, err := filter.Evaluate(defaultDetectorContext, patternVariables)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Matches()).To(ContainElements(match1, match2, match3))
		})
	})

	When("no child filter matches", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Either{
				Children: []filters.Filter{
					&MockFilter{result: nil},
					&MockFilter{result: filters.NewResult()},
				},
			}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("all child filter results are unknown", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Either{
				Children: []filters.Filter{
					&MockFilter{result: nil},
					&MockFilter{result: nil},
				},
			}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})

	When("there are NO child filters", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Either{Children: nil}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("All", func() {
	var filter *filters.All
	var nodes []*tree.Node
	var patternVariables variableshape.Values

	datatype1 := &detectortypes.Detection{RuleID: "dt1"}
	datatype2 := &detectortypes.Detection{RuleID: "dt2"}
	datatype3 := &detectortypes.Detection{RuleID: "dt3"}
	datatype4 := &detectortypes.Detection{RuleID: "dt4"}
	datatype5 := &detectortypes.Detection{RuleID: "dt5"}
	discordantDatatype := &detectortypes.Detection{RuleID: "dtd"}

	var match1, match2, match3, match4, match5, discordantMatch filters.Match

	BeforeEach(func(ctx SpecContext) {
		nodes = parseNodes(ctx, []string{"n1", "n2", "n3", "n4", "n5", "n6", "n7", "n8"})
		patternVariables = []*tree.Node{nodes[0], nil, nil, nil}

		match1 = filters.NewMatch([]*tree.Node{nodes[0], nil, nil, nil}, "", []*detectortypes.Detection{datatype1})
		match2 = filters.NewMatch([]*tree.Node{nil, nodes[2], nodes[4], nil}, "", []*detectortypes.Detection{datatype2})
		match3 = filters.NewMatch([]*tree.Node{nil, nodes[3], nodes[5], nil}, "", []*detectortypes.Detection{datatype3})
		match4 = filters.NewMatch([]*tree.Node{nodes[0], nodes[3], nil, nodes[6]}, "", []*detectortypes.Detection{datatype4})
		match5 = filters.NewMatch([]*tree.Node{nodes[0], nodes[3], nil, nodes[7]}, "", []*detectortypes.Detection{datatype5})
		discordantMatch = filters.NewMatch([]*tree.Node{nodes[1], nil, nil, nil}, "", []*detectortypes.Detection{discordantDatatype})
	})

	When("there is a single child filter with matches", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.All{
				Children: []filters.Filter{
					&MockFilter{result: filters.NewResult(match1, match2)},
				},
			}
		})

		It("returns a result with the child matches", func(ctx SpecContext) {
			result, err := filter.Evaluate(defaultDetectorContext, patternVariables)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Matches()).To(ContainElements(match1, match2))
		})
	})

	When("there are multiple child filters that all match", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.All{
				Children: []filters.Filter{
					&MockFilter{result: filters.NewResult(match1)},
					&MockFilter{result: filters.NewResult(match2, match3)},
					&MockFilter{result: filters.NewResult(match4, match5)},
				},
			}
		})

		It("returns a result with the matches joined by variables", func(ctx SpecContext) {
			result, err := filter.Evaluate(defaultDetectorContext, patternVariables)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.Matches()).To(ContainElements(
				filters.NewMatch(
					[]*tree.Node{nodes[0], nodes[3], nodes[5], nodes[6]},
					"",
					[]*detectortypes.Detection{datatype1, datatype3, datatype4},
				),
				filters.NewMatch(
					[]*tree.Node{nodes[0], nodes[3], nodes[5], nodes[7]},
					"",
					[]*detectortypes.Detection{datatype1, datatype3, datatype5},
				),
			))
		})
	})

	When("one child has NO matches that join to the other matches", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.All{
				Children: []filters.Filter{
					// this is the same example as above, but with the addition of
					// the discordant match
					&MockFilter{result: filters.NewResult(match1)},
					&MockFilter{result: filters.NewResult(match2, match3)},
					&MockFilter{result: filters.NewResult(match4, match5)},
					&MockFilter{result: filters.NewResult(discordantMatch)},
				},
			}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("at least one child filter result is unknown", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.All{
				Children: []filters.Filter{
					&MockFilter{result: filters.NewResult(match1)},
					&MockFilter{result: nil},
				},
			}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})

	When("there are NO child filters", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.All{Children: nil}
		})

		It("returns a result with a single match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})
})

var _ = Describe("FilenameRegex", func() {
	var filter *filters.FilenameRegex
	patternVariables := []*tree.Node{{ID: 42}}

	When("the filename matches the regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.FilenameRegex{Regex: regexp.MustCompile(`foo`)}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the filename does NOT match the regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.FilenameRegex{Regex: regexp.MustCompile(`bar`)}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})
})

var _ = Describe("Rule", func() {
})

var _ = Describe("Values", func() {
	var filter *filters.Values
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "hello", "other")
	})

	When("the variable node's content is equal to one of the filter's values", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Values{Variable: variable, Values: []string{"hello", "other"}}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content does NOT match any filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Values{Variable: variable, Values: []string{"other"}}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})
})

var _ = Describe("Regex", func() {
	var filter *filters.Regex
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "hello", "other")
	})

	When("the variable node's content matches the regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Regex{Variable: variable, Regex: regexp.MustCompile(`l{2}`)}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content does NOT match the regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.Regex{Variable: variable, Regex: regexp.MustCompile(`other`)}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})
})

var _ = Describe("StringLengthLessThan", func() {
	var filter *filters.StringLengthLessThan
	var variable *variableshape.Variable
	var detectorContext detectortypes.Context
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "hello", "other")
		detectorContext = setupStringTest(patternVariables.Node(variable), pointers.String("foo"))
	})

	When("the variable node's string detector value has length less than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringLengthLessThan{Variable: variable, Value: 4}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's string detector value has length greater than or equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringLengthLessThan{Variable: variable, Value: 3}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node has NO string detector value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringLengthLessThan{Variable: variable, Value: 4}
			detectorContext = setupStringTest(patternVariables.Node(variable), nil)
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("StringRegex", func() {
	var filter *filters.StringRegex
	var variable *variableshape.Variable
	var detectorContext detectortypes.Context
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "hello", "other")
		detectorContext = setupStringTest(patternVariables.Node(variable), pointers.String("foo"))
	})

	When("the variable node's string detector value matches the filter regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`o{2}`)}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "foo", nil)),
			))
		})
	})

	When("the variable node's string detector value does NOT match the filter regex", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`bar`)}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node has NO string detector value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.StringRegex{Variable: variable, Regex: regexp.MustCompile(`o{2}`)}
			detectorContext = setupStringTest(patternVariables.Node(variable), nil)
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("EntropyGreaterThan", func() {
	var filter *filters.EntropyGreaterThan
	var variable *variableshape.Variable
	var detectorContext detectortypes.Context
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "hello", "other")
		// entropy("Au+u1hvsvJeEXxky") == 3.75
		detectorContext = setupStringTest(patternVariables.Node(variable), pointers.String("Au+u1hvsvJeEXxky"))
	})

	When("the variable node's content is a string with entropy greater than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.EntropyGreaterThan{Variable: variable, Value: 3.7}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is a string with entropy less than or equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.EntropyGreaterThan{Variable: variable, Value: 3.8}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node is not a string value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.EntropyGreaterThan{Variable: variable, Value: 3.7}
			detectorContext = setupStringTest(patternVariables.Node(variable), nil)
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(detectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("IntegerLessThan", func() {
	var filter *filters.IntegerLessThan
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "9", "other")
	})

	When("the variable node's content is an integer less than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerLessThan{Variable: variable, Value: 10}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is an integer greater than or equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerLessThan{Variable: variable, Value: 9}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node is not an integer value", func() {
		BeforeEach(func(ctx SpecContext) {
			variable, patternVariables = setupContentTest(ctx, "hello", "other")
			filter = &filters.IntegerLessThan{Variable: variable, Value: 10}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("IntegerLessThanOrEqual", func() {
	var filter *filters.IntegerLessThanOrEqual
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "9", "other")
	})

	When("the variable node's content is an integer less than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerLessThanOrEqual{Variable: variable, Value: 10}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})

	})
	When("the variable node's content is an integer equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerLessThanOrEqual{Variable: variable, Value: 9}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is an integer greater than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerLessThanOrEqual{Variable: variable, Value: 8}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node is not an integer value", func() {
		BeforeEach(func(ctx SpecContext) {
			variable, patternVariables = setupContentTest(ctx, "hello", "other")
			filter = &filters.IntegerLessThanOrEqual{Variable: variable, Value: 10}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("IntegerGreaterThan", func() {
	var filter *filters.IntegerGreaterThan
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "9", "other")
	})

	When("the variable node's content is an integer greater than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerGreaterThan{Variable: variable, Value: 8}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is an integer less than or equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerGreaterThan{Variable: variable, Value: 9}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node is not an integer value", func() {
		BeforeEach(func(ctx SpecContext) {
			variable, patternVariables = setupContentTest(ctx, "hello", "other")
			filter = &filters.IntegerGreaterThan{Variable: variable, Value: 8}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("IntegerGreaterThanOrEqual", func() {
	var filter *filters.IntegerGreaterThanOrEqual
	var variable *variableshape.Variable
	var patternVariables variableshape.Values

	BeforeEach(func(ctx SpecContext) {
		variable, patternVariables = setupContentTest(ctx, "9", "other")
	})

	When("the variable node's content is an integer greater than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 8}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is an integer equal to the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 9}
		})

		It("returns a result with a match using the pattern variables", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(
				filters.NewResult(filters.NewMatch(patternVariables, "", nil)),
			))
		})
	})

	When("the variable node's content is an integer less than the filter value", func() {
		BeforeEach(func(ctx SpecContext) {
			filter = &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 10}
		})

		It("returns a result with NO matches", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(Equal(filters.NewResult()))
		})
	})

	When("the variable node is not an integer value", func() {
		BeforeEach(func(ctx SpecContext) {
			variable, patternVariables = setupContentTest(ctx, "hello", "other")
			filter = &filters.IntegerGreaterThanOrEqual{Variable: variable, Value: 8}
		})

		It("returns an unknown result", func(ctx SpecContext) {
			Expect(filter.Evaluate(defaultDetectorContext, patternVariables)).To(BeNil())
		})
	})
})

var _ = Describe("Unknown", func() {
	filter := &filters.Unknown{}

	It("returns an unknown result", func(ctx SpecContext) {
		Expect(filter.Evaluate(defaultDetectorContext, nil)).To(BeNil())
	})
})

func parseNodes(ctx context.Context, content []string) []*tree.Node {
	tree, err := ast.Parse(ctx, ruby.Get(), []byte(strings.Join(content, "\n")))
	Expect(err).To(BeNil())
	return tree.RootNode().NamedChildren()
}

func setupContentTest(ctx context.Context, content, otherContent string) (*variableshape.Variable, variableshape.Values) {
	variableShape := variableshape.NewBuilder().Add("one").Add("two").Build()

	otherVariable, err := variableShape.Variable("one")
	Expect(err).To(BeNil())
	variable, err := variableShape.Variable("two")
	Expect(err).To(BeNil())

	nodes := parseNodes(ctx, []string{otherContent, content})

	patternVariables := variableShape.NewValues()
	patternVariables.Set(otherVariable, nodes[0])
	patternVariables.Set(variable, nodes[1])

	return variable, patternVariables
}

func setupStringTest(node *tree.Node, value *string) detectortypes.Context {
	return &MockDetectorContext{
		filename: defaultDetectorContext.filename,
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

			Fail("unexpected call to MockDetectorContext.scan")
			panic("unreachable")
		},
	}
}
