package patternquery

import (
	"fmt"

	"github.com/bearer/curio/new/language/patternquery/builder"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/new/language/types"
)

type Query struct {
	treeQuery       *tree.Query
	paramToVariable map[string]string
	equalParams     [][]string
	paramToContent  map[string]string
}

func Compile(
	lang types.Language,
	anonymousParentTypes []string,
	input string,
	variables []builder.Variable,
	matchNodeOffset int,
) (*Query, error) {
	builderResult, err := builder.Build(lang, anonymousParentTypes, input, variables, matchNodeOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	treeQuery, err := lang.CompileQuery(builderResult.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %s, %s", err, builderResult.Query)
	}

	return &Query{
		treeQuery:       treeQuery,
		paramToVariable: builderResult.ParamToVariable,
		equalParams:     builderResult.EqualParams,
		paramToContent:  builderResult.ParamToContent,
	}, nil
}

func (query *Query) MatchAt(node *tree.Node) ([]*types.PatternQueryResult, error) {
	treeResults, err := query.treeQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	var results []*types.PatternQueryResult

	for _, treeResult := range treeResults {
		result := query.matchAndTranslateTreeResult(treeResult)
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (query *Query) MatchOnceAt(node *tree.Node) (*types.PatternQueryResult, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	return query.matchAndTranslateTreeResult(treeResult), nil
}

func (query *Query) Close() {
	query.treeQuery.Close()
}

func (query *Query) matchAndTranslateTreeResult(treeResult tree.QueryResult) *types.PatternQueryResult {
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

	return &types.PatternQueryResult{
		MatchNode: getMatchNode(treeResult["match"]),
		Variables: variables,
	}
}

func getMatchNode(node *tree.Node) *tree.Node {
	for {
		parent := node.Parent()
		if parent == nil || parent.StartByte() != node.StartByte() {
			return node
		}

		node = parent
	}
}
