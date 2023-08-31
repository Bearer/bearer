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

type Context interface {
	Filename() string
	Scan(rootNode *tree.Node, detectorID int, scope settings.RuleReferenceScope) ([]*Detection, error)
	// FIXME: remove this
	ScanRule(rootNode *tree.Node, ruleID string, scope settings.RuleReferenceScope) ([]*Detection, error)
}

type Detector interface {
	RuleID() string
	DetectAt(node *tree.Node, detectorContext Context) ([]interface{}, error)
}

type DetectorBase struct{}
