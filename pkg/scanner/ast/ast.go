package ast

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/ruleset"

	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
)

var ExpectedComment = regexp.MustCompile(`\A[^\w]*bearer:expected\s[\w,]+\z`)

func Parse(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
) (*tree.Tree, error) {
	builder, err := parseBuilder(ctx, language, contentBytes, 0)
	if err != nil {
		return nil, err
	}

	return builder.Build(), nil
}

func ParseAndAnalyze(
	ctx context.Context,
	language language.Language,
	ruleSet *ruleset.Set,
	querySet *query.Set,
	contentBytes []byte,
) (*tree.Tree, error) {
	builder, err := parseBuilder(ctx, language, contentBytes, len(ruleSet.Rules()))
	if err != nil {
		return nil, err
	}

	if err := querySet.Query(ctx, builder, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running ast queries: %w", err)
	}

	analyzer := language.NewAnalyzer(builder)
	if err := analyzeNode(ctx, ruleSet, builder, analyzer, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func parseBuilder(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
	ruleCount int,
) (*tree.Builder, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	return tree.NewBuilder(language.SitterLanguage(), contentBytes, sitterTree.RootNode(), ruleCount), nil
}

func analyzeNode(
	ctx context.Context,
	ruleSet *ruleset.Set,
	builder *tree.Builder,
	analyzer language.Analyzer,
	node *sitter.Node,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		var disabledRules []*ruleset.Rule
		var expectedRules []*ruleset.Rule
		for i := 0; i < childCount; i++ {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			disabledRules = addDisabledRules(ruleSet, builder, disabledRules, child)
			expectedRules = addExpectedRules(ruleSet, builder, expectedRules, child)
			if err := analyzeNode(ctx, ruleSet, builder, analyzer, child); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}

func addExpectedRules(
	ruleSet *ruleset.Set,
	builder *tree.Builder,
	expectedRules []*ruleset.Rule,
	node *sitter.Node,
) []*ruleset.Rule {
	if strings.Contains(node.Type(), "comment") {
		nextExpectedRules := expectedRules

		nodeContent := builder.ContentFor(node)
		if ExpectedComment.Match([]byte(nodeContent)) {
			rawRuleIDs := strings.Split(nodeContent, "bearer:expected")[1]

			for _, ruleID := range strings.Split(rawRuleIDs, ",") {
				rule, err := ruleSet.RuleByID(strings.TrimSpace(ruleID))
				if err != nil {
					log.Debug().Msgf("ignoring unknown expected rule '%s': %s", ruleID, err)
					continue
				}

				nextExpectedRules = append(nextExpectedRules, rule)
			}
		}

		return nextExpectedRules
	}

	builder.AddExpectedRules(node, expectedRules)

	return nil
}

func addDisabledRules(
	ruleSet *ruleset.Set,
	builder *tree.Builder,
	disabledRules []*ruleset.Rule,
	node *sitter.Node,
) []*ruleset.Rule {
	if strings.Contains(node.Type(), "comment") {
		nextDisabledRules := disabledRules

		nodeContent := builder.ContentFor(node)
		if strings.Contains(nodeContent, "bearer:disable") {
			rawRuleIDs := strings.Split(nodeContent, "bearer:disable")[1]

			for _, ruleID := range strings.Split(rawRuleIDs, ",") {
				rule, err := ruleSet.RuleByID(strings.TrimSpace(ruleID))
				if err != nil {
					log.Debug().Msgf("ignoring unknown disabled rule '%s': %s", ruleID, err)
					continue
				}

				nextDisabledRules = append(nextDisabledRules, rule)
			}
		}

		return nextDisabledRules
	}

	builder.AddDisabledRules(node, disabledRules)

	return nil
}
