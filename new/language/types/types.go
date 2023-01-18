package types

import "github.com/bearer/curio/new/language/tree"

type PatternQueryResult struct {
	MatchNode *tree.Node
	Variables tree.QueryResult
}

type PatternQuery interface {
	MatchAt(node *tree.Node) ([]*PatternQueryResult, error)
	MatchOnceAt(node *tree.Node) (*PatternQueryResult, error)
	Close()
}

type Language interface {
	Parse(input string) (*tree.Tree, error)
	CompileQuery(input string) (*tree.Query, error)
	CompilePatternQuery(input string) (PatternQuery, error)
}
