package string

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(lang languagetypes.Language) (types.Detector, error) {
	return &stringDetector{}, nil
}

func (detector *stringDetector) Name() string {
	return "string"
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	ruleReferenceType settings.RuleReferenceType,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	switch node.Type() {
	case "string":
		return []interface{}{generictypes.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "template_string":
		return handleTemplateString(node, evaluator)
	case "binary_expression":
		if node.AnonymousChild(0).Content() == "+" {
			return concatenateChildren(node, evaluator)
		}
	}

	return nil, nil
}

func concatenateChildren(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	value := ""
	isLiteral := true

	for i := 0; i < node.ChildCount(); i += 1 {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}

		childValue, childIsLiteral, err := generic.GetStringValue(child, evaluator)
		if err != nil {
			return nil, err
		}

		if childValue == "" && !childIsLiteral {
			childValue = "*"
		}

		value += childValue

		if !childIsLiteral {
			isLiteral = false
		}
	}

	return []interface{}{generictypes.String{
		Value:     value,
		IsLiteral: isLiteral,
	}}, nil
}

func handleTemplateString(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	text := ""
	isLiteral := true

	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		childValue, childIsLiteral, err := generic.GetStringValue(child.Child(1), evaluator)
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

func (detector *stringDetector) Close() {}
