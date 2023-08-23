package query

import (
	"errors"
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

type Set struct {
	mu             sync.RWMutex
	sitterLanguage *sitter.Language
	queries        []Query
	queryByInput   map[string]*Query
	sitterCursor   *sitter.QueryCursor
	sitterQuery    *sitter.Query
}

type Query struct {
	querySet *Set
	id       int
	input    string
}

type Context struct {
	contentBytes []byte
	rootNode     *sitter.Node
	cache        SetResults
}

type Result map[string]*sitter.Node
type NodeResults map[*sitter.Node][]Result
type SetResults map[int]NodeResults

func NewSet(sitterLanguage *sitter.Language) *Set {
	return &Set{
		sitterLanguage: sitterLanguage,
		sitterCursor:   sitter.NewQueryCursor(),
		queryByInput:   make(map[string]*Query),
	}
}

func (querySet *Set) Add(input string) *Query {
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

func (querySet *Set) Query(rootNode *sitter.Node) (SetResults, error) {
	querySet.mu.RLock()
	defer querySet.mu.RUnlock()

	if querySet.sitterQuery == nil {
		return nil, errors.New("query set has not been compiled")
	}

	results := querySet.newResults()
	querySet.sitterCursor.Exec(querySet.sitterQuery, rootNode)

	for {
		match, found := querySet.sitterCursor.NextMatch()
		if !found {
			break
		}

		result := make(Result)
		for _, capture := range match.Captures {
			result[querySet.sitterQuery.CaptureNameForId(capture.Index)] = capture.Node
		}

		resultRoot, rootExists := result["root"]
		if !rootExists {
			return nil, errors.New("missing @root capture in tree sitter query")
		}

		matchNode, matchNodeExists := result["match"]
		if !matchNodeExists {
			matchNode = resultRoot
		}

		results.add(int(match.PatternIndex), matchNode, result)
	}

	return results, nil
}

func (querySet *Set) Compile() error {
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

func (querySet *Set) Close() {
	querySet.sitterCursor.Close()
	querySet.freeSitterQuery()
}

func (queries *Set) freeSitterQuery() {
	if queries.sitterQuery == nil {
		return
	}

	queries.sitterQuery.Close()
	queries.sitterQuery = nil
}

func (querySet *Set) newResults() SetResults {
	results := make(SetResults)

	// make sure all queries are in the map so we don't re-trigger for queries with
	// no results
	for queryID := range querySet.queries {
		results[queryID] = nil
	}

	return results
}

func (query *Query) MatchAt(context *Context, node *sitter.Node) ([]Result, error) {
	inCache := false
	var nodeCache NodeResults
	if context.cache != nil {
		nodeCache, inCache = context.cache[query.id]
	}

	if !inCache {
		results, err := query.querySet.Query(context.rootNode)
		if err != nil {
			return nil, err
		}

		context.cache = results
		nodeCache = results[query.id]
	}

	return nodeCache[node], nil
}

func (query *Query) MatchOnceAt(context *Context, node *sitter.Node) (Result, error) {
	results, err := query.MatchAt(context, node)
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

func (results SetResults) add(queryID int, node *sitter.Node, result Result) {
	nodeResults := results[queryID]
	if nodeResults == nil {
		nodeResults = make(NodeResults)
		results[queryID] = nodeResults
	}

	nodeResults[node] = append(nodeResults[node], result)
}

func NewContext(contentBytes []byte, rootNode *sitter.Node) *Context {
	return &Context{
		contentBytes: contentBytes,
		rootNode:     rootNode,
	}
}

func (context *Context) ContentFor(node *sitter.Node) string {
	return node.Content(context.contentBytes)
}
