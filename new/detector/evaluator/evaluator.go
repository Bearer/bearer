package evaluator

import (
	"fmt"
	"strings"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/implementation"
	langtree "github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"golang.org/x/exp/slices"
)

type evaluator struct {
	langImplementation    implementation.Implementation
	lang                  languagetypes.Language
	detectorSet           types.DetectorSet
	detectionCache        map[langtree.NodeID]map[string][]*types.Detection
	executingDetectors    map[langtree.NodeID][]string
	fileName              string
	rulesDisabledForNodes map[string][]*langtree.Node
}

func New(
	langImplementation implementation.Implementation,
	lang languagetypes.Language,
	detectorSet types.DetectorSet,
	tree *langtree.Tree,
	fileName string,
) types.Evaluator {
	detectionCache := make(map[langtree.NodeID]map[string][]*types.Detection)

	return &evaluator{
		langImplementation:    langImplementation,
		lang:                  lang,
		fileName:              fileName,
		detectorSet:           detectorSet,
		detectionCache:        detectionCache,
		executingDetectors:    make(map[langtree.NodeID][]string),
		rulesDisabledForNodes: mapNodesToDisabledRules(tree.RootNode()),
	}
}

func (evaluator *evaluator) FileName() string {
	return evaluator.fileName
}

func (evaluator *evaluator) Evaluate(
	rootNode *langtree.Node,
	detectorType, sanitizerDetectorType string,
	scope settings.RuleReferenceScope,
	followFlow bool,
) ([]*types.Detection, error) {
	if rootNode == nil {
		return nil, nil
	}

	nestedDetections, err := evaluator.detectorSet.NestedDetections(detectorType)
	if err != nil {
		return nil, err
	}

	var result []*types.Detection
	var nestedMode bool

	if err := rootNode.Walk(func(node *langtree.Node, visitChildren func() error) error {
		if scope == settings.RESULT_SCOPE && !evaluator.langImplementation.ContributesToResult(node) {
			return nil
		}

		if nestedMode && !evaluator.langImplementation.PassthroughNested(node) {
			return nil
		}

		detections, sanitized, err := evaluator.sanitizedNodeDetections(node, detectorType, sanitizerDetectorType, scope)
		if sanitized || err != nil {
			return err
		}

		if followFlow {
			for _, unifiedNode := range node.UnifiedNodes() {
				unifiedNodeDetections, err := evaluator.Evaluate(unifiedNode, detectorType, sanitizerDetectorType, scope, true)
				if err != nil {
					return err
				}

				detections = append(detections, unifiedNodeDetections...)
			}
		}

		result = append(result, detections...)

		if scope != settings.CURSOR_SCOPE {
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

	return result, nil
}

func (evaluator *evaluator) ruleDisabledForNode(ruleId string, node *langtree.Node) bool {
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

func (evaluator *evaluator) sanitizedNodeDetections(
	node *langtree.Node,
	detectorType, sanitizerDetectorType string,
	scope settings.RuleReferenceScope,
) ([]*types.Detection, bool, error) {
	if sanitizerDetectorType != "" {
		sanitizerDetections, err := evaluator.nonUnifiedNodeDetections(node, sanitizerDetectorType, scope)
		if len(sanitizerDetections) != 0 || err != nil {
			return nil, true, err
		}
	}

	detections, err := evaluator.nonUnifiedNodeDetections(node, detectorType, scope)
	return detections, false, err
}

func (evaluator *evaluator) nonUnifiedNodeDetections(
	node *langtree.Node,
	detectorType string,
	scope settings.RuleReferenceScope,
) ([]*types.Detection, error) {
	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		err := evaluator.detectAtNode(node, detectorType, scope)
		if err != nil {
			return nil, err
		}

		nodeDetections = evaluator.detectionCache[node.ID()]
	}

	if detections, ok := nodeDetections[detectorType]; ok {
		return detections, nil
	}

	err := evaluator.detectAtNode(node, detectorType, scope)
	if err != nil {
		return nil, err
	}

	return nodeDetections[detectorType], nil
}

func (evaluator *evaluator) detectAtNode(
	node *langtree.Node,
	detectorType string,
	scope settings.RuleReferenceScope,
) error {
	if evaluator.ruleDisabledForNode(detectorType, node) {
		return nil
	}

	return evaluator.withCycleProtection(node, detectorType, func() error {
		detections, err := evaluator.detectorSet.DetectAt(node, detectorType, scope, evaluator)
		if err != nil {
			return err
		}

		nodeDetections, ok := evaluator.detectionCache[node.ID()]
		if !ok {
			nodeDetections = make(map[string][]*types.Detection)
			evaluator.detectionCache[node.ID()] = nodeDetections
		}

		nodeDetections[detectorType] = detections

		return nil
	})
}

func (evaluator *evaluator) withCycleProtection(node *langtree.Node, detectorType string, body func() error) error {
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
