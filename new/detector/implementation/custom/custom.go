package custom

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/commands/process/settings"
)

type Data struct {
	Pattern   string
	Datatypes []*types.Detection
}

type Pattern struct {
	Pattern string
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
			Pattern: pattern.Pattern,
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
			filtersMatch, datatypeDetections, err := matchAllFilters(result, evaluator, pattern.Filters)
			if err != nil {
				return nil, err
			}

			if !filtersMatch {
				continue
			}

			detections = append(detections, &types.Detection{
				MatchNode: result.MatchNode,
				Data: Data{
					Pattern:   pattern.Pattern,
					Datatypes: datatypeDetections,
				},
			})
		}
	}

	return detections, nil
}

func (detector *customDetector) Close() {
	for _, pattern := range detector.patterns {
		pattern.Query.Close()
	}
}
