package rulescanner

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
)

type scanContext struct {
	cache        *cache.Cache
	queryContext *query.Context
	scope        settings.RuleReferenceScope
	scanner      *Scanner
}

func newContext(scanner *Scanner, cache *cache.Cache, scope settings.RuleReferenceScope) *scanContext {
	return &scanContext{
		cache:        cache,
		queryContext: scanner.queryContext,
		scope:        scope,
		scanner:      scanner,
	}
}

func (context scanContext) Scan(
	rootNode *tree.Node,
	ruleID,
	sanitizerRuleID string,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && context.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	return Scan(
		context.scanner.ctx,
		context.scanner.detectorSet,
		context.scanner.fileName,
		context.scanner.fileStats,
		context.scanner.tree,
		context.scanner.queryContext,
		context.scanner.rulesDisabledForNodes,
		rootNode,
		context.scanner.context.cache,
		effectiveScope,
		ruleID,
		sanitizerRuleID,
	)
}

func (context scanContext) FileName() string {
	return context.scanner.fileName
}

func (context scanContext) QueryContext() *query.Context {
	return context.queryContext
}

func (context scanContext) NodeFromSitter(sitterNode *sitter.Node) *tree.Node {
	return context.scanner.tree.NodeFromSitter(sitterNode)
}

func (context scanContext) QueryMatchAt(query *query.Query, node *tree.Node) ([]detectortypes.QueryResult, error) {
	sitterResults, err := query.MatchAt(context.queryContext, node.SitterNode())
	if err != nil {
		return nil, err
	}

	results := make([]detectortypes.QueryResult, len(sitterResults))

	for i, sitterResult := range sitterResults {
		results[i], err = context.translateResult(sitterResult)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (context scanContext) QueryMatchOnceAt(query *query.Query, node *tree.Node) (detectortypes.QueryResult, error) {
	sitterResult, err := query.MatchOnceAt(context.queryContext, node.SitterNode())
	if err != nil {
		return nil, err
	}

	return context.translateResult(sitterResult)
}

// FIXME: try and remove the translation by caching query results on the ast tree
func (context scanContext) translateResult(sitterResult query.Result) (detectortypes.QueryResult, error) {
	if sitterResult == nil {
		return nil, nil
	}

	result := make(map[string]*tree.Node)

	for name, sitterNode := range sitterResult {
		node := context.NodeFromSitter(sitterNode)
		if node == nil {
			return nil, fmt.Errorf(
				"missing node for sitter node %d:%d:\n%s\n%s",
				sitterNode.StartPoint().Row+1,
				sitterNode.StartPoint().Column+1,
				sitterNode.String(),
				context.QueryContext().ContentFor(sitterNode),
			)
		}

		result[name] = node
	}

	return result, nil
}
