package object

import (
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"

	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
)

type objectDetector struct {
	types.DetectorBase
	// Struct literals
	structExpressionQuery *query.Query
	// Naming (let assignments)
	assignmentQuery *query.Query
	// Projection (field access)
	fieldAccessQuery *query.Query
}

func New(querySet *query.Set) types.Detector {
	// Struct { field: value, ... }
	structExpressionQuery := querySet.Add(`[
		(struct_expression
			name: (_) @struct_name
			body: (field_initializer_list
				(field_initializer
					name: (field_identifier) @field_name
					value: (_) @field_value
				) @field
			)
		) @root
	]`)

	// let user = <object>;
	// let mut user = <object>;
	assignmentQuery := querySet.Add(`[
		(let_declaration
			pattern: (identifier) @name
			value: (_) @value
		) @root
	]`)

	// user.name
	// user.name()
	fieldAccessQuery := querySet.Add(`[
		(field_expression
			value: (_) @object
			field: (field_identifier) @field
		) @root
	]`)

	return &objectDetector{
		structExpressionQuery: structExpressionQuery,
		assignmentQuery:       assignmentQuery,
		fieldAccessQuery:      fieldAccessQuery,
	}
}

func (detector *objectDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinObjectRule
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	// Handle method calls - scan the receiver
	if node.Type() == "call_expression" {
		function := node.ChildByFieldName("function")
		if function != nil {
			detections, err := detectorContext.Scan(function, ruleset.BuiltinObjectRule, traversalstrategy.Cursor)
			if err != nil {
				return nil, err
			}
			results := make([]any, len(detections))
			for i, detection := range detections {
				results[i] = detection.Data
			}
			return results, nil
		}
	}

	detections, err := detector.getAssignment(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getStructExpression(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, detectorContext)
}

func (detector *objectDetector) getStructExpression(node *tree.Node) ([]interface{}, error) {
	results := detector.structExpressionQuery.MatchAt(node)
	if len(results) == 0 {
		return nil, nil
	}

	structName := results[0]["struct_name"].Content()

	var properties []common.Property
	for _, result := range results {
		fieldNameNode := result["field_name"]
		if fieldNameNode == nil {
			continue
		}

		properties = append(properties, common.Property{
			Name: fieldNameNode.Content(),
			Node: result["field"],
		})
	}

	return []interface{}{common.Object{
		Properties: []common.Property{{
			Name: structName,
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

