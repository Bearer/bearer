package object

import (
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/languages/ruby/detectors/common"
	detectorscommon "github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.callsQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		receiverNode := result["receiver"]
		astReceiverNode := receiverNode

		objects, err := detectorscommon.ProjectObject(
			node,
			scanContext,
			astReceiverNode,
			getObjectName(astReceiverNode),
			result["method"].Content(),
			getIsPropertyAccess(astReceiverNode),
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

		objects, err := detectorscommon.ProjectObject(
			node,
			scanContext,
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
	return common.GetLiteralKey(node.Children()[2])
}

func getIsPropertyAccess(objectNode *tree.Node) bool {
	return objectNode.Type() != "call" || objectNode.ChildByFieldName("arguments") == nil
}
