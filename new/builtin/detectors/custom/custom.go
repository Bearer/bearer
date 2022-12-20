package custom

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
	"golang.org/x/exp/slices"
)

// type Data struct {
// 	Name string
// }

type Filter struct {
	Variable         string
	Values           []string
	ContainsDataType bool
}

type Rule struct {
	Pattern string
	Filters []Filter
}

type customDetector struct {
	detectorType string
	patternQuery languagetypes.PatternQuery
	filters      []Filter
}

func New(lang languagetypes.Language, detectorType string, rule Rule) (detector.Detector, error) {
	patternQuery, err := lang.CompilePatternQuery(rule.Pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling pattern: %s", err)
	}

	// TODO: validate filters against pattern

	return &customDetector{
		detectorType: detectorType,
		patternQuery: patternQuery,
		filters:      rule.Filters,
	}, nil
}

func (detector *customDetector) Name() string {
	return detector.detectorType
}

func (detector *customDetector) DetectAt(
	node *tree.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	results, err := detector.patternQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}

	var detections []*detectiontypes.Detection

	for _, result := range results {
		filtersMatch, err := detector.filtersMatch(result, evaluator)
		if err != nil {
			return nil, err
		}

		if !filtersMatch {
			continue
		}

		detections = append(detections, &detectiontypes.Detection{
			MatchNode: node,
			// TODO: include data types
		})
	}

	return detections, nil
}

func (detector *customDetector) filtersMatch(
	result tree.QueryResult,
	evaluator treeevaluatortypes.Evaluator,
) (bool, error) {
	for _, filter := range detector.filters {
		node, ok := result[filter.Variable]
		// shouldn't happen if filters are validated against pattern
		if !ok {
			return false, nil
		}

		if len(filter.Values) != 0 && !slices.Contains(filter.Values, node.Content()) {
			return false, nil
		}

		if filter.ContainsDataType {
			hasDataType, err := evaluator.TreeHasDetection(node, "datatype")
			if err != nil || !hasDataType {
				return false, err
			}
		}
	}

	return true, nil
}

func (detector *customDetector) Close() {
	detector.patternQuery.Close()
}
