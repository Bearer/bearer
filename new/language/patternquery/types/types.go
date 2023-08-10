package types

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type Variable struct {
	NodeTypes  []string
	DummyValue string
	Name       string
}

type PatternQuery interface {
	MatchAt(astContext *tree.QueryContext, node *sitter.Node) ([]*languagetypes.PatternQueryResult, error)
	MatchOnceAt(astContext *tree.QueryContext, node *sitter.Node) (*languagetypes.PatternQueryResult, error)
}
