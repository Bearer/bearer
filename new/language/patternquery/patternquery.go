package patternquery

import (
	"fmt"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery/builder"
	"github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type Query struct {
	treeQuery       *tree.Query
	paramToVariable map[string]string
	equalParams     [][]string
	paramToContent  map[string]map[string]string
}

type RootVariableQuery struct {
	variable *types.Variable
}

func Compile(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	querySet *tree.QuerySet,
	input string,
	focusedVariable string,
) (types.PatternQuery, error) {
	builderResult, err := builder.Build(lang, langImplementation, input, focusedVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	if builderResult.RootVariable != nil {
		log.Trace().Msgf("single variable pattern %s -> %#v", input, *builderResult.RootVariable)
		return &RootVariableQuery{variable: builderResult.RootVariable}, nil
	}

	log.Trace().Msgf("compiled pattern %s -> %s", input, builderResult.Query)

	return &Query{
		treeQuery:       querySet.Add(builderResult.Query),
		paramToVariable: builderResult.ParamToVariable,
		equalParams:     builderResult.EqualParams,
		paramToContent:  builderResult.ParamToContent,
	}, nil
}

func (query *Query) MatchAt(astContext *tree.QueryContext, node *sitter.Node) ([]*languagetypes.PatternQueryResult, error) {
	treeResults, err := query.treeQuery.MatchAt(astContext, node)
	if err != nil {
		return nil, err
	}

	var results []*languagetypes.PatternQueryResult

	for _, treeResult := range treeResults {
		result := query.matchAndTranslateTreeResult(astContext, treeResult)
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (query *Query) MatchOnceAt(astContext *tree.QueryContext, node *sitter.Node) (*languagetypes.PatternQueryResult, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(astContext, node)
	if err != nil {
		return nil, err
	}

	return query.matchAndTranslateTreeResult(astContext, treeResult), nil
}

func (query *Query) matchAndTranslateTreeResult(astContext *tree.QueryContext, treeResult tree.QueryResult) *languagetypes.PatternQueryResult {
	if treeResult == nil {
		return nil
	}

	for _, equalParams := range query.equalParams {
		var equalContent []string
		for _, equalParam := range equalParams {
			if node, exists := treeResult[equalParam]; exists {
				equalContent = append(equalContent, astContext.ContentFor(node))
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

		if content, typeMatched := typedContent[node.Type()]; !typeMatched || astContext.ContentFor(node) != content {
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

func (query *RootVariableQuery) MatchAt(astContext *tree.QueryContext, node *sitter.Node) ([]*languagetypes.PatternQueryResult, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return []*languagetypes.PatternQueryResult{query.resultFor(node)}, nil
}

func (query *RootVariableQuery) MatchOnceAt(astContext *tree.QueryContext, node *sitter.Node) (*languagetypes.PatternQueryResult, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return query.resultFor(node), nil
}

func (query *RootVariableQuery) isCompatibleType(node *sitter.Node) bool {
	if slices.Contains(query.variable.NodeTypes, "_") {
		return true
	}

	return slices.Contains(query.variable.NodeTypes, node.Type())
}

func (query *RootVariableQuery) resultFor(node *sitter.Node) *languagetypes.PatternQueryResult {
	variables := make(tree.QueryResult)
	variables[query.variable.Name] = node

	return &languagetypes.PatternQueryResult{
		MatchNode: node,
		Variables: variables,
	}
}

func (query *RootVariableQuery) Close() {}
