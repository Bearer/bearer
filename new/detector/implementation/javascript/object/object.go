package object

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type objectDetector struct {
	// Gathering properties
	objectPairQuery *tree.Query
	// Naming
	assignmentQuery          *tree.Query
	variableDeclarationQuery *tree.Query
	parentPairQuery          *tree.Query
	// class
	classNameQuery *tree.Query
	// properties
	memberExpressionQuery    *tree.Query
	subscriptExpressionQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { first_name: ..., ... }
	objectPairQuery, err := lang.CompileQuery(`(object (pair) @pair) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling object pair query: %s", err)
	}

	// const user = <object>
	// var user = <object>
	// let user = <object>
	variableDeclarationQuery, err := lang.CompileQuery(`(variable_declarator name: (identifier) @name value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}

	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment_expression left: (identifier) @left right: (_) @right) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}

	// { user: <object> }
	parentPairQuery, err := lang.CompileQuery(`(pair key: (_) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling parent pair query: %s", err)
	}

	// class User
	// end
	classNameQuery, err := lang.CompileQuery(`(class_declaration name: (identifier) @name) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class name query: %s", err)
	}

	// user.name
	memberExpressionQuery, err := lang.CompileQuery(`(member_expression object: (_) @object property: (property_identifier) @property) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling member expression query: %s", err)
	}

	// user[:name]
	subscriptExpressionQuery, err := lang.CompileQuery(`(subscript_expression object: (_) @object index: (string) @index ) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling subscript expression query %s", err)
	}

	return &objectDetector{
		objectPairQuery:          objectPairQuery,
		assignmentQuery:          assignmentQuery,
		variableDeclarationQuery: variableDeclarationQuery,
		parentPairQuery:          parentPairQuery,
		classNameQuery:           classNameQuery,
		memberExpressionQuery:    memberExpressionQuery,
		subscriptExpressionQuery: subscriptExpressionQuery,
	}, nil
}

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	detections, err := detector.getobject(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssigment(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getVariableDeclaration(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getProperties(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.nameParentPairObject(node, evaluator)
}

func (detector *objectDetector) getobject(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	results, err := detector.objectPairQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var properties []*types.Detection
	for _, result := range results {
		nodeProperties, err := evaluator.ForNode(result["pair"], "property", false)
		if err != nil {
			return nil, err
		}

		properties = append(properties, nodeProperties...)
	}

	return []interface{}{generictypes.Object{Properties: properties}}, nil
}

func (detector *objectDetector) getAssigment(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["right"], "object", true)
	if err != nil {
		return nil, err
	}

	var detections []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		if objectData.Name == "" {
			detections = append(detections, generictypes.Object{
				Name:       result["left"].Content(),
				Properties: objectData.Properties,
			})
		}
	}

	return detections, nil
}

func (detector *objectDetector) getVariableDeclaration(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.variableDeclarationQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["value"], "object", true)
	if err != nil {
		return nil, err
	}

	var detections []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		if objectData.Name == "" {
			detections = append(detections, generictypes.Object{
				Name:       result["name"].Content(),
				Properties: objectData.Properties,
			},
			)
		}
	}

	return detections, nil
}

func (detector *objectDetector) getClass(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	result, err := detector.classNameQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	data := generictypes.Object{
		Name:       result["name"].Content(),
		Properties: []*types.Detection{},
	}

	for i := 0; i < node.ChildCount(); i++ {
		detections, err := evaluator.ForNode(node.Child(i), "property", true)
		if err != nil {
			return nil, err
		}
		data.Properties = append(data.Properties, detections...)
	}

	return []interface{}{data}, nil
}

func (detector *objectDetector) nameParentPairObject(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	result, err := detector.parentPairQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["value"], "object", true)
	if err != nil {
		return nil, err
	}

	var detections []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		detections = append(detections, generictypes.Object{
			Name:       result["key"].Content(),
			Properties: objectData.Properties,
		},
		)
	}

	return detections, nil
}

func (detector *objectDetector) Close() {
	// detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.parentPairQuery.Close()
}
