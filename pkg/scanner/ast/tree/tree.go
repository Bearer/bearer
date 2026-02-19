package tree

import (
	"fmt"
	"maps"
	"slices"

	"github.com/bits-and-blooms/bitset"
	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
	"gopkg.in/yaml.v3"
)

type Tree struct {
	contentBytes        []byte
	types               []string
	nodes               []Node
	rootNode            *Node
	sitterToNode        map[*sitter.Node]*Node
	stringFragmentTypes []string
}

type QueryResult map[string]*Node

type Node struct {
	tree *Tree
	ID,
	TypeID int
	fieldName string
	ContentStart,
	ContentEnd Position
	parent *Node
	children,
	dataflowSources,
	aliasOf []*Node
	childrenByField     map[string]*Node
	expectedRules       []string
	disabledRuleIndices *bitset.BitSet
	// FIXME: remove the need for this
	sitterNode   *sitter.Node
	queryResults map[int][]QueryResult
	// FIXME: probably shouldn't be public
	ExecutingDetectors []int
}

type Position struct {
	Byte,
	Line,
	Column int
}

func (tree *Tree) ContentBytes() []byte {
	return tree.contentBytes
}

func (tree *Tree) NodeCount() int {
	return len(tree.nodes)
}

func (tree *Tree) RootNode() *Node {
	return tree.rootNode
}

func (tree *Tree) NodeFromSitter(sitterNode *sitter.Node) *Node {
	return tree.sitterToNode[sitterNode]
}

func (tree *Tree) Nodes() []Node {
	return tree.nodes
}

func (node *Node) Tree() *Tree {
	return node.tree
}

func (node *Node) SitterNode() *sitter.Node {
	return node.sitterNode
}

func (node *Node) Type() string {
	return node.tree.types[node.TypeID]
}

func (node *Node) FieldName() string {
	return node.fieldName
}

func (node *Node) IsNamed() bool {
	// FIXME: don't use the sitter node
	return node.sitterNode.IsNamed()
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) Content() string {
	return string(node.tree.contentBytes[node.ContentStart.Byte:node.ContentEnd.Byte])
}

func (node *Node) Debug() string {
	return fmt.Sprintf(
		"node-%d (%d:%d:%s)",
		node.ID,
		node.ContentStart.Line,
		node.ContentStart.Column,
		node.Type(),
	)
}

func (node *Node) Children() []*Node {
	return node.children
}

func (node *Node) IsMissing() bool {
	return node.sitterNode.IsMissing()
}

func (node *Node) IsError() bool {
	return node.sitterNode.IsError()
}

// FIXME: can we remove this?
func (node *Node) NamedChildren() []*Node {
	namedChildren := make([]*Node, 0, len(node.children))

	for _, child := range node.children {
		if child.IsNamed() {
			namedChildren = append(namedChildren, child)
		}
	}

	return namedChildren
}

func (node *Node) ChildByFieldName(name string) *Node {
	return node.childrenByField[name]
}

func (node *Node) DataflowSources() []*Node {
	return node.dataflowSources
}

func (node *Node) AliasOf() []*Node {
	return node.aliasOf
}

func (node *Node) ExpectedRules() []string {
	return node.expectedRules
}

func (node *Node) RuleDisabled(index int) bool {
	if node.disabledRuleIndices == nil {
		return false
	}

	return node.disabledRuleIndices.Test(uint(index))
}

func (node *Node) QueryResults(queryID int) []QueryResult {
	if node.queryResults == nil {
		return nil
	}

	results := node.queryResults[queryID]
	if len(results) == 0 {
		// Log available query IDs for debugging
		availableIDs := make([]int, 0, len(node.queryResults))
		for id := range node.queryResults {
			availableIDs = append(availableIDs, id)
		}
		log.Trace().Msgf("QueryResults: queryID=%d not found on node=%s, available IDs=%v", queryID, node.Type(), availableIDs)
	}
	return results
}

type nodeDump struct {
	Type            string
	ID              int
	Range           string
	Content         string     `yaml:",omitempty"`
	DataflowSources []int      `yaml:"dataflow_sources,omitempty"`
	AliasOf         []int      `yaml:"alias_of,omitempty"`
	Queries         []int      `yaml:",omitempty"`
	DisabledRules   []int      `yaml:",omitempty"`
	ExpectedRules   []string   `yaml:",omitempty"`
	Children        []nodeDump `yaml:",omitempty"`
}

func (node *Node) Dump() string {
	dump := node.dumpValue()
	yamlDump, err := yaml.Marshal(&dump)
	if err != nil {
		return err.Error()
	}

	return string(yamlDump)
}

func (node *Node) dumpValue() nodeDump {
	childDump := make([]nodeDump, len(node.children))
	for i, child := range node.children {
		childDump[i] = child.dumpValue()
	}

	queries := slices.Sorted(maps.Keys(node.queryResults))

	var disabledRules []int
	if node.disabledRuleIndices != nil {
		for i := 0; i < int(node.disabledRuleIndices.Len()); i++ {
			if node.disabledRuleIndices.Test(uint(i)) {
				disabledRules = append(disabledRules, i)
			}
		}
	}

	var expectedRules []string
	if len(node.expectedRules) > 0 {
		expectedRules = append(expectedRules, node.expectedRules...)
	}

	contentRange := fmt.Sprintf(
		"%d:%d - %d:%d",
		node.ContentStart.Line,
		node.ContentStart.Column,
		node.ContentEnd.Line,
		node.ContentEnd.Column,
	)

	content := ""
	if len(node.children) == 0 && node.Type()[0] != '"' {
		content = node.Content()
	}

	return nodeDump{
		Type:            node.Type(),
		ID:              node.ID,
		Range:           contentRange,
		Content:         content,
		DataflowSources: nodeListToID(node.dataflowSources),
		AliasOf:         nodeListToID(node.aliasOf),
		Children:        childDump,
		Queries:         queries,
		DisabledRules:   disabledRules,
		ExpectedRules:   expectedRules,
	}
}

func nodeListToID(nodes []*Node) []int {
	result := make([]int, len(nodes))

	for i, node := range nodes {
		result[i] = node.ID
	}

	slices.Sort(result)
	return result
}

// EachContentPart iterates over the content parts of a node, calling onText for literal text
// and onChild for child nodes. String fragment types are determined by the language configuration
// stored in the tree.
func (node *Node) EachContentPart(onText func(text string) error, onChild func(child *Node) error) error {
	start := node.ContentStart.Byte
	end := node.ContentEnd.Byte

	emit := func() error {
		if end <= start {
			return nil
		}

		return onText(string(node.tree.contentBytes[start:end]))
	}

	// Use the language-specific string fragment types from the tree
	fragmentTypes := node.tree.stringFragmentTypes

	// Create a map for O(1) lookup
	fragmentTypeMap := make(map[string]bool, len(fragmentTypes))
	for _, t := range fragmentTypes {
		fragmentTypeMap[t] = true
	}

	for _, child := range node.children {
		end = child.ContentStart.Byte
		if err := emit(); err != nil {
			return err
		}

		if child.IsNamed() {
			// Literal string content nodes should emit their text content, not be treated as children.
			// The fragment types are language-specific and configured via Language.StringFragmentTypes()
			childType := child.Type()
			if fragmentTypeMap[childType] {
				if err := onText(child.Content()); err != nil {
					return err
				}
			} else {
				if err := onChild(child); err != nil {
					return err
				}
			}
		}

		start = child.ContentEnd.Byte
	}

	end = node.ContentEnd.Byte
	if err := emit(); err != nil {
		return err
	}

	return nil
}

// FIXME: maybe users of this could work iteratively?
func (node *Node) Walk(visit func(node *Node, visitChildren func() error) error) error {
	visitChildren := func() error {
		for _, child := range node.Children() {
			if err := child.Walk(visit); err != nil {
				return err
			}
		}

		return nil
	}

	return visit(node, visitChildren)
}
