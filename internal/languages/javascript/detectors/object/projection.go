package object

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	objects, err := detector.getMemberExpressionProjections(node, scanContext)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	objects, err = detector.getSubscriptExpressionProjections(node, scanContext)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	objects, err = detector.getCallProjections(node, scanContext)
	if len(objects) != 0 || err != nil {
		return objects, err
	}

	return detector.getObjectDeconstructionProjections(node, scanContext)
}

func (detector *objectDetector) getMemberExpressionProjections(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.memberExpressionQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode, isPropertyAccess := getProjectedObject(scanContext, result["object"])

	objects, err := common.ProjectObject(
		node,
		scanContext,
		objectNode,
		getObjectName(scanContext, objectNode),
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.subscriptExpressionQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode, isPropertyAccess := getProjectedObject(scanContext, result["object"])
	propertyName := getSubscriptProperty(scanContext, result["root"])
	if propertyName == "" {
		return nil, nil
	}

	objects, err := common.ProjectObject(
		node,
		scanContext,
		objectNode,
		getObjectName(scanContext, objectNode),
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.callQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	var properties []common.Property

	functionDetections, err := scanContext.Scan(
		result["function"],
		"object",
		"",
		settings.CURSOR_SCOPE,
	)
	if len(functionDetections) == 0 || err != nil {
		return nil, err
	}

	for _, detection := range functionDetections {
		properties = append(properties, common.Property{Object: detection})
	}

	return []interface{}{common.Object{Properties: properties, IsVirtual: true}}, nil
}

func (detector *objectDetector) getObjectDeconstructionProjections(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.objectDeconstructionQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objectNode := result["value"]
	propertyName := result["match"].Content()
	if propertyName == "" {
		return nil, nil
	}

	objects, err := common.ProjectObject(
		node,
		scanContext,
		objectNode,
		getObjectName(scanContext, objectNode),
		propertyName,
		true,
	)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func getObjectName(scanContext types.ScanContext, objectNode *tree.Node) string {
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
		return getSubscriptProperty(scanContext, objectNode)
	}

	return ""
}

func getSubscriptProperty(scanContext types.ScanContext, node *tree.Node) string {
	indexNode := node.ChildByFieldName("index")
	if indexNode.Type() == "string" {
		return stringutil.StripQuotes(indexNode.Content())
	}

	return ""
}

func getProjectedObject(scanContext types.ScanContext, objectNode *tree.Node) (*tree.Node, bool) {
	if objectNode.Type() == "call_expression" {
		return objectNode.ChildByFieldName("function"), false
	}

	return objectNode, true
}
