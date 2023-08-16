package tree

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"
)

type Builder struct {
	contentBytes    []byte
	types           []string
	nodes           []Node
	rootNodeID      int
	children        map[int][]int
	dataflowSources map[int][]int
	sitterToNodeID  map[*sitter.Node]int
}

func NewBuilder(contentBytes []byte, sitterRootNode *sitter.Node) *Builder {
	builder := &Builder{
		contentBytes:    contentBytes,
		children:        make(map[int][]int),
		dataflowSources: make(map[int][]int),
		sitterToNodeID:  make(map[*sitter.Node]int),
	}

	builder.rootNodeID = builder.addNode(sitterRootNode)

	return builder
}

func (builder *Builder) ContentFor(sitterNode *sitter.Node) string {
	return sitterNode.Content(builder.contentBytes)
}

func (builder *Builder) Dataflow(toNode *sitter.Node, fromNodes ...*sitter.Node) {
	toID := builder.sitterToNodeID[toNode]

	fromIDs := make([]int, len(fromNodes))
	for i, fromNode := range fromNodes {
		fromIDs[i] = builder.sitterToNodeID[fromNode]
	}

	builder.dataflowSources[toID] = append(builder.dataflowSources[toID], fromIDs...)
}

func (builder *Builder) Build() *Tree {
	builder.buildChildren()
	builder.buildDataflowSources()

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
}

func (builder *Builder) buildDataflowSources() {
	totalCount := 0
	for _, sourceIDs := range builder.dataflowSources {
		totalCount += len(sourceIDs)
	}

	dataflowSources := make([]*Node, totalCount)

	offset := 0
	for id := range builder.nodes {
		sourceIDs := builder.dataflowSources[id]
		count := len(sourceIDs)
		if count == 0 {
			continue
		}

		nodeDataflowSources := dataflowSources[offset : offset+count]

		for i, sourceID := range sourceIDs {
			nodeDataflowSources[i] = &builder.nodes[sourceID]
		}

		builder.nodes[id].dataflowSources = nodeDataflowSources
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
