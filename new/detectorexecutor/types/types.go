package types

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Executor interface {
	DetectAt(
		node *language.Node,
		detectorType string,
		evaluator treeevaluatortypes.Evaluator,
	) ([]*detectiontypes.Detection, error)
}
