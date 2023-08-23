package types

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

type Detection struct {
	RuleID    string
	MatchNode *tree.Node
	Data      interface{}
}

type QueryResult map[string]*tree.Node

type ScanContext interface {
	Scan(
		rootNode *tree.Node,
		ruleID,
		sanitizerRuleID string,
		scope settings.RuleReferenceScope,
	) ([]*Detection, error)
	FileName() string
	QueryContext() *query.Context
	NodeFromSitter(sitterNode *sitter.Node) *tree.Node
	QueryMatchAt(query *query.Query, node *tree.Node) ([]QueryResult, error)
	QueryMatchOnceAt(query *query.Query, node *tree.Node) (QueryResult, error)
}

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, scanContext ScanContext) ([]interface{}, error)
}

type DetectorBase struct{}
