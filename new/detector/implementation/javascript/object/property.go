package object

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

func (detector *objectDetector) getProperties(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	var objectParent *tree.Node
	var objectProperty *tree.Node

	results, err := detector.memberExpressionQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		// user.name
		if result["object"].Type() == "identifier" {
			objectParent = result["object"]
			objectProperty = result["property"]
		}

		// address.city.zip
		if result["object"].Type() == "member_expression" {
			memberExpressionProperty := extractFromMemberExpression(result["object"])

			if memberExpressionProperty != nil {
				objectParent = memberExpressionProperty
				objectProperty = result["property"]
			}
		}

		// address["city"].zip
		if result["object"].Type() == "subscript_expression" {
			subscriptExpressionProperty := extractFromSubscriptExpression(result["object"])

			if subscriptExpressionProperty != nil {
				objectParent = subscriptExpressionProperty
				objectProperty = result["property"]
			}
		}
	}

	results, err = detector.subscriptExpressionQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		// address["city"]
		if result["object"].Type() == "identifier" {
			objectParent = result["object"]
			objectProperty = result["index"]
		}

		// address["city"]["zip"]
		if result["object"].Type() == "subscript_expression" {
			subscripExpressionProperty := extractFromSubscriptExpression(result["object"])

			if subscripExpressionProperty != nil {
				objectParent = subscripExpressionProperty
				objectProperty = result["index"]
			}
		}

		//address["city"].zip
		if result["object"].Type() == "member_expression" {
			memberExpressionProperty := extractFromMemberExpression(result["object"])

			if memberExpressionProperty != nil {
				objectParent = memberExpressionProperty
				objectProperty = result["index"]
			}
		}
	}

	if objectParent != nil && objectProperty != nil {
		return []interface{}{generictypes.Object{
			Name: objectParent.Content(),
			Properties: []*types.Detection{
				{
					DetectorType: detector.Name(),
					MatchNode:    node,
					Data: generictypes.Property{
						Name: objectProperty.Content(),
					},
				},
			},
		}}, nil
	}

	return nil, nil
}

func extractFromSubscriptExpression(node *tree.Node) *tree.Node {
	property := node.ChildByFieldName("index")
	arguments := node.ChildByFieldName("arguments")

	if property != nil && property.Type() == "string" && arguments == nil {
		return property
	}

	return nil
}

func extractFromMemberExpression(node *tree.Node) *tree.Node {
	memberExpressionProperty := node.ChildByFieldName("property")
	arguments := node.ChildByFieldName("arguments")

	if memberExpressionProperty.Type() == "property_identifier" && arguments == nil {
		return memberExpressionProperty
	}

	return nil
}
