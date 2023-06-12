package custom

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

type Data struct {
	Pattern       string
	Datatypes     []*types.Detection
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
	scope settings.RuleReferenceScope,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, variableNodes, err := matchAllFilters(
				scope,
				result,
				evaluator,
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
