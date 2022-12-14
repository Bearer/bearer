package language

import (
	"context"
	"log"

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

func (lang *Base) UnifyNodes(a *Node, b *Node) {
	aNodes, aExists := lang.unifiedNodes[a.ID()]
	bNodes, bExists := lang.unifiedNodes[b.ID()]

	if nodesInclude(aNodes, b) {
		// already unified
		return
	}

	if !aExists {
		aNodes = []*Node{a}
	}
	if !bExists {
		bNodes = []*Node{b}
	}

	newNodes := append(aNodes, bNodes...)

	for _, node := range newNodes {
		lang.unifiedNodes[node.ID()] = newNodes
	}
}

func nodesInclude(nodes []*Node, node *Node) bool {
	for _, other := range nodes {
		if other.Equal(node) {
			return true
		}
	}

	return false
}

func (lang *Base) UnifiedNodesFor(node *Node) []*Node {
	nodes, ok := lang.unifiedNodes[node.ID()]
	if !ok {
		return []*Node{node}
	}

	log.Printf("UNIFYIED: %v\n", nodes)
	return nodes
}
