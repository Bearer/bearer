package object

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	classQuery *query.Query
	// Naming
	assignmentQuery *query.Query
	// Projection
	fieldAccessQuery *query.Query
}

func New(querySet *query.Set) types.Detector {
	// user = <object>
	// User user = <object>
	assignmentQuery := querySet.Add(`[
		(assignment_expression left: (identifier) @name right: (_) @value) @root
		(
    	local_variable_declaration (
        	variable_declarator (identifier) @name
            value: (object_creation_expression) @value
        )
    ) @root
	]`)

	// class User {
	//    String name
	//	  String getLevel(){}
	// }
	classQuery := querySet.Add(`
		(class_declaration name: (identifier) @class_name
			(class_body
				[
					(field_declaration (variable_declarator name: (identifier) @name))
					(method_declaration name: (identifier) @name)
				]
			)
		) @root`)

	// user.name
	fieldAccessQuery := querySet.Add(`(field_access object: (_) @object field: (identifier) @field) @root`)

	return &objectDetector{
		assignmentQuery:  assignmentQuery,
		classQuery:       classQuery,
		fieldAccessQuery: fieldAccessQuery,
	}
}

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	detections, err := detector.getAssignment(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, evaluationState)
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.assignmentQuery, node)

	if result == nil || err != nil {
		return nil, err
	}

	rightObjects, err := generic.GetNonVirtualObjects(
		evaluationState,
		result["right"],
	)
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range rightObjects {
		objects = append(objects, generictypes.Object{
			IsVirtual: true,
			Properties: []generictypes.Property{{
				Name:   result["name"].Content(),
				Node:   node,
				Object: object,
			}},
		})
	}

	return objects, nil
}

func (detector *objectDetector) getClass(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	results, err := evaluationState.QueryMatchAt(detector.classQuery, node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	className := results[0]["class_name"].Content()

	var properties []generictypes.Property
	for _, result := range results {
		nameNode := result["name"]

		properties = append(properties, generictypes.Property{
			Name: nameNode.Content(),
			Node: nameNode,
		})
	}

	return []interface{}{generictypes.Object{
		Properties: []generictypes.Property{{
			Name: className,
			Object: &detection.Detection{
				RuleID:    "object",
				MatchNode: node,
				Data: generictypes.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
