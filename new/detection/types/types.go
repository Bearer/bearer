package types

import (
	"github.com/bearer/curio/new/parser"
)

type Detection struct {
	MatchNode  *parser.Node
	ParentNode *parser.Node
	Data       interface{}
}
