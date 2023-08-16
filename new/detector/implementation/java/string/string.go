package string

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
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
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_literal":
		return []interface{}{generictypes.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return generic.ConcatenateChildStrings(node, evaluationState)
		}
	case "assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return generic.ConcatenateAssignEquals(node, evaluationState)
		}
	}

	return nil, nil
}
