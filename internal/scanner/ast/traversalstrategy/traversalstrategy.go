package traversalstrategy

import (
	"fmt"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/util/set"
)

var strategies = map[settings.RuleReferenceScope]*Strategy{
	settings.NESTED_SCOPE: {
		nextNodes: func(node *tree.Node) []*tree.Node {
			return append(node.Children(), node.AliasOf()...)
		},
	},
	settings.NESTED_STRICT_SCOPE: {
		nextNodes: func(node *tree.Node) []*tree.Node {
			return node.Children()
		},
	},
	settings.RESULT_SCOPE: {
		nextNodes: func(node *tree.Node) []*tree.Node {
			return append(node.AliasOf(), node.DataflowSources()...)
		},
	},
	settings.CURSOR_SCOPE: {
		nextNodes: func(node *tree.Node) []*tree.Node {
			return node.AliasOf()
		},
	},
	settings.CURSOR_STRICT_SCOPE: {},
}

type Strategy struct {
	nextNodes func(node *tree.Node) []*tree.Node
}

func Get(scope settings.RuleReferenceScope) (*Strategy, error) {
	strategy, exists := strategies[scope]
	if !exists {
		return nil, fmt.Errorf("unknown scope '%s'", scope)
	}

	return strategy, nil
}

func (strategy *Strategy) Traverse(rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	if strategy.nextNodes == nil {
		_, err := visit(rootNode)
		return err
	}

	next := make([]*tree.Node, 0, 1000)
	nodes := make([]*tree.Node, 0, 1000)
	nodes = append(nodes, rootNode)
	seen := set.New[*tree.Node]()

	for {
		if len(nodes) == 0 {
			break
		}

		for _, node := range nodes {
			if !seen.Add(node) {
				continue
			}

			stopTraversal, err := visit(node)
			if err != nil {
				return err
			}

			if stopTraversal {
				continue
			}

			next = append(next, strategy.nextNodes(node)...)
		}

		old := nodes
		nodes = next
		// allow memory to be re-used
		next = old[:0]
	}

	return nil
}
