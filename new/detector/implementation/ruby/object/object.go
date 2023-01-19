package object

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type Data struct {
	Name       string
	Properties []*types.Detection
}

type objectDetector struct {
	// Gathering properties
	hashPairQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	parentPairQuery *tree.Query
	// class
	classNameQuery *tree.Query
	// properties
	propertiesQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { first_name: ..., ... }
	hashPairQuery, err := lang.CompileQuery(`(hash (pair) @pair) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling hash pair query: %s", err)
	}

	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment left: (identifier) @left right: (_) @right) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling assignment query: %s", err)
	}
	// { user: <object> }
	parentPairQuery, err := lang.CompileQuery(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling parent pair query: %s", err)
	}

	// class User
	// end
	classNameQuery, err := lang.CompileQuery(`(class name: (constant) @name) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class name query: %s", err)
	}

	// user.name
	propertiesQuery, err := lang.CompileQuery(`(call receiver: (_) @receiver method: (identifier) @method) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class name query: %s", err)
	}

	return &objectDetector{
		hashPairQuery:   hashPairQuery,
		assignmentQuery: assignmentQuery,
		parentPairQuery: parentPairQuery,
		classNameQuery:  classNameQuery,
		propertiesQuery: propertiesQuery,
	}, nil
}

func (detector *objectDetector) Name() string {
	return "object"
}

func (detector *objectDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	detections, err := detector.getHash(node, evaluator)
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

	detections, err = detector.getProperties(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.nameParentPairObject(node, evaluator)
}

func (detector *objectDetector) getProperties(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	results, err := detector.propertiesQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	for _, result := range results {

		if result["receiver"].Type() == "identifier" {
			return []*types.Detection{{
				MatchNode:   node,
				ContextNode: node,
				Data: Data{
					Name: result["receiver"].Content(),
					Properties: []*types.Detection{
						{
							MatchNode: result["root"],
							Data: Data{
								Name: result["method"].Content(),
							},
						},
					},
				},
			}}, nil
		}

		if result["receiver"].Type() == "call" {
			childMethodNode := result["receiver"].ChildByFieldName("method")

			return []*types.Detection{{
				MatchNode:   node,
				ContextNode: node,
				Data: Data{
					Name: childMethodNode.Content(),
					Properties: []*types.Detection{
						{
							MatchNode: result["root"],
							Data: Data{
								Name: result["method"].Content(),
							},
						},
					},
				},
			}}, nil
		}
	}

	return nil, nil
}

func (detector *objectDetector) getHash(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	results, err := detector.hashPairQuery.MatchAt(node)
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
		Data:      Data{Properties: properties},
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
		objectData := object.Data.(Data)

		if objectData.Name == "" {
			detections = append(detections, &types.Detection{
				MatchNode: node,
				Data: Data{
					Name:       result["left"].Content(),
					Properties: objectData.Properties,
				},
			})
		} else { // FIXME: should we remove this case?
			detections = append(detections, &types.Detection{
				MatchNode: node,
				Data: Data{
					Name: result["left"].Content(),
					Properties: []*types.Detection{{
						MatchNode: object.MatchNode,
						Data: Data{
							Name: objectData.Name,
						},
					}},
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

	data := Data{
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
		objectData := object.Data.(Data)

		detections = append(detections, &types.Detection{
			MatchNode: node,
			Data: Data{
				Name:       result["key"].Content(),
				Properties: objectData.Properties,
			},
		})
	}

	return detections, nil
}

func (detector *objectDetector) Close() {
	detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.parentPairQuery.Close()
}
