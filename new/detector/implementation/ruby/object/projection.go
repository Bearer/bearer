package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.callsQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		receiverNode := result["receiver"]

		objects, err := generic.ProjectObject(
			node,
			evaluator,
			receiverNode,
			getObjectName(receiverNode),
			result["method"].Content(),
			getIsPropertyAccess(receiverNode),
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	result, err = detector.elementReferenceQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]
		propertyName := getElementProperty(result["root"])
		if propertyName == "" {
			return nil, nil
		}

		objects, err := generic.ProjectObject(
			node,
			evaluator,
			objectNode,
			getObjectName(objectNode),
			propertyName,
			getIsPropertyAccess(objectNode),
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	return nil, nil
}

func getObjectName(objectNode *tree.Node) string {
	// user.name or user["name"]
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// @user.name or @user["name"]
	if objectNode.Type() == "instance_variable" {
		return objectNode.Content()[1:]
	}

	// address.city.zip or address.city["zip"]
	if objectNode.Type() == "call" {
		return objectNode.ChildByFieldName("method").Content()
	}

	// address["city"].zip or address["city"]["zip"]
	if objectNode.Type() == "element_reference" {
		return getElementProperty(objectNode)
	}

	return ""
}

func getElementProperty(node *tree.Node) string {
	return common.GetLiteralKey(node.NamedChild(1))
}

func getIsPropertyAccess(objectNode *tree.Node) bool {
	return objectNode.Type() != "call" || objectNode.ChildByFieldName("arguments") == nil
}
