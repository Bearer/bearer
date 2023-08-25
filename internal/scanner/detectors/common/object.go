package common

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

type Object struct {
	Properties []Property
	// IsVirtual describes whether this object actually exists, or has
	// been detected as part of a variable name
	IsVirtual bool
}

type Property struct {
	Name   string
	Node   *tree.Node
	Object *types.Detection
}

func GetNonVirtualObjects(
	detectorContext types.Context,
	node *tree.Node,
) ([]*types.Detection, error) {
	detections, err := detectorContext.Scan(node, "object", settings.CURSOR_SCOPE)
	if err != nil {
		return nil, err
	}

	var result []*types.Detection
	for _, detection := range detections {
		data := detection.Data.(Object)
		if !data.IsVirtual {
			result = append(result, detection)
		}
	}

	return result, nil
}

func ProjectObject(
	node *tree.Node,
	detectorContext types.Context,
	objectNode *tree.Node,
	objectName,
	propertyName string,
	isPropertyAccess bool,
) ([]interface{}, error) {
	var result []interface{}

	if isPropertyAccess {
		objectDetections, err := GetNonVirtualObjects(detectorContext, objectNode)
		if err != nil {
			return nil, err
		}

		for _, objectDetection := range objectDetections {
			objectData := objectDetection.Data.(Object)

			for _, property := range objectData.Properties {
				if property.Name == propertyName && property.Object != nil {
					result = append(result, property.Object.Data)
					result = append(result, Object{
						Properties: []Property{{
							Name: propertyName,
							Object: &types.Detection{
								RuleID:    "object",
								MatchNode: node,
								Data:      property.Object.Data,
							},
						}},
						IsVirtual: true,
					})
				}
			}
		}
	}

	if objectName != "" {
		result = append(result, Object{
			Properties: []Property{{
				Name: objectName,
				Object: &types.Detection{
					RuleID:    "object",
					MatchNode: node,
					Data: Object{
						Properties: []Property{{Name: propertyName}},
						IsVirtual:  true,
					},
				},
			}},
			IsVirtual: true,
		})
	}

	return result, nil
}
