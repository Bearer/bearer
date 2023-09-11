package rulescanner

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/cache"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Scanner struct {
	ctx            context.Context
	detectorSet    detectorset.Set
	filename       string
	stats          *stats.FileStats
	traversalCache *traversalstrategy.Cache
	cache          *cache.Cache
}

func New(
	ctx context.Context,
	detectorSet detectorset.Set,
	filename string,
	stats *stats.FileStats,
	traversalCache *traversalstrategy.Cache,
	cache *cache.Cache,
) *Scanner {
	return &Scanner{
		ctx:            ctx,
		detectorSet:    detectorSet,
		filename:       filename,
		stats:          stats,
		traversalCache: traversalCache,
		cache:          cache,
	}
}

func (scanner *Scanner) Scan(
	rootNode *tree.Node,
	rule *ruleset.Rule,
	traversalStrategy traversalstrategy.Strategy,
) ([]*detectortypes.Detection, error) {
	if scanner.stats != nil {
		startTime := time.Now()
		defer scanner.stats.Rule(rule.ID(), startTime)
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan start at %s [%s]",
			rule.ID(),
			rootNode.Debug(),
			traversalStrategy.Scope(),
		)
	}

	var detections []*detectortypes.Detection
	if err := traversalStrategy.Traverse(scanner.traversalCache, rootNode, func(node *tree.Node) (bool, error) {
		if scanner.ctx.Err() != nil {
			return false, scanner.ctx.Err()
		}

		result, err := scanner.detectAtNode(rule, node)
		if result == nil || err != nil {
			return false, err
		}

		detections = append(detections, result.Detections...)
		return result.Sanitized, nil
	}); err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan end at %s [%s]: %d detections",
			rule.ID(),
			rootNode.Debug(),
			traversalStrategy.Scope(),
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *Scanner) Filename() string {
	return scanner.filename
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

		scanner.cache.Put(node, rule, nil)
		return nil, nil
	}

	result, err := scanner.detectorSet.DetectAt(node, rule, scanner)
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

func (scanner *Scanner) ruleDisabledForNode(rule *ruleset.Rule, node *tree.Node) bool {
	for current := node; current != nil; current = current.Parent() {
		if slices.Contains(node.DisabledRuleIndices(), rule.Index()) {
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
