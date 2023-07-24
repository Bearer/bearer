package tree

import (
	"errors"
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

type QuerySet struct {
	mu             sync.RWMutex
	sitterLanguage *sitter.Language
	queries        []Query
	queryByInput   map[string]*Query
	sitterCursor   *sitter.QueryCursor
	sitterQuery    *sitter.Query
}

type Query struct {
	querySet *QuerySet
	id       int
	input    string
}

type QueryResult map[string]*Node
type NodeResults map[NodeID][]QueryResult
type QuerySetResults map[int]NodeResults

func NewQuerySet(sitterLanguage *sitter.Language) *QuerySet {
	return &QuerySet{
		sitterLanguage: sitterLanguage,
		sitterCursor:   sitter.NewQueryCursor(),
		queryByInput:   make(map[string]*Query),
	}
}

func (querySet *QuerySet) Add(input string) *Query {
	querySet.mu.Lock()
	defer querySet.mu.Unlock()

	if query := querySet.queryByInput[input]; query != nil {
		return query
	}

	id := len(querySet.queries)
	querySet.queries = append(querySet.queries, Query{
		querySet: querySet,
		id:       id,
		input:    input,
	})

	querySet.freeSitterQuery()

	query := &querySet.queries[id]
	querySet.queryByInput[input] = query
	return query
}

func (querySet *QuerySet) Query(tree *Tree) (QuerySetResults, error) {
	querySet.mu.RLock()
	defer querySet.mu.RUnlock()

	if querySet.sitterQuery == nil {
		return nil, errors.New("query set has not been compiled")
	}

	results := querySet.newResults()
	querySet.sitterCursor.Exec(querySet.sitterQuery, tree.RootNode().sitterNode)

	for {
		match, found := querySet.sitterCursor.NextMatch()
		if !found {
			break
		}

		result := make(QueryResult)
		for _, capture := range match.Captures {
			result[querySet.sitterQuery.CaptureNameForId(capture.Index)] = tree.wrap(capture.Node)
		}

		resultRoot, rootExists := result["root"]
		if !rootExists {
			return nil, errors.New("missing @root capture in tree sitter query")
		}

		matchNode, matchNodeExists := result["match"]
		if !matchNodeExists {
			matchNode = resultRoot
		}

		results.add(int(match.PatternIndex), matchNode.ID(), result)
	}

	return results, nil
}

func (querySet *QuerySet) Compile() error {
	querySet.mu.Lock()
	defer querySet.mu.Unlock()

	if querySet.sitterQuery != nil {
		return nil
	}

	var s strings.Builder

	for _, query := range querySet.queries {
		s.WriteString(query.input)
		s.WriteString("\n")
	}

	sitterQuery, err := sitter.NewQuery([]byte(s.String()), querySet.sitterLanguage)
	if err != nil {
		return err
	}

	querySet.sitterQuery = sitterQuery

	return nil
}

func (querySet *QuerySet) Close() {
	querySet.sitterCursor.Close()
	querySet.freeSitterQuery()
}

func (queries *QuerySet) freeSitterQuery() {
	if queries.sitterQuery == nil {
		return
	}

	queries.sitterQuery.Close()
	queries.sitterQuery = nil
}

func (querySet *QuerySet) newResults() QuerySetResults {
	results := make(QuerySetResults)

	// make sure all queries are in the map so we don't re-trigger for queries with
	// no results
	for queryID := range querySet.queries {
		results[queryID] = nil
	}

	return results
}

func (query *Query) MatchAt(node *Node) ([]QueryResult, error) {
	inCache := false
	var nodeCache NodeResults
	if node.tree.queryCache != nil {
		nodeCache, inCache = node.tree.queryCache[query.id]
	}

	if !inCache {
		results, err := query.querySet.Query(node.tree)
		if err != nil {
			return nil, err
		}

		node.tree.queryCache = results
		nodeCache = results[query.id]
	}

	return nodeCache[node.ID()], nil
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

func (results QuerySetResults) add(queryID int, nodeID NodeID, result QueryResult) {
	nodeResults := results[queryID]
	if nodeResults == nil {
		nodeResults = make(NodeResults)
		results[queryID] = nodeResults
	}

	nodeResults[nodeID] = append(nodeResults[nodeID], result)
}
