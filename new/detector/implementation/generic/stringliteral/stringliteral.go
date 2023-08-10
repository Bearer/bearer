package stringliteral

import (
	"github.com/bearer/bearer/new/detector/types"
	langtree "github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type stringLiteralDetector struct {
	types.DetectorBase
}

func New(querySet *langtree.QuerySet) (types.Detector, error) {
	return &stringLiteralDetector{}, nil
}

func (detector *stringLiteralDetector) Name() string {
	return "string_literal"
}

func (detector *stringLiteralDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	detections, err := evaluationState.Evaluate(node, "string", "", settings.CURSOR_SCOPE, false)
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
