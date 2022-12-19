package datatypes

import (
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	Name string
}

type datatypesDetector struct{}

func New(lang languagetypes.Language) (detector.Detector, error) {
	return &datatypesDetector{}, nil
}

func (detector *datatypesDetector) Type() string {
	return "datatypes"
}

func (detector *datatypesDetector) DetectAt(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	objectDetections, err := evaluator.NodeDetections(node, "objects")
	if err != nil {
		return nil, err
	}

	var result []*detectiontypes.Detection

	for _, object := range objectDetections {
		data := object.Data
		if data["name"] != "user" {
			continue
		}

		for _, property := range data["properties"].([]*detectiontypes.Detection) {
			propertyData := property.Data
			if propertyData["name"] == "first_name" {
				result = append(result, &detectiontypes.Detection{
					ContextNode: node,
					MatchNode:   property.MatchNode,
					Data:        map[string]interface{}{"name": "Person Name"},
				})
			}
		}
	}

	return result, nil
}

func (detector *datatypesDetector) Close() {}
