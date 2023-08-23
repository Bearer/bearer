package string

import (
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

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
	case "string_content":
		return []interface{}{common.String{
			Value:     node.Content(),
			IsLiteral: true,
		}}, nil
	case "interpolation", "string":
		return common.ConcatenateChildStrings(node, scanContext)
	case "binary":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, scanContext)
		}
	case "operator_assignment":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, scanContext)
		}
	}

	return nil, nil
}
