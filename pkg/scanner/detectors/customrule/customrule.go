package customrule

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/scanner/variableshape"

	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/filters"
	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery"
	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/types"
)

type Pattern struct {
	Index   int
	Pattern string
	Query   patternquery.Query
	Filter  filters.Filter
}

type Detector struct {
	detectortypes.DetectorBase
	rule     *ruleset.Rule
	patterns []Pattern
}

func New(
	language language.Language,
	ruleSet *ruleset.Set,
	variableShapeSet *variableshape.Set,
	querySet *query.Set,
	rule *ruleset.Rule,
) (detectortypes.Detector, error) {
	variableShape := variableShapeSet.Shape(rule)

	var compiledPatterns []Pattern
	for i, pattern := range rule.Patterns() {
		patternQuery, err := patternquery.Compile(
			language,
			querySet,
			rule.ID(),
			i,
			pattern.Pattern,
			pattern.Focus,
			variableShape,
		)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}

		filter, err := translateFiltersTop(ruleSet, variableShapeSet, variableShapeSet.Shape(rule), pattern.Filters)
		if err != nil {
			return nil, err
		}

		compiledPatterns = append(compiledPatterns, Pattern{
			Index:   i,
			Pattern: pattern.Pattern,
			Query:   patternQuery,
			Filter:  filter,
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
	if log.Trace().Enabled() {
		log.Trace().Msgf("DetectAt: rule=%s, node=%s", detector.rule.ID(), node.Debug())
	}
	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		if log.Trace().Enabled() {
			if len(results) != 0 {
				log.Trace().Msgf("pattern %s matched (without filters) at node=%s", pattern.Query.ID(), node.Debug())
			} else {
				log.Trace().Msgf("pattern %s did not match at node=%s", pattern.Query.ID(), node.Debug())
			}
		}

		for _, result := range results {
			filterResult, err := pattern.Filter.Evaluate(detectorContext, result.Variables)
			if err != nil {
				return nil, err
			}
			if filterResult == nil || len(filterResult.Matches()) == 0 {
				log.Trace().Msg("filters didn't match")
				continue
			}

			for _, match := range filterResult.Matches() {
				detectionsData = append(detectionsData, types.Data{
					Pattern:   pattern.Pattern,
					Datatypes: match.DatatypeDetections(),
					Variables: match.Variables(),
					Value:     match.Value(),
				})
			}

			log.Trace().Msg("filters matched")
		}
	}

	return detectionsData, nil
}
