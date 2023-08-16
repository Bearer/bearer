package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/tree"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.callsQuery, node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		receiverNode := result["receiver"]
		astReceiverNode := receiverNode

		objects, err := generic.ProjectObject(
			node,
			evaluationState,
			astReceiverNode,
			getObjectName(evaluationState, astReceiverNode),
			result["method"].Content(),
			getIsPropertyAccess(astReceiverNode),
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	result, err = evaluationState.QueryMatchOnceAt(detector.elementReferenceQuery, node)
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
			evaluationState,
			objectNode,
			getObjectName(evaluationState, objectNode),
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

func getObjectName(
	evaluationState types.EvaluationState,
	objectNode *tree.Node,
) string {
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
	return common.GetLiteralKey(node.Children()[2])
}

func getIsPropertyAccess(objectNode *tree.Node) bool {
	return objectNode.Type() != "call" || objectNode.ChildByFieldName("arguments") == nil
}
