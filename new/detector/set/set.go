package set

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
)

type set struct {
	detectors map[string]types.Detector
}

func New(detectors []types.Detector) (types.DetectorSet, error) {
	detectorMap := make(map[string]types.Detector)

	for _, detector := range detectors {
		name := detector.Name()

		if _, existing := detectorMap[name]; existing {
			return nil, fmt.Errorf("duplicate detector '%s'", name)
		}

		detectorMap[name] = detector
	}

	return &set{
		detectors: detectorMap,
	}, nil
}

func (set *set) NestedDetections(detectorType string) (bool, error) {
	detector, err := set.lookupDetector(detectorType)
	if err != nil {
		return false, err
	}

	return detector.NestedDetections(), nil
}

func (set *set) DetectAt(
	rootNode, node *tree.Node,
	detectorType string,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	detector, err := set.lookupDetector(detectorType)
	if err != nil {
		return nil, err
	}

	detectionsData, err := detector.DetectAt(rootNode, node, evaluator)
	if err != nil {
		return nil, err
	}

	detections := make([]*types.Detection, len(detectionsData))
	for i, data := range detectionsData {
		detections[i] = &types.Detection{
			DetectorType: detectorType,
			MatchNode:    node,
			Data:         data,
		}
	}

	return detections, nil
}

func (set *set) lookupDetector(detectorType string) (types.Detector, error) {
	detector, ok := set.detectors[detectorType]
	if !ok {
		return nil, fmt.Errorf("detector type '%s' not registered", detectorType)
	}

	return detector, nil
}
