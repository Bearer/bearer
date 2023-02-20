package queries

import (
	"github.com/bearer/bearer/pkg/parser"
)

type ChildMatch interface {
	Match(*parser.Node) *parser.Node
}
