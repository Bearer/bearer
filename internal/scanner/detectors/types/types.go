package types

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

type Detection struct {
	RuleID    string
	MatchNode *tree.Node
	Data      interface{}
}

type ScanContext interface {
	FileName() string
	Scan(
		rootNode *tree.Node,
		ruleID,
		sanitizerRuleID string,
		scope settings.RuleReferenceScope,
	) ([]*Detection, error)
}

type Detector interface {
	Name() string
	DetectAt(node *tree.Node, scanContext ScanContext) ([]interface{}, error)
}

type DetectorBase struct{}
