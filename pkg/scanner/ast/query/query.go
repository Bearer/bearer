package query

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
)

type Set struct {
	languageID     string
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

func NewSet(languageID string, sitterLanguage *sitter.Language) *Set {
	return &Set{
		languageID:     languageID,
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

		if log.Trace().Enabled() {
			log.Trace().Msgf("storing query result: PatternIndex=%d (querySet has %d queries), matchNode=%s, root=%s", match.PatternIndex, len(querySet.queries), matchNode.Type(), resultRoot.Type())
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

	if log.Trace().Enabled() {
		log.Trace().Msgf("%s queries:\n%s", querySet.languageID, querySet.dump())
	}

	sitterQuery, err := sitter.NewQuery([]byte(s.String()), querySet.sitterLanguage)
	if err != nil {
		return fmt.Errorf("%w\n\n%s", err, s.String())
	}

	querySet.sitterQuery = sitterQuery

	return nil
}

type dumpValue struct {
	ID    int
	Input string
}

func (querySet *Set) dump() string {
	queries := make([]dumpValue, len(querySet.queries))

	for i, query := range querySet.queries {
		queries[i].ID = query.id
		queries[i].Input = query.input
	}

	yamlQueries, err := yaml.Marshal(queries)
	if err != nil {
		return err.Error()
	}

	return string(yamlQueries)
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
	if log.Trace().Enabled() {
		log.Trace().Msgf("MatchAt: query.id=%d, node id=%d type=%s content=%q", query.id, node.ID, node.Type(), node.Content())
	}
	results := node.QueryResults(query.id)
	if log.Trace().Enabled() {
		if len(results) == 0 {
			log.Trace().Msgf("MatchAt: query.id=%d, node id=%d, found=0 results", query.id, node.ID)
		} else {
			log.Trace().Msgf("MatchAt: query.id=%d, node id=%d, found=%d results", query.id, node.ID, len(results))
		}
	}
	return results
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
