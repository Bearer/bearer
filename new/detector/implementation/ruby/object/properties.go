package object

import (
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	"github.com/bearer/curio/new/detector/implementation/ruby/common"
)

func (detector *objectDetector) getProperties(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	var objectParent string
	var objectProperty string

	result, err := detector.callsQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		// user.name
		if result["receiver"].Type() == "identifier" {
			objectParent = result["receiver"].Content()
			objectProperty = result["method"].Content()
		}

		// @user.name
		if result["receiver"].Type() == "instance_variable" {
			objectParent = result["receiver"].Content()[1:]
			objectProperty = result["method"].Content()
		}

		// address.city.zip
		if result["receiver"].Type() == "call" {
			if callProperty := extractFromCall(result["receiver"]); callProperty != nil {
				objectParent = callProperty.Content()
				objectProperty = result["method"].Content()
			}
		}

		// address[:city].zip
		if result["receiver"].Type() == "element_reference" {
			elementReferenceProperty := extractFromElementReference(result["receiver"])

			if elementReferenceProperty != nil {
				objectParent = common.GetLiteralKey(elementReferenceProperty)
				objectProperty = result["method"].Content()
			}
		}
	}

	result, err = detector.elementReferenceQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		key := common.GetLiteralKey(result["key"])
		if key == "" {
			return nil, nil
		}

		// address[:city]
		if result["object"].Type() == "identifier" {
			objectParent = result["object"].Content()
			objectProperty = key
		}

		// @address[:city]
		if result["object"].Type() == "instance_variable" {
			objectParent = result["object"].Content()[1:]
			objectProperty = key
		}

		// address[:city][:zip].international
		if result["object"].Type() == "element_reference" {
			elementReferenceProperty := extractFromElementReference(result["object"])

			if elementReferenceProperty != nil {
				objectParent = common.GetLiteralKey(elementReferenceProperty)
				objectProperty = key
			}
		}

		// address[:city].zip
		if result["object"].Type() == "call" {
			callProperty := extractFromCall(result["object"])

			if callProperty != nil {
				objectParent = callProperty.Content()
				objectProperty = key
			}
		}
	}

	if objectParent != "" && objectProperty != "" {
		return []interface{}{
			generictypes.Object{
				Name: objectParent,
				Properties: []*types.Detection{
					{
						DetectorType: detector.Name(),
						MatchNode:    node,
						Data: generictypes.Property{
							Name: objectProperty,
						},
					},
				},
			},
		}, nil
	}

	return nil, nil
}

func extractFromElementReference(node *tree.Node) *tree.Node {
	for i := 0; i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Type() == "simple_symbol" || child.Type() == "string" {
			return child
		}
	}

	return nil
}

func extractFromCall(node *tree.Node) *tree.Node {
	callProperty := node.ChildByFieldName("method")
	arguments := node.ChildByFieldName("arguments")

	if callProperty.Type() == "identifier" && arguments == nil {
		return callProperty
	}

	return nil
}
