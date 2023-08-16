package object

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	objectPairQuery *query.Query
	classQuery      *query.Query
	// Naming
	assignmentQuery *query.Query
	// Projection
	memberExpressionQuery     *query.Query
	subscriptExpressionQuery  *query.Query
	callQuery                 *query.Query
	objectDeconstructionQuery *query.Query
	spreadElementQuery        *query.Query
}

func New(querySet *query.Set) types.Detector {
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
	detections, err := detector.getObject(node, evaluationState)
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

func (detector *objectDetector) getObject(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	var properties []generictypes.Property
	spreadResults, err := evaluationState.QueryMatchAt(detector.spreadElementQuery, node)
	if err != nil {
		return nil, err
	}

	for _, spreadResult := range spreadResults {
		detections, err := evaluationState.Evaluate(
			spreadResult["identifier"],
			"object",
			"",
			settings.CURSOR_SCOPE,
			true,
		)

		if err != nil {
			return nil, err
		}
		for _, detection := range detections {
			properties = append(properties, detection.Data.(generictypes.Object).Properties...)
		}
	}

	results, err := evaluationState.QueryMatchAt(detector.objectPairQuery, node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	for _, result := range results {
		var name string
		key := result["key"]
		keyContent := key.Content()

		switch key.Type() {
		case "string": // {"user": "admin_user"}
			name = stringutil.StripQuotes(keyContent)
		case "property_identifier": // { user: "admin_user"}
			name = keyContent
		}

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
		methodName := result["method_name"].Content()
		if methodName == "constructor" {
			for _, param := range result["params"].Children() {
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
				RuleID:    "object",
				MatchNode: node,
				Data: generictypes.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
