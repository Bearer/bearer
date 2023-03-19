package custom

import (
	"fmt"
	"strings"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	soufflequery "github.com/bearer/bearer/pkg/souffle/query"
	"github.com/bearer/bearer/pkg/util/output"
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
	types.DetectorBase
	detectorType string
	patterns     []Pattern
}

func New(
	lang languagetypes.Language,
	detectorType string,
	patterns []settings.RulePattern,
) (types.Detector, error) {
	var compiledPatterns []Pattern
	for _, pattern := range patterns {
		var patternQuery languagetypes.PatternQuery
		if !strings.HasPrefix(detectorType, "ruby_") {
			var err error
			patternQuery, err = lang.CompilePatternQuery(detectorType, pattern.Pattern)
			if err != nil {
				return nil, fmt.Errorf("error compiling pattern: %s", err)
			}
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
	queryContext *soufflequery.QueryContext,
) ([]interface{}, error) {
	var detectionsData []interface{}

	for i, pattern := range detector.patterns {
		var results []*languagetypes.PatternQueryResult

		if queryContext == nil {
			var err error
			results, err = pattern.Query.MatchAt(node)
			if err != nil {
				return nil, err
			}
		} else {
			results = queryContext.MatchAt(idgenerator.PatternId(detector.detectorType, i), node)
		}

		for _, result := range results {
			output.StdErrLogger().Msgf("found pattern %s at: %s", idgenerator.PatternId(detector.detectorType, i), result.MatchNode.Content())
			filtersMatch, datatypeDetections, err := matchAllFilters(result, evaluator, pattern.Filters)
			if err != nil {
				return nil, err
			}

			if !filtersMatch {
				continue
			}

			detectionsData = append(detectionsData, Data{
				Pattern:   pattern.Pattern,
				Datatypes: datatypeDetections,
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
