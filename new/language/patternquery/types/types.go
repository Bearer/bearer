package types

import (
	sitter "github.com/smacker/go-tree-sitter"

	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/query"
)

type Variable struct {
	NodeTypes  []string
	DummyValue string
	Name       string
}

type PatternQuery interface {
	MatchAt(astContext *query.Context, node *sitter.Node) ([]*languagetypes.PatternQueryResult, error)
	MatchOnceAt(astContext *query.Context, node *sitter.Node) (*languagetypes.PatternQueryResult, error)
}
