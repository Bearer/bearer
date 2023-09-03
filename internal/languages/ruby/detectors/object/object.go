package object

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"

	"github.com/bearer/bearer/internal/languages/ruby/detectors/common"
	detectorscommon "github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	hashPairQuery        *query.Query
	keywordArgumentQuery *query.Query
	classQuery           *query.Query
	// Naming
	assignmentQuery *query.Query
	// Projection
	callsQuery            *query.Query
	elementReferenceQuery *query.Query
}

func New(querySet *query.Set) types.Detector {
	// { first_name: ..., ... }
	hashPairQuery := querySet.Add(`(hash (pair key: (_) @key value: (_) @value) @pair) @root`)

	// call(first_name: ...)
	keywordArgumentQuery := querySet.Add(`(argument_list (pair key: (_) @key value: (_) @value) @match) @root`)

	// user = <object>
	assignmentQuery := querySet.Add(`(assignment left: (identifier) @name right: (_) @value) @root`)

	// class User
	//   attr_accessor :name
	//
	//   def get_first_name()
	//   end
	// end
	classQuery := querySet.Add(`
		(class name: (constant) @class_name
			[
				(call arguments: (argument_list (simple_symbol) @name))
				(method name: (identifier) @name)
			]
		) @root`)

	// user.name
	callsQuery := querySet.Add(`(call receiver: (_) @receiver method: (identifier) @method) @root`)

	// user[:name]
	elementReferenceQuery := querySet.Add(`(element_reference object: (_) @object . (_) @key . ) @root`)

	return &objectDetector{
		hashPairQuery:         hashPairQuery,
		keywordArgumentQuery:  keywordArgumentQuery,
		assignmentQuery:       assignmentQuery,
		classQuery:            classQuery,
		callsQuery:            callsQuery,
		elementReferenceQuery: elementReferenceQuery,
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

	detections, err = detector.getKeywordArgument(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, detectorContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, detectorContext)
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

	var properties []detectorscommon.Property
	for _, result := range results {
		pairNode := result["pair"]

		name := common.GetLiteralKey(result["key"])
		if name == "" {
			continue
		}

		propertyObjects, err := detectorContext.Scan(result["value"], ruleset.BuiltinObjectRule, traversalstrategy.Cursor)
		if err != nil {
			return nil, err
		}

		if len(propertyObjects) == 0 {
			properties = append(properties, detectorscommon.Property{
				Name: name,
				Node: pairNode,
			})

			continue
		}

		for _, propertyObject := range propertyObjects {
			properties = append(properties, detectorscommon.Property{
				Name:   name,
				Node:   pairNode,
				Object: propertyObject,
			})
		}
	}

	return []interface{}{detectorscommon.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getKeywordArgument(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	result, err := detector.keywordArgumentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	name := common.GetLiteralKey(result["key"])
	if name == "" {
		return nil, nil
	}

	propertyObjects, err := detectorContext.Scan(result["value"], ruleset.BuiltinObjectRule, traversalstrategy.Cursor)
	if err != nil {
		return nil, err
	}

	var properties []detectorscommon.Property

	if len(propertyObjects) == 0 {
		properties = append(properties, detectorscommon.Property{
			Name: name,
			Node: node,
		})
	}

	for _, propertyObject := range propertyObjects {
		properties = append(properties, detectorscommon.Property{
			Name:   name,
			Node:   node,
			Object: propertyObject,
		})
	}

	return []interface{}{detectorscommon.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	valueObjects, err := detectorscommon.GetNonVirtualObjects(detectorContext, result["value"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range valueObjects {
		objects = append(objects, detectorscommon.Object{
			IsVirtual: true,
			Properties: []detectorscommon.Property{{
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
	detectorContext types.Context,
) ([]interface{}, error) {
	results := detector.classQuery.MatchAt(node)
	if len(results) == 0 {
		return nil, nil
	}

	className := results[0]["class_name"].Content()

	var properties []detectorscommon.Property
	for _, result := range results {
		nameNode := result["name"]
		name := nameNode.Content()

		if nameNode.Type() == "simple_symbol" {
			name = name[1:]
		}

		if name != "initialize" {
			properties = append(properties, detectorscommon.Property{
				Name: name,
				Node: nameNode,
			})
		}
	}

	return []interface{}{detectorscommon.Object{
		Properties: []detectorscommon.Property{{
			Name: className,
			Object: &types.Detection{
				RuleID:    ruleset.BuiltinObjectRule.ID(),
				MatchNode: node,
				Data: detectorscommon.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
