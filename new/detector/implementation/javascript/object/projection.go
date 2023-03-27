package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"
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
		objectNode, isPropertyAccess := getProjectedObject(result["object"])

		objects, err := generic.ProjectObject(
			node,
			evaluator,
			objectNode,
			getObjectName(objectNode),
			result["property"].Content(),
			isPropertyAccess,
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
		objectNode, isPropertyAccess := getProjectedObject(result["object"])
		propertyName := getSubscriptProperty(result["root"])
		if propertyName == "" {
			return nil, nil
		}

		objects, err := generic.ProjectObject(
			node,
			evaluator,
			objectNode,
			getObjectName(objectNode),
			propertyName,
			isPropertyAccess,
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	result, err = detector.callQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		var properties []generictypes.Property

		functionDetections, err := evaluator.ForTree(result["function"], "object", true)
		if err != nil {
			return nil, err
		}

		for _, detection := range functionDetections {
			properties = append(properties, generictypes.Property{Object: detection})
		}

		if len(properties) != 0 {
			return []interface{}{generictypes.Object{Properties: properties}}, nil
		}
	}

	return nil, nil
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

func getProjectedObject(objectNode *tree.Node) (*tree.Node, bool) {
	if objectNode.Type() == "call_expression" {
		return objectNode.ChildByFieldName("function"), false
	}

	return objectNode, true
}
