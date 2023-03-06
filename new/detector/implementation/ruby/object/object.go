package object

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	hashPairQuery *tree.Query
	classQuery    *tree.Query
	// Naming
	assignmentQuery *tree.Query
	// Projection
	callsQuery            *tree.Query
	elementReferenceQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { first_name: ..., ... }
	hashPairQuery, err := lang.CompileQuery(`(hash (pair key: (_) @key value: (_) @value)) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling hash pair query: %s", err)
	}

	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment left: (identifier) @name right: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}

	// class User
	//   attr_accessor :name
	//
	//   def get_first_name()
	//   end
	// end
	classQuery, err := lang.CompileQuery(`
		(class name: (constant) @class_name
			[
				(call arguments: (argument_list (simple_symbol) @name))
				(method name: (identifier) @name)
			]
		) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class query: %s", err)
	}

	// user.name
	callsQuery, err := lang.CompileQuery(`(call receiver: (_) @receiver method: (identifier) @method) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling call query: %s", err)
	}

	// user[:name]
	elementReferenceQuery, err := lang.CompileQuery(`(element_reference object: (_) @object (_) @key) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling element reference query %s", err)
	}

	return &objectDetector{
		hashPairQuery:         hashPairQuery,
		assignmentQuery:       assignmentQuery,
		classQuery:            classQuery,
		callsQuery:            callsQuery,
		elementReferenceQuery: elementReferenceQuery,
	}, nil
}

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) NestedDetections() bool {
	return false
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	detections, err := detector.getHash(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, evaluator)
}

func (detector *objectDetector) getHash(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	results, err := detector.hashPairQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	var properties []generictypes.Property
	for _, result := range results {
		name := common.GetLiteralKey(result["key"])
		if name == "" {
			continue
		}

		propertyObjects, err := generic.GetNonVirtualObjects(evaluator, result["value"])
		if err != nil {
			return nil, err
		}

		if len(propertyObjects) == 0 {
			properties = append(properties, generictypes.Property{
				Name: name,
			})

			continue
		}

		for _, propertyObject := range propertyObjects {
			properties = append(properties, generictypes.Property{
				Name:   name,
				Object: propertyObject,
			})
		}
	}

	return []interface{}{generictypes.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	valueObjects, err := generic.GetNonVirtualObjects(evaluator, result["value"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range valueObjects {
		objects = append(objects, generictypes.Object{
			IsVirtual: true,
			Properties: []generictypes.Property{{
				Name:   result["name"].Content(),
				Object: object,
			}},
		})
	}

	return objects, nil
}

func (detector *objectDetector) getClass(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	results, err := detector.classQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	className := results[0]["class_name"].Content()

	var properties []generictypes.Property
	for _, result := range results {
		name := result["name"].Content()

		if result["name"].Type() == "simple_symbol" {
			name = name[1:]
		}

		if name != "initialize" {
			properties = append(properties, generictypes.Property{Name: name})
		}
	}

	return []interface{}{generictypes.Object{
		Properties: []generictypes.Property{{
			Name: className,
			Object: &types.Detection{
				DetectorType: "object",
				MatchNode:    node,
				Data: generictypes.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}

func (detector *objectDetector) Close() {
	detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.classQuery.Close()
	detector.callsQuery.Close()
	detector.elementReferenceQuery.Close()
}
