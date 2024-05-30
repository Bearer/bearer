package traversalstrategy

import (
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
)

var (
	Nested       = &nestedStrategy{}
	NestedStrict = &nestedStrictStrategy{}
	Result       = &resultStrategy{}
	Cursor       = &cursorStrategy{}
	CursorStrict = &cursorStrictStrategy{}
)

type nestedStrategy struct{}

var nestedTraverse = makeTraverse(func(next *[]*tree.Node, node *tree.Node) {
	*next = append(*next, node.Children()...)
	*next = append(*next, node.AliasOf()...)
})

func (strategy *nestedStrategy) Scope() settings.RuleReferenceScope {
	return settings.NESTED_SCOPE
}

func (strategy *nestedStrategy) Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	return nestedTraverse(cache, rootNode, visit)
}

type nestedStrictStrategy struct{}

var nestedStrictTraverse = makeTraverse(func(next *[]*tree.Node, node *tree.Node) {
	*next = append(*next, node.Children()...)
})

func (strategy *nestedStrictStrategy) Scope() settings.RuleReferenceScope {
	return settings.NESTED_STRICT_SCOPE
}

func (strategy *nestedStrictStrategy) Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	return nestedStrictTraverse(cache, rootNode, visit)
}

type resultStrategy struct{}

var resultTraverse = makeTraverse(func(next *[]*tree.Node, node *tree.Node) {
	*next = append(*next, node.AliasOf()...)
	*next = append(*next, node.DataflowSources()...)
})

func (strategy *resultStrategy) Scope() settings.RuleReferenceScope {
	return settings.RESULT_SCOPE
}

func (strategy *resultStrategy) Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	return resultTraverse(cache, rootNode, visit)
}

type cursorStrategy struct{}

var cursorTraverse = makeTraverse(func(next *[]*tree.Node, node *tree.Node) {
	*next = append(*next, node.AliasOf()...)
})

func (strategy *cursorStrategy) Scope() settings.RuleReferenceScope {
	return settings.CURSOR_SCOPE
}

func (strategy *cursorStrategy) Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	return cursorTraverse(cache, rootNode, visit)
}

type cursorStrictStrategy struct{}

func (strategy *cursorStrictStrategy) Scope() settings.RuleReferenceScope {
	return settings.CURSOR_STRICT_SCOPE
}

func (strategy *cursorStrictStrategy) Traverse(cache *Cache, rootNode *tree.Node, visit func(node *tree.Node) (bool, error)) error {
	_, err := visit(rootNode)
	return err
}
