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

	// Process string content by iterating over children and extracting literal text between them
	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		var childValue string
		var childIsLiteral bool

		switch child.Type() {
		case "escape_sequence":
			// Handle line continuation inside a string
			if child.Content() == "\\\n" || child.Content() == "\\\r\n" {
				childValue = ""
			} else {
				// Use Python-specific unescaping which handles unknown escapes
				childValue = stringutil.UnescapePython(child.Content())
			}
			childIsLiteral = true
		case "interpolation":
			// Interpolation like {var} - get the value from the expression inside
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
		default:
			// Other child types - try to get their string value
			var err error
			childValue, childIsLiteral, err = common.GetStringValue(child, detectorContext)
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

	if err != nil {
		return nil, err
	}

	return []interface{}{common.String{
		Value:     text,
		IsLiteral: isLiteral,
	}}, nil
}

func concatenateChildStrings(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	return common.ConcatenateChildStrings(node, detectorContext, "string_start", "string_end")
}
