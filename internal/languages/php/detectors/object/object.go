package object

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/ruleset"
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
	// $user = new <object>;
	assignmentQuery := querySet.Add(`[
		(assignment_expression left: (variable_name) @name right: (_) @value) @root
	]`)

	// class User {
	//   public $name;
	// 	 public $gender;
	//   function set_name($name) {
	//     $this->name = $name;
	//   }
	// }
	classQuery := querySet.Add(`
	(
		class_declaration
			name: (name) @class_name
			body: (
				declaration_list [
						(property_declaration (property_element (variable_name) @name ))
						(method_declaration name: (name) @name)
					]
			)
	) @root`)

	// $user->name;
	// $user->name();
	fieldAccessQuery := querySet.Add(`[
		(member_access_expression object: (_) @object name: (name) @field) @root
		(member_call_expression object: (_) @object name: (name) @field) @root
	]`)

	return &objectDetector{
		assignmentQuery:  assignmentQuery,
		classQuery:       classQuery,
		fieldAccessQuery: fieldAccessQuery,
	}
}

func (detector *objectDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinObjectRule
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	detections, err := detector.getAssignment(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, detectorContext)
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)

	if result == nil || err != nil {
		return nil, err
	}

	rightObjects, err := common.GetNonVirtualObjects(
		detectorContext,
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
				RuleID:    ruleset.BuiltinObjectRule.ID(),
				MatchNode: node,
				Data: common.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
