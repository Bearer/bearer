package custom

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"golang.org/x/exp/slices"
)

type Data struct {
	Pattern       string
	Datatypes     []*detection.Detection
	VariableNodes map[string]*tree.Node
}

type Pattern struct {
	Pattern string
	Query   languagetypes.PatternQuery
	Filters []settings.PatternFilter
}

type customDetector struct {
	types.DetectorBase
	detectorType string
	patterns     []Pattern
	rules        map[string]*settings.Rule
}

func New(
	lang languagetypes.Language,
	detectorType string,
	patterns []settings.RulePattern,
	rules map[string]*settings.Rule,
) (types.Detector, error) {
	var compiledPatterns []Pattern
	for _, pattern := range patterns {
		patternQuery, err := lang.CompilePatternQuery(pattern.Pattern, pattern.Focus)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}

		sortFilters(pattern.Filters)

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
		rules:        rules,
	}, nil
}

func (detector *customDetector) Name() string {
	return detector.detectorType
}

func (detector *customDetector) DetectAt(
	node *tree.Node,
	evaluationState types.EvaluationState,
) ([]interface{}, error) {
	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, variableNodes, err := matchAllFilters(
				evaluationState,
				result,
				pattern.Filters,
				detector.rules,
			)
			if err != nil {
				return nil, err
			}

			if !filtersMatch {
				continue
			}

			detectionsData = append(detectionsData, Data{
				Pattern:       pattern.Pattern,
				Datatypes:     datatypeDetections,
				VariableNodes: variableNodes,
			})
		}
	}

	return detectionsData, nil
}

func (detector *customDetector) Close() {
	for _, pattern := range detector.patterns {
		pattern.Query.Close()
	}
}

func sortFilters(filters []settings.PatternFilter) {
	slices.SortFunc(filters, func(a, b settings.PatternFilter) bool {
		return scoreFilter(a) < scoreFilter(b)
	})

	for i := range filters {
		sortFilter(&filters[i])
	}
}

func sortFilter(filter *settings.PatternFilter) {
	switch {
	case len(filter.Either) != 0:
		sortFilters(filter.Either)
	case filter.Not != nil:
		sortFilter(filter.Not)
	}
}

func scoreFilter(filter settings.PatternFilter) int {
	if filter.Regex != nil ||
		len(filter.Values) != 0 ||
		filter.LengthLessThan != nil ||
		filter.LessThan != nil ||
		filter.LessThanOrEqual != nil ||
		filter.GreaterThan != nil ||
		filter.GreaterThanOrEqual != nil ||
		filter.FilenameRegex != nil {
		return 1
	}

	if filter.Detection == "datatype" {
		return 5
	}

	if filter.StringRegex != nil ||
		filter.Detection != "" && filter.Scope == settings.CURSOR_SCOPE {
		return 2
	}

	if filter.Detection != "" && filter.Scope == settings.RESULT_SCOPE {
		return 3
	}

	if filter.Detection != "" && filter.Scope == settings.NESTED_SCOPE {
		return 4
	}

	if filter.Not != nil {
		return scoreFilter(*filter.Not)
	}

	if len(filter.Either) != 0 {
		max := 0

		for _, subFilter := range filter.Either {
			if subScore := scoreFilter(subFilter); subScore > max {
				max = subScore
			}
		}

		return max
	}

	panic(fmt.Sprintf("unknown filter %#v", filter))
}
