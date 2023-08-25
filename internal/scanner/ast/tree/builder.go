package tree

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"
)

type Builder struct {
	contentBytes []byte
	types        []string
	nodes        []Node
	rootNodeID   int
	children,
	dataflowSources,
	aliasOf map[int][]int
	sitterRootNode *sitter.Node
	sitterToNodeID map[*sitter.Node]int
}

func NewBuilder(contentBytes []byte, sitterRootNode *sitter.Node) *Builder {
	builder := &Builder{
		contentBytes:    contentBytes,
		children:        make(map[int][]int),
		dataflowSources: make(map[int][]int),
		aliasOf:         make(map[int][]int),
		sitterRootNode:  sitterRootNode,
		sitterToNodeID:  make(map[*sitter.Node]int),
	}

	builder.rootNodeID = builder.addNode(sitterRootNode)

	return builder
}

func (builder *Builder) SitterRootNode() *sitter.Node {
	return builder.sitterRootNode
}

func (builder *Builder) LastChild(node *sitter.Node) *sitter.Node {
	childCount := int(node.ChildCount())
	if childCount == 0 {
		return nil
	}

	return node.Child(childCount - 1)
}

func (builder *Builder) ChildrenFor(node *sitter.Node) []*sitter.Node {
	childCount := int(node.ChildCount())
	children := make([]*sitter.Node, childCount)

	for i := 0; i < childCount; i++ {
		children[i] = node.Child(i)
	}

	return children
}
func (builder *Builder) ChildrenExcept(node, excludedNode *sitter.Node) []*sitter.Node {
	childCount := int(node.ChildCount())
	children := make([]*sitter.Node, 0, childCount)

	for i := 0; i < childCount; i++ {
		if child := node.Child(i); child != excludedNode {
			children = append(children, child)
		}
	}

	return children
}

func (builder *Builder) ContentFor(node *sitter.Node) string {
	return node.Content(builder.contentBytes)
}

func (builder *Builder) Dataflow(toNode *sitter.Node, fromNodes ...*sitter.Node) {
	toID := builder.sitterToNodeID[toNode]

	builder.dataflowSources[toID] = append(
		builder.dataflowSources[toID],
		builder.sitterToNodeIDs(fromNodes)...,
	)
}

func (builder *Builder) Alias(toNode *sitter.Node, fromNodes ...*sitter.Node) {
	toID := builder.sitterToNodeID[toNode]

	builder.aliasOf[toID] = append(
		builder.aliasOf[toID],
		builder.sitterToNodeIDs(fromNodes)...,
	)
}

func (builder *Builder) AddDisabledRules(sitterNode *sitter.Node, ruleIDs []string) {
	node := &builder.nodes[builder.sitterToNodeID[sitterNode]]
	node.disabledRuleIDs = append(node.disabledRuleIDs, ruleIDs...)
}

func (builder *Builder) sitterToNodeIDs(nodes []*sitter.Node) []int {
	ids := make([]int, len(nodes))

	for i, node := range nodes {
		ids[i] = builder.sitterToNodeID[node]
	}

	return ids
}

func (builder *Builder) QueryResult(queryID int, sitterNode *sitter.Node, result map[string]*sitter.Node) {
	node := &builder.nodes[builder.sitterToNodeID[sitterNode]]

	if node.queryResults == nil {
		node.queryResults = make(map[int][]QueryResult)
	}

	node.queryResults[queryID] = append(node.queryResults[queryID], builder.translateNodeMap(result))
}

func (builder *Builder) Build() *Tree {
	builder.buildChildren()
	builder.buildDataflowSources()
	builder.buildAliasOf()

	tree := &Tree{
		contentBytes: builder.contentBytes,
		types:        builder.types,
		nodes:        builder.nodes,
		rootNode:     &builder.nodes[builder.rootNodeID],
		sitterToNode: builder.buildSitterToNode(),
	}

	for i := range tree.nodes {
		tree.nodes[i].tree = tree
		tree.nodes[i].parent = tree.sitterToNode[tree.nodes[i].sitterNode.Parent()]
	}

	return tree
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
		sitterNode: sitterNode,
		ID:         id,
		TypeID:     builder.internType(sitterType),
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

func (builder *Builder) buildChildren() {
	builder.buildAdjacencyList(builder.children, func(node *Node, children []*Node) {
		node.children = children
	})
}

func (builder *Builder) buildDataflowSources() {
	builder.buildAdjacencyList(builder.dataflowSources, func(node *Node, dataflowSources []*Node) {
		node.dataflowSources = dataflowSources
	})
}

func (builder *Builder) buildAliasOf() {
	builder.buildAdjacencyList(builder.aliasOf, func(node *Node, aliasOf []*Node) {
		node.aliasOf = aliasOf
	})
}

func (builder *Builder) buildAdjacencyList(
	nodeToAdjacencyIDs map[int][]int,
	assignToNode func(node *Node, adjacentNodes []*Node),
) {
	totalCount := 0
	for _, adjacentIDs := range nodeToAdjacencyIDs {
		totalCount += len(adjacentIDs)
	}

	// use a single backing slice for memory-local traversal
	store := make([]*Node, totalCount)

	offset := 0
	for id := range builder.nodes {
		adjacentIDs := nodeToAdjacencyIDs[id]
		count := len(adjacentIDs)
		if count == 0 {
			continue
		}

		// this shares memory with the store
		adjacentNodes := store[offset : offset+count]

		for i, adjacentID := range adjacentIDs {
			adjacentNodes[i] = &builder.nodes[adjacentID]
		}

		assignToNode(&builder.nodes[id], adjacentNodes)
		offset += count
	}
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

func (builder *Builder) translateNodeMap(sitterMap map[string]*sitter.Node) map[string]*Node {
	result := make(map[string]*Node)

	for name, sitterNode := range sitterMap {
		result[name] = &builder.nodes[builder.sitterToNodeID[sitterNode]]
	}

	return result
}
