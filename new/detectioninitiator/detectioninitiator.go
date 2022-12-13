package detectioninitiator

import (
	detectortypes "github.com/bearer/curio/new/detection/types"
	executortypes "github.com/bearer/curio/new/detectionexecutor/types"
	"github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/parser"
)

type detectionInitiator struct {
	executor       executortypes.DetectionExecutor
	detectionCache map[parser.NodeID]map[string]*detectortypes.Detection
}

func New(executor executortypes.DetectionExecutor, tree *parser.Tree) types.TreeDetectionInitiator {
	detectionCache := make(map[parser.NodeID]map[string]*detectortypes.Detection)

	return &detectionInitiator{
		executor:       executor,
		detectionCache: detectionCache,
	}
}

func (initiator *detectionInitiator) TreeDetections(rootNode *parser.Node, detectorType string) ([]*detectortypes.Detection, error) {
	var result []*detectortypes.Detection

	if err := rootNode.Walk(func(node *parser.Node) error {
		detection, err := initiator.NodeDetection(node, detectorType)
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

func (initiator *detectionInitiator) NodeDetection(
	node *parser.Node,
	detectorType string,
) (*detectortypes.Detection, error) {
	nodeDetections, ok := initiator.detectionCache[node.ID()]
	if !ok {
		err := initiator.detectAtNode(node, detectorType)
		if err != nil {
			return nil, err
		}

		nodeDetections = initiator.detectionCache[node.ID()]
	}

	if detection, ok := nodeDetections[detectorType]; ok {
		return detection, nil
	}

	initiator.detectAtNode(node, detectorType)
	return nodeDetections[detectorType], nil
}

func (initiator *detectionInitiator) TreeHasDetection(rootNode *parser.Node, detectorType string) (bool, error) {
	hasDetection := false

	if err := rootNode.Walk(func(node *parser.Node) error {
		var err error
		hasDetection, err = initiator.NodeHasDetection(node, detectorType)
		if err != nil {
			return err
		}

		if hasDetection {
			return parser.ErrTerminateWalk
		}

		return nil
	}); err != nil {
		return false, err
	}

	return hasDetection, nil
}

func (initiator *detectionInitiator) NodeHasDetection(node *parser.Node, detectorType string) (bool, error) {
	detection, err := initiator.NodeDetection(node, detectorType)
	if err != nil {
		return false, err
	}

	return detection != nil, nil
}

func (initiator *detectionInitiator) detectAtNode(node *parser.Node, detectorType string) error {
	detection, err := initiator.executor.DetectAt(node, detectorType, initiator)
	if err != nil {
		return err
	}

	nodeDetections, ok := initiator.detectionCache[node.ID()]
	if !ok {
		nodeDetections = make(map[string]*detectortypes.Detection)
		initiator.detectionCache[node.ID()] = nodeDetections
	}

	nodeDetections[detectorType] = detection

	return nil
}
