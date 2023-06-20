package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := detector.fieldAccessQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]

		objects, err := generic.ProjectObject(
			node,
			evaluationState,
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
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// address.city.zip
	if objectNode.Type() == "field_access" {
		return objectNode.ChildByFieldName("field").Content()
	}

	// address["city"].zip or address["city"]["zip"]
	if objectNode.Type() == "method_invocation" {
		return objectNode.ChildByFieldName("name").Content()
	}

	return ""
}
