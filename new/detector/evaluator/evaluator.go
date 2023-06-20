package evaluator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bearer/bearer/new/detector/detection"
	cachepkg "github.com/bearer/bearer/new/detector/evaluator/cache"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"
	langtree "github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Evaluator struct {
	ctx                   context.Context
	langImplementation    implementation.Implementation
	lang                  languagetypes.Language
	detectorSet           types.DetectorSet
	stats                 *stats.Stats
	executingDetectors    map[langtree.NodeID][]string
	fileName              string
	rulesDisabledForNodes map[string][]*langtree.Node
}

func New(
	ctx context.Context,
	langImplementation implementation.Implementation,
	lang languagetypes.Language,
	detectorSet types.DetectorSet,
	tree *langtree.Tree,
	fileName string,
	stats *stats.Stats,
) *Evaluator {
	return &Evaluator{
		ctx:                   ctx,
		langImplementation:    langImplementation,
		lang:                  lang,
		fileName:              fileName,
		detectorSet:           detectorSet,
		stats:                 stats,
		executingDetectors:    make(map[langtree.NodeID][]string),
		rulesDisabledForNodes: mapNodesToDisabledRules(tree.RootNode()),
	}
}

func (evaluator *Evaluator) Evaluate(
	rootNode *langtree.Node,
	detectorType, sanitizerDetectorType string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
	followFlow bool,
) ([]*detection.Detection, error) {
	if rootNode == nil {
		return nil, nil
	}

	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"evaluate start: %d:%d:%s:\n%s",
			rootNode.StartLineNumber(),
			rootNode.StartColumnNumber(),
			rootNode.Type(),
			rootNode.Content(),
		)
	}

	key := cachepkg.NewKey(rootNode, detectorType, scope, followFlow)

	if detections, cached := cache.Get(key); cached {
		if evaluator.stats != nil {
			evaluator.stats.Record(detectorType, startTime)
		}
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"evaluate end: %d:%d:%s: %d detections (cached)",
				rootNode.StartLineNumber(),
				rootNode.StartColumnNumber(),
				rootNode.Type(),
				len(detections),
			)
		}

		return detections, nil
	}

	nestedDetections, err := evaluator.detectorSet.NestedDetections(detectorType)
	if err != nil {
		return nil, err
	}

	var result []*detection.Detection
	var nestedMode bool

	if err := rootNode.Walk(func(node *langtree.Node, visitChildren func() error) error {
		if evaluator.ctx.Err() != nil {
			return evaluator.ctx.Err()
		}

		if scope == settings.RESULT_SCOPE && !evaluator.langImplementation.ContributesToResult(node) {
			return nil
		}

		if nestedMode && !evaluator.langImplementation.PassthroughNested(node) {
			return nil
		}

		detections, sanitized, err := evaluator.sanitizedNodeDetections(node, detectorType, sanitizerDetectorType, cache, scope)
		if sanitized || err != nil {
			return err
		}

		if followFlow {
			for _, unifiedNode := range node.UnifiedNodes() {
				unifiedNodeDetections, err := evaluator.Evaluate(unifiedNode, detectorType, sanitizerDetectorType, cache, scope, true)
				if err != nil {
					return err
				}

				detections = append(detections, unifiedNodeDetections...)
			}
		}

		result = append(result, detections...)

		if scope != settings.CURSOR_SCOPE && !evaluator.langImplementation.IsMatchLeaf(node) {
			parentNestedMode := nestedMode

			if len(detections) != 0 && nestedDetections {
				nestedMode = true
			}

			err = visitChildren()
			nestedMode = parentNestedMode
		}

		return err
	}); err != nil {
		return nil, err
	}

	cache.Put(key, result)

	if evaluator.stats != nil {
		evaluator.stats.Record(detectorType, startTime)
	}
	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"evaluate end: %d:%d:%s: %d detections",
			rootNode.StartLineNumber(),
			rootNode.StartColumnNumber(),
			rootNode.Type(),
			len(result),
		)
	}

	return result, nil
}

func (evaluator *Evaluator) ruleDisabledForNode(ruleId string, node *langtree.Node) bool {
	nodesToIgnore := evaluator.rulesDisabledForNodes[ruleId]
	if nodesToIgnore == nil {
		return false
	}

	// check node
	for _, ignoredNode := range nodesToIgnore {
		if ignoredNode.Equal(node) {
			return true
		}
	}

	// check node ancestors
	parent := node.Parent()
	for parent != nil {
		for _, ignoredNode := range nodesToIgnore {
			if ignoredNode.Equal(parent) {
				return true
			}
		}

		parent = parent.Parent()
	}

	return false
}

func mapNodesToDisabledRules(rootNode *langtree.Node) map[string][]*langtree.Node {
	res := make(map[string][]*langtree.Node)
	var disabledRules []string
	err := rootNode.Walk(func(node *langtree.Node, visitChildren func() error) error {
		if node.Type() == "comment" {
			// reset rules skipped array
			disabledRules = []string{}

			nodeContent := node.Content()
			if strings.Contains(nodeContent, "bearer:disable") {
				ruleIdsStr := strings.Split(nodeContent, "bearer:disable")[1]

				for _, ruleId := range strings.Split(ruleIdsStr, ",") {
					disabledRules = append(disabledRules, strings.TrimSpace(ruleId))
				}
			}

			return visitChildren()
		}

		// add rules skipped and node to result map
		for _, ruleId := range disabledRules {
			res[ruleId] = append(res[ruleId], node)
		}

		// reset rules skipped array
		disabledRules = []string{}
		return visitChildren()
	})

	// walk itself shouldn't trigger an error, and we aren't creating any
	if err != nil {
		panic(err)
	}

	return res
}

func (evaluator *Evaluator) sanitizedNodeDetections(
	node *langtree.Node,
	detectorType, sanitizerDetectorType string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, bool, error) {
	if sanitizerDetectorType != "" {
		sanitizerDetections, err := evaluator.detectAtNode(node, sanitizerDetectorType, cache, settings.DefaultScope)
		if len(sanitizerDetections) != 0 || err != nil {
			return nil, true, err
		}
	}

	detections, err := evaluator.detectAtNode(node, detectorType, cache, scope)
	return detections, false, err
}

func (evaluator *Evaluator) detectAtNode(
	node *langtree.Node,
	detectorType string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, error) {
	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node start: %s at %d:%d:%s\n%s",
			detectorType,
			node.StartLineNumber(),
			node.StartColumnNumber(),
			node.Type(),
			node.Content(),
		)
	}
	key := cachepkg.NewKey(node, detectorType, settings.CURSOR_SCOPE, false)

	if detections, cached := cache.Get(key); cached {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %d:%d:%s: %d detections (cached)",
				detectorType,
				node.StartLineNumber(),
				node.StartColumnNumber(),
				node.Type(),
				len(detections),
			)
		}

		return detections, nil
	}

	if evaluator.ruleDisabledForNode(detectorType, node) {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %d:%d:%s: rule disabled",
				detectorType,
				node.StartLineNumber(),
				node.StartColumnNumber(),
				node.Type(),
			)
		}

		cache.Put(key, nil)
		return nil, nil
	}

	var detections []*detection.Detection
	if err := evaluator.withCycleProtection(node, detectorType, func() (err error) {
		state := evaluationState{
			cache:     cache,
			scope:     scope,
			evaluator: evaluator,
		}
		detections, err = evaluator.detectorSet.DetectAt(node, detectorType, state)
		cache.Put(key, detections)
		return
	}); err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"detect at node end: %s at %d:%d:%s: %d detections",
			detectorType,
			node.StartLineNumber(),
			node.StartColumnNumber(),
			node.Type(),
			len(detections),
		)
	}

	return detections, nil
}

func (evaluator *Evaluator) withCycleProtection(node *langtree.Node, detectorType string, body func() error) error {
	nodeID := node.ID()

	executingDetectors := evaluator.executingDetectors[nodeID]
	if slices.Contains(evaluator.executingDetectors[nodeID], detectorType) {
		return fmt.Errorf(
			"cycle found in detector usage: [%s > %s]\nnode type: %s, content:\n%s",
			strings.Join(executingDetectors, " > "),
			detectorType,
			node.Type(),
			node.Content(),
		)
	}

	evaluator.executingDetectors[nodeID] = append(evaluator.executingDetectors[nodeID], detectorType)

	if err := body(); err != nil {
		return err
	}

	if len(evaluator.executingDetectors[nodeID]) == 1 {
		delete(evaluator.executingDetectors, nodeID)
	} else {
		executingDetectors := evaluator.executingDetectors[nodeID]
		evaluator.executingDetectors[nodeID] = executingDetectors[:len(executingDetectors)-1]
	}

	return nil
}

type evaluationState struct {
	cache     *cachepkg.Cache
	scope     settings.RuleReferenceScope
	evaluator *Evaluator
}

func (state evaluationState) Evaluate(
	rootNode *tree.Node,
	detectorType,
	sanitizerDetectorType string,
	scope settings.RuleReferenceScope,
	followFlow bool,
) ([]*detection.Detection, error) {
	effectiveScope := scope
	if effectiveScope == settings.NESTED_SCOPE && state.scope == settings.RESULT_SCOPE {
		effectiveScope = settings.RESULT_SCOPE
	}

	return state.evaluator.Evaluate(
		rootNode,
		detectorType,
		sanitizerDetectorType,
		state.cache,
		effectiveScope,
		followFlow,
	)
}

func (state evaluationState) FileName() string {
	return state.evaluator.fileName
}
