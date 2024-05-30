package insecureurl

import (
	"regexp"

	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
)

type insecureURLDetector struct {
	types.DetectorBase
}

var insecureUrlPattern = regexp.MustCompile(`^http:`)
var localhostInsecureUrlPattern = regexp.MustCompile(`^http://(localhost|127.0.0.1)`)

func New(querySet *query.Set) types.Detector {
	return &insecureURLDetector{}
}

func (detector *insecureURLDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinInsecureURLRule
}

func (detector *insecureURLDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	detections, err := detectorContext.Scan(node, ruleset.BuiltinStringRule, traversalstrategy.CursorStrict)
	if err != nil {
		return nil, err
	}

	for _, detection := range detections {
		value := detection.Data.(common.String).Value
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
