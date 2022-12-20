package detector

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language/tree"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, evaluator treeevaluatortypes.Evaluator) ([]*detectiontypes.Detection, error)
	Close()
}
