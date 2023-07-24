package types

import (
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type Variable struct {
	NodeTypes  []string
	DummyValue string
	Name       string
}

type PatternQuery interface {
	MatchAt(node *tree.Node) ([]*languagetypes.PatternQueryResult, error)
	MatchOnceAt(node *tree.Node) (*languagetypes.PatternQueryResult, error)
}
