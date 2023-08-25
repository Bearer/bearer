package patternquery

import (
	"fmt"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"

	astquery "github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/language"

	"github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
)

type Query interface {
	ID() string
	MatchAt(node *tree.Node) ([]*Result, error)
	MatchOnceAt(node *tree.Node) (*Result, error)
}

type query struct {
	id              string
	input           string
	treeQuery       *astquery.Query
	paramToVariable map[string]string
	equalParams     [][]string
	paramToContent  map[string]map[string]string
}

type rootVariableQuery struct {
	id       string
	variable *language.PatternVariable
}

type Result struct {
	MatchNode *tree.Node
	Variables tree.QueryResult
}

func Compile(
	language language.Language,
	querySet *astquery.Set,
	ruleID string,
	patternIndex int,
	input string,
	focusedVariable string,
) (Query, error) {
	builderResult, err := builder.Build(language, input, focusedVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	id := fmt.Sprintf("%s[%d]", ruleID, patternIndex)

	if builderResult.RootVariable != nil {
		log.Trace().Msgf("single variable pattern %s: %s -> %#v", id, input, *builderResult.RootVariable)
		return &rootVariableQuery{id: id, variable: builderResult.RootVariable}, nil
	}

	query := &query{
		id:              id,
		input:           input,
		treeQuery:       querySet.Add(builderResult.Query),
		paramToVariable: builderResult.ParamToVariable,
		equalParams:     builderResult.EqualParams,
		paramToContent:  builderResult.ParamToContent,
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("compiled pattern:\n%s", query.dump())
	}

	return query, nil
}

type dumpValue struct {
	ID              string
	Pattern         string
	TreeQueryID     int                          `yaml:"tree_query_id"`
	ParamToVariable map[string]string            `yaml:"param_to_variable,omitempty"`
	ParamToContent  map[string]map[string]string `yaml:"param_to_content,omitempty"`
	EqualParams     [][]string                   `yaml:"equal_params,omitempty"`
}

func (query *query) dump() string {
	yamlQuery, err := yaml.Marshal(&dumpValue{
		ID:              query.id,
		Pattern:         query.input,
		TreeQueryID:     query.treeQuery.ID(),
		ParamToVariable: query.paramToVariable,
		ParamToContent:  query.paramToContent,
		EqualParams:     query.equalParams,
	})
	if err != nil {
		return err.Error()
	}

	return string(yamlQuery)
}

func (query *query) ID() string {
	return query.id
}

func (query *query) MatchAt(node *tree.Node) ([]*Result, error) {
	treeResults := query.treeQuery.MatchAt(node)

	var results []*Result
	for _, treeResult := range treeResults {
		if result := query.matchAndTranslateTreeResult(treeResult); result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (query *query) MatchOnceAt(node *tree.Node) (*Result, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	return query.matchAndTranslateTreeResult(treeResult), nil
}

func (query *query) matchAndTranslateTreeResult(treeResult tree.QueryResult) *Result {
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

	return &Result{
		MatchNode: treeResult["match"],
		Variables: variables,
	}
}

func (query *rootVariableQuery) ID() string {
	return query.id
}

func (query *rootVariableQuery) MatchAt(node *tree.Node) ([]*Result, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return []*Result{query.resultFor(node)}, nil
}

func (query *rootVariableQuery) MatchOnceAt(node *tree.Node) (*Result, error) {
	if !query.isCompatibleType(node) {
		return nil, nil
	}

	return query.resultFor(node), nil
}

func (query *rootVariableQuery) isCompatibleType(node *tree.Node) bool {
	if slices.Contains(query.variable.NodeTypes, "_") {
		return true
	}

	return slices.Contains(query.variable.NodeTypes, node.Type())
}

func (query *rootVariableQuery) resultFor(node *tree.Node) *Result {
	variables := make(tree.QueryResult)
	variables[query.variable.Name] = node

	return &Result{
		MatchNode: node,
		Variables: variables,
	}
}
