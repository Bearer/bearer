package types

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/parser"
)

type DetectionExecutor interface {
	RegisterDetector(detector detector.Detector) error
	DetectAt(
		node *parser.Node,
		detectorType string,
		initiator initiatortypes.TreeDetectionInitiator,
	) (*detectiontypes.Detection, error)
}
