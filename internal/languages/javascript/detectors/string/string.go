package string

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Name() string {
	return "string"
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	switch node.Type() {
	case "string":
		return []interface{}{common.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "template_string":
		return handleTemplateString(node, scanContext)
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, scanContext)
		}
	case "augmented_assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, scanContext)
		}
	}

	return nil, nil
}

func handleTemplateString(node *tree.Node, scanContext types.ScanContext) ([]interface{}, error) {
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
			childValue, childIsLiteral, err = common.GetStringValue(namedChildren[0], scanContext)
			if err != nil {
				return err
			}
		}

		if childValue == "" && !childIsLiteral {
			childValue = "*"
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
