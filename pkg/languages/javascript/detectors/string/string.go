package string

import (
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/util/stringutil"

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
	case "string", "template_string":
		return common.ConcatenateChildStrings(node, detectorContext)
	case "string_fragment":
		return common.Literal(node.Content()), nil
	case "escape_sequence":
		// Use JavaScript-specific unescaping which handles unknown escapes
		// (like \s, \d) by stripping the backslash, matching JavaScript behavior
		value := stringutil.UnescapeJavaScript(node.Content())
		return common.Literal(value), nil
	case "template_substitution":
		return common.GetStringData(node.NamedChildren()[0], detectorContext)
	case "binary_expression":
		switch node.Children()[1].Content() {
		case "+":
			return common.ConcatenateChildStrings(node, detectorContext)
		case "||":
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
	case "augmented_assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}
