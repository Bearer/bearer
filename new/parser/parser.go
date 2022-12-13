package parser

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	input      []byte
	sitterTree *sitter.Tree
}

func Parse(language *sitter.Language, input string) (*Tree, error) {
	inputBytes := []byte(input)

	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language)

	sitterTree, err := parser.ParseCtx(context.Background(), nil, inputBytes)
	if err != nil {
		return nil, err
	}

	return &Tree{
		input:      inputBytes,
		sitterTree: sitterTree,
	}, nil
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
