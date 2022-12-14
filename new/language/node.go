package language

import (
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
)

type Node struct {
	tree       *Tree
	sitterNode *sitter.Node
}

type NodeID *sitter.Node

var ErrTerminateWalk = errors.New("terminate walk")

func (node *Node) ID() NodeID {
	return node.sitterNode
}

func (node *Node) Equal(other *Node) bool {
	if other == nil {
		return false
	}

	return node.sitterNode.Equal(other.sitterNode)
}

func (node *Node) Type() string {
	return node.sitterNode.Type()
}

func (node *Node) Content() string {
	return node.sitterNode.Content(node.tree.input)
}

func (node *Node) Parent() *Node {
	return node.tree.wrap(node.sitterNode.Parent())
}

func (node *Node) ChildByFieldName(name string) *Node {
	return node.tree.wrap(node.sitterNode.ChildByFieldName(name))
}

func (node *Node) Walk(visit func(node *Node) error) error {
	cursor := sitter.NewTreeCursor(node.sitterNode)
	defer cursor.Close()

	for {
		if cursor.CurrentNode().IsNamed() {
			if err := visit(node.tree.wrap(cursor.CurrentNode())); err != nil {
				if err == ErrTerminateWalk {
					return nil
				}

				return err
			}
		}

		if cursor.GoToFirstChild() || cursor.GoToNextSibling() {
			continue
		}

		for {
			if !cursor.GoToParent() {
				// Reached the root again
				return nil
			}

			if cursor.GoToNextSibling() {
				break
			}
		}
	}
}
