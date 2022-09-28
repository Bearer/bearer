package queries

import (
	"github.com/bearer/curio/pkg/parser"
)

type ChildMatch interface {
	Match(*parser.Node) *parser.Node
}
