package tree

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	content      []byte
	types        []string
	nodes        []Node
	rootNode     *Node
	children     []*Node
	sitterToNode map[*sitter.Node]*Node
}

type Node struct {
	tree *Tree
	ID,
	TypeID int
	ContentStart,
	ContentEnd Position
	children []*Node
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

func (node *Node) Content() string {
	return string(node.tree.content[node.ContentStart.Byte:node.ContentEnd.Byte])
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

// FIXME: remove this
func (node *Node) EachContentPart(onText func(text string) error, onChild func(child *Node) error) error {
	start := node.ContentStart.Byte
	end := start

	emit := func() error {
		if end <= start {
			return nil
		}

		return onText(string(node.tree.content[start:end]))
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
