package objects

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/parser"
)

type objectsDetector struct {
	lang *language.Language
}

func New(lang *language.Language) detector.Detector {
	return &objectsDetector{}
}

func (detector *objectsDetector) Name() string {
	return "objects"
}

func (detector *objectsDetector) DetectAt(
	node *parser.Node,
	initiator initiatortypes.TreeDetectionInitiator,
) (*detectiontypes.Detection, error) {
	return nil, nil
}
