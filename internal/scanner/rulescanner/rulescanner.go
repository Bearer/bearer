package rulescanner

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/util/set"

	"github.com/bearer/bearer/internal/scanner/cache"
)

type Scanner struct {
	fileContext *filecontext.Context
	rootNode    *tree.Node
	context,
	sanitizerContext *detectorContext
	ruleID,
	sanitizerRuleID string
}

func (scanner *Scanner) Scan() ([]*detectortypes.Detection, error) {
	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan start at %s [%s]",
			scanner.ruleID,
			scanner.rootNode.Debug(),
			scanner.context.scope,
		)
	}

	var detections []*detectortypes.Detection
	var err error

	switch scanner.context.scope {
	case settings.NESTED_SCOPE:
		detections, err = scanner.scanDescendantsAndAliases()
	case settings.NESTED_STRICT_SCOPE:
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

	scanner.fileContext.RuleStats(scanner.ruleID, startTime)

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan end at %s [%s]: %d detections",
			scanner.ruleID,
			scanner.rootNode.Debug(),
			scanner.context.scope,
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *Scanner) scanDescendantsAndAliases() ([]*detectortypes.Detection, error) {
	seen := set.New[*tree.Node]()

	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		var result []*tree.Node

		for _, child := range node.Children() {
			if seen.Add(child) {
				result = append(result, child)
			}
		}
		for _, alias := range node.AliasOf() {
			if seen.Add(alias) {
				result = append(result, alias)
			}
		}

		return result
	})
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
	next := make([]*tree.Node, 0, 1000)
	nodes := make([]*tree.Node, 0, 1000)
	nodes = append(nodes, scanner.rootNode)

	var detections []*detectortypes.Detection

	for {
		if len(nodes) == 0 {
			break
		}

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

		old := nodes
		nodes = next
		// allow memory to be re-used
		next = old[:0]
	}

	return detections, nil
}

func (scanner *Scanner) ruleDisabledForNode(ruleID string, node *tree.Node) bool {
	for current := node; current != nil; current = current.Parent() {
		if slices.Contains(current.DisabledRuleIDs(), ruleID) {
			return true
		}
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
	context *detectorContext,
	node *tree.Node,
	ruleID string,
) ([]*detectortypes.Detection, error) {
	if scanner.fileContext.Err() != nil {
		return nil, scanner.fileContext.Err()
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", ruleID, node.Debug())
	}

	if detections, cached := context.cache.Get(node, ruleID); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: %d detections (cached)",
				ruleID,
				node.Debug(),
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
				node.Debug(),
			)
		}

		context.cache.Put(node, ruleID, nil)
		return nil, nil
	}

	detections, err := scanner.detectWithoutCycles(node, ruleID, context)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %s: %d detections",
			ruleID,
			node.Debug(),
			len(detections),
		)
	}

	context.cache.Put(node, ruleID, detections)
	return detections, nil
}

func (scanner *Scanner) detectWithoutCycles(
	node *tree.Node,
	ruleID string,
	context *detectorContext,
) ([]*detectortypes.Detection, error) {
	if slices.Contains(node.ExecutingRules, ruleID) {
		return nil, fmt.Errorf(
			"cycle found during rule evaluation at %s: [%s > %s]",
			node.Debug(),
			strings.Join(node.ExecutingRules, " > "),
			ruleID,
		)
	}

	node.ExecutingRules = append(node.ExecutingRules, ruleID)
	detections, err := scanner.fileContext.DetectAt(node, ruleID, context)
	node.ExecutingRules = node.ExecutingRules[:len(node.ExecutingRules)-1]

	return detections, err
}

func Scan(
	fileContext *filecontext.Context,
	cache *cache.Cache,
	scope settings.RuleReferenceScope,
	ruleID string,
	rootNode *tree.Node,
) ([]*detectortypes.Detection, error) {
	sanitizerRuleID := ""
	if rule, ok := fileContext.Rules()[ruleID]; ok {
		sanitizerRuleID = rule.SanitizerRuleID
	}

	scanner := &Scanner{
		fileContext:     fileContext,
		ruleID:          ruleID,
		sanitizerRuleID: sanitizerRuleID,
		rootNode:        rootNode,
	}

	scanner.context = newContext(fileContext, cache, scope)
	scanner.sanitizerContext = newContext(fileContext, cache, settings.DefaultScope)

	return scanner.Scan()
}
