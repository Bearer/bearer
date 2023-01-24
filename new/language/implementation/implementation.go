package implementation

import (
	sitter "github.com/smacker/go-tree-sitter"

	patternquerytypes "github.com/bearer/curio/new/language/patternquery/types"
	"github.com/bearer/curio/new/language/tree"
)

type Implementation interface {
	SitterLanguage() *sitter.Language
	AnalyzeFlow(rootNode *tree.Node) error
	ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error)
	FindPatternMatchNode(input []byte) [][]int
	FindPatternUnanchoredPoints(input []byte) [][]int
	AnonymousPatternNodeParentTypes() []string
	PatternIsAnchored(node *tree.Node) bool
}
