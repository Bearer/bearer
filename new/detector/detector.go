package detector

import (
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
	"github.com/open-policy-agent/opa/ast"
)

type Detector interface {
	Type() string
	DetectAt(node *language.Node, evaluator treeevaluatortypes.Evaluator) (*ast.Array, error)
	Close()
}
