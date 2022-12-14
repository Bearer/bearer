package types

import "github.com/bearer/curio/new/language"

type Detection struct {
	MatchNode   *language.Node
	ContextNode *language.Node
	Data        interface{}
}
