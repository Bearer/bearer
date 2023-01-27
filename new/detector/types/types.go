package types

import (
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/pkg/util/file"
)

type Detection struct {
	MatchNode   *tree.Node
	ContextNode *tree.Node
	Data        interface{}
}

type Evaluator interface {
	ForTree(rootNode *tree.Node, detectorType string) ([]*Detection, error)
	ForNode(node *tree.Node, detectorType string) ([]*Detection, error)
	TreeHas(rootNode *tree.Node, detectorType string) (bool, error)
	NodeHas(node *tree.Node, detectorType string) (bool, error)
	FileName() string
}

type DetectorSet interface {
	DetectAt(
		node *tree.Node,
		detectorType string,
		evaluator Evaluator,
	) ([]*Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, evaluator Evaluator) ([]*Detection, error)
	Close()
}

type CompositionDetection struct {
	RuleName   string
	File       *file.FileInfo
	Detections []*Detection
}

type Composition interface {
	DetectFromFile(file *file.FileInfo) ([]*CompositionDetection, error)
	Close()
}
