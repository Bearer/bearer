package properties

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/parser"
)

type propertiesDetector struct {
	lang *language.Language
}

func New(lang *language.Language) detector.Detector {
	return &propertiesDetector{}
}

func (detector *propertiesDetector) Name() string {
	return "properties"
}

func (detector *propertiesDetector) DetectAt(
	node *parser.Node,
	initiator initiatortypes.TreeDetectionInitiator,
) (*detectiontypes.Detection, error) {
	return nil, nil
}
