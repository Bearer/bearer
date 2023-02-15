package custom

import (
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/commands/process/settings"
)

func matchFilter(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	filter settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	if filter.Not != nil {
		match, _, err := matchFilter(result, evaluator, *filter.Not)
		if match == nil {
			return nil, nil, err
		}
		return boolPointer(!*match), nil, err
	}

	if len(filter.Either) != 0 {
		return matchEitherFilters(result, evaluator, filter.Either)
	}

	node, ok := result.Variables[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return nil, nil, nil
	}

	if filter.Detection != "" {
		return matchDetectionFilter(result, evaluator, node, filter.Detection)
	}

	return matchContentFilter(filter, node.Content()), nil, nil
}

func matchAllFilters(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(result, evaluator, filter)
		if matched == nil || !*matched || err != nil {
			return false, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, nil
}

func matchEitherFilters(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection
	oneMatched := false
	oneNotMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(result, evaluator, subFilter)
		if err != nil {
			return nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDatatypeDetections...)
		oneMatched = oneMatched || (subMatch != nil && *subMatch)
		oneNotMatched = oneNotMatched || (subMatch != nil && !*subMatch)
	}

	if oneMatched {
		return boolPointer(true), datatypeDetections, nil
	}

	if oneNotMatched {
		return boolPointer(false), nil, nil
	}

	return nil, nil, nil
}

func matchDetectionFilter(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	node *tree.Node,
	detectorType string,
) (*bool, []*types.Detection, error) {
	if detectorType == "datatype" {
		detections, err := evaluator.ForTree(node, "datatype", true)

		return boolPointer(len(detections) != 0), detections, err
	}

	hasDetection, err := evaluator.TreeHas(node, detectorType)
	return boolPointer(hasDetection), nil, err
}

func matchContentFilter(filter settings.PatternFilter, content string) *bool {
	if len(filter.Values) != 0 && !slices.Contains(filter.Values, content) {
		return boolPointer(false)
	}

	if filter.Regex != nil {
		return boolPointer(filter.Regex.MatchString(content))
	}

	if filter.LessThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil
		}

		if value >= *filter.LessThan {
			return boolPointer(false)
		}
	}

	if filter.LessThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil
		}

		if value > *filter.LessThanOrEqual {
			return boolPointer(false)
		}
	}

	if filter.GreaterThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil
		}

		if value <= *filter.GreaterThan {
			return boolPointer(false)
		}
	}

	if filter.GreaterThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil
		}

		if value < *filter.GreaterThanOrEqual {
			return boolPointer(false)
		}
	}

	return boolPointer(true)
}

func boolPointer(value bool) *bool {
	return &value
}
