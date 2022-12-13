package detector

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/parser"
)

type Detector interface {
	Name() string
	DetectAt(node *parser.Node, initiator initiatortypes.TreeDetectionInitiator) (*detectiontypes.Detection, error)
}
