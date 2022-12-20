package patternquery

import (
	"log"

	"github.com/bearer/curio/new/language/patternquery/builder"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/new/language/types"
)

type Query struct {
	treeQuery       *tree.Query
	paramToVariable map[string]string
}

func Compile(lang types.Language, input string, variables []builder.Variable) (*Query, error) {
	builderResult, err := builder.Build(lang, input, variables)
	if err != nil {
		return nil, err
	}

	log.Printf("translated query from:\n\n%s\n\nto:\n\n%#v\n\n", input, builderResult)

	treeQuery, err := lang.CompileQuery(builderResult.Query)
	if err != nil {
		return nil, err
	}

	return &Query{treeQuery: treeQuery, paramToVariable: builderResult.ParamToVariable}, nil
}

func (query *Query) MatchAt(node *tree.Node) ([]tree.QueryResult, error) {
	treeResults, err := query.treeQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	results := make([]tree.QueryResult, len(treeResults))

	for i, treeResult := range treeResults {
		results[i] = query.translateTreeResult(treeResult)
	}

	return results, nil
}

func (query *Query) MatchOnceAt(node *tree.Node) (tree.QueryResult, error) {
	treeResult, err := query.treeQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	return query.translateTreeResult(treeResult), nil
}

func (query *Query) Close() {
	query.treeQuery.Close()
}

func (query *Query) translateTreeResult(treeResult tree.QueryResult) tree.QueryResult {
	result := make(tree.QueryResult)

	for paramName, node := range treeResult {
		result[query.paramToVariable[paramName]] = node
	}

	return result
}
