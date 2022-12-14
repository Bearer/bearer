package detectorexecutor

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/detectorexecutor/types"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type detectorExecutor struct {
	lang          languagetypes.Language
	detectorStack []string
	detectors     map[string]detector.Detector
}

func New(lang languagetypes.Language, detectors []detector.Detector) types.Executor {

	return &detectorExecutor{
		lang:      lang,
		detectors: makeDetectorMap(detectors),
	}
}

func makeDetectorMap(detectors []detector.Detector) map[string]detector.Detector {
	result := make(map[string]detector.Detector)

	for _, detector := range detectors {
		result[detector.Name()] = detector
	}

	return result
}

func (executor *detectorExecutor) DetectAt(
	node *language.Node,
	detectorType string,
	evaluator treeevaluatortypes.Evaluator,
) (*detectiontypes.Detection, error) {
	detector, ok := executor.detectors[detectorType]
	if !ok {
		return nil, fmt.Errorf("detector type '%s' not registered", detectorType)
	}

	if slices.Contains(executor.detectorStack, detectorType) {
		return nil, fmt.Errorf(
			"cycle found in detector usage: [%s > %s]",
			strings.Join(executor.detectorStack, " > "),
			detectorType,
		)
	}

	executor.detectorStack = append(executor.detectorStack, detectorType)

	detection, err := detector.DetectAt(node, evaluator)
	if err != nil {
		return nil, err
	}

	executor.detectorStack = executor.detectorStack[:len(executor.detectorStack)-1]

	return detection, nil
}
