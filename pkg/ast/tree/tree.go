package tree

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	types        []string
	nodes        []Node
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
}

type Position struct {
	Byte,
	Line,
	Column int
}

func (node *Node) Type() string {
	return node.tree.types[node.TypeID]
}

func (node *Node) Child(index int) *Node {
	return node.children[index]
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
