package detectorexecutor

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/detectorexecutor/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type detectorExecutor struct {
	lang          languagetypes.Language
	detectorStack map[tree.NodeID][]string
	detectors     map[string]detector.Detector
}

func New(lang languagetypes.Language, detectors []detector.Detector) (types.Executor, error) {
	detectorMap, err := makeDetectorMap(detectors)
	if err != nil {
		return nil, err
	}

	return &detectorExecutor{
		lang:          lang,
		detectorStack: make(map[tree.NodeID][]string),
		detectors:     detectorMap,
	}, nil
}

func makeDetectorMap(detectors []detector.Detector) (map[string]detector.Detector, error) {
	result := make(map[string]detector.Detector)

	for _, detector := range detectors {
		name := detector.Name()

		if _, existing := result[name]; existing {
			return nil, fmt.Errorf("duplicate detector '%s'", name)
		}

		result[name] = detector
	}

	return result, nil
}

func (executor *detectorExecutor) DetectAt(
	node *tree.Node,
	detectorType string,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	detector, ok := executor.detectors[detectorType]
	if !ok {
		return nil, fmt.Errorf("detector type '%s' not registered", detectorType)
	}

	nodeID := node.ID()
	executingDetectors := executor.detectorStack[nodeID]

	if slices.Contains(executingDetectors, detectorType) {
		return nil, fmt.Errorf(
			"cycle found in detector usage: [%s > %s]\nnode type: %s, content:\n%s",
			strings.Join(executingDetectors, " > "),
			detectorType,
			node.Type(),
			node.Content(),
		)
	}

	executor.detectorStack[nodeID] = append(executor.detectorStack[nodeID], detectorType)

	detections, err := detector.DetectAt(node, evaluator)
	if err != nil {
		return nil, err
	}

	if len(executor.detectorStack[nodeID]) == 1 {
		delete(executor.detectorStack, nodeID)
	} else {
		executor.detectorStack[nodeID] = executor.detectorStack[nodeID][:len(executor.detectorStack[nodeID])-1]
	}

	return detections, nil
}
