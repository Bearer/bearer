package patternquery

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	astquery "github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/language"

	"github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
)

type Query interface {
	MatchAt(astContext *astquery.Context, node *sitter.Node) ([]*Result, error)
	MatchOnceAt(astContext *astquery.Context, node *sitter.Node) (*Result, error)
}

type query struct {
	treeQuery       *astquery.Query
	paramToVariable map[string]string
	equalParams     [][]string
	paramToContent  map[string]map[string]string
}

type rootVariableQuery struct {
	variable *language.PatternVariable
}

type Result struct {
	MatchNode *sitter.Node
	Variables astquery.Result
}

func Compile(
	language language.Language,
	querySet *astquery.Set,
	input string,
	focusedVariable string,
) (Query, error) {
	builderResult, err := builder.Build(language, input, focusedVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	if builderResult.RootVariable != nil {
		log.Trace().Msgf("single variable pattern %s -> %#v", input, *builderResult.RootVariable)
		return &rootVariableQuery{variable: builderResult.RootVariable}, nil
	}

	log.Trace().Msgf("compiled pattern %s -> %s", input, builderResult.Query)

	return &query{
		treeQuery:       querySet.Add(builderResult.Query),
		paramToVariable: builderResult.ParamToVariable,
		equalParams:     builderResult.EqualParams,
		paramToContent:  builderResult.ParamToContent,
	}, nil
}

func (query *query) MatchAt(astContext *astquery.Context, node *sitter.Node) ([]*Result, error) {
	treeResults, err := query.treeQuery.MatchAt(astContext, node)
	if err != nil {
		return nil, err
	}

	var results []*Result
	for _, treeResult := range treeResults {
		result := query.matchAndTranslateTreeResult(astContext, treeResult)
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (query *query) MatchOnceAt(astContext *astquery.Context, node *sitter.Node) (*Result, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(astContext, node)
	if err != nil {
		return nil, err
	}

	return query.matchAndTranslateTreeResult(astContext, treeResult), nil
}

func (query *query) matchAndTranslateTreeResult(astContext *astquery.Context, treeResult astquery.Result) *Result {
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

	variables := make(astquery.Result)

	for paramName, node := range treeResult {
		variableName := query.paramToVariable[paramName]
		if variableName != "" {
			variables[variableName] = node
		}
	}

	return &Result{
		MatchNode: treeResult["match"],
		Variables: variables,
	}
}

func (query *rootVariableQuery) MatchAt(astContext *astquery.Context, node *sitter.Node) ([]*Result, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return []*Result{query.resultFor(node)}, nil
}

func (query *rootVariableQuery) MatchOnceAt(astContext *astquery.Context, node *sitter.Node) (*Result, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return query.resultFor(node), nil
}

func (query *rootVariableQuery) isCompatibleType(node *sitter.Node) bool {
	if slices.Contains(query.variable.NodeTypes, "_") {
		return true
	}

	return slices.Contains(query.variable.NodeTypes, node.Type())
}

func (query *rootVariableQuery) resultFor(node *sitter.Node) *Result {
	variables := make(astquery.Result)
	variables[query.variable.Name] = node

	return &Result{
		MatchNode: node,
		Variables: variables,
	}
}
