package language

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery"
	"github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
)

type Language struct {
	implementation implementation.Implementation
}

func New(implementation implementation.Implementation) *Language {
	return &Language{implementation: implementation}
}

func (lang *Language) Parse(ctx context.Context, contentBytes []byte) (*tree.Tree, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(lang.implementation.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	builder := tree.NewBuilder(contentBytes, sitterTree.RootNode())
	analyzer := lang.implementation.NewAnalyzer(builder)

	if err := lang.analyzeNode(ctx, analyzer, sitterTree.RootNode()); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func (lang *Language) analyzeNode(ctx context.Context, analyzer implementation.Analyzer, node *sitter.Node) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		for i := 0; i < childCount; i++ {
			if err := lang.analyzeNode(ctx, analyzer, node.Child(i)); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}

func (lang *Language) CompilePatternQuery(querySet *query.Set, input, focusedVariable string) (types.PatternQuery, error) {
	return patternquery.Compile(lang, lang.implementation.Pattern(), querySet, input, focusedVariable)
}
