package types

import (
	detectortypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/parser"
)

type TreeDetectionInitiator interface {
	TreeDetections(rootNode *parser.Node, detectorType string) ([]*detectortypes.Detection, error)
	NodeDetection(node *parser.Node, detectorType string) (*detectortypes.Detection, error)
	TreeHasDetection(rootNode *parser.Node, detectorType string) (bool, error)
	NodeHasDetection(node *parser.Node, detectorType string) (bool, error)
}
