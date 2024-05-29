package string

import (
	"fmt"

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
	case "string":
		return common.ConcatenateChildStrings(node, detectorContext)
	case "string_fragment":
		return []interface{}{common.String{
			Value:     node.Content(),
			IsLiteral: true,
		}}, nil
	case "escape_sequence":
		value, err := stringutil.Unescape(node.Content())
		if err != nil {
			return nil, fmt.Errorf("failed to decode escape sequence: %w", err)
		}

		return []interface{}{common.String{
			Value:     value,
			IsLiteral: true,
		}}, nil
	case "template_string":
		return handleTemplateString(node, detectorContext)
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

func handleTemplateString(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	text := ""
	isLiteral := true

	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		var childValue string
		var childIsLiteral bool
		namedChildren := child.NamedChildren()

		if len(namedChildren) == 0 {
			childValue = ""
			childIsLiteral = true
		} else {
			var err error
			childValue, childIsLiteral, err = common.GetStringValue(namedChildren[0], detectorContext)
			if err != nil {
				return err
			}
		}

		if childValue == "" && !childIsLiteral {
			childValue = common.NonLiteralValue
		}

		text += childValue

		if !childIsLiteral {
			isLiteral = false
		}

		return nil
	})

	return []interface{}{common.String{
		Value:     text,
		IsLiteral: isLiteral,
	}}, err
}
