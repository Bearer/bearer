package string

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "assignment_statement":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	case "interpreted_string_literal", "raw_string_literal":
		value := stringutil.StripQuotes(node.Content())

		return []interface{}{common.String{
			Value:     value,
			IsLiteral: true,
		}}, nil
	}

	return nil, nil
}
