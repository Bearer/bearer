package object

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	classQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	// Projection
	fieldAccessQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment_expression left: (identifier) @name right: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}

	// class User {
	//    String name
	//	  String getLevel(){}
	// }
	classQuery, err := lang.CompileQuery(`
		(class_declaration name: (identifier) @class_name
			(class_body
				[
					(field_declaration (variable_declarator name: (identifier) @name))
					(method_declaration name: (identifier) @name)
				]
			)
		) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class query: %s", err)
	}

	// user.name
	fieldAccessQuery, err := lang.CompileQuery(`(field_access object: (_) @object field: (identifier) @field) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling call query: %s", err)
	}

	return &objectDetector{
		assignmentQuery:  assignmentQuery,
		classQuery:       classQuery,
		fieldAccessQuery: fieldAccessQuery,
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
	ruleReferenceType settings.RuleReferenceScope,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	log.Debug().Msgf("node is %s", node.Debug())

	detections, err := detector.getAssignment(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.getProjections(node, evaluator)
}

func (detector *objectDetector) getAssignment(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	rightObjects, err := generic.GetNonVirtualObjects(evaluator, result["right"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range rightObjects {
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

func (detector *objectDetector) getClass(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	results, err := detector.classQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	className := results[0]["class_name"].Content()

	var properties []generictypes.Property
	for _, result := range results {
		nameNode := result["name"]

		properties = append(properties, generictypes.Property{
			Name: nameNode.Content(),
			Node: nameNode,
		})
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
	detector.classQuery.Close()
	detector.assignmentQuery.Close()
	detector.fieldAccessQuery.Close()
}
