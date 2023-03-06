package generic

import (
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
)

func GetNonVirtualObjects(evaluator types.Evaluator, node *tree.Node) ([]*types.Detection, error) {
	detections, err := evaluator.ForNode(node, "object", true)
	if err != nil {
		return nil, err
	}

	var result []*types.Detection
	for _, detection := range detections {
		data := detection.Data.(generictypes.Object)
		if !data.IsVirtual {
			result = append(result, detection)
		}
	}

	return result, nil
}

func ProjectObject(
	node *tree.Node,
	evaluator types.Evaluator,
	objectNode *tree.Node,
	objectName,
	propertyName string,
	isPropertyAccess bool,
) ([]interface{}, error) {
	var result []interface{}

	if isPropertyAccess {
		objectDetections, err := GetNonVirtualObjects(evaluator, objectNode)
		if err != nil {
			return nil, err
		}

		for _, objectDetection := range objectDetections {
			objectData := objectDetection.Data.(generictypes.Object)

			for _, property := range objectData.Properties {
				if property.Name == propertyName && property.Object != nil {
					result = append(result, property.Object.Data)
					result = append(result, generictypes.Object{
						Properties: []generictypes.Property{{
							Name: propertyName,
							Object: &types.Detection{
								DetectorType: "object",
								MatchNode:    node,
								Data:         property.Object.Data,
							},
						}},
						IsVirtual: true,
					})
				}
			}
		}
	}

	if objectName != "" {
		result = append(result, generictypes.Object{
			Properties: []generictypes.Property{{
				Name: objectName,
				Object: &types.Detection{
					DetectorType: "object",
					MatchNode:    node,
					Data: generictypes.Object{
						Properties: []generictypes.Property{{Name: propertyName}},
						IsVirtual:  true,
					},
				},
			}},
			IsVirtual: true,
		})
	}

	return result, nil
}
