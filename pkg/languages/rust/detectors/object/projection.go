package object

import (
	"github.com/bearer/bearer/pkg/scanner/ast/tree"

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

	return nil, nil
}

func getObjectName(objectNode *tree.Node) string {
	// user.name
	// user.name()
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// self.name
	if objectNode.Type() == "self" {
		return "self"
	}

	return ""
}

