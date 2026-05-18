package language

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
)

type PatternBase struct{}

func (*PatternBase) IsLeaf(node *tree.Node) bool {
	return false
}

func (*PatternBase) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*PatternBase) IsRoot(node *tree.Node) bool {
	return true
}

func (*PatternBase) ShouldSkipNode(node *tree.Node) bool {
	return false
}

func (*PatternBase) IsContainer(node *tree.Node) bool {
	return false
}

func (*PatternBase) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	return dummyValue
}

func (*PatternBase) AnonymousParentTypes() []string {
	return nil
}

func (*PatternBase) AdjustInput(input string) string {
	return input
}

func (*PatternBase) FixupMissing(node *tree.Node) string {
	return ""
}

func (*PatternBase) IsVariable(node *tree.Node, dummyValue string) bool {
	return strings.EqualFold(node.Content(), dummyValue)
}

func (*PatternBase) FieldNameFor(sitterLanguage *sitter.Language, node *tree.Node) string {
	return FieldNameForWorkaround(sitterLanguage, node)
}

// FieldNameForWorkaround resolves a node's field name by iterating the
// language's field names and asking the parent for the child bound to each.
// This is a workaround until https://github.com/tree-sitter/tree-sitter/pull/2104
// is released. It misidentifies children whose field name is shared with a
// sibling — languages that have such nodes should use FieldNameForChild.
func FieldNameForWorkaround(sitterLanguage *sitter.Language, node *tree.Node) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	for i := 1; ; i++ {
		name := sitterLanguage.FieldName(i)
		if name == "" {
			return ""
		}

		if parent.ChildByFieldName(name) == node {
			return name
		}
	}
}

// FieldNameForChild resolves a node's field name using tree-sitter's canonical
// FieldNameForChild API, which always matches what scan-time queries match
// against. Use this for languages where multiple children may share a field
// name.
func FieldNameForChild(node *tree.Node) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	sitterParent := parent.SitterNode()
	sitterNode := node.SitterNode()
	for i := 0; i < int(sitterParent.ChildCount()); i++ {
		if sitterParent.Child(i).Equal(sitterNode) {
			return sitterParent.FieldNameForChild(i)
		}
	}
	return ""
}
