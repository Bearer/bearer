package query

import (
	"fmt"
	"reflect"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/ast/languages/ruby"
	"github.com/bearer/bearer/pkg/util/souffle"
	relationwriter "github.com/bearer/bearer/pkg/util/souffle/writer/relation"
)

type QueryContext struct {
	souffle            *souffle.Souffle
	language           *ruby.Language
	langImplementation implementation.Implementation
	cache              map[string]map[uint32][]*languagetypes.PatternQueryResult
	nodeIdGenerator    *idgenerator.NodeIdGenerator
	tree               *tree.Tree
	patternVariables   map[string][]string
	patternTypes       map[string]reflect.Type
}

func NewContext(
	souffle *souffle.Souffle,
	language *ruby.Language,
	langImplementation implementation.Implementation,
	tree *tree.Tree,
	input []byte,
	patternVariables map[string][]string,
) (*QueryContext, error) {
	context := &QueryContext{
		souffle:            souffle,
		language:           language,
		langImplementation: langImplementation,
		cache:              make(map[string]map[uint32][]*languagetypes.PatternQueryResult),
		nodeIdGenerator:    idgenerator.NewNodeIdGenerator(),
		tree:               tree,
		patternVariables:   patternVariables,
		patternTypes:       makePatternTypes(patternVariables),
	}

	return context, context.run(tree.RootNode().SitterNode(), input)
}

func (context *QueryContext) MatchAt(patternId string, node *tree.Node) []*languagetypes.PatternQueryResult {
	nodeCache, ruleExists := context.cache[patternId]
	if !ruleExists {
		return nil
	}

	r := nodeCache[context.nodeIdGenerator.Get(node.SitterNode())]

	return r
}

func (context *QueryContext) put(nodeId uint32, patternId string, result *languagetypes.PatternQueryResult) {
	// log.Error().Msgf("putting result for %d: %#v", patternId, result)

	nodeCache, ruleExists := context.cache[patternId]
	if !ruleExists {
		nodeCache = make(map[uint32][]*languagetypes.PatternQueryResult)
		context.cache[patternId] = nodeCache
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

	if err := context.readMatches(); err != nil {
		return fmt.Errorf("error reading matches: %w", err)
	}

	return nil
}

func (context *QueryContext) readMatches() error {
	for patternId, typ := range context.patternTypes {
		// FIXME: Define dynamic names in a common place
		relation, err := context.souffle.Relation(fmt.Sprintf("Rule_Match_%s", patternId))
		if err != nil {
			// FIXME: need to support all rules
			// log.Error().Msgf("pattern relation error: %w", err)
			continue
			// return err
		}

		iterator := relation.NewIterator()
		defer iterator.Close()

		// log.Error().Msgf("output count: %d", relation.Size())

		variableNames := context.patternVariables[patternId]

		for i := 0; iterator.HasNext(); i++ {
			match := reflect.New(typ)
			matchValue := reflect.Indirect(match)
			if err := context.souffle.Unmarshal(match.Interface(), iterator.GetNext()); err != nil {
				return fmt.Errorf("failed to read tuple %d: %w", i, err)
			}

			matchNodeId := uint32(matchValue.Field(0).Uint())
			matchNode := context.resultNode(matchNodeId)

			variables := make(map[string]*tree.Node)

			for i := 1; i < typ.NumField(); i++ {
				variables[variableNames[i-1]] = context.resultNode(uint32(matchValue.Field(i).Uint()))
			}

			context.put(matchNodeId, patternId, &languagetypes.PatternQueryResult{
				MatchNode: matchNode,
				Variables: variables,
			})
		}
	}

	return nil
}

func (context *QueryContext) resultNode(nodeId uint32) *tree.Node {
	return context.tree.Wrap(context.nodeIdGenerator.InverseLookup(nodeId))
}

func makePatternTypes(patternVariables map[string][]string) map[string]reflect.Type {
	uint32_t := reflect.TypeOf(uint32(0))
	result := make(map[string]reflect.Type)

	for patternId, variables := range patternVariables {
		fields := []reflect.StructField{{Name: "Match", Type: uint32_t}}

		for i := range variables {
			fields = append(fields, reflect.StructField{Name: fmt.Sprintf("Variable%d", i), Type: uint32_t})
		}

		result[patternId] = reflect.StructOf(fields)
	}

	return result
}
