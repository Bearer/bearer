package insecureurl

import (
	"regexp"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type insecureURLDetector struct {
	types.DetectorBase
}

var insecureUrlPattern = regexp.MustCompile(`^http:`)
var localhostInsecureUrlPattern = regexp.MustCompile(`^http://(localhost|127.0.0.1)`)

func New(lang languagetypes.Language) (types.Detector, error) {
	return &insecureURLDetector{}, nil
}

func (detector *insecureURLDetector) Name() string {
	return "insecure_url"
}

func (detector *insecureURLDetector) DetectAt(
	node *tree.Node,
	_ settings.RuleReferenceScope,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	detections, err := evaluator.Evaluate(node, "string", "", settings.CURSOR_SCOPE, false)
	if err != nil {
		return nil, err
	}

	for _, detection := range detections {
		value := detection.Data.(generictypes.String).Value
		if insecureUrlPattern.MatchString(value) {
			if localhostInsecureUrlPattern.MatchString(value) {
				// ignore insecure local URLs
				continue
			}

			return []interface{}{nil}, nil
		}
	}

	return nil, nil
}

func (detector *insecureURLDetector) Close() {}
