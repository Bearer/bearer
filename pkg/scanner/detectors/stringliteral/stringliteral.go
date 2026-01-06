package stringliteral

import (
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
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
	// string_fragment, string_content, and string_value are parts of template/interpolated strings,
	// not standalone literals. They should not be considered "safe" string literals
	// because the parent string may contain dynamic content.
	// Different languages use different node types for literal content:
	// - string_fragment: JavaScript, TypeScript, TSX, Java
	// - string_content: Python, Ruby
	// - string_value: PHP
	nodeType := node.Type()
	if nodeType == "string_fragment" || nodeType == "string_content" || nodeType == "string_value" {
		return nil, nil
	}

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
