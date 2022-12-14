package objects

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	Properties []detectiontypes.Detection
}

type objectsDetector struct {
	pairQuery *language.Query
}

func New(lang languagetypes.Language) (detector.Detector, error) {
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
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
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
		property, err := evaluator.NodeDetection(result["pair"], "properties")
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

func (detector *objectsDetector) Close() {
	detector.pairQuery.Close()
}
