package object

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	"github.com/rs/zerolog/log"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type objectDetector struct {
	// Gathering properties
	objectPairQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	parentPairQuery *tree.Query
	// class
	classNameQuery *tree.Query
	// properties
	// callsQuery            *tree.Query
	// elementReferenceQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	tree, err := lang.Parse(`{"user": "name"}`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile tree %s", err)
	}
	log.Debug().Msgf(tree.RootNode().Debug())

	// { first_name: ..., ... }
	objectPairQuery, err := lang.CompileQuery(`(object (pair) @pair) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling object pair query: %s", err)
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

	// // user.name
	// callsQuery, err := lang.CompileQuery(`(call receiver: (_) @receiver method: (identifier) @method) @root`)
	// if err != nil {
	// 	return nil, fmt.Errorf("error compiling call query: %s", err)
	// }

	// // user[:name]
	// elementReferenceQuery, err := lang.CompileQuery(`(element_reference object: (_) @object (simple_symbol) @simple_symbol) @root`)
	// if err != nil {
	// 	return nil, fmt.Errorf("error compiling element reference query %s", err)
	// }

	return &objectDetector{
		objectPairQuery: objectPairQuery,
		assignmentQuery: assignmentQuery,
		parentPairQuery: parentPairQuery,
		classNameQuery:  classNameQuery,
		// callsQuery:            callsQuery,
		// elementReferenceQuery: elementReferenceQuery,
	}, nil
}

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	detections, err := detector.getobject(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getAssigment(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.getClass(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.nameParentPairObject(node, evaluator)
}

func (detector *objectDetector) getobject(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	results, err := detector.objectPairQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var properties []*types.Detection
	for _, result := range results {
		nodeProperties, err := evaluator.ForNode(result["pair"], "property")
		if err != nil {
			return nil, err
		}

		properties = append(properties, nodeProperties...)
	}

	return []*types.Detection{{
		MatchNode: node,
		Data:      generictypes.Object{Properties: properties},
	}}, nil
}

func (detector *objectDetector) getAssigment(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["right"], "object")
	if err != nil {
		return nil, err
	}

	var detections []*types.Detection
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		if objectData.Name == "" {
			detections = append(detections, &types.Detection{
				MatchNode: node,
				Data: generictypes.Object{
					Name:       result["left"].Content(),
					Properties: objectData.Properties,
				},
			})
		}
	}

	return detections, nil
}

func (detector *objectDetector) getClass(node *tree.Node, evaluator types.Evaluator) ([]*types.Detection, error) {
	result, err := detector.classNameQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	data := generictypes.Object{
		Name:       result["name"].Content(),
		Properties: []*types.Detection{},
	}

	for i := 0; i < node.ChildCount(); i++ {
		detections, err := evaluator.ForNode(node.Child(i), "property")
		if err != nil {
			return nil, err
		}
		data.Properties = append(data.Properties, detections...)
	}

	return []*types.Detection{{
		MatchNode:   node,
		ContextNode: node,
		Data:        data,
	}}, nil
}

func (detector *objectDetector) nameParentPairObject(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	result, err := detector.parentPairQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.ForNode(result["value"], "object")
	if err != nil {
		return nil, err
	}

	var detections []*types.Detection
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		detections = append(detections, &types.Detection{
			MatchNode: node,
			Data: generictypes.Object{
				Name:       result["key"].Content(),
				Properties: objectData.Properties,
			},
		})
	}

	return detections, nil
}

func (detector *objectDetector) Close() {
	// detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.parentPairQuery.Close()
}
