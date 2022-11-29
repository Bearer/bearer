package customdetector

import (
	"github.com/bearer/curio/pkg/parser"
)

func (detector *Detector) IsMatchAnything(node *parser.Node) bool {
	return false
}
