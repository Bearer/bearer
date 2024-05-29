package types

import (
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
)

type Detection struct {
	RuleID    string
	MatchNode *tree.Node
	Data      interface{}
}

type Context interface {
	Filename() string
	Scan(
		rootNode *tree.Node,
		rule *ruleset.Rule,
		traversalStrategy traversalstrategy.Strategy,
	) ([]*Detection, error)
}

type Detector interface {
	Rule() *ruleset.Rule
	DetectAt(node *tree.Node, detectorContext Context) ([]interface{}, error)
	DetectExpectedAt(node *tree.Node, detectorContext Context) ([]interface{}, error)
}

type DetectorBase interface {
	DetectExpectedAt(node *tree.Node, detectorContext Context) ([]interface{}, error)
}
