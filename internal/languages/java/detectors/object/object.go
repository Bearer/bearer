package object

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	detections, err := detector.getAssignment(node, scanContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, scanContext)
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)

	if result == nil || err != nil {
		return nil, err
	}

	rightObjects, err := common.GetNonVirtualObjects(
		scanContext,
		result["value"],
	)
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range rightObjects {
		objects = append(objects, common.Object{
			IsVirtual: true,
			Properties: []common.Property{{
				Name:   result["name"].Content(),
				Node:   node,
				Object: object,
			}},
		})
	}

	return objects, nil
}

func (detector *objectDetector) getClass(node *tree.Node) ([]interface{}, error) {
	results := detector.classQuery.MatchAt(node)
	if len(results) == 0 {
		return nil, nil
	}

	className := results[0]["class_name"].Content()

	var properties []common.Property
	for _, result := range results {
		nameNode := result["name"]

		properties = append(properties, common.Property{
			Name: nameNode.Content(),
			Node: nameNode,
		})
	}

	return []interface{}{common.Object{
		Properties: []common.Property{{
			Name: className,
			Object: &types.Detection{
				RuleID:    "object",
				MatchNode: node,
				Data: common.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
