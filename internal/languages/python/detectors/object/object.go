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
	hashPairQuery *query.Query
	classQuery    *query.Query
	// arrayCreationQuery *query.Query
	// Naming
	assignmentQuery *query.Query
	// Projection
	fieldAccessQuery *query.Query
	subscriptQuery   *query.Query
}

func New(querySet *query.Set) types.Detector {
	// { "foo": "bar" }
	hashPairQuery := querySet.Add(`(dictionary (pair key: (_) @key value: (_) @value) @pair) @root`)

	// user = <object>
	assignmentQuery := querySet.Add(`[
		(assignment left: (identifier) @name right: (_) @value) @root
	]`)

	// class User:
	// 	def __init__(self, name='', gender=''):
	// 			self.name = name
	// 			self.gender = gender
	classQuery := querySet.Add(`
	(
		class_definition
				name: (identifier) @class_name
				body: (block (function_definition
					name: (identifier) @method.name
						parameters: (
							parameters [
									(identifier) @name
									(default_parameter (identifier) @name)
							]
						)
					)
				)
				) @root
				`)

	// user.name
	// user.name()
	fieldAccessQuery := querySet.Add(`[
		(attribute object: (_) @object attribute: (identifier) @field) @root
	]`)

	// user["uuid"]
	subscriptQuery := querySet.Add(`
			(subscript value: (_) @object subscript: (_) @key) @root
	`)

	return &objectDetector{
		hashPairQuery: hashPairQuery,
		classQuery:    classQuery,
		// arrayCreationQuery:       arrayCreationQuery,
		assignmentQuery:  assignmentQuery,
		fieldAccessQuery: fieldAccessQuery,
		subscriptQuery:   subscriptQuery,
	}
}

func (detector *objectDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinObjectRule
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	detections, err := detector.getHash(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, detectorContext)
}

func (detector *objectDetector) getHash(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	results := detector.hashPairQuery.MatchAt(node)
	if len(results) == 0 {
		return nil, nil
	}

	var properties []common.Property
	for _, result := range results {
		pairNode := result["pair"]

		name := result["key"].Content()
		if name == "" {
			continue
		}

		propertyObjects, err := detectorContext.Scan(result["value"], ruleset.BuiltinObjectRule, traversalstrategy.Cursor)
		if err != nil {
			return nil, err
		}

		if len(propertyObjects) == 0 {
			properties = append(properties, common.Property{
				Name: name,
				Node: pairNode,
			})

			continue
		}

		for _, propertyObject := range propertyObjects {
			properties = append(properties, common.Property{
				Name:   name,
				Node:   pairNode,
				Object: propertyObject,
			})
		}
	}

	return []interface{}{common.Object{Properties: properties}}, nil
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

		if result["method.name"].Content() != "__init__" {
			continue
		}

		if result["name"].Content() == "self" {
			continue
		}

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
