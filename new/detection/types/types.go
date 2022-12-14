package types

import "github.com/bearer/curio/new/language"

type Detection struct {
	MatchNode  *language.Node
	ParentNode *language.Node
	Data       interface{}
}
