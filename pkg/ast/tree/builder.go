package tree

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"
)

type Builder struct {
	types          []string
	nodes          []Node
	rootNodeID     int
	children       map[int][]int
	sitterToNodeID map[*sitter.Node]int
}

func NewBuilder(sitterRootNode *sitter.Node) *Builder {
	builder := &Builder{
		children:       make(map[int][]int),
		sitterToNodeID: make(map[*sitter.Node]int),
	}

	builder.rootNodeID = builder.addNode(sitterRootNode)

	return builder
}

func (builder *Builder) AddPatternMatch(sitterNode *sitter.Node, match struct{}) {
	// node := &builder.nodes[builder.sitterToNodeID[sitterNode]]
	// node.patternMatches = append(node.patternMatches, match)
}

func (builder *Builder) Build() *Node {
	tree := &Tree{}
	tree.types = builder.types
	tree.nodes = builder.nodes
	tree.children = builder.buildChildren()
	// tree.patternMatches = builder.buildPatternMatches()
	tree.sitterToNode = builder.buildSitterToNode()

	for i := range tree.nodes {
		tree.nodes[i].tree = tree
	}

	return &tree.nodes[builder.rootNodeID]
}

func (builder *Builder) addNode(sitterNode *sitter.Node) int {
	id := len(builder.nodes)
	builder.sitterToNodeID[sitterNode] = id

	startPoint := sitterNode.StartPoint()
	endPoint := sitterNode.EndPoint()

	sitterType := sitterNode.Type()
	if !sitterNode.IsNamed() {
		sitterType = `"` + sitterType + `"`
	}

	builder.nodes = append(builder.nodes, Node{
		ID:     id,
		TypeID: builder.internType(sitterType),
		ContentStart: Position{
			Byte:   int(sitterNode.StartByte()),
			Line:   int(startPoint.Row) + 1,
			Column: int(startPoint.Column) + 1,
		},
		ContentEnd: Position{
			Byte:   int(sitterNode.EndByte()),
			Line:   int(endPoint.Row) + 1,
			Column: int(endPoint.Column) + 1,
		},
	})

	builder.children[id] = builder.addChildren(id, sitterNode)

	return id
}

func (builder *Builder) addChildren(parentID int, sitterNode *sitter.Node) []int {
	childCount := int(sitterNode.ChildCount())
	if childCount == 0 {
		return nil
	}

	children := make([]int, childCount)
	for i := 0; i < childCount; i++ {
		children[i] = builder.addNode(sitterNode.Child(i))
	}

	return children
}

func (builder *Builder) buildChildren() []*Node {
	totalCount := 0
	for _, childIDs := range builder.children {
		totalCount += len(childIDs)
	}

	children := make([]*Node, totalCount)

	offset := 0
	for id := range builder.nodes {
		childIDs := builder.children[id]
		count := len(childIDs)
		if count == 0 {
			continue
		}

		nodeChildren := children[offset : offset+count]

		for i, childID := range childIDs {
			nodeChildren[i] = &builder.nodes[childID]
		}

		builder.nodes[id].children = nodeChildren
		offset += count
	}

	return children
}

func (builder *Builder) buildSitterToNode() map[*sitter.Node]*Node {
	result := make(map[*sitter.Node]*Node)

	for sitterNode, ID := range builder.sitterToNodeID {
		result[sitterNode] = &builder.nodes[ID]
	}

	return result
}

func (builder *Builder) internType(nodeType string) int {
	id := slices.Index(builder.types, nodeType)
	if id != -1 {
		return id
	}

	id = len(builder.types)
	builder.types = append(builder.types, nodeType)
	return id
}
