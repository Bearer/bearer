package string

import (
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/ruleset"

	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_value":
		return []interface{}{common.String{
			Value:     node.Content(),
			IsLiteral: true,
		}}, nil
	case "string", "encapsed_string":
		return common.ConcatenateChildStrings(node, detectorContext)
	case "binary_expression":
		if node.Children()[1].Content() == "." {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "augmented_assignment_expression":
		if node.Children()[1].Content() == ".=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}
