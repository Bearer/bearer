package types

import (
	"context"

	"github.com/bearer/bearer/new/language/tree"
)

type PatternQueryResult struct {
	MatchNode *tree.Node
	Variables tree.QueryResult
}

type PatternQuery interface {
	MatchAt(node *tree.Node) ([]*PatternQueryResult, error)
	MatchOnceAt(node *tree.Node) (*PatternQueryResult, error)
}

type Language interface {
	Parse(ctx context.Context, input string) (*tree.Tree, error)
	NewQuerySet() *tree.QuerySet
	CompilePatternQuery(querySet *tree.QuerySet, input, focusedVariable string) (PatternQuery, error)
}
