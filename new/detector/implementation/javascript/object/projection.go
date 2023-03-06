package object

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.memberExpressionQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]

		objects, err := projectObject(
			node,
			evaluator,
			objectNode,
			result["property"].Content(),
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	result, err = detector.subscriptExpressionQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]
		propertyName := getSubscriptProperty(result["root"])
		if propertyName == "" {
			return nil, nil
		}

		objects, err := projectObject(
			node,
			evaluator,
			objectNode,
			propertyName,
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	return nil, nil
}

func projectObject(
	node *tree.Node,
	evaluator types.Evaluator,
	objectNode *tree.Node,
	propertyName string,
) ([]interface{}, error) {
	var result []interface{}

	objectDetections, err := getNonVirtualObjects(evaluator, objectNode)
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

	if objectName := getObjectName(objectNode); objectName != "" {
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

func getObjectName(objectNode *tree.Node) string {
	// user.name or user["name"]
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// address.city.zip or address.city["zip"]
	if objectNode.Type() == "member_expression" {
		return objectNode.ChildByFieldName("property").Content()
	}

	// address["city"].zip or address["city"]["zip"]
	if objectNode.Type() == "subscript_expression" {
		return getSubscriptProperty(objectNode)
	}

	return ""
}

func getSubscriptProperty(node *tree.Node) string {
	indexNode := node.ChildByFieldName("index")
	if indexNode.Type() == "string" {
		return stringutil.StripQuotes(indexNode.Content())
	}

	return ""
}
