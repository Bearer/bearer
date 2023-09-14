package stringliteral

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

type stringLiteralDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringLiteralDetector{}
}

func (detector *stringLiteralDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringLiteralRule
}

func (detector *stringLiteralDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	detections, err := detectorContext.Scan(node, ruleset.BuiltinStringRule, traversalstrategy.CursorStrict)
	if err != nil {
		return nil, err
	}

	for _, detection := range detections {
		data := detection.Data.(common.String)
		if data.IsLiteral {
			if len(data.Value) > 0 {
				return []interface{}{nil}, nil
			}
		}
	}

	return nil, nil
}
