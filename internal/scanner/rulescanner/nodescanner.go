package rulescanner

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/util/set"
)

type nodeScanner struct {
	fileContext *filecontext.Context
	rootNode    *tree.Node
	ruleScanner *Scanner
	ruleID      string
}

func (scanner *nodeScanner) Scan() ([]*detectortypes.Detection, error) {
	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan start at %s [%s]",
			scanner.ruleID,
			scanner.rootNode.Debug(),
			scanner.ruleScanner.scope,
		)
	}

	var detections []*detectortypes.Detection
	var err error

	switch scanner.ruleScanner.scope {
	case settings.NESTED_SCOPE:
		detections, err = scanner.scanDescendantsAndAliases()
	case settings.NESTED_STRICT_SCOPE:
		detections, err = scanner.scanDescendants()
	case settings.RESULT_SCOPE:
		detections, err = scanner.scanDataflowAndAliases()
	case settings.CURSOR_SCOPE:
		detections, err = scanner.scanAliases()
	case settings.CURSOR_STRICT_SCOPE:
		detections, err = scanner.scanNode()
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
			scanner.ruleScanner.scope,
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *nodeScanner) scanNode() ([]*detectortypes.Detection, error) {
	result, err := scanner.detectAtNode(scanner.rootNode)
	if err != nil {
		return nil, err
	}

	return result.Detections, nil
}

func (scanner *nodeScanner) scanDescendantsAndAliases() ([]*detectortypes.Detection, error) {
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

func (scanner *nodeScanner) scanDescendants() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.Children()
	})
}

func (scanner *nodeScanner) scanDataflowAndAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return append(node.DataflowSources(), node.AliasOf()...)
	})
}
func (scanner *nodeScanner) scanAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.AliasOf()
	})
}

func (scanner *nodeScanner) detectWithNext(
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
			nodeResult, err := scanner.detectAtNode(node)
			if err != nil {
				return nil, err
			}
			if nodeResult.Sanitized {
				continue
			}

			detections = append(detections, nodeResult.Detections...)
			next = append(next, getNext(node)...)
		}

		old := nodes
		nodes = next
		// allow memory to be re-used
		next = old[:0]
	}

	return detections, nil
}

func (scanner *nodeScanner) ruleDisabledForNode(node *tree.Node) bool {
	for current := node; current != nil; current = current.Parent() {
		if slices.Contains(current.DisabledRuleIDs(), scanner.ruleID) {
			return true
		}
	}

	return false
}

func (scanner *nodeScanner) detectAtNode(node *tree.Node) (*detectorset.Result, error) {
	if scanner.fileContext.Err() != nil {
		return nil, scanner.fileContext.Err()
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", scanner.ruleID, node.Debug())
	}

	if result, cached := scanner.ruleScanner.cache.Get(node, scanner.ruleID); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: %s (cached)",
				scanner.ruleID,
				node.Debug(),
				traceResultText(result),
			)
		}

		return result, nil
	}

	if scanner.ruleDisabledForNode(node) {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: rule disabled",
				scanner.ruleID,
				node.Debug(),
			)
		}

		scanner.ruleScanner.cache.Put(node, scanner.ruleID, &detectorset.Result{})
		return nil, nil
	}

	result, err := scanner.detectWithoutCycles(node)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %s: %s",
			scanner.ruleID,
			node.Debug(),
			traceResultText(result),
		)
	}

	scanner.ruleScanner.cache.Put(node, scanner.ruleID, result)
	return result, nil
}

func (scanner *nodeScanner) detectWithoutCycles(node *tree.Node) (*detectorset.Result, error) {
	if slices.Contains(node.ExecutingRules, scanner.ruleID) {
		return nil, fmt.Errorf(
			"cycle found during rule evaluation at %s: [%s > %s]",
			node.Debug(),
			strings.Join(node.ExecutingRules, " > "),
			scanner.ruleID,
		)
	}

	node.ExecutingRules = append(node.ExecutingRules, scanner.ruleID)
	result, err := scanner.fileContext.DetectAt(node, scanner.ruleID, scanner.ruleScanner)
	node.ExecutingRules = node.ExecutingRules[:len(node.ExecutingRules)-1]

	return result, err
}

func traceResultText(result *detectorset.Result) string {
	if result.Sanitized {
		return "sanitized"
	}

	return fmt.Sprintf("%d detections", len(result.Detections))
}
