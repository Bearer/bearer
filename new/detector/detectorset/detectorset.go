package detectorset

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
)

type detectorSet struct {
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

	return &detectorSet{
		detectors: detectorMap,
	}, nil
}

func (detectorSet *detectorSet) DetectAt(
	node *tree.Node,
	detectorType string,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	detector, ok := detectorSet.detectors[detectorType]
	if !ok {
		return nil, fmt.Errorf("detector type '%s' not registered", detectorType)
	}

	return detector.DetectAt(node, evaluator)
}
