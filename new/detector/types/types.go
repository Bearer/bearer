package types

import (
	"context"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/file"
)

type EvaluationState interface {
	Evaluate(
		rootNode *tree.Node,
		detectorType,
		sanitizerDetectorType string,
		scope settings.RuleReferenceScope,
		followFlow bool,
	) ([]*detection.Detection, error)
	FileName() string
}

type DetectorSet interface {
	NestedDetections(detectorType string) (bool, error)
	DetectAt(
		node *tree.Node,
		detectorType string,
		evaluationState EvaluationState,
	) ([]*detection.Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, evaluationState EvaluationState) ([]interface{}, error)
	NestedDetections() bool
	Close()
}

type DetectorBase struct{}

func (*DetectorBase) NestedDetections() bool {
	return true
}

type Composition interface {
	DetectFromFile(ctx context.Context, file *file.FileInfo) ([]*detection.Detection, error)
	DetectFromFileWithTypes(
		ctx context.Context,
		file *file.FileInfo,
		detectorTypes, sharedDetectorTypes []string,
	) ([]*detection.Detection, error)
	Close()
}
