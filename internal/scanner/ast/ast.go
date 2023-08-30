package ast

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/language"

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
	if err := analyzeNode(ctx, builder, analyzer, builder.SitterRootNode()); err != nil {
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
	builder *tree.Builder,
	analyzer language.Analyzer,
	node *sitter.Node,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		var disabledRules []string
		for i := 0; i < childCount; i++ {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			disabledRules = addDisabledRules(builder, disabledRules, node)
			if err := analyzeNode(ctx, builder, analyzer, child); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}

func addDisabledRules(builder *tree.Builder, disabledRules []string, node *sitter.Node) []string {
	if node.Type() == "comment" {
		nextDisabledRules := disabledRules

		nodeContent := builder.ContentFor(node)
		if strings.Contains(nodeContent, "bearer:disable") {
			ruleIdsStr := strings.Split(nodeContent, "bearer:disable")[1]

			for _, ruleId := range strings.Split(ruleIdsStr, ",") {
				nextDisabledRules = append(nextDisabledRules, strings.TrimSpace(ruleId))
			}
		}

		return nextDisabledRules
	}

	builder.AddDisabledRules(node, disabledRules)

	return nil
}
