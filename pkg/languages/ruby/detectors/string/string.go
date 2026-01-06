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
	case "string_content":
		return common.Literal(node.Content()), nil
	case "interpolation", "string":
		return common.ConcatenateChildStrings(node, detectorContext)
	case "binary":
		children := node.Children()
		if len(children) > 1 {
			switch children[1].Content() {
			case "+":
				return common.ConcatenateChildStrings(node, detectorContext)
			case "||", "or":
				leftData, err := common.GetStringData(node.ChildByFieldName("left"), detectorContext)
				if err != nil {
					return nil, err
				}

				rightData, err := common.GetStringData(node.ChildByFieldName("right"), detectorContext)
				if err != nil {
					return nil, err
				}

				return append(leftData, rightData...), nil
			}
		}
	case "operator_assignment":
		children := node.Children()
		if len(children) > 1 && children[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}
