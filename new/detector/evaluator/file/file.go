package file

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/detector/detection"
	cachepkg "github.com/bearer/bearer/new/detector/evaluator/cache"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	asttree "github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

type FileEvaluator struct {
	ctx                   context.Context
	tree                  *asttree.Tree
	detectorSet           types.DetectorSet
	fileStats             *stats.FileStats
	fileName              string
	rulesDisabledForNodes map[string][]*asttree.Node
	queryContext          *query.Context
}

func New(
	ctx context.Context,
	detectorSet types.DetectorSet,
	tree *asttree.Tree,
	fileName string,
	fileStats *stats.FileStats,
) *FileEvaluator {
	return &FileEvaluator{
		ctx:                   ctx,
		tree:                  tree,
		fileName:              fileName,
		detectorSet:           detectorSet,
		fileStats:             fileStats,
		rulesDisabledForNodes: mapNodesToDisabledRules(tree.RootNode()),
		queryContext:          query.NewContext(tree.ContentBytes(), tree.RootNode().SitterNode()),
	}
}

func (evaluator *FileEvaluator) Evaluate(
	rootNode *asttree.Node,
	ruleID string,
	sanitizerRuleID string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
	// FIXME: support or remove followFlow
	followFlow bool,
) ([]*detection.Detection, error) {
	// FIXME: remove this and fix the caller that's passing nil
	if rootNode == nil {
		return nil, nil
	}

	startTime := time.Now()

	if log.Trace().Enabled() {
		log.Trace().Msgf("evaluate start: %s at %s", ruleID, rootNode.Debug(true))
	}

	key := cachepkg.NewKey(rootNode, ruleID, scope, followFlow)

	if detections, cached := cache.Get(key); cached {
		evaluator.fileStats.Rule(ruleID, startTime)

		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"evaluate end: %s at %s: %d detections (cached)",
				ruleID,
				rootNode.Debug(false),
				len(detections),
			)
		}

		return detections, nil
	}

	var detections []*detection.Detection
	var err error

	switch scope {
	case settings.NESTED_SCOPE:
		detections, err = evaluator.evalAtDescendents(ruleID, sanitizerRuleID, rootNode, cache, scope)
	case settings.RESULT_SCOPE:
		detections, err = evaluator.evalAtDataflowSources(ruleID, sanitizerRuleID, rootNode, cache, scope)
	case settings.CURSOR_SCOPE:
		// FIXME: variable lookups
		detections, _, err = evaluator.sanitizedNodeDetections(rootNode, ruleID, sanitizerRuleID, cache, scope)
	}

	if err != nil {
		return nil, err
	}

	cache.Put(key, detections)

	evaluator.fileStats.Rule(ruleID, startTime)

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"evaluate end: %s at %s: %d detections",
			ruleID,
			rootNode.Debug(false),
			len(detections),
		)
	}

	return detections, nil
}

func (evaluator *FileEvaluator) evalAtDescendents(
	ruleID,
	sanitizerRuleID string,
	rootNode *asttree.Node,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, error) {
	nodes := []*asttree.Node{rootNode}

	var detections []*detection.Detection

	for {
		if len(nodes) == 0 {
			break
		}

		var next []*asttree.Node

		for _, node := range nodes {
			nodeDetections, sanitized, err := evaluator.sanitizedNodeDetections(node, ruleID, sanitizerRuleID, cache, scope)
			if err != nil {
				return nil, err
			}
			if sanitized {
				continue
			}

			detections = append(detections, nodeDetections...)
			next = append(next, node.Children()...)
		}

		nodes = next
	}

	return detections, nil
}

func (evaluator *FileEvaluator) evalAtDataflowSources(
	ruleID,
	sanitizerRuleID string,
	rootNode *asttree.Node,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, error) {
	nodes := []*asttree.Node{rootNode}

	var detections []*detection.Detection

	for {
		if len(nodes) == 0 {
			break
		}

		var next []*asttree.Node

		for _, node := range nodes {
			nodeDetections, sanitized, err := evaluator.sanitizedNodeDetections(node, ruleID, sanitizerRuleID, cache, scope)
			if err != nil {
				return nil, err
			}
			if sanitized {
				continue
			}

			detections = append(detections, nodeDetections...)
			next = append(next, node.DataflowSources()...)
		}

		nodes = next
	}

	return detections, nil
}

func (evaluator *FileEvaluator) ruleDisabledForNode(ruleId string, node *asttree.Node) bool {
	nodesToIgnore := evaluator.rulesDisabledForNodes[ruleId]
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

func mapNodesToDisabledRules(rootNode *asttree.Node) map[string][]*asttree.Node {
	res := make(map[string][]*asttree.Node)
	var disabledRules []string
	err := rootNode.Walk(func(node *asttree.Node, visitChildren func() error) error {
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

func (evaluator *FileEvaluator) sanitizedNodeDetections(
	node *asttree.Node,
	ruleID, sanitizerRuleID string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, bool, error) {
	if sanitizerRuleID != "" {
		sanitizerDetections, err := evaluator.detectAtNode(node, sanitizerRuleID, cache, settings.DefaultScope)
		if len(sanitizerDetections) != 0 || err != nil {
			return nil, true, err
		}
	}

	detections, err := evaluator.detectAtNode(node, ruleID, cache, scope)
	return detections, false, err
}

func (evaluator *FileEvaluator) detectAtNode(
	node *asttree.Node,
	ruleID string,
	cache *cachepkg.Cache,
	scope settings.RuleReferenceScope,
) ([]*detection.Detection, error) {
	if log.Trace().Enabled() {
		log.Trace().Msgf("detect at node start: %s at %s", ruleID, node.Debug(true))
	}

	key := cachepkg.NewKey(node, ruleID, settings.CURSOR_SCOPE, false)

	if detections, cached := cache.Get(key); cached {
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

	if evaluator.ruleDisabledForNode(ruleID, node) {
		if log.Trace().Enabled() {
			log.Trace().Msgf(
				"detect at node end: %s at %s: rule disabled",
				ruleID,
				node.Debug(false),
			)
		}

		cache.Put(key, nil)
		return nil, nil
	}

	var detections []*detection.Detection
	if err := evaluator.withCycleProtection(node, ruleID, func() (err error) {
		state := evaluationState{
			cache:        cache,
			queryContext: evaluator.queryContext,
			scope:        scope,
			evaluator:    evaluator,
		}
		detections, err = evaluator.detectorSet.DetectAt(node, ruleID, state)
		cache.Put(key, detections)
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

func (evaluator *FileEvaluator) withCycleProtection(node *asttree.Node, ruleID string, body func() error) error {
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

type evaluationState struct {
	cache        *cachepkg.Cache
	queryContext *query.Context
	scope        settings.RuleReferenceScope
	evaluator    *FileEvaluator
}

func (state evaluationState) Evaluate(
	rootNode *asttree.Node,
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

func (state evaluationState) QueryContext() *query.Context {
	return state.queryContext
}

func (state evaluationState) NodeFromSitter(sitterNode *sitter.Node) *asttree.Node {
	return state.evaluator.tree.NodeFromSitter(sitterNode)
}

func (state evaluationState) QueryMatchAt(query *query.Query, node *asttree.Node) ([]types.QueryResult, error) {
	sitterResults, err := query.MatchAt(state.queryContext, node.SitterNode())
	if err != nil {
		return nil, err
	}

	results := make([]types.QueryResult, len(sitterResults))

	for i, sitterResult := range sitterResults {
		results[i] = state.translateResult(sitterResult)
	}

	return results, nil
}

func (state evaluationState) QueryMatchOnceAt(query *query.Query, node *asttree.Node) (types.QueryResult, error) {
	sitterResult, err := query.MatchOnceAt(state.queryContext, node.SitterNode())
	if err != nil {
		return nil, err
	}

	return state.translateResult(sitterResult), nil
}

// FIXME: try and remove the translation by caching query results on the ast tree
func (state evaluationState) translateResult(sitterResult query.Result) types.QueryResult {
	result := make(map[string]*asttree.Node)

	for name, sitterNode := range sitterResult {
		result[name] = state.NodeFromSitter(sitterNode)
	}

	return result
}
