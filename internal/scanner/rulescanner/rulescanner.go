package rulescanner

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

type Scanner struct {
	fileContext *filecontext.Context
	cache       *cache.Cache
	scope       settings.RuleReferenceScope
}

func New(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	scope settings.RuleReferenceScope,
) *Scanner {
	return &Scanner{
		fileContext: fileContext,
		cache:       cache,
		scope:       scope,
	}
}

func (scanner *Scanner) Scan(
	rootNode *tree.Node,
	rule *ruleset.Rule,
	scope settings.RuleReferenceScope,
) ([]*detectortypes.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && scanner.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	nextScanner := scanner
	if nextScanner.scope != effectiveScope {
		nextScanner = New(scanner.fileContext, scanner.cache, effectiveScope)
	}

	return nextScanner.scan(rule, rootNode)
}

func (scanner *Scanner) scan(rule *ruleset.Rule, rootNode *tree.Node) ([]*detectortypes.Detection, error) {
	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan start at %s [%s]",
			rule.ID(),
			rootNode.Debug(),
			scanner.scope,
		)
	}

	var detections []*detectortypes.Detection
	traversalStrategy, err := traversalstrategy.Get(scanner.scope)
	if err != nil {
		return nil, err
	}

	if err := traversalStrategy.Traverse(rootNode, func(node *tree.Node) (bool, error) {
		result, err := scanner.detectAtNode(rule, node)
		detections = append(detections, result.Detections...)
		return result.Sanitized, err
	}); err != nil {
		return nil, err
	}

	scanner.fileContext.RuleStats(rule, startTime)

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan end at %s [%s]: %d detections",
			rule.ID(),
			rootNode.Debug(),
			scanner.scope,
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *Scanner) Filename() string {
	return scanner.fileContext.Filename()
}

func (scanner *Scanner) detectAtNode(rule *ruleset.Rule, node *tree.Node) (*detectorset.Result, error) {
	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", rule.ID(), node.Debug())
	}

	if result, cached := scanner.cache.Get(node, rule); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: %s (cached)",
				rule.ID(),
				node.Debug(),
				traceResultText(result),
			)
		}

		return result, nil
	}

	if scanner.ruleDisabledForNode(rule, node) {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: rule disabled",
				rule.ID(),
				node.Debug(),
			)
		}

		scanner.cache.Put(node, rule, &detectorset.Result{})
		return nil, nil
	}

	result, err := scanner.detectWithoutCycles(rule, node)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %s: %s",
			rule.ID(),
			node.Debug(),
			traceResultText(result),
		)
	}

	scanner.cache.Put(node, rule, result)
	return result, nil
}

func (scanner *Scanner) detectWithoutCycles(rule *ruleset.Rule, node *tree.Node) (*detectorset.Result, error) {
	if slices.Contains(node.ExecutingDetectors, rule.Index()) {
		executingRules := make([]string, len(node.ExecutingDetectors))
		for i, ruleIndex := range node.ExecutingDetectors {
			executingRules[i] = scanner.fileContext.Rule(ruleIndex).ID()
		}

		return nil, fmt.Errorf(
			"cycle found during rule evaluation at %s: [%s > %s]",
			node.Debug(),
			strings.Join(executingRules, " > "),
			rule.ID(),
		)
	}

	node.ExecutingDetectors = append(node.ExecutingDetectors, rule.Index())
	result, err := scanner.fileContext.DetectAt(node, rule, scanner)
	node.ExecutingDetectors = node.ExecutingDetectors[:len(node.ExecutingDetectors)-1]

	return result, err
}

func (scanner *Scanner) ruleDisabledForNode(rule *ruleset.Rule, node *tree.Node) bool {
	for current := node; current != nil; current = current.Parent() {
		if slices.Contains(current.DisabledRuleIndices(), rule.Index()) {
			return true
		}
	}

	return false
}

func traceResultText(result *detectorset.Result) string {
	if result.Sanitized {
		return "sanitized"
	}

	return fmt.Sprintf("%d detections", len(result.Detections))
}

func ScanTopLevelRule(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	tree *tree.Tree,
	rule *ruleset.Rule,
) ([]*detectortypes.Detection, error) {
	context := New(fileContext, cache, settings.NESTED_STRICT_SCOPE)
	return context.Scan(tree.RootNode(), rule, settings.NESTED_STRICT_SCOPE)
}
