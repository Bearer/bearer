package string

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"

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
	case "string_content":
		return []interface{}{generictypes.String{Value: node.Content()}}, nil
	case "interpolation", "string":
		return concatenateChildren(node, evaluator)
	case "binary":
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

func (detector *stringDetector) Close() {}
