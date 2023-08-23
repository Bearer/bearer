package ast

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/language"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

func Parse(ctx context.Context, language language.Language, contentBytes []byte) (*tree.Tree, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	builder := tree.NewBuilder(contentBytes, sitterTree.RootNode())
	analyzer := language.NewAnalyzer(builder)

	if err := analyzeNode(ctx, analyzer, sitterTree.RootNode()); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func analyzeNode(ctx context.Context, analyzer language.Analyzer, node *sitter.Node) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		for i := 0; i < childCount; i++ {
			if err := analyzeNode(ctx, analyzer, node.Child(i)); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}
