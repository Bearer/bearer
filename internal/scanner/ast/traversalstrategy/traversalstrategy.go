package traversalstrategy

import (
	"fmt"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bits-and-blooms/bitset"
)

var (
	Nested = &Strategy{
		scope: settings.NESTED_SCOPE,
		nextNodes: func(node *tree.Node) []*tree.Node {
			return append(node.Children(), node.AliasOf()...)
		},
	}

	NestedStrict = &Strategy{
		scope: settings.NESTED_STRICT_SCOPE,
		nextNodes: func(node *tree.Node) []*tree.Node {
			return node.Children()
		},
	}

	Result = &Strategy{
		scope: settings.RESULT_SCOPE,
		nextNodes: func(node *tree.Node) []*tree.Node {
			return append(node.AliasOf(), node.DataflowSources()...)
		},
	}

	Cursor = &Strategy{
		scope: settings.CURSOR_SCOPE,
		nextNodes: func(node *tree.Node) []*tree.Node {
			return node.AliasOf()
		},
	}

	CursorStrict = &Strategy{
		scope:     settings.CURSOR_STRICT_SCOPE,
		nextNodes: nil,
	}
)

type Strategy struct {
	scope     settings.RuleReferenceScope
	nextNodes func(node *tree.Node) []*tree.Node
}

func Get(scope settings.RuleReferenceScope) (*Strategy, error) {
	switch scope {
	case settings.NESTED_SCOPE:
		return Nested, nil
	case settings.NESTED_STRICT_SCOPE:
		return NestedStrict, nil
	case settings.RESULT_SCOPE:
		return Result, nil
	case settings.CURSOR_SCOPE:
		return Cursor, nil
	case settings.CURSOR_STRICT_SCOPE:
		return CursorStrict, nil
	default:
		return nil, fmt.Errorf("unknown scope '%s'", scope)
	}
}

func (strategy *Strategy) Scope() settings.RuleReferenceScope {
	return strategy.scope
}

func (strategy *Strategy) Traverse(rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	if strategy.nextNodes == nil {
		_, err := visit(rootNode)
		return err
	}

	next := make([]*tree.Node, 0, 1000)
	nodes := make([]*tree.Node, 0, 1000)
	nodes = append(nodes, rootNode)
	seen := bitset.New(uint(rootNode.Tree().NodeCount()))

	for {
		if len(nodes) == 0 {
			break
		}

		for _, node := range nodes {
			bit := uint(node.ID)
			if seen.Test(bit) {
				continue
			}
			seen.Set(bit)

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
