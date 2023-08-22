package stringliteral

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type stringLiteralDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringLiteralDetector{}
}

func (detector *stringLiteralDetector) Name() string {
	return "string_literal"
}

func (detector *stringLiteralDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	detections, err := evaluationState.Evaluate(node, "string", "", settings.CURSOR_STRICT_SCOPE)
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
