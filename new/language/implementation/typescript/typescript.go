package typescript

import (
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"

	javascriptImplementation "github.com/bearer/bearer/new/language/implementation/javascript"

	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type typescriptImplementation struct {
	javascript implementation.Implementation
}

func Get() implementation.Implementation {
	return &typescriptImplementation{
		javascript: javascriptImplementation.Get(),
	}
}

func (implementation *typescriptImplementation) SitterLanguage() *sitter.Language {
	return typescript.GetLanguage()
}

func (implementation *typescriptImplementation) AnalyzeFlow(rootNode *tree.Node) error {
	return implementation.javascript.AnalyzeFlow(rootNode)
}

// TODO: See if anything needs to be added here
func (implementation *typescriptImplementation) ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error) {
	return implementation.javascript.ExtractPatternVariables(input)
}

// TODO: See if anything needs to be added here
func (implementation *typescriptImplementation) AnonymousPatternNodeParentTypes() []string {
	return implementation.javascript.AnonymousPatternNodeParentTypes()
}

// TODO: See if anything needs to be added here
func (implementation *typescriptImplementation) FindPatternMatchNode(input []byte) [][]int {
	return implementation.javascript.FindPatternMatchNode(input)
}

// TODO: See if anything needs to be added here
func (implementation *typescriptImplementation) FindPatternUnanchoredPoints(input []byte) [][]int {
	return implementation.javascript.FindPatternUnanchoredPoints(input)
}

func (implementation *typescriptImplementation) PatternMatchNodeContainerTypes() []string {
	return implementation.javascript.PatternMatchNodeContainerTypes()
}

func (implementation *typescriptImplementation) PatternLeafContentTypes() []string {
	types := []string{"string_fragment"}
	javascriptTypes := implementation.javascript.PatternLeafContentTypes()
	for _, v := range javascriptTypes {
		if v != "string" {
			types = append(types, v)
		}
	}

	return types
}

func (implementation *typescriptImplementation) PatternIsAnchored(node *tree.Node) (bool, bool) {
	return implementation.javascript.PatternIsAnchored(node)
}

func (implementation *typescriptImplementation) IsRootOfRuleQuery(node *tree.Node) bool {
	return implementation.javascript.IsRootOfRuleQuery(node)
}

func (implementation *typescriptImplementation) PatternNodeTypes(node *tree.Node) []string {
	return implementation.javascript.PatternNodeTypes(node)
}

func (implementation *typescriptImplementation) TranslatePatternContent(fromNodeType, toNodeType, content string) string {
	return implementation.javascript.TranslatePatternContent(fromNodeType, toNodeType, content)
}
