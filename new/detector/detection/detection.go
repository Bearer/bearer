package detection

import "github.com/bearer/bearer/pkg/ast/tree"

type Detection struct {
	RuleID    string
	MatchNode *tree.Node
	Data      interface{}
}
