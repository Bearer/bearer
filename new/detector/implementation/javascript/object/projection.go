package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	objects, err := detector.getMemberExpressionProjections(node, evaluationState)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	objects, err = detector.getSubscriptExpressionProjections(node, evaluationState)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	objects, err = detector.getCallProjections(node, evaluationState)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	return detector.getObjectDeconstructionProjections(node, evaluationState)
}

func (detector *objectDetector) getMemberExpressionProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.memberExpressionQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode, isPropertyAccess := getProjectedObject(evaluationState, result["object"])

	objects, err := generic.ProjectObject(
		node,
		evaluationState,
		objectNode,
		getObjectName(evaluationState, objectNode),
		result["property"].Content(),
		isPropertyAccess,
	)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (detector *objectDetector) getSubscriptExpressionProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.subscriptExpressionQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode, isPropertyAccess := getProjectedObject(evaluationState, result["object"])
	propertyName := getSubscriptProperty(evaluationState, result["root"])
	if propertyName == "" {
		return nil, nil
	}

	objects, err := generic.ProjectObject(
		node,
		evaluationState,
		objectNode,
		getObjectName(evaluationState, objectNode),
		propertyName,
		isPropertyAccess,
	)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (detector *objectDetector) getCallProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.callQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	var properties []generictypes.Property

	functionDetections, err := evaluationState.Evaluate(
		result["function"],
		"object",
		"",
		settings.NESTED_SCOPE,
		true,
	)
	if len(functionDetections) == 0 || err != nil {
		return nil, err
	}

	for _, detection := range functionDetections {
		properties = append(properties, generictypes.Property{Object: detection})
	}

	return []interface{}{generictypes.Object{Properties: properties, IsVirtual: true}}, nil
}

func (detector *objectDetector) getObjectDeconstructionProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.objectDeconstructionQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode := result["value"]
	propertyName := result["match"].Content()
	if propertyName == "" {
		return nil, nil
	}

	objects, err := generic.ProjectObject(
		node,
		evaluationState,
		objectNode,
		getObjectName(evaluationState, objectNode),
		propertyName,
		true,
	)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func getObjectName(evaluationState types.EvaluationState, objectNode *tree.Node) string {
	// user.name or user["name"]
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// address.city.zip or address.city["zip"]
	if objectNode.Type() == "member_expression" {
		return evaluationState.NodeFromSitter(objectNode.SitterNode().ChildByFieldName("property")).Content()
	}

	// address["city"].zip or address["city"]["zip"]
	if objectNode.Type() == "subscript_expression" {
		return getSubscriptProperty(evaluationState, objectNode)
	}

	return ""
}

func getSubscriptProperty(evaluationState types.EvaluationState, node *tree.Node) string {
	indexNode := evaluationState.NodeFromSitter(node.SitterNode().ChildByFieldName("index"))
	if indexNode.Type() == "string" {
		return stringutil.StripQuotes(indexNode.Content())
	}

	return ""
}

func getProjectedObject(evaluationState types.EvaluationState, objectNode *tree.Node) (*tree.Node, bool) {
	if objectNode.Type() == "call_expression" {
		return evaluationState.NodeFromSitter(objectNode.SitterNode().ChildByFieldName("function")), false
	}

	return objectNode, true
}
