package treeevaluator

import (
	detectortypes "github.com/bearer/curio/new/detection/types"
	detectorexecutortypes "github.com/bearer/curio/new/detectorexecutor/types"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/new/treeevaluator/types"
)

type treeEvaluator struct {
	lang           languagetypes.Language
	executor       detectorexecutortypes.Executor
	detectionCache map[language.NodeID]map[string][]*detectortypes.Detection
}

func New(
	lang languagetypes.Language,
	executor detectorexecutortypes.Executor,
	tree *language.Tree,
) types.Evaluator {
	detectionCache := make(map[language.NodeID]map[string][]*detectortypes.Detection)

	return &treeEvaluator{
		lang:           lang,
		executor:       executor,
		detectionCache: detectionCache,
	}
}

func (evaluator *treeEvaluator) TreeDetections(rootNode *language.Node, detectorType string) ([]*detectortypes.Detection, error) {
	var result []*detectortypes.Detection

	if err := rootNode.Walk(func(node *language.Node) error {
		detections, err := evaluator.NodeDetections(node, detectorType)
		if err != nil {
			return err
		}

		result = append(result, detections...)

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (evaluator *treeEvaluator) NodeDetections(
	node *language.Node,
	detectorType string,
) ([]*detectortypes.Detection, error) {
	var detections []*detectortypes.Detection

	for _, unifiedNode := range evaluator.lang.UnifiedNodesFor(node) {
		unifiedNodeDetections, err := evaluator.nonUnifiedNodeDetections(unifiedNode, detectorType)
		if err != nil {
			return nil, err
		}

		detections = append(detections, unifiedNodeDetections...)
	}

	return detections, nil
}

func (evaluator *treeEvaluator) nonUnifiedNodeDetections(
	node *language.Node,
	detectorType string,
) ([]*detectortypes.Detection, error) {
	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		err := evaluator.detectAtNode(node, detectorType)
		if err != nil {
			return nil, err
		}

		nodeDetections = evaluator.detectionCache[node.ID()]
	}

	if detections, ok := nodeDetections[detectorType]; ok {
		return detections, nil
	}

	evaluator.detectAtNode(node, detectorType)
	return nodeDetections[detectorType], nil
}

func (evaluator *treeEvaluator) TreeHasDetection(rootNode *language.Node, detectorType string) (bool, error) {
	hasDetection := false

	if err := rootNode.Walk(func(node *language.Node) error {
		var err error
		hasDetection, err = evaluator.NodeHasDetection(node, detectorType)
		if err != nil {
			return err
		}

		if hasDetection {
			return language.ErrTerminateWalk
		}

		return nil
	}); err != nil {
		return false, err
	}

	return hasDetection, nil
}

func (evaluator *treeEvaluator) NodeHasDetection(node *language.Node, detectorType string) (bool, error) {
	detections, err := evaluator.NodeDetections(node, detectorType)
	if err != nil {
		return false, err
	}

	return len(detections) != 0, nil
}

func (evaluator *treeEvaluator) detectAtNode(node *language.Node, detectorType string) error {
	detections, err := evaluator.executor.DetectAt(node, detectorType, evaluator)
	if err != nil {
		return err
	}

	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		nodeDetections = make(map[string][]*detectortypes.Detection)
		evaluator.detectionCache[node.ID()] = nodeDetections
	}

	nodeDetections[detectorType] = detections

	return nil
}
