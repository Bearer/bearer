package patternquery

import (
	"fmt"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery/builder"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/rs/zerolog/log"
)

type Query struct {
	treeQuery       *tree.Query
	paramToVariable map[string]string
	equalParams     [][]string
	paramToContent  map[string]map[string]string
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
		treeQuery:       treeQuery,
		paramToVariable: builderResult.ParamToVariable,
		equalParams:     builderResult.EqualParams,
		paramToContent:  builderResult.ParamToContent,
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

func (query *Query) Variables() []string {
	variables := set.New[string]()

	for _, variable := range query.paramToVariable {
		variables.Add(variable)
	}

	return variables.Items()
}

func (query *Query) Close() {
	query.treeQuery.Close()
}

func (query *Query) matchAndTranslateTreeResult(treeResult tree.QueryResult, rootNode *tree.Node) *languagetypes.PatternQueryResult {
	if treeResult == nil {
		return nil
	}

	for _, equalParams := range query.equalParams {
		var equalContent []string
		for _, equalParam := range equalParams {
			if node, exists := treeResult[equalParam]; exists {
				equalContent = append(equalContent, node.Content())
			}
		}

		if len(equalContent) < 2 {
			continue
		}

		value := equalContent[0]
		for _, content := range equalContent[1:] {
			if content != value {
				return nil
			}
		}
	}

	for param, typedContent := range query.paramToContent {
		node, exists := treeResult[param]
		if !exists {
			continue
		}

		if content, typeMatched := typedContent[node.Type()]; !typeMatched || node.Content() != content {
			return nil
		}
	}

	variables := make(tree.QueryResult)

	for paramName, node := range treeResult {
		variableName := query.paramToVariable[paramName]
		if variableName != "" {
			variables[variableName] = node
		}
	}

	return &languagetypes.PatternQueryResult{
		MatchNode: treeResult["match"],
		Variables: variables,
	}
}
