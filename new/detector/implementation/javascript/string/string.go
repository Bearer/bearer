package string

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
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
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	switch node.Type() {
	case "string":
		return []interface{}{generictypes.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "template_string":
		return handleTemplateString(node, evaluationState)
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return generic.ConcatenateChildStrings(node, evaluationState)
		}
	case "augmented_assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return generic.ConcatenateAssignEquals(node, evaluationState)
		}
	}

	return nil, nil
}

func handleTemplateString(node *tree.Node, evaluationState types.EvaluationState) ([]interface{}, error) {
	text := ""
	isLiteral := true

	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		childValue, childIsLiteral, err := generic.GetStringValue(child.Children()[1], evaluationState)
		if err != nil {
			return err
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

	return []interface{}{generictypes.String{
		Value:     text,
		IsLiteral: isLiteral,
	}}, err
}
