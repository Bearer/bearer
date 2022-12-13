package parser

import (
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
)

type Node struct {
	tree       *Tree
	sitterNode *sitter.Node
}

type NodeID *sitter.Node

var TerminateWalk = errors.New("terminate walk")

func (node *Node) ID() NodeID {
	return node.sitterNode
}

func (node *Node) Walk(visit func(*Node) error) error {
	cursor := sitter.NewTreeCursor(node.sitterNode)
	defer cursor.Close()

	for {
		if cursor.CurrentNode().IsNamed() {
			if err := visit(node.tree.wrap(cursor.CurrentNode())); err != nil {
				if err == TerminateWalk {
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
