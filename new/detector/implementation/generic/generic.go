package generic

import (
	"fmt"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

func GetNonVirtualObjects(
	evaluator types.Evaluator,
	node *tree.Node,
) ([]*types.Detection, error) {
	detections, err := evaluator.Evaluate(node, "object", "", settings.CURSOR_SCOPE, true)
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

func GetStringValue(node *tree.Node, evaluator types.Evaluator) (string, bool, error) {
	detections, err := evaluator.Evaluate(node, "string", "", settings.CURSOR_SCOPE, true)
	if err != nil {
		return "", false, err
	}

	switch len(detections) {
	case 0:
		return "", false, nil
	case 1:
		childString := detections[0].Data.(generictypes.String)

		return childString.Value, childString.IsLiteral, nil
	default:
		return "", false, fmt.Errorf("expected single string detection but got %d", len(detections))
	}
}

func ConcatenateChildStrings(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	value := ""
	isLiteral := true

	for i := 0; i < node.ChildCount(); i += 1 {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}

		childValue, childIsLiteral, err := GetStringValue(child, evaluator)
		if err != nil {
			return nil, err
		}

		if childValue == "" && !childIsLiteral {
			childValue = "*"
		}

		value += childValue

		if !childIsLiteral {
			isLiteral = false
		}
	}

	return []interface{}{generictypes.String{
		Value:     value,
		IsLiteral: isLiteral,
	}}, nil
}
