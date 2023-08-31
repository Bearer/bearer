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

type scanner struct {
	fileContext *filecontext.Context
	context     *context
	rootNode    *tree.Node
	detectorID  int
}

func (scanner *scanner) Scan() ([]*detectortypes.Detection, error) {
	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan start at %s [%s]",
			scanner.fileContext.RuleIDFor(scanner.detectorID),
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
		detections, err = scanner.scanNode()
	}

	if err != nil {
		return nil, err
	}

	scanner.fileContext.RuleStats(scanner.detectorID, startTime)

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"rule %s scan end at %s [%s]: %d detections",
			scanner.fileContext.RuleIDFor(scanner.detectorID),
			scanner.rootNode.Debug(),
			scanner.context.scope,
			len(detections),
		)
	}

	return detections, nil
}

func (scanner *scanner) scanNode() ([]*detectortypes.Detection, error) {
	result, err := scanner.detectAtNode(scanner.rootNode)
	if err != nil {
		return nil, err
	}

	return result.Detections, nil
}

func (scanner *scanner) scanDescendantsAndAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return append(node.Children(), node.AliasOf()...)
	})
}

func (scanner *scanner) scanDescendants() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.Children()
	})
}

func (scanner *scanner) scanDataflowAndAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return append(node.DataflowSources(), node.AliasOf()...)
	})
}
func (scanner *scanner) scanAliases() ([]*detectortypes.Detection, error) {
	return scanner.detectWithNext(func(node *tree.Node) []*tree.Node {
		return node.AliasOf()
	})
}

func (scanner *scanner) detectWithNext(
	getNext func(node *tree.Node) []*tree.Node,
) ([]*detectortypes.Detection, error) {
	next := make([]*tree.Node, 0, 1000)
	nodes := make([]*tree.Node, 0, 1000)
	nodes = append(nodes, scanner.rootNode)

	var detections []*detectortypes.Detection

	seen := set.New[*tree.Node]()

	for {
		if len(nodes) == 0 {
			break
		}

		for _, node := range nodes {
			if !seen.Add(node) {
				continue
			}

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

func (scanner *scanner) ruleDisabledForNode(node *tree.Node) bool {
	for current := node; current != nil; current = current.Parent() {
		// FIXME: we should use the detector id
		if slices.Contains(current.DisabledRuleIDs(), scanner.ruleID()) {
			return true
		}
	}

	return false
}

func (scanner *scanner) detectAtNode(node *tree.Node) (*detectorset.Result, error) {
	if scanner.fileContext.Err() != nil {
		return nil, scanner.fileContext.Err()
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", scanner.ruleID(), node.Debug())
	}

	if result, cached := scanner.context.cache.Get(node, scanner.detectorID); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: %s (cached)",
				scanner.ruleID(),
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
				scanner.ruleID(),
				node.Debug(),
			)
		}

		scanner.context.cache.Put(node, scanner.detectorID, &detectorset.Result{})
		return nil, nil
	}

	result, err := scanner.detectWithoutCycles(node)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %s: %s",
			scanner.ruleID(),
			node.Debug(),
			traceResultText(result),
		)
	}

	scanner.context.cache.Put(node, scanner.detectorID, result)
	return result, nil
}

func (scanner *scanner) detectWithoutCycles(node *tree.Node) (*detectorset.Result, error) {
	if slices.Contains(node.ExecutingDetectors, scanner.detectorID) {
		executingRules := make([]string, len(node.ExecutingDetectors))
		for i, detectorID := range node.ExecutingDetectors {
			executingRules[i] = scanner.fileContext.RuleIDFor(detectorID)
		}

		return nil, fmt.Errorf(
			"cycle found during rule evaluation at %s: [%s > %s]",
			node.Debug(),
			strings.Join(executingRules, " > "),
			scanner.ruleID(),
		)
	}

	node.ExecutingDetectors = append(node.ExecutingDetectors, scanner.detectorID)
	result, err := scanner.fileContext.DetectAt(node, scanner.detectorID, scanner.context)
	node.ExecutingDetectors = node.ExecutingDetectors[:len(node.ExecutingDetectors)-1]

	return result, err
}

func (scanner *scanner) ruleID() string {
	return scanner.fileContext.RuleIDFor(scanner.detectorID)
}

func traceResultText(result *detectorset.Result) string {
	if result.Sanitized {
		return "sanitized"
	}

	return fmt.Sprintf("%d detections", len(result.Detections))
}
