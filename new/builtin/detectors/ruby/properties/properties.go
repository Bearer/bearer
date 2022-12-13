package properties

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/parser"
)

type Data struct {
	Name string
}

type propertiesDetector struct {
	pairQuery *parser.Query
}

func New(lang *language.Language) (detector.Detector, error) {
	pairQuery, err := lang.CompileQuery(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	return &propertiesDetector{pairQuery: pairQuery}, nil
}

func (detector *propertiesDetector) Name() string {
	return "properties"
}

func (detector *propertiesDetector) DetectAt(
	node *parser.Node,
	initiator initiatortypes.TreeDetectionInitiator,
) (*detectiontypes.Detection, error) {
	result, err := detector.pairQuery.MatchAtOnce(node)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	return &detectiontypes.Detection{
		MatchNode: node,
		Data:      Data{Name: result["key"].Content()},
	}, nil
}

func (detector *propertiesDetector) Close() {
	detector.pairQuery.Close()
}
