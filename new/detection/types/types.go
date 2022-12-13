package types

import (
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/parser"
)

type Detection struct {
	MatchNode  *parser.Node
	ParentNode *parser.Node
	Data       interface{}
}

type QueryHandlerCallback func(captures map[string]*parser.Node, initiator initiatortypes.TreeDetectionInitiator) (*Detection, error)

type QueryHandler struct {
	Query    string
	Callback func(captures map[string]*parser.Node) error
}

type Detector interface {
	Name() string
	Queries() []QueryHandler
}
