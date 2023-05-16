package tree

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Node struct {
	tree       *Tree
	sitterNode *sitter.Node
}

type NodeID *sitter.Node

func (node *Node) Debug() string {
	return node.sitterNode.String()
}

func (node *Node) ID() NodeID {
	return node.sitterNode
}

func (node *Node) Equal(other *Node) bool {
	if other == nil {
		return false
	}

	return node.sitterNode.Equal(other.sitterNode)
}

func (node *Node) Type() string {
	return node.sitterNode.Type()
}

func (node *Node) Content() string {
	return node.sitterNode.Content(node.tree.input)
}

func (node *Node) StartByte() int {
	return int(node.sitterNode.StartByte())
}

func (node *Node) EndByte() int {
	return int(node.sitterNode.EndByte())
}

func (node *Node) LineNumber() int {
	return int(node.sitterNode.StartPoint().Row + 1)
}

func (node *Node) ColumnNumber() int {
	return int(node.sitterNode.StartPoint().Column + 1)
}

func (node *Node) Parent() *Node {
	return node.tree.wrap(node.sitterNode.Parent())
}

func (node *Node) ChildCount() int {
	return int(node.sitterNode.ChildCount())
}

func (node *Node) NamedChildCount() int {
	count := 0
	for i := 0; i < int(node.sitterNode.ChildCount()); i++ {
		if node.sitterNode.Child(i).IsNamed() {
			count++
		}
	}

	return count
}

func (node *Node) Child(i int) *Node {
	return node.tree.wrap(node.sitterNode.Child(i))
}

func (node *Node) NamedChild(i int) *Node {
	count := 0
	for j := 0; j < int(node.sitterNode.ChildCount()); j++ {
		child := node.sitterNode.Child(j)

		if child.IsNamed() {
			if count == i {
				return node.tree.wrap(child)
			}

			count++
		}
	}

	return nil
}

func (node *Node) EachContentPart(onText func(text string) error, onChild func(child *Node) error) error {
	n := int(node.sitterNode.ChildCount())

	start := node.sitterNode.StartByte()
	end := start

	emit := func() error {
		if end <= start {
			return nil
		}

		return onText(string(node.tree.input[start:end]))
	}

	for i := 0; i < n; i++ {
		child := node.sitterNode.Child(i)
		end = child.StartByte()

		if err := emit(); err != nil {
			return err
		}

		if child.IsNamed() {
			if err := onChild(node.tree.wrap(child)); err != nil {
				return err
			}
		}

		start = child.EndByte()
		end = start
	}

	if err := emit(); err != nil {
		return err
	}

	return nil
}

func (node *Node) AnonymousChild(i int) *Node {
	n := int(node.sitterNode.ChildCount())
	k := 0

	for j := 0; j < n; j++ {
		child := node.sitterNode.Child(j)
		if child.IsNamed() {
			continue
		}

		if k == i {
			return node.tree.wrap(child)
		}

		k += 1
	}

	return nil
}

func (node *Node) ChildByFieldName(name string) *Node {
	return node.tree.wrap(node.sitterNode.ChildByFieldName(name))
}

func (node *Node) IsNamed() bool {
	return node.sitterNode.IsNamed()
}

func (node *Node) IsError() bool {
	return node.sitterNode.IsError()
}

func (node *Node) Walk(visit func(node *Node, visitChildren func() error) error) error {
	visitChildren := func() error {
		for i := 0; i < node.ChildCount(); i += 1 {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			if err := child.Walk(visit); err != nil {
				return err
			}
		}

		return nil
	}

	return visit(node, visitChildren)
}

func (node *Node) UnifyWith(earlierNode *Node) {
	node.tree.unifyNodes(node, earlierNode)
}

func (node *Node) UnifiedNodes() []*Node {
	return node.tree.unifiedNodesFor(node)
}
