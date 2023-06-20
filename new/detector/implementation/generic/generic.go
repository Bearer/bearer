package generic

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/detection"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

func GetNonVirtualObjects(
	evaluationState types.EvaluationState,
	node *tree.Node,
) ([]*detection.Detection, error) {
	detections, err := evaluationState.Evaluate(node, "object", "", settings.CURSOR_SCOPE, true)
	if err != nil {
		return nil, err
	}

	var result []*detection.Detection
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
	evaluationState types.EvaluationState,
	objectNode *tree.Node,
	objectName,
	propertyName string,
	isPropertyAccess bool,
) ([]interface{}, error) {
	var result []interface{}

	if isPropertyAccess {
		objectDetections, err := GetNonVirtualObjects(evaluationState, objectNode)
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
							Object: &detection.Detection{
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
				Object: &detection.Detection{
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

func GetStringValue(node *tree.Node, evaluationState types.EvaluationState) (string, bool, error) {
	detections, err := evaluationState.Evaluate(node, "string", "", settings.CURSOR_SCOPE, true)
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
		return "", false, fmt.Errorf(
			"expected single string detection but got %d for %s",
			len(detections),
			node.Debug(true),
		)
	}
}

func ConcatenateChildStrings(node *tree.Node, evaluationState types.EvaluationState) ([]interface{}, error) {
	value := ""
	isLiteral := true

	for i := 0; i < node.ChildCount(); i += 1 {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}

		childValue, childIsLiteral, err := GetStringValue(child, evaluationState)
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

func ConcatenateAssignEquals(node *tree.Node, evaluationState types.EvaluationState) ([]interface{}, error) {
	unifiedNodes := node.ChildByFieldName("left").UnifiedNodes()
	if len(unifiedNodes) == 0 {
		return nil, nil
	}
	if len(unifiedNodes) != 1 {
		return nil, fmt.Errorf("expected exactly one unified `+=` node but got %d", len(unifiedNodes))
	}

	left, leftIsLiteral, err := GetStringValue(unifiedNodes[0], evaluationState)
	if err != nil {
		return nil, err
	}

	right, rightIsLiteral, err := GetStringValue(node.ChildByFieldName("right"), evaluationState)
	if err != nil {
		return nil, err
	}

	if left == "" && !leftIsLiteral {
		left = "*"

		// No detection when neither parts are a string
		if right == "" && !rightIsLiteral {
			return nil, nil
		}
	}

	if right == "" && !rightIsLiteral {
		right = "*"
	}

	return []interface{}{generictypes.String{
		Value:     left + right,
		IsLiteral: leftIsLiteral && rightIsLiteral,
	}}, nil
}
