package string

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

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
) ([]*types.Detection, error) {
	switch node.Type() {
	case "string_content":
		return []*types.Detection{{
			MatchNode: node,
			Data:      generictypes.String{Value: node.Content()},
		}}, nil
	case "interpolation", "string":
		return concatenateChildren(node, evaluator)
	case "binary":
		if node.AnonymousChild(0).Content() == "+" {
			return concatenateChildren(node, evaluator)
		}
	}

	return nil, nil
}

func concatenateChildren(node *tree.Node, evaluator types.Evaluator) ([]*types.Detection, error) {
	value := ""

	for i := 0; i < node.ChildCount(); i += 1 {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}

		detections, err := evaluator.ForNode(child, "string")
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

	return []*types.Detection{{
		MatchNode: node,
		Data:      generictypes.String{Value: value},
	}}, nil
}

func (detector *stringDetector) Close() {}
