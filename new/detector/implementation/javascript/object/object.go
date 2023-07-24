package object

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	objectPairQuery *tree.Query
	classQuery      *tree.Query
	// Naming
	assignmentQuery *tree.Query
	// Projection
	memberExpressionQuery     *tree.Query
	subscriptExpressionQuery  *tree.Query
	callQuery                 *tree.Query
	objectDeconstructionQuery *tree.Query
	spreadElementQuery        *tree.Query
}

func New(querySet *tree.QuerySet) (types.Detector, error) {
	// { first_name: ..., ... }
	objectPairQuery := querySet.Add(`(object (pair key: (_) @key value: (_) @value) @pair) @root`)

	// user = <object>
	// const user = <object>
	// var user = <object>
	// let user = <object>
	assignmentQuery := querySet.Add(`[
		(assignment_expression left: (identifier) @name right: (_) @value)
		(variable_declarator name: (identifier) @name value: (_) @value)
	] @root`)

	// const { user } = <object>
	// let { user } = <object>
	// var { user } = <object>
	objectDeconstructionQuery := querySet.Add(`(variable_declarator name: (object_pattern (shorthand_property_identifier_pattern) @match) value: (_) @value) @root`)

	// { ...user, foo: "bar" }
	spreadElementQuery := querySet.Add(`(object (spread_element (identifier) @identifier)) @root`)

	// class User {
	//   constructor(name, surname) {}
	//   GetName() {}
	// }
	classQuery := querySet.Add(`
		(class_declaration
		  name: (type_identifier) @class_name
      body: (class_body
        (method_definition name: (property_identifier) @method_name (formal_parameters) @params)
      )
    ) @root`)

	// user.name
	memberExpressionQuery := querySet.Add(`(member_expression object: (_) @object property: (property_identifier) @property) @root`)

	// user[:name]
	subscriptExpressionQuery := querySet.Add(`(subscript_expression object: (_) @object index: (string) @index ) @root`)

	callQuery := querySet.Add(`(call_expression function: (_) @function) @root`)

	return &objectDetector{
		objectPairQuery:           objectPairQuery,
		assignmentQuery:           assignmentQuery,
		spreadElementQuery:        spreadElementQuery,
		objectDeconstructionQuery: objectDeconstructionQuery,
		classQuery:                classQuery,
		memberExpressionQuery:     memberExpressionQuery,
		subscriptExpressionQuery:  subscriptExpressionQuery,
		callQuery:                 callQuery,
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
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	detections, err := detector.getObject(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, evaluationState)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, evaluationState)
}

func (detector *objectDetector) getObject(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	var properties []generictypes.Property
	spreadResults, err := detector.spreadElementQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	for _, spreadResult := range spreadResults {
		detections, err := evaluationState.Evaluate(spreadResult["identifier"], "object", "", settings.CURSOR_SCOPE, true)

		if err != nil {
			return nil, err
		}
		for _, detection := range detections {
			properties = append(properties, detection.Data.(generictypes.Object).Properties...)
		}
	}

	results, err := detector.objectPairQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	for _, result := range results {
		var name string
		key := result["key"]

		switch key.Type() {
		case "string": // {"user": "admin_user"}
			name = stringutil.StripQuotes(key.Content())
		case "property_identifier": // { user: "admin_user"}
			name = key.Content()
		}

		if name == "" {
			continue
		}

		propertyObjects, err := evaluationState.Evaluate(result["value"], "object", "", settings.NESTED_SCOPE, true)
		if err != nil {
			return nil, err
		}

		pairNode := result["pair"]

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

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
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

func (detector *objectDetector) getClass(node *tree.Node) ([]interface{}, error) {
	results, err := detector.classQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	className := results[0]["class_name"].Content()

	var properties []generictypes.Property
	for _, result := range results {
		methodName := result["method_name"].Content()
		if methodName == "constructor" {
			params := result["params"]

			for i := 0; i < params.ChildCount(); i++ {
				param := params.Child(i)
				if param.Type() != "identifier" {
					continue
				}

				properties = append(properties, generictypes.Property{Name: param.Content()})
			}
		} else {
			properties = append(properties, generictypes.Property{Name: methodName})
		}
	}

	return []interface{}{generictypes.Object{
		Properties: []generictypes.Property{{
			Name: className,
			Object: &detection.Detection{
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
}
