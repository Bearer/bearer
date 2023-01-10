package custom

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/commands/process/settings"
)

type Data struct {
	Datatypes []*types.Detection
}

type Pattern struct {
	Query   languagetypes.PatternQuery
	Filters []settings.PatternFilter
}

type customDetector struct {
	detectorType string
	patterns     []Pattern
}

func New(lang languagetypes.Language, detectorType string, patterns []settings.RulePattern) (types.Detector, error) {
	var compiledPatterns []Pattern
	for _, pattern := range patterns {
		patternQuery, err := lang.CompilePatternQuery(pattern.Pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}
		compiledPatterns = append(compiledPatterns, Pattern{
			Query:   patternQuery,
			Filters: pattern.Filters,
		})

		// TODO: validate filters against pattern
	}

	return &customDetector{
		detectorType: detectorType,
		patterns:     compiledPatterns,
	}, nil
}

func (detector *customDetector) Name() string {
	return detector.detectorType
}

func (detector *customDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	var detections []*types.Detection
	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, err := detector.matchFilters(result, evaluator, pattern.Filters)
			if err != nil {
				return nil, err
			}

			if !filtersMatch {
				continue
			}

			detections = append(detections, &types.Detection{
				MatchNode: node,
				Data: Data{
					Datatypes: datatypeDetections,
				},
			})
		}
	}

	return detections, nil
}

func (detector *customDetector) matchFilters(
	result tree.QueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection

	for _, filter := range filters {
		node, ok := result[filter.Variable]
		// shouldn't happen if filters are validated against pattern
		if !ok {
			return false, nil, nil
		}

		if len(filter.Values) != 0 && !slices.Contains(filter.Values, node.Content()) {
			return false, nil, nil
		}

		if filter.Detection == "datatype" {
			filterDetections, err := evaluator.ForTree(node, "datatype")
			if err != nil {
				return false, nil, err
			}

			datatypeDetections = append(datatypeDetections, filterDetections...)
		} else if filter.Detection != "" {
			hasDetection, err := evaluator.TreeHas(node, filter.Detection)
			if err != nil || !hasDetection {
				return false, nil, err
			}
		}
	}

	return true, datatypeDetections, nil
}

func (detector *customDetector) Close() {
	for _, pattern := range detector.patterns {
		pattern.Query.Close()
	}
}
