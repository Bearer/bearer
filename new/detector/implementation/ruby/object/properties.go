package object

import (
	"github.com/bearer/curio/new/detector/implementation/ruby/property"
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
)

func (detector *objectDetector) getProperties(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	var objectParent *tree.Node
	var objectProperty *tree.Node

	results, err := detector.callsQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		// user.name || @user.name
		if result["receiver"].Type() == "identifier" || result["receiver"].Type() == "instance_variable" {
			objectParent = result["receiver"]
			objectProperty = result["method"]
		}

		// address.city.zip
		if result["receiver"].Type() == "call" {
			callProperty := extractFromCall(result["receiver"])

			if callProperty != nil {
				objectParent = callProperty
				objectProperty = result["method"]
			}
		}

		// address[:city].zip
		if result["receiver"].Type() == "element_reference" {
			elementReferenceProperty := extractFromElementReference(result["object"])

			if elementReferenceProperty != nil {
				objectParent = elementReferenceProperty
				objectProperty = result["method"]
			}
		}
	}

	results, err = detector.elementReferenceQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		// address[:city] || @address[:city]
		if result["object"].Type() == "identifier" || result["object"].Type() == "instance_variable" {
			objectParent = result["object"]
			objectProperty = result["simple_symbol"]
		}

		// address[:city][:zip].international
		if result["object"].Type() == "element_reference" {
			elementReferenceProperty := extractFromElementReference(result["object"])

			if elementReferenceProperty != nil {
				objectParent = elementReferenceProperty
				objectProperty = result["simple_symbol"]
			}
		}

		// address[:city].zip
		if result["object"].Type() == "call" {
			callProperty := extractFromCall(result["object"])

			if callProperty != nil {
				objectParent = callProperty
				objectProperty = result["simple_symbol"]
			}
		}
	}

	if objectParent != nil && objectProperty != nil {
		return []*types.Detection{{
			MatchNode:   node,
			ContextNode: node,
			Data: Data{
				Name: objectParent.Content(),
				Properties: []*types.Detection{
					{
						MatchNode: node,
						Data: property.Data{
							Name: objectProperty.Content(),
						},
					},
				},
			},
		}}, nil
	}

	return nil, nil
}

func extractFromElementReference(node *tree.Node) *tree.Node {
	for i := 0; i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Type() == "simple_symbol" {
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
