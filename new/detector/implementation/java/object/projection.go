package object

import (
	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/tree"
)

func (detector *objectDetector) getProjections(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.fieldAccessQuery, node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		objectNode := result["object"]

		objects, err := generic.ProjectObject(
			node,
			evaluationState,
			objectNode,
			getObjectName(evaluationState, objectNode),
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

func getObjectName(evaluationState types.EvaluationState, objectNode *tree.Node) string {
	// user.name
	if objectNode.Type() == "identifier" {
		return objectNode.Content()
	}

	// address.city.zip
	if objectNode.Type() == "field_access" {
		// FIXME: implement field names
		return objectNode.ChildByFieldName("field").Content()
	}

	// address["city"].zip or address["city"]["zip"]
	if objectNode.Type() == "method_invocation" {
		return objectNode.ChildByFieldName("name").Content()
	}

	return ""
}
