package query

import (
	"fmt"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/ast/languages/ruby"
	"github.com/bearer/bearer/pkg/souffle/relationtypes"
	"github.com/bearer/bearer/pkg/util/souffle"
	relationwriter "github.com/bearer/bearer/pkg/util/souffle/writer/relation"
	sitter "github.com/smacker/go-tree-sitter"
)

type QueryContext struct {
	souffle            *souffle.Souffle
	language           *ruby.Language
	langImplementation implementation.Implementation
	cache              map[string]map[uint32][]*languagetypes.PatternQueryResult
	nodeIdGenerator    *idgenerator.NodeIdGenerator
	tree               *tree.Tree
}

func NewContext(
	souffle *souffle.Souffle,
	language *ruby.Language,
	langImplementation implementation.Implementation,
	tree *tree.Tree,
	input []byte,
) (*QueryContext, error) {
	context := &QueryContext{
		souffle:            souffle,
		language:           language,
		langImplementation: langImplementation,
		cache:              make(map[string]map[uint32][]*languagetypes.PatternQueryResult),
		nodeIdGenerator:    idgenerator.NewNodeIdGenerator(),
		tree:               tree,
	}

	return context, context.run(tree.RootNode().SitterNode(), input)
}

func (context *QueryContext) MatchAt(ruleName string, patternIndex int, node *tree.Node) []*languagetypes.PatternQueryResult {
	nodeCache, ruleExists := context.cache[patternKey(ruleName, patternIndex)]
	if !ruleExists {
		return nil
	}

	return nodeCache[context.nodeIdGenerator.Get(node.SitterNode())]
}

func (context *QueryContext) put(nodeId uint32, ruleName string, patternIndex int, result *languagetypes.PatternQueryResult) {
	nodeCache, ruleExists := context.cache[patternKey(ruleName, patternIndex)]
	if !ruleExists {
		nodeCache = make(map[uint32][]*languagetypes.PatternQueryResult)
		context.cache[ruleName] = nodeCache
	}

	nodeCache[nodeId] = append(nodeCache[nodeId], result)
}

func (context *QueryContext) run(rootNode *sitter.Node, input []byte) error {
	// FIXME: this should go and we should write structs from sourcefacts
	writer := relationwriter.New(context.souffle.Program())

	if err := context.language.WriteSourceFacts(input, rootNode, context.nodeIdGenerator, writer); err != nil {
		return fmt.Errorf("fact generation error: %w", err)
	}

	context.souffle.Run()

	return context.readMatches()
}

func (context *QueryContext) readMatches() error {
	relation, err := context.souffle.Relation("Rule_Match")
	if err != nil {
		return err
	}

	iterator := relation.NewIterator()
	defer iterator.Close()

	for i := 0; iterator.HasNext(); i++ {
		var match relationtypes.Rule_Match
		if err := context.souffle.Unmarshal(&match, iterator.GetNext()); err != nil {
			return fmt.Errorf("failed to read tuple %d for Rule_Match: %w", i, err)
		}

		context.put(match.Node, match.RuleName, int(match.PatternIndex), &languagetypes.PatternQueryResult{
			MatchNode: context.tree.Wrap(context.nodeIdGenerator.InverseLookup(match.Node)),
			Variables: make(map[string]*tree.Node),
		})
	}

	return nil
}

func patternKey(ruleName string, patternIndex int) string {
	return fmt.Sprintf("%s:%d", ruleName, patternIndex)
}
