package types

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language/tree"
)

type Evaluator interface {
	TreeDetections(rootNode *tree.Node, detectorType string) ([]*detectiontypes.Detection, error)
	NodeDetections(node *tree.Node, detectorType string) ([]*detectiontypes.Detection, error)
	TreeHasDetection(rootNode *tree.Node, detectorType string) (bool, error)
	NodeHasDetection(node *tree.Node, detectorType string) (bool, error)
}
