package customrule

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/ruleset"

	"github.com/bearer/bearer/internal/scanner/detectors/customrule/filters"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule/types"
)

type Pattern struct {
	Index   int
	Pattern string
	Query   patternquery.Query
	Filters []filters.Filter
}

type Detector struct {
	detectortypes.DetectorBase
	rule     *ruleset.Rule
	patterns []Pattern
}

func New(
	language language.Language,
	ruleSet *ruleset.Set,
	querySet *query.Set,
	rule *ruleset.Rule,
) (detectortypes.Detector, error) {
	var compiledPatterns []Pattern
	for i, pattern := range rule.Patterns() {
		patternQuery, err := patternquery.Compile(language, querySet, rule.ID(), i, pattern.Pattern, pattern.Focus)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}

		filters, err := translateFilters(ruleSet, pattern.Filters)
		if err != nil {
			return nil, err
		}

		compiledPatterns = append(compiledPatterns, Pattern{
			Index:   i,
			Pattern: pattern.Pattern,
			Query:   patternQuery,
			Filters: filters,
		})
	}

	return &Detector{
		patterns: compiledPatterns,
		rule:     rule,
	}, nil
}

func (detector *Detector) Rule() *ruleset.Rule {
	return detector.rule
}

func (detector *Detector) DetectAt(
	node *tree.Node,
	detectorContext detectortypes.Context,
) ([]interface{}, error) {
	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		if log.Trace().Enabled() && len(results) != 0 {
			log.Trace().Msgf("pattern %s matched (without filters)", pattern.Query.ID())
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, variableNodes, err := filters.Match(
				detectorContext,
				result.Variables,
				pattern.Filters,
			)
			if err != nil {
				return nil, err
			}

			if !filtersMatch {
				log.Trace().Msg("filters didn't match")
				continue
			}

			detectionsData = append(detectionsData, types.Data{
				Pattern:       pattern.Pattern,
				Datatypes:     datatypeDetections,
				VariableNodes: variableNodes,
			})

			log.Trace().Msg("filters matched")
		}
	}

	return detectionsData, nil
}
