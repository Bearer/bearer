package types

import (
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/file"
)

type Detection struct {
	DetectorType string
	MatchNode    *tree.Node
	Data         interface{}
}

type Evaluator interface {
	Evaluate(
		rootNode *tree.Node,
		detectorType,
		sanitizerDetectorType string,
		scope settings.RuleReferenceScope,
		followFlow bool,
	) ([]*Detection, error)
	FileName() string
}

type DetectorSet interface {
	NestedDetections(detectorType string) (bool, error)
	DetectAt(
		node *tree.Node,
		detectorType string,
		ruleReferenceType settings.RuleReferenceScope,
		evaluator Evaluator,
	) ([]*Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(
		node *tree.Node,
		ruleReferenceType settings.RuleReferenceScope,
		evaluator Evaluator,
	) ([]interface{}, error)
	NestedDetections() bool
	Close()
}

type DetectorBase struct{}

func (*DetectorBase) NestedDetections() bool {
	return true
}

type Composition interface {
	DetectFromFile(file *file.FileInfo) ([]*Detection, error)
	DetectFromFileWithTypes(file *file.FileInfo, detectorTypes []string) ([]*Detection, error)
	Close()
}
