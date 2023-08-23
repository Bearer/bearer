package string

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
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

func (detector *stringDetector) Name() string {
	return "string"
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	scanContext types.ScanContext,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_literal":
		return []interface{}{common.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, scanContext)
		}
	case "assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, scanContext)
		}
	}

	return nil, nil
}
