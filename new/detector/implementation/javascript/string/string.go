package string

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"
	"github.com/rs/zerolog/log"

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
	evaluator types.Evaluator,
) ([]interface{}, error) {
	switch node.Type() {
	case "string":
		return []interface{}{generictypes.String{Value: stringutil.StripQuotes(node.Content())}}, nil
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

	for i := 0; i < node.ChildCount(); i += 1 {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}

		childValue, err := getStringValue(child, evaluator)
		if err != nil {
			return nil, err
		}

		value += childValue
	}

	return []interface{}{generictypes.String{Value: value}}, nil
}

func handleTemplateString(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	text := ""

	err := node.EachContentPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		childValue, err := getStringValue(child.Child(1), evaluator)
		if err != nil {
			return err
		}

		text += childValue

		return nil
	})

	log.Debug().Msgf("node is %s", node.Debug())

	return []interface{}{generictypes.String{Value: text}}, err
}

func getStringValue(node *tree.Node, evaluator types.Evaluator) (string, error) {
	detections, err := evaluator.ForNode(node, "string", true)
	if err != nil {
		return "", err
	}

	switch len(detections) {
	case 0:
		return "*", nil
	case 1:
		return detections[0].Data.(generictypes.String).Value, nil
	default:
		return "", fmt.Errorf("expected single string detection but got %d", len(detections))
	}
}

func (detector *stringDetector) Close() {}
