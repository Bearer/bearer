package object

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type objectDetector struct {
	types.DetectorBase
	// Gathering properties
	objectPairQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	parentPairQuery *tree.Query
	// Variables
	variableDeclarationQuery  *tree.Query
	objectDeconstructionQuery *tree.Query
	// class
	classNameQuery   *tree.Query
	constructorQuery *tree.Query
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

	// const { user } = <object>
	// let { user } = <object>
	// var { user } = <object>
	objectDeconstructionQuery, err := lang.CompileQuery(`(variable_declarator name:(object_pattern (shorthand_property_identifier_pattern) @match ) value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling object deconstruction query: %s", err)
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

	// new User()
	constructorQuery, err := lang.CompileQuery(`(new_expression constructor: (identifier) @name) @root`)
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
		objectPairQuery:           objectPairQuery,
		assignmentQuery:           assignmentQuery,
		variableDeclarationQuery:  variableDeclarationQuery,
		objectDeconstructionQuery: objectDeconstructionQuery,
		parentPairQuery:           parentPairQuery,
		classNameQuery:            classNameQuery,
		constructorQuery:          constructorQuery,
		memberExpressionQuery:     memberExpressionQuery,
		subscriptExpressionQuery:  subscriptExpressionQuery,
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

	detections, err = detector.getObjectDeconstruction(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getConstructor(node, evaluator)
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

func (detector *objectDetector) getClass(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	result, err := detector.classNameQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	data := generictypes.Object{
		Name:       result["name"].Content(),
		Properties: []*types.Detection{},
	}

	body := node.ChildByFieldName("body")

	for i := 0; i < body.ChildCount(); i++ {
		detections, err := evaluator.ForNode(body.Child(i), "property", true)
		if err != nil {
			return nil, err
		}
		data.Properties = append(data.Properties, detections...)
	}

	return []interface{}{data}, nil
}

func (detector *objectDetector) getConstructor(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	result, err := detector.constructorQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	data := generictypes.Object{
		Name: result["name"].Content(),
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

		objectName := result["key"].Content()
		objectNameNode := result["key"]
		if objectNameNode.Type() == "string" {
			objectName = stringutil.StripQuotes(objectName)
		}

		detections = append(detections, generictypes.Object{
			Name:       objectName,
			Properties: objectData.Properties,
		},
		)
	}

	return detections, nil
}

func (detector *objectDetector) Close() {
	detector.objectPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.variableDeclarationQuery.Close()
	detector.objectDeconstructionQuery.Close()
	detector.parentPairQuery.Close()
	detector.classNameQuery.Close()
	detector.memberExpressionQuery.Close()
	detector.subscriptExpressionQuery.Close()
}
