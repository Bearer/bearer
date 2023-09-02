package object

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
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

func (detector *objectDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinObjectRule
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	detections, err := detector.getObject(node, detectorContext)
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

func (detector *objectDetector) getObject(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	var properties []common.Property

	spreadResults := detector.spreadElementQuery.MatchAt(node)
	for _, spreadResult := range spreadResults {
		detections, err := detectorContext.Scan(
			spreadResult["identifier"],
			ruleset.BuiltinObjectRule,
			settings.CURSOR_SCOPE,
		)
		if err != nil {
			return nil, err
		}

		for _, detection := range detections {
			properties = append(properties, detection.Data.(common.Object).Properties...)
		}
	}

	results := detector.objectPairQuery.MatchAt(node)
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

		propertyObjects, err := detectorContext.Scan(result["value"], ruleset.BuiltinObjectRule, settings.CURSOR_SCOPE)
		if err != nil {
			return nil, err
		}

		pairNode := result["pair"]

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

	if len(properties) == 0 {
		return nil, nil
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

	valueObjects, err := common.GetNonVirtualObjects(detectorContext, result["value"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range valueObjects {
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
		methodName := result["method_name"].Content()
		if methodName == "constructor" {
			for _, param := range result["params"].Children() {
				if param.Type() != "identifier" {
					continue
				}

				properties = append(properties, common.Property{Name: param.Content()})
			}
		} else {
			properties = append(properties, common.Property{Name: methodName})
		}
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
