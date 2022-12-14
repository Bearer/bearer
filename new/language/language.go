package language

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type Base struct {
	sitterLanguage *sitter.Language
	analyzeFlow    func(rootNode *Node)
}

func New(sitterLanguage *sitter.Language, analyzeFlow func(rootNode *Node)) Base {
	return Base{
		sitterLanguage: sitterLanguage,
		analyzeFlow:    analyzeFlow,
	}
}

func (lang *Base) Parse(input string) (*Tree, error) {
	inputBytes := []byte(input)

	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(lang.sitterLanguage)

	sitterTree, err := parser.ParseCtx(context.Background(), nil, inputBytes)
	if err != nil {
		return nil, err
	}

	tree := &Tree{
		input:        inputBytes,
		sitterTree:   sitterTree,
		unifiedNodes: make(map[NodeID][]*Node),
	}

	lang.analyzeFlow(tree.RootNode())

	return tree, nil
}

func (lang *Base) CompileQuery(input string) (*Query, error) {
	sitterQuery, err := sitter.NewQuery([]byte(input), lang.sitterLanguage)
	if err != nil {
		return nil, err
	}

	return &Query{sitterQuery: sitterQuery}, nil
}
