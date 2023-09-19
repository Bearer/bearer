package ast

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/ruleset"

	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

func Parse(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
) (*tree.Tree, error) {
	builder, err := parseBuilder(ctx, language, contentBytes)
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
	builder, err := parseBuilder(ctx, language, contentBytes)
	if err != nil {
		return nil, err
	}

	if err := querySet.Query(ctx, builder, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running ast queries: %w", err)
	}

	analyzer := language.NewAnalyzer(builder)
	var disabledRules []*ruleset.Rule
	if err := analyzeNode(ctx, ruleSet, builder, analyzer, builder.SitterRootNode(), disabledRules); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func parseBuilder(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
) (*tree.Builder, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	return tree.NewBuilder(contentBytes, sitterTree.RootNode()), nil
}

func analyzeNode(
	ctx context.Context,
	ruleSet *ruleset.Set,
	builder *tree.Builder,
	analyzer language.Analyzer,
	node *sitter.Node,
	disabledRules []*ruleset.Rule,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		for i := 0; i < childCount; i++ {
			child := node.Child(i)
			nextChild := node.Child(i + 1)
			if !child.IsNamed() {
				continue
			}

			if nextDisabledRules := addDisabledRules(ruleSet, builder, disabledRules, child); nextDisabledRules != nil {
				if err := analyzeNode(ctx, ruleSet, builder, analyzer, nextChild, nextDisabledRules); err != nil {
					return err
				}
			} else {
				builder.AddDisabledRules(child, disabledRules)

				if err := analyzeNode(ctx, ruleSet, builder, analyzer, child, disabledRules); err != nil {
					return err
				}
			}

		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}

func addDisabledRules(
	ruleSet *ruleset.Set,
	builder *tree.Builder,
	disabledRules []*ruleset.Rule,
	node *sitter.Node,
) []*ruleset.Rule {
	if node.Type() == "comment" {
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

	return nil
}
