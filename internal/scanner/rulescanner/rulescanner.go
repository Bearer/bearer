package rulescanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"

	"github.com/bearer/bearer/internal/scanner/cache"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Scanner struct {
	ctx                   context.Context
	tree                  *tree.Tree
	detectorSet           detectorset.Set
	fileStats             *stats.FileStats
	fileName              string
	rulesDisabledForNodes map[string][]*tree.Node
	queryContext          *query.Context
	rootNode              *tree.Node
	context,
	sanitizerContext *scanContext
	ruleID,
	sanitizerRuleID string
}

func new(
	ctx context.Context,
	detectorSet detectorset.Set,
	fileName string,
	fileStats *stats.FileStats,
	tree *tree.Tree,
	queryContext *query.Context,
	rulesDisabledForNodes map[string][]*tree.Node,
	rootNode *tree.Node,
	cache *cache.Cache,
	scope settings.RuleReferenceScope,
	ruleID,
	sanitizerRuleID string,
) *Scanner {
	scanner := &Scanner{
		ctx:                   ctx,
		tree:                  tree,
		fileName:              fileName,
		detectorSet:           detectorSet,
		fileStats:             fileStats,
		rulesDisabledForNodes: rulesDisabledForNodes,
		queryContext:          queryContext,
		rootNode:              rootNode,
		ruleID:                ruleID,
		sanitizerRuleID:       sanitizerRuleID,
	}

	scanner.context = newContext(scanner, cache, scope)
	scanner.sanitizerContext = newContext(scanner, cache, settings.DefaultScope)

	return scanner
}

func (scanner *Scanner) Scan() ([]*detectortypes.Detection, error) {
	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf("file scan start: %s at %s", scanner.ruleID, scanner.rootNode.Debug(true))
	}

	var detections []*detectortypes.Detection
	var err error

	switch scanner.context.scope {
	case settings.NESTED_SCOPE:
		detections, err = scanner.scanDescendants()
	case settings.RESULT_SCOPE:
		detections, err = scanner.scanDataflowAndAliases()
	case settings.CURSOR_SCOPE:
		detections, err = scanner.scanAliases()
	case settings.CURSOR_STRICT_SCOPE:
		detections, _, err = scanner.sanitizedNodeDetections(scanner.rootNode)
	}

	if err != nil {
		return nil, err
	}

	scanner.fileStats.Rule(scanner.ruleID, startTime)

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"file scan end: %s at %s: %d detections",
			scanner.ruleID,
			scanner.rootNode.Debug(false),
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *Scanner) scanDescendants() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.Children()
	})
}

func (scanner *Scanner) scanDataflowAndAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return append(node.DataflowSources(), node.AliasOf()...)
	})
}
func (scanner *Scanner) scanAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.AliasOf()
	})
}

func (scanner *Scanner) detectWithNext(
	getNext func(node *tree.Node) []*tree.Node,
) ([]*detectortypes.Detection, error) {
	nodes := []*tree.Node{scanner.rootNode}

	var detections []*detectortypes.Detection

	for {
		if len(nodes) == 0 {
			break
		}

		var next []*tree.Node

		for _, node := range nodes {
			nodeDetections, sanitized, err := scanner.sanitizedNodeDetections(node)
			if err != nil {
				return nil, err
			}
			if sanitized {
				continue
			}

			detections = append(detections, nodeDetections...)
			next = append(next, getNext(node)...)
		}

		nodes = next
	}

	return detections, nil
}

func (scanner *Scanner) ruleDisabledForNode(ruleId string, node *tree.Node) bool {
	nodesToIgnore := scanner.rulesDisabledForNodes[ruleId]
	if nodesToIgnore == nil {
		return false
	}

	// check node
	for _, ignoredNode := range nodesToIgnore {
		if ignoredNode == node {
			return true
		}
	}

	// check node ancestors
	parent := node.Parent()
	for parent != nil {
		for _, ignoredNode := range nodesToIgnore {
			if ignoredNode == parent {
				return true
			}
		}

		parent = parent.Parent()
	}

	return false
}

func (scanner *Scanner) sanitizedNodeDetections(node *tree.Node) ([]*detectortypes.Detection, bool, error) {
	if scanner.sanitizerRuleID != "" {
		sanitizerDetections, err := scanner.detectAtNode(scanner.sanitizerContext, node, scanner.sanitizerRuleID)
		if len(sanitizerDetections) != 0 || err != nil {
			return nil, true, err
		}
	}

	detections, err := scanner.detectAtNode(scanner.context, node, scanner.ruleID)
	return detections, false, err
}

func (scanner *Scanner) detectAtNode(
	context *scanContext,
	node *tree.Node,
	ruleID string,
) ([]*detectortypes.Detection, error) {
	if scanner.ctx.Err() != nil {
		return nil, scanner.ctx.Err()
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", ruleID, node.Debug(true))
	}

	if detections, cached := context.cache.Get(node, ruleID); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: %d detections (cached)",
				ruleID,
				node.Debug(false),
				len(detections),
			)
		}

		return detections, nil
	}

	if scanner.ruleDisabledForNode(ruleID, node) {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: rule disabled",
				ruleID,
				node.Debug(false),
			)
		}

		context.cache.Put(node, ruleID, nil)
		return nil, nil
	}

	var detections []*detectortypes.Detection
	if err := scanner.withCycleProtection(node, ruleID, func() (err error) {
		detections, err = scanner.detectorSet.DetectAt(node, ruleID, context)
		context.cache.Put(node, ruleID, detections)
		return
	}); err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %s: %d detections",
			ruleID,
			node.Debug(false),
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *Scanner) withCycleProtection(node *tree.Node, ruleID string, body func() error) error {
	if slices.Contains(node.ExecutingRules, ruleID) {
		return fmt.Errorf(
			"cycle found during rule evaluation: [%s > %s]\nnode: %s",
			strings.Join(node.ExecutingRules, " > "),
			ruleID,
			node.Debug(true),
		)
	}

	node.ExecutingRules = append(node.ExecutingRules, ruleID)

	if err := body(); err != nil {
		return err
	}

	node.ExecutingRules = node.ExecutingRules[:len(node.ExecutingRules)-1]

	return nil
}

func Scan(
	ctx context.Context,
	detectorSet detectorset.Set,
	fileName string,
	fileStats *stats.FileStats,
	tree *tree.Tree,
	queryContext *query.Context,
	rulesDisabledForNodes map[string][]*tree.Node,
	rootNode *tree.Node,
	cache *cache.Cache,
	scope settings.RuleReferenceScope,
	ruleID,
	sanitizerRuleID string,
) ([]*detectortypes.Detection, error) {
	return new(
		ctx,
		detectorSet,
		fileName,
		fileStats,
		tree,
		queryContext,
		rulesDisabledForNodes,
		rootNode,
		cache,
		scope,
		ruleID,
		sanitizerRuleID,
	).Scan()
}
