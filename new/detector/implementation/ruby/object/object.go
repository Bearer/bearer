package object

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	soufflequery "github.com/bearer/bearer/pkg/souffle/query"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type objectDetector struct {
	types.DetectorBase
	// Gathering properties
	hashPairQuery *tree.Query
	// Naming
	assignmentQuery *tree.Query
	parentPairQuery *tree.Query
	// class
	classNameQuery *tree.Query
	// properties
	callsQuery            *tree.Query
	elementReferenceQuery *tree.Query
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
	parentPairQuery, err := lang.CompileQuery(`(pair key: (_) @key value: (_) @value) @root`)
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
	callsQuery, err := lang.CompileQuery(`(call receiver: (_) @receiver method: (identifier) @method) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling call query: %s", err)
	}

	// user[:name]
	elementReferenceQuery, err := lang.CompileQuery(`(element_reference object: (_) @object (_) @key) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling element reference query %s", err)
	}

	return &objectDetector{
		hashPairQuery:         hashPairQuery,
		assignmentQuery:       assignmentQuery,
		parentPairQuery:       parentPairQuery,
		classNameQuery:        classNameQuery,
		callsQuery:            callsQuery,
		elementReferenceQuery: elementReferenceQuery,
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
	queryContext *soufflequery.QueryContext,
) ([]interface{}, error) {
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

func (detector *objectDetector) getHash(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	results, err := detector.hashPairQuery.MatchAt(node)
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

	var detectionsData []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		if objectData.Name == "" {
			detectionsData = append(detectionsData, generictypes.Object{
				Name:       result["left"].Content(),
				Properties: objectData.Properties,
			})
		}
	}

	return detectionsData, nil
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

	key := common.GetLiteralKey(result["key"])
	if key == "" {
		return nil, nil
	}

	objects, err := evaluator.ForNode(result["value"], "object", true)
	if err != nil {
		return nil, err
	}

	var detectionsData []interface{}
	for _, object := range objects {
		objectData := object.Data.(generictypes.Object)

		detectionsData = append(detectionsData, generictypes.Object{
			Name:       key,
			Properties: objectData.Properties,
		})
	}

	return detectionsData, nil
}

func (detector *objectDetector) Close() {
	detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.parentPairQuery.Close()
	detector.classNameQuery.Close()
	detector.callsQuery.Close()
	detector.elementReferenceQuery.Close()
}
