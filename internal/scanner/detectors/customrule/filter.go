package customrule

import (
	"slices"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

func matchFilter(
	detectorContext types.Context,
	variables,
	joinedVariables map[string]*tree.Node,
	filter settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	if filter.Not != nil {
		match, _, err := matchFilter(detectorContext, variables, joinedVariables, *filter.Not)
		if match == nil {
			return nil, nil, err
		}
		return boolPointer(!*match), nil, err
	}

	if len(filter.Either) != 0 {
		return matchEitherFilters(detectorContext, variables, joinedVariables, filter.Either)
	}

	if filter.FilenameRegex != nil {
		return boolPointer(filter.FilenameRegex.MatchString(detectorContext.Filename())), nil, nil
	}

	node, ok := variables[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return nil, nil, nil
	}

	if filter.Detection != "" {
		return matchDetectionFilter(
			detectorContext,
			variables,
			joinedVariables,
			node,
			filter,
		)
	}

	matched, err := matchContentFilter(detectorContext, filter, node)
	return matched, nil, err
}

func matchAllFilters(
	detectorContext types.Context,
	variables map[string]*tree.Node,
	filters []settings.PatternFilter,
) (bool, []*types.Detection, map[string]*tree.Node, error) {
	var datatypeDetections []*types.Detection

	joinedVariables := make(map[string]*tree.Node)
	for name, node := range variables {
		joinedVariables[name] = node
	}

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(detectorContext, variables, joinedVariables, filter)
		if matched == nil || !*matched || err != nil {
			return false, nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, joinedVariables, nil
}

func matchEitherFilters(
	detectorContext types.Context,
	variables,
	joinedVariables map[string]*tree.Node,
	filters []settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection
	oneMatched := false
	oneNotMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(detectorContext, variables, joinedVariables, subFilter)
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
	detectorContext types.Context,
	variables,
	joinedVariables map[string]*tree.Node,
	node *tree.Node,
	filter settings.PatternFilter,
) (*bool, []*types.Detection, error) {
	ruleID := filter.Detection
	detections, err := detectorContext.Scan(node, ruleID, filter.Scope)
	if err != nil {
		return nil, nil, err
	}

	if ruleID == "datatype" {
		return boolPointer(len(detections) != 0), detections, nil
	}

	var datatypeDetections []*types.Detection
	ignoredVariables := getIgnoredVariables(detections)
	foundDetection := false

	for _, detection := range detections {
		data, ok := detection.Data.(Data)
		if !ok { // Built-in detector
			foundDetection = true
			log.Trace().Msg("detection match (built-in)")
			continue
		}

		filtersMatch, _, _, err := matchAllFilters(detectorContext, data.VariableNodes, filter.Filters)
		if err != nil {
			return nil, nil, err
		}
		if !filtersMatch {
			log.Trace().Msg("detection filters do not match")
			continue
		}

		variablesMatch := true
		for name, node := range data.VariableNodes {
			if existingNode, existing := joinedVariables[name]; existing {
				if existingNode != node {
					variablesMatch = false
					break
				}
			}
		}

		if !variablesMatch {
			log.Trace().Msg("detection variable mismatch")
			continue
		}

		foundDetection = true
		for name, node := range data.VariableNodes {
			if _, ignored := ignoredVariables[name]; !ignored {
				joinedVariables[name] = node
			}
		}

		datatypeDetections = append(datatypeDetections, data.Datatypes...)
		log.Trace().Msg("detection match")
	}

	return boolPointer(foundDetection), datatypeDetections, nil
}

func matchContentFilter(
	detectorContext types.Context,
	filter settings.PatternFilter,
	node *tree.Node,
) (*bool, error) {
	content := node.Content()

	if len(filter.Values) != 0 {
		return boolPointer(slices.Contains(filter.Values, content)), nil
	}

	if filter.Regex != nil {
		return boolPointer(filter.Regex.MatchString(content)), nil
	}

	if filter.LengthLessThan != nil {
		strValue, _, err := common.GetStringValue(node, detectorContext)
		if err != nil || strValue == "" {
			return nil, err
		}

		if len(strValue) >= *filter.LengthLessThan {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.StringRegex != nil {
		value, isLiteral, err := common.GetStringValue(node, detectorContext)
		if err != nil || (value == "" && !isLiteral) {
			return nil, err
		}

		return boolPointer(filter.StringRegex.MatchString(value)), nil
	}

	if filter.LessThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil, nil
		}

		if value >= *filter.LessThan {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.LessThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil, nil
		}

		if value > *filter.LessThanOrEqual {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.GreaterThan != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil, nil
		}

		if value <= *filter.GreaterThan {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.GreaterThanOrEqual != nil {
		value, err := strconv.Atoi(content)
		if err != nil {
			return nil, nil
		}

		if value < *filter.GreaterThanOrEqual {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	log.Debug().Msgf("unknown filter: %#v", filter)
	return nil, nil
}

func boolPointer(value bool) *bool {
	return &value
}

func getIgnoredVariables(detections []*types.Detection) map[string]struct{} {
	ignoredVariables := make(map[string]struct{})
	seenNodes := make(map[string]*tree.Node)

	for _, detection := range detections {
		data, ok := detection.Data.(Data)
		if !ok {
			continue
		}

		for name, node := range data.VariableNodes {
			seenNode := seenNodes[name]
			if seenNode != nil && seenNode != node {
				ignoredVariables[name] = struct{}{}
			}

			seenNodes[name] = node
		}
	}

	return ignoredVariables
}
