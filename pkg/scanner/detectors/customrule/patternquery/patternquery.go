package patternquery

import (
	"fmt"
	"slices"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	astquery "github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/variableshape"

	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery/builder"
)

type Query interface {
	ID() string
	MatchAt(node *tree.Node) ([]*Result, error)
	MatchOnceAt(node *tree.Node) (*Result, error)
}

type query struct {
	id                   string
	input                string
	treeQuery            *astquery.Query
	paramToShapeVariable map[string]*variableshape.Variable
	equalParams          [][]string
	paramToContent       map[string]map[string]string
	variableShape        *variableshape.Shape
}

type rootVariableQuery struct {
	id            string
	variable      *language.PatternVariable
	shapeVariable *variableshape.Variable
	variableShape *variableshape.Shape
}

type Result struct {
	MatchNode *tree.Node
	Variables variableshape.Values
}

func Compile(
	language language.Language,
	querySet *astquery.Set,
	ruleID string,
	patternIndex int,
	input string,
	focusedVariable string,
	variableShape *variableshape.Shape,
) (Query, error) {
	builderResult, err := builder.Build(language, input, focusedVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %s", err)
	}

	id := fmt.Sprintf("%s[%d]", ruleID, patternIndex)

	if builderResult.RootVariable != nil {
		log.Trace().Msgf("single variable pattern %s: %s -> %#v", id, input, *builderResult.RootVariable)

		shapeVariable, err := variableShape.Variable(builderResult.RootVariable.Name)
		if err != nil {
			return nil, err
		}

		return &rootVariableQuery{
			id:            id,
			variable:      builderResult.RootVariable,
			shapeVariable: shapeVariable,
			variableShape: variableShape,
		}, nil
	}

	paramToShapeVariable := make(map[string]*variableshape.Variable)
	for param, variableName := range builderResult.ParamToVariable {
		shapeVariable, err := variableShape.Variable(variableName)
		if err != nil {
			return nil, err
		}

		paramToShapeVariable[param] = shapeVariable
	}

	query := &query{
		id:                   id,
		input:                input,
		treeQuery:            querySet.Add(builderResult.Query),
		paramToShapeVariable: paramToShapeVariable,
		equalParams:          builderResult.EqualParams,
		paramToContent:       builderResult.ParamToContent,
		variableShape:        variableShape,
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
	paramToVariableName := make(map[string]string)
	for param, variable := range query.paramToShapeVariable {
		paramToVariableName[param] = variable.Name()
	}

	yamlQuery, err := yaml.Marshal(&dumpValue{
		ID:              query.id,
		Pattern:         query.input,
		TreeQueryID:     query.treeQuery.ID(),
		ParamToVariable: paramToVariableName,
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

	variables := query.variableShape.NewValues()

	for paramName, node := range treeResult {
		if variable := query.paramToShapeVariable[paramName]; variable != nil {
			variables.Set(variable, node)
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
	variables := query.variableShape.NewValues()
	variables.Set(query.shapeVariable, node)

	return &Result{
		MatchNode: node,
		Variables: variables,
	}
}
