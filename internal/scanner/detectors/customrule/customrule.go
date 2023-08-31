package customrule

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"

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
	ruleID              string
	sanitizerDetectorID int
	patterns            []Pattern
}

func New(
	language language.Language,
	querySet *query.Set,
	detectorIDByRuleID map[string]int,
	ruleID,
	sanitizerRuleID string,
	patterns []settings.RulePattern,
) (detectortypes.Detector, error) {
	var compiledPatterns []Pattern
	for i, pattern := range patterns {
		patternQuery, err := patternquery.Compile(language, querySet, ruleID, i, pattern.Pattern, pattern.Focus)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", err)
		}

		filters, err := translateFilters(detectorIDByRuleID, pattern.Filters)
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

	sanitizerDetectorID := -1
	if sanitizerRuleID != "" {
		sanitizerDetectorID = detectorIDByRuleID[sanitizerRuleID]
	}

	return &Detector{
		ruleID:              ruleID,
		sanitizerDetectorID: sanitizerDetectorID,
		patterns:            compiledPatterns,
	}, nil
}

func (detector *Detector) RuleID() string {
	return detector.ruleID
}

func (detector *Detector) SanitizerDetectorID() int {
	return detector.sanitizerDetectorID
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
