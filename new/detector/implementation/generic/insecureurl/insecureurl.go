package insecureurl

import (
	"regexp"

	stringdetector "github.com/bearer/curio/new/detector/implementation/ruby/string"
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
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
		value := detection.Data.(stringdetector.Data).Value

		if insecureUrlPattern.MatchString(value) {
			return []*types.Detection{{MatchNode: node}}, nil
		}
	}

	return nil, nil
}

func (detector *insecureURLDetector) Close() {}
