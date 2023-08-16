package types

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
)

type PatternQueryResult struct {
	MatchNode *sitter.Node
	Variables query.Result
}

type PatternQuery interface {
	MatchAt(context *query.Context, node *sitter.Node) ([]*PatternQueryResult, error)
	MatchOnceAt(context *query.Context, node *sitter.Node) (*PatternQueryResult, error)
}

type Language interface {
	Parse(ctx context.Context, content string) (*tree.Tree, error)
	NewQuerySet() *query.Set
	CompilePatternQuery(querySet *query.Set, input, focusedVariable string) (PatternQuery, error)
}
