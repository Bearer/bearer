package rulescanner

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
)

type detectorContext struct {
	fileContext *filecontext.Context
	cache       *cache.Cache
	scope       settings.RuleReferenceScope
}

func newContext(fileContext *filecontext.Context, cache *cache.Cache, scope settings.RuleReferenceScope) *detectorContext {
	return &detectorContext{
		fileContext: fileContext,
		cache:       cache,
		scope:       scope,
	}
}

func (context detectorContext) Scan(
	rootNode *tree.Node,
	ruleID string,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && context.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	return Scan(context.fileContext, context.cache, effectiveScope, ruleID, rootNode)
}

func (context detectorContext) Filename() string {
	return context.fileContext.Filename()
}
