package query

import (
	"context"
	"errors"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

type Set struct {
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

func NewSet(sitterLanguage *sitter.Language) *Set {
	return &Set{
		sitterLanguage: sitterLanguage,
		sitterCursor:   sitter.NewQueryCursor(),
		queryByInput:   make(map[string]*Query),
	}
}

func (querySet *Set) Add(input string) *Query {
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

func (querySet *Set) Query(ctx context.Context, builder *tree.Builder, rootNode *sitter.Node) error {
	if querySet.sitterQuery == nil {
		return errors.New("query set has not been compiled")
	}

	querySet.sitterCursor.Exec(querySet.sitterQuery, rootNode)

	captureNames := make(map[uint32]string)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		match, found := querySet.sitterCursor.NextMatch()
		if !found {
			break
		}

		result := make(map[string]*sitter.Node)
		for _, capture := range match.Captures {
			name := captureNames[capture.Index]
			if name == "" {
				name = querySet.sitterQuery.CaptureNameForId(capture.Index)
				captureNames[capture.Index] = name
			}

			result[name] = capture.Node
		}

		resultRoot, rootExists := result["root"]
		if !rootExists {
			return errors.New("missing @root capture in tree sitter query")
		}

		matchNode, matchNodeExists := result["match"]
		if !matchNodeExists {
			matchNode = resultRoot
		}

		builder.QueryResult(int(match.PatternIndex), matchNode, result)
	}

	return nil
}

func (querySet *Set) Compile() error {
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

func (query *Query) ID() int {
	return query.id
}

func (query *Query) MatchAt(node *tree.Node) []tree.QueryResult {
	return node.QueryResults(query.id)
}

func (query *Query) MatchOnceAt(node *tree.Node) (tree.QueryResult, error) {
	results := query.MatchAt(node)
	if len(results) > 1 {
		return nil, errors.New("query returned more than one result")
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}
