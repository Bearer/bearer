package nodeid

import (
	"strconv"

	"github.com/bearer/bearer/internal/parser"
	"github.com/google/uuid"
)

type Map struct {
	tree      *parser.Tree
	values    map[parser.NodeID]string
	generator Generator
}

func New(tree *parser.Tree, generator Generator) *Map {
	return &Map{
		tree:      tree,
		generator: generator,
		values:    make(map[parser.NodeID]string),
	}
}

func (finder *Map) Annotate() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck
		finder.values[child.ID()] = finder.generator.GenerateId()
		return nil
	})
}

func (finder *Map) ValueForNode(id parser.NodeID) string {
	return finder.values[id]
}

type Generator interface {
	GenerateId() string
}

type IntGenerator struct {
	Counter int
}

func (generator *IntGenerator) GenerateId() string {
	generator.Counter++
	return strconv.Itoa(generator.Counter)
}

type UUIDGenerator struct {
}

func (generator *UUIDGenerator) GenerateId() string {
	return uuid.NewString()
}
