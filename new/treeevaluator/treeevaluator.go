package treeevaluator

import (
	detectortypes "github.com/bearer/curio/new/detection/types"
	detectorexecutortypes "github.com/bearer/curio/new/detectorexecutor/types"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/treeevaluator/types"
)

type treeEvaluator struct {
	executor       detectorexecutortypes.Executor
	detectionCache map[language.NodeID]map[string]*detectortypes.Detection
}

func New(executor detectorexecutortypes.Executor, tree *language.Tree) types.Evaluator {
	detectionCache := make(map[language.NodeID]map[string]*detectortypes.Detection)

	return &treeEvaluator{
		executor:       executor,
		detectionCache: detectionCache,
	}
}

func (evaluator *treeEvaluator) TreeDetections(rootNode *language.Node, detectorType string) ([]*detectortypes.Detection, error) {
	var result []*detectortypes.Detection

	if err := rootNode.Walk(func(node *language.Node) error {
		detection, err := evaluator.NodeDetection(node, detectorType)
		if err != nil {
			return err
		}

		if detection != nil {
			result = append(result, detection)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (evaluator *treeEvaluator) NodeDetection(
	node *language.Node,
	detectorType string,
) (*detectortypes.Detection, error) {
	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		err := evaluator.detectAtNode(node, detectorType)
		if err != nil {
			return nil, err
		}

		nodeDetections = evaluator.detectionCache[node.ID()]
	}

	if detection, ok := nodeDetections[detectorType]; ok {
		return detection, nil
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
	detection, err := evaluator.NodeDetection(node, detectorType)
	if err != nil {
		return false, err
	}

	return detection != nil, nil
}

func (evaluator *treeEvaluator) detectAtNode(node *language.Node, detectorType string) error {
	detection, err := evaluator.executor.DetectAt(node, detectorType, evaluator)
	if err != nil {
		return err
	}

	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		nodeDetections = make(map[string]*detectortypes.Detection)
		evaluator.detectionCache[node.ID()] = nodeDetections
	}

	nodeDetections[detectorType] = detection

	return nil
}
