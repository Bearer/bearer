package types

import (
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
	"github.com/open-policy-agent/opa/ast"
)

type Executor interface {
	DetectAt(
		node *language.Node,
		detectorType string,
		evaluator treeevaluatortypes.Evaluator,
	) (*ast.Array, error)
}
