package stringliteral

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type stringLiteralDetector struct {
	types.DetectorBase
}

func New(lang languagetypes.Language) (types.Detector, error) {
	return &stringLiteralDetector{}, nil
}

func (detector *stringLiteralDetector) Name() string {
	return "string_literal"
}

func (detector *stringLiteralDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	detections, err := evaluator.ForNode(node, "string", false)
	if err != nil {
		return nil, err
	}

	for _, detection := range detections {
		if detection.Data.(generictypes.String).IsLiteral {
			if len(detection.Data.(generictypes.String).Value) > 0 {
				return []interface{}{nil}, nil
			}
		}
	}

	return nil, nil
}

func (detector *stringLiteralDetector) Close() {}
