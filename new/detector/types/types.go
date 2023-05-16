package types

import (
	"crypto/sha256"

	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/file"
)

type Detection struct {
	DetectorType string
	MatchNode    *tree.Node
	Data         interface{}
}

type EvaluationContext []EvaluationScope

func (context EvaluationContext) Cursor() *tree.Node {
	return context[len(context)-1].Cursor
}

func (context EvaluationContext) CursorKey() string {
	hash := sha256.New()

	for _, scope := range context {
		hash.Write([]byte(scope.Cursor.IDString()))
		hash.Write([]byte(","))
	}

	return string(hash.Sum(nil))
}

type EvaluationScope struct {
	Root   *tree.Node
	Cursor *tree.Node
}

type Evaluator interface {
	ForTree(rootNode *tree.Node, detectorType string, followFlow bool) ([]*Detection, error)
	ForNode(node *tree.Node, detectorType string, followFlow bool) ([]*Detection, error)
	TreeHas(rootNode *tree.Node, detectorType string) (bool, error)
	NodeHas(node *tree.Node, detectorType string) (bool, error)
	FileName() string
}

type DetectorSet interface {
	NestedDetections(detectorType string) (bool, error)
	DetectAt(
		evaluationContext EvaluationContext,
		detectorType string,
		evaluator Evaluator,
	) ([]*Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(evaluationContext EvaluationContext, evaluator Evaluator) ([]interface{}, error)
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
