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

func (detector *objectsDetector) Type() string {
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
		Data:      map[string]interface{}{"properties": properties},
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
		objectData := object.Data

		if objectData["name"] == "" {
			detections = append(detections, &detectiontypes.Detection{
				MatchNode: node,
				Data: map[string]interface{}{
					"name":       result["left"].Content(),
					"properties": objectData["properties"],
				},
			})
		} else {
			detections = append(detections, &detectiontypes.Detection{
				MatchNode: node,
				Data: map[string]interface{}{
					"name": result["left"].Content(),
					"properties": []*detectiontypes.Detection{{
						MatchNode: object.MatchNode,
						Data: map[string]interface{}{
							"name": objectData["name"],
						},
					}},
				},
			})
		}
	}

	return detections, nil
}

func (detector *objectsDetector) nameParentPairObject(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	result, err := detector.parentPairQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	objects, err := evaluator.NodeDetections(result["value"], "objects")
	if err != nil {
		return nil, err
	}

	var detections []*detectiontypes.Detection
	for _, object := range objects {
		objectData := object.Data

		detections = append(detections, &detectiontypes.Detection{
			MatchNode: node,
			Data: map[string]interface{}{
				"name":       result["key"].Content(),
				"properties": objectData["properties"],
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
