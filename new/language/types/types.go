package types

import "github.com/bearer/curio/new/language"

type Language interface {
	Parse(input string) (*language.Tree, error)
	CompileQuery(input string) (*language.Query, error)
	// FIXME: should be part of Tree
	UnifiedNodesFor(node *language.Node) []*language.Node
}
