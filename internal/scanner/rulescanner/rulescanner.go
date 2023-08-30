package rulescanner

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
)

type Scanner struct {
	fileContext *filecontext.Context
	cache       *cache.Cache
	scope       settings.RuleReferenceScope
}

func New(fileContext *filecontext.Context, cache *cache.Cache, scope settings.RuleReferenceScope) *Scanner {
	return &Scanner{
		fileContext: fileContext,
		cache:       cache,
		scope:       scope,
	}
}

func (scanner *Scanner) Scan(
	rootNode *tree.Node,
	ruleID string,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && scanner.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	nodeScanner := &nodeScanner{
		fileContext: scanner.fileContext,
		ruleScanner: New(scanner.fileContext, scanner.cache, effectiveScope),
		ruleID:      ruleID,
		rootNode:    rootNode,
	}

	return nodeScanner.Scan()
}

func (scanner *Scanner) Filename() string {
	return scanner.fileContext.Filename()
}

func ScanTopLevelRule(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	tree *tree.Tree,
	ruleID string,
) ([]*detectortypes.Detection, error) {
	scanner := New(fileContext, cache, settings.NESTED_STRICT_SCOPE)
	return scanner.Scan(tree.RootNode(), ruleID, settings.NESTED_STRICT_SCOPE)
}
