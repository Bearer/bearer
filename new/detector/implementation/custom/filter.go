package custom

import (
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/pkg/commands/process/settings"
)

func matchFilter(
	result tree.QueryResult,
	evaluator types.Evaluator,
	filter settings.PatternFilter,
) (bool, []*types.Detection, error) {
	if len(filter.Or) != 0 {
		return matchOrFilter(result, evaluator, filter.Or)
	}

	node, ok := result[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return false, nil, nil
	}

	if filter.Detection != "" {
		return matchDetectionFilter(result, evaluator, node, filter.Detection)
	}

	return matchContentFilter(filter, node.Content()), nil, nil
}

func matchAndFilter(
	result tree.QueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(result, evaluator, filter)
		if !matched || err != nil {
			return false, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, nil
}

func matchOrFilter(
	result tree.QueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection
	oneMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(result, evaluator, subFilter)
		if err != nil {
			return false, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDatatypeDetections...)
		oneMatched = oneMatched || subMatch
	}

	if !oneMatched {
		return false, nil, nil
	}

	return true, datatypeDetections, nil
}

func matchDetectionFilter(
	result tree.QueryResult,
	evaluator types.Evaluator,
	node *tree.Node,
	detectorType string,
) (bool, []*types.Detection, error) {
	if detectorType == "datatype" {
		detections, err := evaluator.ForTree(node, "datatype")

		return len(detections) != 0, detections, err
	}

	hasDetection, err := evaluator.TreeHas(node, detectorType)
	return hasDetection, nil, err
}

func matchContentFilter(filter settings.PatternFilter, content string) bool {
	if len(filter.Values) != 0 && !slices.Contains(filter.Values, content) {
		return false
	}

	if filter.LessThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return false
		}

		if value >= *filter.LessThan {
			return false
		}
	}

	if filter.LessThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return false
		}

		if value > *filter.LessThanOrEqual {
			return false
		}
	}

	if filter.GreaterThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return false
		}

		if value <= *filter.GreaterThan {
			return false
		}
	}

	if filter.GreaterThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return false
		}

		if value < *filter.GreaterThanOrEqual {
			return false
		}
	}

	return true
}
