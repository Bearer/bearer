package detectionexecutor

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detectionexecutor/types"
	initiatortypes "github.com/bearer/curio/new/detectioninitiator/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/parser"
)

type detectionExecutor struct {
	lang          *language.Language
	detectorStack []string
	detectors     map[string]detector.Detector
}

func New(lang *language.Language) types.DetectionExecutor {
	detectors := make(map[string]detector.Detector)

	return &detectionExecutor{
		lang:      lang,
		detectors: detectors,
	}
}

func (executor *detectionExecutor) RegisterDetector(detector detector.Detector) error {
	name := detector.Name()

	if _, alreadyRegistered := executor.detectors[name]; alreadyRegistered {
		return fmt.Errorf("detector '%s' is already registered", name)
	}

	executor.detectors[name] = detector

	return nil
}

func (executor *detectionExecutor) DetectAt(
	node *parser.Node,
	detectorType string,
	initiator initiatortypes.TreeDetectionInitiator,
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

	detection, err := detector.DetectAt(node, initiator)
	if err != nil {
		return nil, err
	}

	executor.detectorStack = executor.detectorStack[:len(executor.detectorStack)-1]

	return detection, nil
}
