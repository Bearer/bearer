package types

import "github.com/bearer/curio/new/language/tree"

type Detection struct {
	MatchNode   *tree.Node
	ContextNode *tree.Node
	Data        interface{}
}
