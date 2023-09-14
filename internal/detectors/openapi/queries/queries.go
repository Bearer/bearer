package queries

import (
	"github.com/bearer/bearer/internal/parser"
)

type ChildMatch interface {
	Match(*parser.Node) *parser.Node
}
