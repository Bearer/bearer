package treeevaluator

import (
	detectorexecutortypes "github.com/bearer/curio/new/detectorexecutor/types"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/new/treeevaluator/types"
	"github.com/open-policy-agent/opa/ast"
)

type treeEvaluator struct {
	lang           languagetypes.Language
	executor       detectorexecutortypes.Executor
	detectionCache map[language.NodeID]map[string]*ast.Array
}

func New(
	lang languagetypes.Language,
	executor detectorexecutortypes.Executor,
	tree *language.Tree,
) types.Evaluator {
	detectionCache := make(map[language.NodeID]map[string]*ast.Array)

	return &treeEvaluator{
		lang:           lang,
		executor:       executor,
		detectionCache: detectionCache,
	}
}

func (evaluator *treeEvaluator) TreeDetections(rootNode *language.Node, detectorType string) (*ast.Array, error) {
	var resultTerms []*ast.Term

	if err := rootNode.Walk(func(node *language.Node) error {
		detections, err := evaluator.nonUnifiedNodeDetections(node, detectorType)
		if err != nil {
			return err
		}

		resultTerms = append(resultTerms, collectTerms(detections)...)

		for _, unifiedNode := range node.UnifiedNodes() {
			unifiedNodeDetections, err := evaluator.TreeDetections(unifiedNode, detectorType)
			if err != nil {
				return err
			}

			resultTerms = append(resultTerms, collectTerms(unifiedNodeDetections)...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ast.NewArray(resultTerms...), nil
}

func (evaluator *treeEvaluator) NodeDetections(
	node *language.Node,
	detectorType string,
) (*ast.Array, error) {
	detections, err := evaluator.nonUnifiedNodeDetections(node, detectorType)
	if err != nil {
		return nil, err
	}

	resultTerms := collectTerms(detections)

	for _, unifiedNode := range node.UnifiedNodes() {
		unifiedNodeDetections, err := evaluator.nonUnifiedNodeDetections(unifiedNode, detectorType)
		if err != nil {
			return nil, err
		}

		resultTerms = append(resultTerms, collectTerms(unifiedNodeDetections)...)
	}

	return ast.NewArray(resultTerms...), nil
}

func (evaluator *treeEvaluator) nonUnifiedNodeDetections(
	node *language.Node,
	detectorType string,
) (*ast.Array, error) {
	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		err := evaluator.detectAtNode(node, detectorType)
		if err != nil {
			return nil, err
		}

		nodeDetections = evaluator.detectionCache[node.ID()]
	}

	if detections, ok := nodeDetections[detectorType]; ok {
		return detections, nil
	}

	evaluator.detectAtNode(node, detectorType)
	return nodeDetections[detectorType], nil
}

func (evaluator *treeEvaluator) TreeHasDetection(rootNode *language.Node, detectorType string) (bool, error) {
	hasDetection := false

	if err := rootNode.Walk(func(node *language.Node) error {
		var err error
		hasDetection, err = evaluator.NodeHasDetection(node, detectorType)
		if err != nil {
			return err
		}

		if hasDetection {
			return language.ErrTerminateWalk
		}

		return nil
	}); err != nil {
		return false, err
	}

	return hasDetection, nil
}

func (evaluator *treeEvaluator) NodeHasDetection(node *language.Node, detectorType string) (bool, error) {
	detections, err := evaluator.NodeDetections(node, detectorType)
	if err != nil {
		return false, err
	}

	return detections.Len() != 0, nil
}

func (evaluator *treeEvaluator) detectAtNode(node *language.Node, detectorType string) error {
	detections, err := evaluator.executor.DetectAt(node, detectorType, evaluator)
	if err != nil {
		return err
	}

	nodeDetections, ok := evaluator.detectionCache[node.ID()]
	if !ok {
		nodeDetections = make(map[string]*ast.Array)
		evaluator.detectionCache[node.ID()] = nodeDetections
	}

	nodeDetections[detectorType] = detections

	return nil
}

func collectTerms(array *ast.Array) []*ast.Term {
	result := make([]*ast.Term, array.Len())

	for i := 0; i < array.Len(); i += 1 {
		result[i] = array.Elem(i)
	}

	return result
}
