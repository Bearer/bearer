package tree

import (
	"errors"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

var maxQueryID = 0

type Query struct {
	sitterQuery *sitter.Query
	id          int
	input       string
}

type QueryResult map[string]*Node

var mu sync.Mutex = sync.Mutex{}

func CompileQuery(sitterLanguage *sitter.Language, input string) (*Query, error) {
	sitterQuery, err := sitter.NewQuery([]byte(input), sitterLanguage)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	id := maxQueryID
	maxQueryID += 1
	defer mu.Unlock()

	return &Query{sitterQuery: sitterQuery, id: id, input: input}, nil
}

// Revisit if https://github.com/tree-sitter/tree-sitter/issues/1212 gets implemented
func (query *Query) MatchAt(node *Node) ([]QueryResult, error) {
	if _, inCache := node.tree.queryCache[query.id]; !inCache {
		results, err := query.resultsFor(node.tree)
		if err != nil {
			return nil, err
		}

		node.tree.queryCache[query.id] = results
	}

	return node.tree.queryCache[query.id][node.ID()], nil
}

func (query *Query) resultsFor(tree *Tree) (map[NodeID][]QueryResult, error) {
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(query.sitterQuery, tree.RootNode().sitterNode)

	nodeResults := make(map[NodeID][]QueryResult)

	for {
		match, found := cursor.NextMatch()
		if !found {
			break
		}

		result := make(QueryResult)
		for _, capture := range match.Captures {
			result[query.sitterQuery.CaptureNameForId(capture.Index)] = tree.wrap(capture.Node)
		}

		resultRoot, rootExists := result["root"]
		if !rootExists {
			return nil, errors.New("missing @root capture in tree sitter query")
		}

		matchNode, matchNodeExists := result["match"]
		if !matchNodeExists {
			matchNode = resultRoot
		}

		nodeResults[matchNode.ID()] = append(nodeResults[matchNode.ID()], result)
	}

	return nodeResults, nil
}

func (query *Query) MatchOnceAt(node *Node) (QueryResult, error) {
	results, err := query.MatchAt(node)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}
	if len(results) > 1 {
		return nil, errors.New("query returned more than one result")
	}

	return results[0], nil
}

func (query *Query) Close() {
	query.sitterQuery.Close()
}
