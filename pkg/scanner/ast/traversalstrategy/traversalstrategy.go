package traversalstrategy

import (
	"fmt"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bits-and-blooms/bitset"
)

type Strategy interface {
	Scope() settings.RuleReferenceScope
	Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error
}

func Get(scope settings.RuleReferenceScope) (Strategy, error) {
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

type Cache struct {
	nodeCount int
	allocated []*data
}

type data struct {
	seen *bitset.BitSet
	nodes,
	next []*tree.Node
}

func NewCache(nodeCount int) *Cache {
	return &Cache{nodeCount: nodeCount}
}

func (cache *Cache) get() *data {
	if len(cache.allocated) == 0 {
		return &data{
			seen:  bitset.New(uint(cache.nodeCount)),
			nodes: make([]*tree.Node, 0, 1000),
			next:  make([]*tree.Node, 0, 1000),
		}
	}

	index := len(cache.allocated) - 1
	data := cache.allocated[index]
	cache.allocated = cache.allocated[:index]
	return data
}

func (cache *Cache) put(data *data) {
	// same buffer but zero length
	data.nodes = data.nodes[:0]
	data.next = data.next[:0]
	data.seen.ClearAll()

	cache.allocated = append(cache.allocated, data)
}

func makeTraverse(appendNext func(next *[]*tree.Node, node *tree.Node)) func(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	return func(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
		data := cache.get()
		defer cache.put(data)

		data.nodes = append(data.nodes, rootNode)

		for len(data.nodes) != 0 {

			for _, node := range data.nodes {
				bit := uint(node.ID)
				if data.seen.Test(bit) {
					continue
				}
				data.seen.Set(bit)

				stopTraversal, err := visit(node)
				if err != nil {
					return err
				}

				if stopTraversal {
					continue
				}

				appendNext(&data.next, node)
			}

			old := data.nodes
			data.nodes = data.next
			// allow memory to be re-used
			data.next = old[:0]
		}

		return nil
	}
}
