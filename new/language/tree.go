package language

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	input      []byte
	sitterTree *sitter.Tree
}

func (tree *Tree) RootNode() *Node {
	return tree.wrap(tree.sitterTree.RootNode())
}

func (tree *Tree) Close() {
	tree.sitterTree.Close()
}

func (tree *Tree) wrap(sitterNode *sitter.Node) *Node {
	return &Node{tree: tree, sitterNode: sitterNode}
}
