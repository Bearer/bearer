package patternquery

import (
	"fmt"

	"github.com/bearer/curio/new/language/implementation"
	"github.com/bearer/curio/new/language/patternquery/builder"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Query struct {
	langImplementation implementation.Implementation
	treeQuery          *tree.Query
	paramToVariable    map[string]string
	equalParams        [][]string
	paramToContent     map[string]string
}

func Compile(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	input string,
) (*Query, error) {
	builderResult, err := builder.Build(lang, langImplementation, input)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	treeQuery, err := lang.CompileQuery(builderResult.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %s, %s -> %s", err, input, builderResult.Query)
	}

	log.Debug().Msgf("compiled pattern %s -> %s", input, builderResult.Query)

	return &Query{
		langImplementation: langImplementation,
		treeQuery:          treeQuery,
		paramToVariable:    builderResult.ParamToVariable,
		equalParams:        builderResult.EqualParams,
		paramToContent:     builderResult.ParamToContent,
	}, nil
}

func (query *Query) MatchAt(node *tree.Node) ([]*languagetypes.PatternQueryResult, error) {
	treeResults, err := query.treeQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	var results []*languagetypes.PatternQueryResult

	for _, treeResult := range treeResults {
		result := query.matchAndTranslateTreeResult(treeResult, node)
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (query *Query) MatchOnceAt(node *tree.Node) (*languagetypes.PatternQueryResult, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	return query.matchAndTranslateTreeResult(treeResult, node), nil
}

func (query *Query) Close() {
	query.treeQuery.Close()
}

func (query *Query) matchAndTranslateTreeResult(treeResult tree.QueryResult, rootNode *tree.Node) *languagetypes.PatternQueryResult {
	if treeResult == nil {
		return nil
	}

	for _, equalParams := range query.equalParams {
		value := treeResult[equalParams[0]].Content()

		for _, param := range equalParams[1:] {
			if treeResult[param].Content() != value {
				return nil
			}
		}
	}

	for param, content := range query.paramToContent {
		if treeResult[param].Content() != content {
			return nil
		}
	}

	variables := make(tree.QueryResult)

	for paramName, node := range treeResult {
		variables[query.paramToVariable[paramName]] = node
	}

	return &languagetypes.PatternQueryResult{
		MatchNode: query.getMatchNode(treeResult["match"], rootNode),
		Variables: variables,
	}
}

func (query *Query) getMatchNode(node *tree.Node, rootNode *tree.Node) *tree.Node {
	for {
		parent := node.Parent()
		if parent == nil ||
			parent.StartByte() != node.StartByte() ||
			slices.Contains(query.langImplementation.PatternMatchNodeContainerTypes(), parent.Type()) ||
			node.Equal(rootNode) {
			return node
		}

		node = parent
	}
}
