package objects

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/parser"
)

type Data struct {
	Properties []detectiontypes.Detection
}

type objectsDetector struct {
	pairQuery *parser.Query
}

func New(lang *language.Language) (detector.Detector, error) {
	pairQuery, err := lang.CompileQuery(`(hash (pair) @pair) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	return &objectsDetector{pairQuery: pairQuery}, nil
}

func (detector *objectsDetector) Name() string {
	return "objects"
}

func (detector *objectsDetector) DetectAt(
	node *parser.Node,
	initiator initiatortypes.TreeDetectionInitiator,
) (*detectiontypes.Detection, error) {
	results, err := detector.pairQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var properties []detectiontypes.Detection
	for _, result := range results {
		property, err := initiator.NodeDetection(result["pair"], "properties")
		if err != nil {
			return nil, err
		}

		if property != nil {
			properties = append(properties, *property)
		}
	}

	return &detectiontypes.Detection{
		MatchNode: node,
		Data:      Data{Properties: properties},
	}, nil
}
