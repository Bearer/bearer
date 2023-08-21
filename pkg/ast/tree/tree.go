package tree

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	contentBytes []byte
	types        []string
	nodes        []Node
	rootNode     *Node
	sitterToNode map[*sitter.Node]*Node
}

type Node struct {
	tree *Tree
	ID,
	TypeID int
	ContentStart,
	ContentEnd Position
	parent *Node
	children,
	dataflowSources,
	aliasOf []*Node
	// FIXME: remove the need for this
	sitterNode *sitter.Node
	// FIXME: probably shouldn't be public
	ExecutingRules []string
}

type Position struct {
	Byte,
	Line,
	Column int
}

func (tree *Tree) ContentBytes() []byte {
	return tree.contentBytes
}

func (tree *Tree) RootNode() *Node {
	return tree.rootNode
}

func (tree *Tree) NodeFromSitter(sitterNode *sitter.Node) *Node {
	return tree.sitterToNode[sitterNode]
}

func (node *Node) SitterNode() *sitter.Node {
	return node.sitterNode
}

func (node *Node) Type() string {
	return node.tree.types[node.TypeID]
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) Content() string {
	return string(node.tree.contentBytes[node.ContentStart.Byte:node.ContentEnd.Byte])
}

func (node *Node) Debug(includeContent bool) string {
	content := ""
	if includeContent {
		content = ":\n" + node.Content()
	}

	return fmt.Sprintf(
		"%d:%d:%s%s",
		node.ContentStart.Line,
		node.ContentStart.Column,
		node.Type(),
		content,
	)
}

func (node *Node) Children() []*Node {
	return node.children
}

// FIXME: can we remove this?
func (node *Node) NamedChildren() []*Node {
	var namedChildren []*Node

	for _, child := range node.children {
		// FIXME: don't use the sitter node
		if child.sitterNode.IsNamed() {
			namedChildren = append(namedChildren, child)
		}
	}

	return namedChildren
}

func (node *Node) ChildByFieldName(name string) *Node {
	// FIXME: don't use the sitter node
	return node.tree.sitterToNode[node.sitterNode.ChildByFieldName(name)]
}

// FIXME: this is only used by tests
func (node *Node) NodeAndDescendentIDs() []int {
	var result []int

	next := []int{node.ID}
	for {
		if len(next) == 0 {
			break
		}

		result = append(result, next...)

		var newNext []int
		for _, id := range next {
			for _, child := range node.tree.nodes[id].children {
				newNext = append(newNext, child.ID)
			}
		}

		next = newNext
	}

	return result
}

func (node *Node) DataflowSources() []*Node {
	return node.dataflowSources
}

func (node *Node) AliasOf() []*Node {
	return node.aliasOf
}

func (node *Node) Dump() string {
	var s strings.Builder

	s.WriteString("(")         //nolint:errcheck
	s.WriteString(node.Type()) //nolint:errcheck

	if len(node.children) != 0 {
		for _, child := range node.children {
			s.WriteString(" ")          //nolint:errcheck
			s.WriteString(child.Dump()) //nolint:errcheck
		}
	}

	s.WriteString(")") //nolint:errcheck
	return s.String()
}

func (node *Node) Dump2() string {
	var s strings.Builder

	s.WriteString("(")         //nolint:errcheck
	s.WriteString(node.Type()) //nolint:errcheck

	if len(node.dataflowSources) != 0 {
		for _, child := range node.dataflowSources {
			s.WriteString(" ")           //nolint:errcheck
			s.WriteString(child.Dump2()) //nolint:errcheck
		}
	}

	s.WriteString(")") //nolint:errcheck
	return s.String()
}

// FIXME: remove this
func (node *Node) EachContentPart(onText func(text string) error, onChild func(child *Node) error) error {
	start := node.ContentStart.Byte
	end := start

	emit := func() error {
		if end <= start {
			return nil
		}

		return onText(string(node.tree.contentBytes[start:end]))
	}

	for _, child := range node.children {
		end = child.ContentStart.Byte

		if err := emit(); err != nil {
			return err
		}

		if child.SitterNode().IsNamed() {
			if err := onChild(child); err != nil {
				return err
			}
		}

		start = child.ContentEnd.Byte
		end = start
	}

	if err := emit(); err != nil {
		return err
	}

	return nil
}

// FIXME: maybe users of this could work iteratively?
func (node *Node) Walk(visit func(node *Node, visitChildren func() error) error) error {
	visitChildren := func() error {
		for _, child := range node.Children() {
			if err := child.Walk(visit); err != nil {
				return err
			}
		}

		return nil
	}

	return visit(node, visitChildren)
}
