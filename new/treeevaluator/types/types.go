package types

import (
	"github.com/bearer/curio/new/language"
	"github.com/open-policy-agent/opa/ast"
)

type Evaluator interface {
	TreeDetections(rootNode *language.Node, detectorType string) (*ast.Array, error)
	NodeDetections(node *language.Node, detectorType string) (*ast.Array, error)
	TreeHasDetection(rootNode *language.Node, detectorType string) (bool, error)
	NodeHasDetection(node *language.Node, detectorType string) (bool, error)
}
