package rulescanner

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
)

type scanContext struct {
	cache   *cache.Cache
	scope   settings.RuleReferenceScope
	scanner *Scanner
}

func newContext(scanner *Scanner, cache *cache.Cache, scope settings.RuleReferenceScope) *scanContext {
	return &scanContext{
		cache:   cache,
		scope:   scope,
		scanner: scanner,
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
