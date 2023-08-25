package customrule

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"

	"github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery"
)

type Data struct {
	Pattern       string
	Datatypes     []*types.Detection
	VariableNodes map[string]*tree.Node
}

type Pattern struct {
	Pattern string
	Query   patternquery.Query
	Filters []settings.PatternFilter
}

type customDetector struct {
	types.DetectorBase
	ruleID   string
	patterns []Pattern
}

func New(
	language language.Language,
	querySet *query.Set,
	ruleID string,
	patterns []settings.RulePattern,
) (types.Detector, error) {
	var compiledPatterns []Pattern
	for _, pattern := range patterns {
		patternQuery, err := patternquery.Compile(language, querySet, pattern.Pattern, pattern.Focus)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}

		sortFilters(pattern.Filters)

		compiledPatterns = append(compiledPatterns, Pattern{
			Pattern: pattern.Pattern,
			Query:   patternQuery,
			Filters: pattern.Filters,
		})
	}

	return &customDetector{
		ruleID:   ruleID,
		patterns: compiledPatterns,
	}, nil
}

func (detector *customDetector) RuleID() string {
	return detector.ruleID
}

func (detector *customDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, variableNodes, err := matchAllFilters(
				detectorContext,
				result.Variables,
				pattern.Filters,
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
