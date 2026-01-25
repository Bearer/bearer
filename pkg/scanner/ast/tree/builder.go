package tree

import (
	"slices"

	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bits-and-blooms/bitset"
	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

type Builder struct {
	contentBytes        []byte
	types               []string
	nodes               []Node
	rootNodeID          int
	children,
	dataflowSources,
	aliasOf             map[int][]int
	childrenByField     map[int]map[string]int
	sitterRootNode      *sitter.Node
	sitterToNodeID      map[*sitter.Node]int
	fieldNames          []string
	ruleCount           int
	stringFragmentTypes []string
}

func NewBuilder(
	sitterLanguage *sitter.Language,
	contentBytes []byte,
	sitterRootNode *sitter.Node,
	ruleCount int,
	stringFragmentTypes []string,
) *Builder {
	var fieldNames []string
	for i := 1; ; i++ {
		name := sitterLanguage.FieldName(i)
		if name == "" {
			break
		}

		fieldNames = append(fieldNames, name)
	}

	builder := &Builder{
		contentBytes:        contentBytes,
		nodes:               make([]Node, 0, 1000),
		children:            make(map[int][]int),
		dataflowSources:     make(map[int][]int),
		aliasOf:             make(map[int][]int),
		childrenByField:     make(map[int]map[string]int),
		sitterRootNode:      sitterRootNode,
		sitterToNodeID:      make(map[*sitter.Node]int),
		fieldNames:          fieldNames,
		ruleCount:           ruleCount,
		stringFragmentTypes: stringFragmentTypes,
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

func (builder *Builder) AddExpectedRules(sitterNode *sitter.Node, rules []*ruleset.Rule) {
	if len(rules) == 0 {
		return
	}

	builder.addExpectedRulesForNode(builder.sitterToNodeID[sitterNode], rules)
}

func (builder *Builder) AddDisabledRules(sitterNode *sitter.Node, rules []*ruleset.Rule) {
	if len(rules) == 0 {
		return
	}

	builder.addDisabledRulesForNode(builder.sitterToNodeID[sitterNode], rules)
}

func (builder *Builder) addExpectedRulesForNode(nodeID int, rules []*ruleset.Rule) {
	node := &builder.nodes[nodeID]

	for _, rule := range rules {
		node.expectedRules = append(node.expectedRules, rule.ID())
	}
}

func (builder *Builder) addDisabledRulesForNode(nodeID int, rules []*ruleset.Rule) {
	node := &builder.nodes[nodeID]
	if node.disabledRuleIndices == nil {
		node.disabledRuleIndices = bitset.New(uint(builder.ruleCount))
	}

	for _, rule := range rules {
		node.disabledRuleIndices.Set(uint(rule.Index()))
	}

	for _, childID := range builder.children[nodeID] {
		builder.addDisabledRulesForNode(childID, rules)
	}
}

func (builder *Builder) sitterToNodeIDs(nodes []*sitter.Node) []int {
	ids := make([]int, len(nodes))

	for i, node := range nodes {
		ids[i] = builder.sitterToNodeID[node]
	}

	return ids
}

func (builder *Builder) QueryResult(queryID int, sitterNode *sitter.Node, result map[string]*sitter.Node) {
	nodeID := builder.sitterToNodeID[sitterNode]
	node := &builder.nodes[nodeID]

	if node.queryResults == nil {
		node.queryResults = make(map[int][]QueryResult)
	}

	translatedResult := builder.translateNodeMap(result)

	// Deduplicate: skip if an equivalent result already exists for this queryID
	// Two results are considered equivalent if they have the same root node
	for _, existing := range node.queryResults[queryID] {
		if existingRoot, ok := existing["root"]; ok {
			if newRoot, ok := translatedResult["root"]; ok {
				if existingRoot == newRoot {
					if log.Trace().Enabled() {
						log.Trace().Msgf("QueryResult: skipping duplicate queryID=%d on node id=%d", queryID, nodeID)
					}
					return
				}
			}
		}
	}

	if log.Trace().Enabled() {
		content := builder.contentBytes[sitterNode.StartByte():sitterNode.EndByte()]
		nodeType := builder.types[node.TypeID]
		log.Trace().Msgf("QueryResult: storing queryID=%d on node id=%d type=%s content=%q", queryID, nodeID, nodeType, string(content))
	}

	node.queryResults[queryID] = append(node.queryResults[queryID], translatedResult)
}

func (builder *Builder) Build() *Tree {
	builder.buildChildren()
	builder.buildChildrenByField()
	builder.buildDataflowSources()
	builder.buildAliasOf()

	tree := &Tree{
		contentBytes:        builder.contentBytes,
		types:               builder.types,
		nodes:               builder.nodes,
		rootNode:            &builder.nodes[builder.rootNodeID],
		sitterToNode:        builder.buildSitterToNode(),
		stringFragmentTypes: builder.stringFragmentTypes,
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
		sitterNode:      sitterNode,
		childrenByField: make(map[string]*Node),
		ID:              id,
		TypeID:          builder.internType(sitterType),
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

	builder.children[id], builder.childrenByField[id] = builder.addChildren(id, sitterNode)

	return id
}

func (builder *Builder) addChildren(parentID int, sitterNode *sitter.Node) ([]int, map[string]int) {
	sitterParent := builder.nodes[parentID].sitterNode

	childCount := int(sitterNode.ChildCount())
	if childCount == 0 {
		return nil, make(map[string]int)
	}

	children := make([]int, childCount)
	childrenByField := make(map[string]int)

	for i := 0; i < childCount; i++ {
		sitterChild := sitterNode.Child(i)
		childID := builder.addNode(sitterChild)
		children[i] = childID

		var fieldName string
		for _, candidateName := range builder.fieldNames {
			if candidate := sitterParent.ChildByFieldName(candidateName); candidate != nil && sitterChild.Equal(candidate) {
				fieldName = candidateName
				break
			}
		}

		if fieldName != "" {
			builder.nodes[childID].fieldName = fieldName
			childrenByField[fieldName] = childID
		}
	}

	return children, childrenByField
}

func (builder *Builder) buildChildren() {
	builder.buildAdjacencyList(builder.children, func(node *Node, children []*Node) {
		node.children = children
	})
}

func (builder *Builder) buildChildrenByField() {
	for id := range builder.nodes {
		node := &builder.nodes[id]

		for name, childID := range builder.childrenByField[id] {
			node.childrenByField[name] = &builder.nodes[childID]
		}
	}
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
