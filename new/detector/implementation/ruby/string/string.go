package string

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *tree.QuerySet) (types.Detector, error) {
	return &stringDetector{}, nil
}

func (detector *stringDetector) Name() string {
	return "string"
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_content":
		return []interface{}{generictypes.String{
			Value:     node.Content(),
			IsLiteral: true,
		}}, nil
	case "interpolation", "string":
		return generic.ConcatenateChildStrings(node, evaluationState)
	case "binary":
		if node.AnonymousChild(0).Content() == "+" {
			return generic.ConcatenateChildStrings(node, evaluationState)
		}
	case "operator_assignment":
		if node.AnonymousChild(0).Content() == "+=" {
			return generic.ConcatenateAssignEquals(node, evaluationState)
		}
	}

	return nil, nil
}

func (detector *stringDetector) Close() {}
