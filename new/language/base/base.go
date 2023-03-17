package base

import (
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/new/language/types"
	soufflequery "github.com/bearer/bearer/pkg/souffle/query"
)

type Language struct {
	implementation implementation.Implementation
	souffle        bool
}

func New(souffle bool, implementation implementation.Implementation) *Language {
	return &Language{souffle: souffle, implementation: implementation}
}

func (lang *Language) Implementation() implementation.Implementation {
	return lang.implementation
}

func (lang *Language) Parse(input string) (*tree.Tree, error) {
	tree, err := tree.Parse(lang.implementation.SitterLanguage(), input)
	if err != nil {
		return nil, err
	}

	if err := lang.implementation.AnalyzeFlow(tree.RootNode()); err != nil {
		return nil, err
	}

	return tree, nil
}

func (lang *Language) CompileQuery(input string) (*tree.Query, error) {
	return tree.CompileQuery(lang.implementation.SitterLanguage(), input)
}

func (lang *Language) CompilePatternQuery(ruleName, input string) (types.PatternQuery, error) {
	if lang.souffle {
		return soufflequery.NewQuery(ruleName), nil
	}

	return patternquery.Compile(lang, lang.implementation, input)
}
