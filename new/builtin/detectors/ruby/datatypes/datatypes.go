package datatypes

import (
	"github.com/bearer/curio/new/builtin/detectors/ruby/objects"
	"github.com/bearer/curio/new/builtin/detectors/ruby/properties"
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

func (detector *datatypesDetector) Name() string {
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
		data := object.Data.(objects.Data)
		if data.Name != "user" {
			continue
		}

		for _, property := range data.Properties {
			propertyData := property.Data.(properties.Data)
			if propertyData.Name == "first_name" {
				result = append(result, &detectiontypes.Detection{
					ContextNode: node,
					MatchNode:   property.MatchNode,
					Data:        Data{Name: "Person Name"},
				})
			}
		}
	}

	return result, nil
}

func (detector *datatypesDetector) Close() {}
