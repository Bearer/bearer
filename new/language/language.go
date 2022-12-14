package language

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type Base struct {
	sitterLanguage *sitter.Language
	unifiedNodes   map[NodeID][]*Node
}

func New(sitterLanguage *sitter.Language) Base {
	return Base{
		sitterLanguage: sitterLanguage,
		unifiedNodes:   make(map[NodeID][]*Node),
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

	return &Tree{
		input:      inputBytes,
		sitterTree: sitterTree,
	}, nil
}

func (lang *Base) CompileQuery(input string) (*Query, error) {
	sitterQuery, err := sitter.NewQuery([]byte(input), lang.sitterLanguage)
	if err != nil {
		return nil, err
	}

	return &Query{sitterQuery: sitterQuery}, nil
}

func (lang *Base) UnifyNodes(laterNode *Node, earlierNode *Node) {
	if laterNode.Equal(earlierNode) {
		return
	}

	existingUnifiedNodes := lang.unifiedNodes[laterNode.ID()]

	for _, other := range existingUnifiedNodes {
		if other.Equal(earlierNode) {
			// already unified
			return
		}
	}

	lang.unifiedNodes[laterNode.ID()] = append(existingUnifiedNodes, earlierNode)
}

func (lang *Base) UnifiedNodesFor(node *Node) []*Node {
	return lang.unifiedNodes[node.ID()]
}
