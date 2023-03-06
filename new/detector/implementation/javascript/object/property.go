package object

import (
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
)

func getNonVirtualObjects(evaluator types.Evaluator, node *tree.Node) ([]*types.Detection, error) {
	detections, err := evaluator.ForNode(node, "object", true)
	if err != nil {
		return nil, err
	}

	var result []*types.Detection
	for _, detection := range detections {
		data := detection.Data.(generictypes.Object)
		if !data.IsVirtual {
			result = append(result, detection)
		}
	}

	return result, nil
}
