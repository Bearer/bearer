package insecureurl

import (
	"regexp"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type insecureURLDetector struct{}

var insecureUrlPattern = regexp.MustCompile(`^http:`)

func New(lang languagetypes.Language) (types.Detector, error) {
	return &insecureURLDetector{}, nil
}

func (detector *insecureURLDetector) Name() string {
	return "insecure_url"
}

func (detector *insecureURLDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	detections, err := evaluator.ForNode(node, "string")
	if err != nil {
		return nil, err
	}

	for _, detection := range detections {
		value := detection.Data.(generictypes.String).Value

		if insecureUrlPattern.MatchString(value) {
			return []*types.Detection{{MatchNode: node}}, nil
		}
	}

	return nil, nil
}

func (detector *insecureURLDetector) Close() {}
