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
	detectorType    string
	patterns        []Pattern
	sanitizerRuleID string
}

func New(
	lang languagetypes.Language,
	detectorType string,
	patterns []settings.RulePattern,
	sanitizerRuleID string,
) (types.Detector, error) {
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
		detectorType:    detectorType,
		patterns:        compiledPatterns,
		sanitizerRuleID: sanitizerRuleID,
	}, nil
}

func (detector *customDetector) Name() string {
	return detector.detectorType
}

func (detector *customDetector) DetectAt(
	rootNode *tree.Node,
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	sanitized, err := detector.isSanitized(rootNode, node, evaluator)
	if err != nil {
		return nil, fmt.Errorf("error running sanitizer: %w", err)
	}

	if sanitized {
		return nil, nil
	}

	var detectionsData []interface{}

	for _, pattern := range detector.patterns {
		results, err := pattern.Query.MatchAt(node)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			filtersMatch, datatypeDetections, variableNodes, err := matchAllFilters(result, evaluator, pattern.Filters)
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

func (detector *customDetector) isSanitized(rootNode, node *tree.Node, evaluator types.Evaluator) (bool, error) {
	if detector.sanitizerRuleID == "" {
		return false, nil
	}

	for ancestor := node; !ancestor.Equal(rootNode); ancestor = ancestor.Parent() {
		sanitized, err := evaluator.NodeHas(ancestor, detector.sanitizerRuleID)
		if err != nil {
			return false, err
		}

		if sanitized {
			return true, nil
		}
	}

	return false, nil
}

func (detector *customDetector) Close() {
	for _, pattern := range detector.patterns {
		pattern.Query.Close()
	}
}
