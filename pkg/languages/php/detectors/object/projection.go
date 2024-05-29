package object

import (
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	result, err := detector.fieldAccessQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]
		objects, err := common.ProjectObject(
			node,
			detectorContext,
			objectNode,
			getObjectName(objectNode),
			result["field"].Content(),
			true,
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
		propertyName := getElementProperty(result["key"])
		if propertyName == "" {
			return nil, nil
		}

		objects, err := common.ProjectObject(
			node,
			detectorContext,
			objectNode,
			getObjectName(objectNode),
			propertyName,
			false,
		)
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	return nil, nil
}

func getObjectName(objectNode *tree.Node) string {
	switch objectNode.Type() {
	// $user->name()
	// $user->name
	// user->name
	case "variable_name", "name":
		return objectNode.Content()
	// $user->foo->name
	// $user->foo()->name
	case "member_access_expression", "member_call_expression":
		return objectNode.ChildByFieldName("name").Content()
	}

	return ""
}

func getElementProperty(node *tree.Node) string {
	switch node.Type() {
	case "encapsed_string":
		return stringutil.StripQuotes(node.Content())
	default:
		return node.Content()
	}
}
