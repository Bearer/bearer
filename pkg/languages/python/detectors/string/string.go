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
	case "string_content":
		return handleTemplateString(node, detectorContext)
	case "concatenated_string", "string":
		return concatenateChildStrings(node, detectorContext)
	case "binary_operator":
		if node.Children()[1].Content() == "+" {
			return concatenateChildStrings(node, detectorContext)
		}
	case "boolean_operator":
		if node.Children()[1].Content() == "or" {
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
	case "augmented_assignment":
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

		switch {
		case child.Type() == "escape_sequence":
			// tree sitter parser doesn't handle line continuation inside a string
			if child.Content() == "\\\n" || child.Content() == "\\\r\n" {
				childValue = ""
			} else {
				value, err := stringutil.Unescape(child.Content())
				if err != nil {
					return fmt.Errorf("failed to decode escape sequence '%s': %w", child.Content(), err)
				}

				childValue = value
			}

			childIsLiteral = true
		case len(namedChildren) == 0:
			childValue = ""
			childIsLiteral = true
		default:
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

func concatenateChildStrings(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	return common.ConcatenateChildStrings(node, detectorContext, "string_start", "string_end")
}
