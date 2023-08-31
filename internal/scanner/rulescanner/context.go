package rulescanner

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
)

type context struct {
	fileContext *filecontext.Context
	cache       *cache.Cache
	scope       settings.RuleReferenceScope
}

func newContext(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	scope settings.RuleReferenceScope,
) *context {
	return &context{
		fileContext: fileContext,
		cache:       cache,
		scope:       scope,
	}
}

func (context *context) Scan(
	rootNode *tree.Node,
	detectorID int,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && context.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	ruleScanner := &scanner{
		fileContext: context.fileContext,
		context:     newContext(context.fileContext, context.cache, effectiveScope),
		detectorID:  detectorID,
		rootNode:    rootNode,
	}

	return ruleScanner.Scan()
}

// FIXME: remove this
func (context *context) ScanRule(
	rootNode *tree.Node,
	ruleID string,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	return context.Scan(rootNode, context.fileContext.DetectorIDFor(ruleID), scope)
}

func (context *context) Filename() string {
	return context.fileContext.Filename()
}

func ScanTopLevelRule(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	tree *tree.Tree,
	detectorID int,
) ([]*detectortypes.Detection, error) {
	context := newContext(fileContext, cache, settings.NESTED_STRICT_SCOPE)
	return context.Scan(tree.RootNode(), detectorID, settings.NESTED_STRICT_SCOPE)
}
