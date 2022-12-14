package types

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language"
)

type Evaluator interface {
	TreeDetections(rootNode *language.Node, detectorType string) ([]*detectiontypes.Detection, error)
	NodeDetection(node *language.Node, detectorType string) (*detectiontypes.Detection, error)
	TreeHasDetection(rootNode *language.Node, detectorType string) (bool, error)
	NodeHasDetection(node *language.Node, detectorType string) (bool, error)
}
