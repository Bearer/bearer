package object

import (
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/rs/zerolog/log"
)

func (detector *objectDetector) getVariableDeclaration(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.variableDeclarationQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["value"], "object", true)
	if err != nil {
		return nil, err
	}

	var detections []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		if objectData.Name == "" {
			detections = append(detections, generictypes.Object{
				Name:       result["name"].Content(),
				Properties: objectData.Properties,
			},
			)
		}
	}

	return detections, nil
}

func (detector *objectDetector) getObjectDeconstruction(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.objectDeconstructionQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["value"], "object", true)
	if err != nil {
		return nil, err
	}

	var detections []interface{}

	if len(objects) != 1 {
		return detections, nil
	}

	object := objects[0]
	objectData := object.Data.(generictypes.Object)

	name := result["match"]

	for i := 0; i < name.ChildCount(); i++ {
		child := name.Child(i)
		if child.Type() != "shorthand_property_identifier_pattern" {
			log.Debug().Msgf("child type is %s", child.Type())
			continue
		}

		detections = append(detections, generictypes.Object{
			Name: objectData.Name,
			Properties: []*types.Detection{
				{
					DetectorType: detector.Name(),
					MatchNode:    node,
					Data: generictypes.Property{
						Name: child.Content(),
					},
				},
			},
		})
	}

	return detections, nil
}
