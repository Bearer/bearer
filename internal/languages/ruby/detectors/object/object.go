package object

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

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

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	detections, err := detector.getHash(node, scanContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getKeywordArgument(node, scanContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssignment(node, scanContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, scanContext)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, scanContext)
}

func (detector *objectDetector) getHash(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	results, err := scanContext.QueryMatchAt(detector.hashPairQuery, node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	var properties []detectorscommon.Property
	for _, result := range results {
		pairNode := result["pair"]

		name := common.GetLiteralKey(result["key"])
		if name == "" {
			continue
		}

		propertyObjects, err := scanContext.Scan(
			result["value"],
			"object",
			"",
			settings.CURSOR_SCOPE,
		)
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := scanContext.QueryMatchOnceAt(detector.keywordArgumentQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	name := common.GetLiteralKey(result["key"])
	if name == "" {
		return nil, nil
	}

	propertyObjects, err := scanContext.Scan(
		result["value"],
		"object",
		"",
		settings.CURSOR_SCOPE,
	)
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	result, err := scanContext.QueryMatchOnceAt(detector.assignmentQuery, node)
	if result == nil || err != nil {
		return nil, err
	}

	valueObjects, err := detectorscommon.GetNonVirtualObjects(scanContext, result["value"])
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
	scanContext types.ScanContext,
) ([]interface{}, error) {
	results, err := scanContext.QueryMatchAt(detector.classQuery, node)
	if len(results) == 0 || err != nil {
		return nil, err
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
				RuleID:    "object",
				MatchNode: node,
				Data: detectorscommon.Object{
					Properties: properties,
				},
			},
		}},
	}}, nil
}
