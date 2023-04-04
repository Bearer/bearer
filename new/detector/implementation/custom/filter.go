package custom

import (
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

func matchFilter(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	variableNodes map[string]*tree.Node,
	filter settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	if filter.Not != nil {
		match, _, err := matchFilter(result, evaluator, variableNodes, *filter.Not)
		if match == nil {
			return nil, nil, err
		}
		return boolPointer(!*match), nil, err
	}

	if len(filter.Either) != 0 {
		return matchEitherFilters(result, evaluator, variableNodes, filter.Either)
	}

	node, ok := result.Variables[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return nil, nil, nil
	}

	if filter.Detection != "" {
		return matchDetectionFilter(
			result,
			evaluator,
			variableNodes,
			node,
			filter.Detection,
			filter.Contains == nil || *filter.Contains,
		)
	}

	return matchContentFilter(filter, node.Content()), nil, nil
}

func matchAllFilters(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, map[string]*tree.Node, error) {
	var datatypeDetections []*types.Detection

	variableNodes := make(map[string]*tree.Node)
	for name, node := range result.Variables {
		variableNodes[name] = node
	}

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(result, evaluator, variableNodes, filter)
		if matched == nil || !*matched || err != nil {
			return false, nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, variableNodes, nil
}

func matchEitherFilters(
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	variableNodes map[string]*tree.Node,
	filters []settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection
	oneMatched := false
	oneNotMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(result, evaluator, variableNodes, subFilter)
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
	variableNodes map[string]*tree.Node,
	node *tree.Node,
	detectorType string,
	contains bool,
) (*bool, []*types.Detection, error) {
	var evaluateDetections func(*tree.Node, string, bool) ([]*types.Detection, error)
	if contains {
		evaluateDetections = evaluator.ForTree
	} else {
		evaluateDetections = evaluator.ForNode
	}

	if detectorType == "datatype" {
		detections, err := evaluateDetections(node, "datatype", true)

		return boolPointer(len(detections) != 0), detections, err
	}

	detections, err := evaluateDetections(node, detectorType, true)

	var datatypeDetections []*types.Detection

	foundDetection := false
	for _, detection := range detections {
		data, ok := detection.Data.(Data)
		if !ok { // Built-in detector
			foundDetection = true
			continue
		}

		variablesMatch := true
		for name, node := range data.VariableNodes {
			if existingNode, existing := variableNodes[name]; existing {
				if !existingNode.Equal(node) {
					variablesMatch = false
					break
				}
			}
		}

		if !variablesMatch {
			continue
		}

		foundDetection = true
		for name, node := range data.VariableNodes {
			variableNodes[name] = node
		}

		datatypeDetections = append(datatypeDetections, data.Datatypes...)
	}

	return boolPointer(foundDetection), datatypeDetections, err
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
