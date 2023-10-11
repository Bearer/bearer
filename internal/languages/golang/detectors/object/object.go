package object

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	objectQuery *query.Query
	// Naming
	assignmentQuery *query.Query
	// Projection
	fieldAccessQuery *query.Query
}

func New(querySet *query.Set) types.Detector {
	// user = User{}
	// user := User{}
	assignmentQuery := querySet.Add(`[
		(short_var_declaration left: (expression_list . (identifier) @name) right: (_) @value) @root
		(assignment_statement left: (expression_list . (identifier) @name) right: (_) @value) @root
	]`)

	// user.name
	// user.name()
	fieldAccessQuery := querySet.Add(`[
		(selector_expression operand: (_) @object field: (_) @field) @root
	]`)

	// User{ Name: "", Foo: ""}
	// type User struct {
	// 	Name string
	// }
	objectQuery := querySet.Add(`[
		(composite_literal
			type:
				(type_identifier) @object_name
			body:
				(literal_value
					(keyed_element . (_) @name)
			)) @root
	]`)
	// func (x User) FullName() string {} ->
	// (method_declaration
	// 	receiver: (parameter_list (parameter_declaration type: (type_identifier) @object_name))
	// 	name: (field_identifier) @name
	// ) @root
	// (type_declaration
	// 	(type_spec
	// 		name: (type_identifier) @object_name
	// 		type: (struct_type (field_declaration_list (field_declaration name: (field_identifier) @name )))
	// )) @root

	return &objectDetector{
		objectQuery:      objectQuery,
		assignmentQuery:  assignmentQuery,
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
	if node.Type() == "call_expression" {
		detections, err := detectorContext.Scan(node.ChildByFieldName("function"), ruleset.BuiltinObjectRule, traversalstrategy.Cursor)
		if err != nil {
			return nil, err
		}
		results := make([]any, len(detections))
		for i, detection := range detections {
			results[i] = detection.Data
		}

		return results, nil
	}

	detections, err := detector.getAssignment(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getObject(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, detectorContext)
}

func (detector *objectDetector) getObject(node *tree.Node) ([]interface{}, error) {
	results := detector.objectQuery.MatchAt(node)
	if len(results) == 0 {
		return nil, nil
	}

	className := results[0]["object_name"].Content()

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
