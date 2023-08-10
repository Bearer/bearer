package types

import (
	"context"

	"github.com/bearer/bearer/new/language/tree"
	sitter "github.com/smacker/go-tree-sitter"
)

type PatternQueryResult struct {
	MatchNode *sitter.Node
	Variables tree.QueryResult
}

type PatternQuery interface {
	MatchAt(context *tree.QueryContext, node *sitter.Node) ([]*PatternQueryResult, error)
	MatchOnceAt(context *tree.QueryContext, node *sitter.Node) (*PatternQueryResult, error)
}

type Language interface {
	Parse(ctx context.Context, input string) (*tree.Tree, error)
	NewQuerySet() *tree.QuerySet
	CompilePatternQuery(querySet *tree.QuerySet, input, focusedVariable string) (PatternQuery, error)
}
