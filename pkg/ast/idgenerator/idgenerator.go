package idgenerator

import sitter "github.com/smacker/go-tree-sitter"

type Generator struct {
	next uint32
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (generator *Generator) Get() uint32 {
	next := generator.next
	generator.next++
	return next
}

type NodeIdGenerator struct {
	ids         map[*sitter.Node]uint32
	inverseIds  map[uint32]*sitter.Node
	idGenerator *Generator
}

func NewNodeIdGenerator() *NodeIdGenerator {
	return &NodeIdGenerator{
		ids:         make(map[*sitter.Node]uint32),
		inverseIds:  make(map[uint32]*sitter.Node),
		idGenerator: NewGenerator(),
	}
}

func (generator *NodeIdGenerator) Get(node *sitter.Node) uint32 {
	if id, cached := generator.ids[node]; cached {
		return id
	}

	id := generator.idGenerator.Get()
	generator.ids[node] = id
	generator.inverseIds[id] = node
	return id
}

func (generator *NodeIdGenerator) InverseLookup(nodeId uint32) *sitter.Node {
	return generator.inverseIds[nodeId]
}
