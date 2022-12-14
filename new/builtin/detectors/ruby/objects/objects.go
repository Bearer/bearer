package objects

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	Name       string
	Properties []*detectiontypes.Detection
}

type objectsDetector struct {
	// Gathering properties
	hashPairQuery *language.Query
	// Naming
	assignmentQuery *language.Query
	parentPairQuery *language.Query
}

func New(lang languagetypes.Language) (detector.Detector, error) {
	// { first_name: ..., ... }
	hashPairQuery, err := lang.CompileQuery(`(hash (pair) @pair) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling hash pair query: %s", err)
	}

	// user = <object>
	assignmentQuery, err := lang.CompileQuery(`(assignment left: (identifier) @left right: (_) @right) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling hash pair query: %s", err)
	}
	// { user: <object> }
	parentPairQuery, err := lang.CompileQuery(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling hash pair query: %s", err)
	}

	return &objectsDetector{
		hashPairQuery:   hashPairQuery,
		assignmentQuery: assignmentQuery,
		parentPairQuery: parentPairQuery,
	}, nil
}

func (detector *objectsDetector) Name() string {
	return "objects"
}

func (detector *objectsDetector) DetectAt(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	detections, err := detector.gatherProperties(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	detections, err = detector.nameAssignedObject(node, evaluator)
	if len(detections) != 0 || err != nil {
		return detections, err
	}

	return detector.nameParentPairObject(node, evaluator)
}

func (detector *objectsDetector) gatherProperties(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	results, err := detector.hashPairQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var properties []*detectiontypes.Detection
	for _, result := range results {
		nodeProperties, err := evaluator.NodeDetections(result["pair"], "properties")
		if err != nil {
			return nil, err
		}

		properties = append(properties, nodeProperties...)
	}

	return []*detectiontypes.Detection{{
		MatchNode: node,
		Data:      Data{Properties: properties},
	}}, nil
}

func (detector *objectsDetector) nameAssignedObject(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.NodeDetections(result["right"], "objects")
	if err != nil {
		return nil, err
	}

	var detections []*detectiontypes.Detection
	for _, object := range objects {
		objectData := object.Data.(Data)

		detections = append(detections, &detectiontypes.Detection{
			MatchNode: node,
			Data: Data{
				Name:       result["left"].Content(),
				Properties: objectData.Properties,
			},
		})
	}

	return detections, nil
}

func (detector *objectsDetector) nameParentPairObject(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	result, err := detector.assignmentQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.NodeDetections(result["value"], "objects")
	if err != nil {
		return nil, err
	}

	var detections []*detectiontypes.Detection
	for _, object := range objects {
		objectData := object.Data.(Data)

		detections = append(detections, &detectiontypes.Detection{
			MatchNode: node,
			Data: Data{
				Name:       result["key"].Content(),
				Properties: objectData.Properties,
			},
		})
	}

	return detections, nil
}

func (detector *objectsDetector) Close() {
	detector.hashPairQuery.Close()
	detector.assignmentQuery.Close()
	detector.parentPairQuery.Close()
}
