package types

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/file"
)

type QueryResult map[string]*tree.Node

type EvaluationState interface {
	Evaluate(
		rootNode *tree.Node,
		ruleID,
		sanitizerRuleID string,
		scope settings.RuleReferenceScope,
		followFlow bool,
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
	NestedDetections() bool
}

type DetectorBase struct{}

func (*DetectorBase) NestedDetections() bool {
	return true
}

type Composition interface {
	DetectFromFile(ctx context.Context, fileStats *stats.FileStats, file *file.FileInfo) ([]*detection.Detection, error)
	DetectFromFileWithTypes(
		ctx context.Context,
		fileStats *stats.FileStats,
		file *file.FileInfo,
		detectorTypes, sharedDetectorTypes []string,
	) ([]*detection.Detection, error)
	Close()
}
