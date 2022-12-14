package language

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	input        []byte
	sitterTree   *sitter.Tree
	unifiedNodes map[NodeID][]*Node
}

func (tree *Tree) RootNode() *Node {
	return tree.wrap(tree.sitterTree.RootNode())
}

func (tree *Tree) Close() {
	tree.sitterTree.Close()
}

func (tree *Tree) wrap(sitterNode *sitter.Node) *Node {
	if sitterNode == nil {
		return nil
	}

	return &Node{tree: tree, sitterNode: sitterNode}
}

func (tree *Tree) unifyNodes(laterNode *Node, earlierNode *Node) {
	if laterNode.Equal(earlierNode) {
		return
	}

	existingUnifiedNodes := tree.unifiedNodes[laterNode.ID()]

	for _, other := range existingUnifiedNodes {
		if other.Equal(earlierNode) {
			// already unified
			return
		}
	}

	tree.unifiedNodes[laterNode.ID()] = append(existingUnifiedNodes, earlierNode)
}

func (tree *Tree) unifiedNodesFor(node *Node) []*Node {
	return tree.unifiedNodes[node.ID()]
}
