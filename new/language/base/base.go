package base

import (
	"context"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/new/language/types"
)

type Language struct {
	implementation implementation.Implementation
}

func New(implementation implementation.Implementation) *Language {
	return &Language{implementation: implementation}
}

func (lang *Language) Parse(ctx context.Context, input string) (*tree.Tree, error) {
	tree, err := tree.Parse(ctx, lang.implementation.SitterLanguage(), input)
	if err != nil {
		return nil, err
	}

	if err := lang.implementation.AnalyzeFlow(ctx, tree.RootNode()); err != nil {
		return nil, err
	}

	return tree, nil
}

func (lang *Language) NewQuerySet() *tree.QuerySet {
	return tree.NewQuerySet(lang.implementation.SitterLanguage())
}

func (lang *Language) CompilePatternQuery(querySet *tree.QuerySet, input, focusedVariable string) (types.PatternQuery, error) {
	return patternquery.Compile(lang, lang.implementation, querySet, input, focusedVariable)
}
