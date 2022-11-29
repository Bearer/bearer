package customdetector

import (
	"strings"

	"github.com/bearer/curio/pkg/parser"
)

func (detector *Detector) IsMatchAnything(node *parser.Node) bool {
	return strings.Index(node.Content(), "Var_Anything") == 0
}
