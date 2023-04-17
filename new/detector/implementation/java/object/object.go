package object

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type objectDetector struct {
	types.DetectorBase
	// Base
	classNameQuery     *tree.Query
	classPropertyQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	// Projection
	fieldAcessQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment_expression left: (identifier) @name right: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}

	// class User {
	//
	// }
	classNameQuery, err := lang.CompileQuery(`
		(class_declaration name: (identifier) @class_name
			(class_body) @class_body
		) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class name query: %s", err)
	}

	// class User {
	//    String name
	//	  String getLevel(){}
	// }
	classPropertyQuery, err := lang.CompileQuery(
		`(class_body
			[
			  (field_declaration (variable_declarator name: (identifier) @name))
			  (method_declaration name: (identifier) @name)
			]
		)@root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class property query: %s", err)
	}

	// user.name
	fieldAcessQuery, err := lang.CompileQuery(`(field_access object: (_) @object field: (identifier) @field) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling call query: %s", err)
	}

	return &objectDetector{
		assignmentQuery:    assignmentQuery,
		classNameQuery:     classNameQuery,
		classPropertyQuery: classPropertyQuery,
		fieldAcessQuery:    fieldAcessQuery,
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
	log.Debug().Msgf("node is %s", node.Debug())

	detections, err := detector.getAssignment(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClassName(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClassProperties(node, evaluator)
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

func (detector *objectDetector) getClassName(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	results, err := detector.classNameQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	result := results[0]

	className := result["class_name"].Content()

	properties, err := generic.GetNonVirtualObjects(evaluator, result["right"])
	if err != nil {
		return nil, err
	}

	var objects []interface{}
	for _, object := range properties {
		objects = append(objects, generictypes.Object{
			IsVirtual: false,
			Properties: []generictypes.Property{{
				Name:   className,
				Node:   node,
				Object: object,
			}},
		})
	}

	return objects, nil
}

func (detector *objectDetector) getClassProperties(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	results, err := detector.classPropertyQuery.MatchAt(node)
	if len(results) == 0 || err != nil {
		return nil, err
	}

	var objects []interface{}

	for _, result := range results {
		objects = append(objects, generictypes.Object{
			IsVirtual: false,
			Properties: []generictypes.Property{{
				Name: result["name"].Content(),
				Node: node,
				Object: &types.Detection{
					DetectorType: "object",
					MatchNode:    result["name"],
				},
			}},
		})
	}

	return objects, nil
}

func (detector *objectDetector) Close() {
	detector.classNameQuery.Close()
	detector.classPropertyQuery.Close()
	detector.assignmentQuery.Close()
	detector.fieldAcessQuery.Close()
}
