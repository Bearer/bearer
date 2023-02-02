package string

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/pkg/util/stringutil"
	"github.com/rs/zerolog/log"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type stringDetector struct{}

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
	default:
		log.Debug().Msgf("got node: %s, with content: %s", node.Type(), node.Content())
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

		detections, err := evaluator.ForNode(child, "string", true)
		if err != nil {
			return nil, err
		}

		switch len(detections) {
		case 0:
			value += "*"
		case 1:
			value += detections[0].Data.(generictypes.String).Value
		default:
			return nil, fmt.Errorf("expected single string detection but got %d", len(detections))
		}
	}

	return []interface{}{generictypes.String{Value: value}}, nil
}

func handleTemplateString(node *tree.Node, evaluator types.Evaluator) ([]interface{}, error) {
	text := ""

	node.EachPart(func(partText string) error {
		text += partText
		return nil
	}, func(child *tree.Node) error {
		text += "*"
		return nil
	})

	return []interface{}{generictypes.String{Value: text}}, nil
}

func (detector *stringDetector) Close() {}
