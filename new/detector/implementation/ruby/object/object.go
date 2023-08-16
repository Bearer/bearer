package object

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
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

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) NestedDetections() bool {
	return false
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	detections, err := detector.getHash(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getKeywordArgument(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, evaluationState)
}

func (detector *objectDetector) getHash(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	results, err := evaluationState.QueryMatchAt(detector.hashPairQuery, node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	var properties []generictypes.Property
	for _, result := range results {
		pairNode := result["pair"]

		name := common.GetLiteralKey(result["key"])
		if name == "" {
			continue
		}

		propertyObjects, err := evaluationState.Evaluate(
			result["value"],
			"object",
			"",
			settings.NESTED_SCOPE,
			true,
		)
		if err != nil {
			return nil, err
		}

		if len(propertyObjects) == 0 {
			properties = append(properties, generictypes.Property{
				Name: name,
				Node: pairNode,
			})

			continue
		}

		for _, propertyObject := range propertyObjects {
			properties = append(properties, generictypes.Property{
				Name:   name,
				Node:   pairNode,
				Object: propertyObject,
			})
		}
	}

	return []interface{}{generictypes.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getKeywordArgument(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.keywordArgumentQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	name := common.GetLiteralKey(result["key"])
	if name == "" {
		return nil, nil
	}

	propertyObjects, err := evaluationState.Evaluate(
		result["value"],
		"object",
		"",
		settings.NESTED_SCOPE,
		true,
	)
	if err != nil {
		return nil, err
	}

	var properties []generictypes.Property

	if len(propertyObjects) == 0 {
		properties = append(properties, generictypes.Property{
			Name: name,
			Node: node,
		})
	}

	for _, propertyObject := range propertyObjects {
		properties = append(properties, generictypes.Property{
			Name:   name,
			Node:   node,
			Object: propertyObject,
		})
	}

	return []interface{}{generictypes.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := evaluationState.QueryMatchOnceAt(detector.assignmentQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	valueObjects, err := generic.GetNonVirtualObjects(evaluationState, result["value"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range valueObjects {
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
		name := nameNode.Content()

		if nameNode.Type() == "simple_symbol" {
			name = name[1:]
		}

		if name != "initialize" {
			properties = append(properties, generictypes.Property{
				Name: name,
				Node: nameNode,
			})
		}
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
