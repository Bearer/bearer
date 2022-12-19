package detector

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Detector interface {
	Type() string
	DetectAt(node *language.Node, evaluator treeevaluatortypes.Evaluator) ([]*detectiontypes.Detection, error)
	Close()
}
