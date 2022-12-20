package datatype

import (
	objectdetector "github.com/bearer/curio/new/builtin/detectors/ruby/object"
	propertydetector "github.com/bearer/curio/new/builtin/detectors/ruby/property"
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	Name string
}

type datatypeDetector struct{}

func New(lang languagetypes.Language) (detector.Detector, error) {
	return &datatypeDetector{}, nil
}

func (detector *datatypeDetector) Name() string {
	return "datatype"
}

func (detector *datatypeDetector) DetectAt(
	node *tree.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	objectDetections, err := evaluator.NodeDetections(node, "object")
	if err != nil {
		return nil, err
	}

	var result []*detectiontypes.Detection

	for _, object := range objectDetections {
		data := object.Data.(objectdetector.Data)
		if data.Name != "user" {
			continue
		}

		for _, property := range data.Properties {
			propertyData := property.Data.(propertydetector.Data)
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

func (detector *datatypeDetector) Close() {}
