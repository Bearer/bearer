package implementation

import (
	sitter "github.com/smacker/go-tree-sitter"

	patternquerybuilder "github.com/bearer/curio/new/language/patternquery/builder"
	"github.com/bearer/curio/new/language/tree"
)

type Implementation interface {
	SitterLanguage() *sitter.Language
	AnalyzeFlow(rootNode *tree.Node)
	ExtractPatternVariables(input string) (string, []patternquerybuilder.Variable, error)
	AnonymousPatternNodeParentTypes() []string
}
