package base

import (
	"fmt"

	"github.com/bearer/curio/new/language/implementation"
	"github.com/bearer/curio/new/language/patternquery"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/new/language/types"
)

type Language struct {
	implementation implementation.Implementation
}

func New(implementation implementation.Implementation) *Language {
	return &Language{implementation: implementation}
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

func (lang *Language) CompilePatternQuery(input string) (types.PatternQuery, error) {
	inputWithoutVariables, params, err := lang.implementation.ExtractPatternVariables(input)
	if err != nil {
		return nil, fmt.Errorf("error processing variables: %s", err)
	}

	inputWithoutMatchNode, matchNodeOffset, err := lang.implementation.ExtractPatternMatchNode(inputWithoutVariables)
	if err != nil {
		return nil, fmt.Errorf("error processing match node: %s", err)
	}

	return patternquery.Compile(
		lang,
		lang.implementation.AnonymousPatternNodeParentTypes(),
		inputWithoutMatchNode,
		params,
		matchNodeOffset,
	)
}
