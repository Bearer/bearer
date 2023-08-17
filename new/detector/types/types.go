package types

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

type QueryResult map[string]*tree.Node

type EvaluationState interface {
	Evaluate(
		rootNode *tree.Node,
		ruleID,
		sanitizerRuleID string,
		scope settings.RuleReferenceScope,
	) ([]*detection.Detection, error)
	FileName() string
	QueryContext() *query.Context
	NodeFromSitter(sitterNode *sitter.Node) *tree.Node
	QueryMatchAt(query *query.Query, node *tree.Node) ([]QueryResult, error)
	QueryMatchOnceAt(query *query.Query, node *tree.Node) (QueryResult, error)
}

type DetectorSet interface {
	BuiltinAndSharedRuleIDs() []string
	TopLevelRuleIDs() []string
	DetectAt(
		node *tree.Node,
		ruleID string,
		evaluationState EvaluationState,
	) ([]*detection.Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, evaluationState EvaluationState) ([]interface{}, error)
}

type DetectorBase struct{}
