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

func (lang *Language) Parse(ctx context.Context, content string) (*tree.Tree, error) {
	contentBytes := []byte(content)

	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(lang.implementation.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	builder := tree.NewBuilder(content, sitterTree.RootNode())

	if err := lang.implementation.AnalyzeTree(ctx, sitterTree.RootNode(), builder); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func (lang *Language) NewQuerySet() *query.Set {
	return query.NewSet(lang.implementation.SitterLanguage())
}

func (lang *Language) CompilePatternQuery(querySet *query.Set, input, focusedVariable string) (types.PatternQuery, error) {
	return patternquery.Compile(lang, lang.implementation, querySet, input, focusedVariable)
}
