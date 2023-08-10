package custom

import (
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

func matchFilter(
	evaluationState types.EvaluationState,
	variables,
	joinedVariables map[string]*sitter.Node,
	filter settings.PatternFilter,
	rules map[string]*settings.Rule,
) (*bool, []*detection.Detection, error) {
	if filter.Not != nil {
		match, _, err := matchFilter(evaluationState, variables, joinedVariables, *filter.Not, rules)
		if match == nil {
			return nil, nil, err
		}
		return boolPointer(!*match), nil, err
	}

	if len(filter.Either) != 0 {
		return matchEitherFilters(evaluationState, variables, joinedVariables, filter.Either, rules)
	}

	if filter.FilenameRegex != nil {
		return boolPointer(filter.FilenameRegex.MatchString(evaluationState.FileName())), nil, nil
	}

	node, ok := variables[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return nil, nil, nil
	}

	if filter.Detection != "" {
		return matchDetectionFilter(
			evaluationState,
			variables,
			joinedVariables,
			node,
			filter,
			rules,
		)
	}

	matched, err := matchContentFilter(evaluationState, filter, node)
	return matched, nil, err
}

func matchAllFilters(
	evaluationState types.EvaluationState,
	variables map[string]*sitter.Node,
	filters []settings.PatternFilter,
	rules map[string]*settings.Rule,
) (bool, []*detection.Detection, map[string]*sitter.Node, error) {
	var datatypeDetections []*detection.Detection

	joinedVariables := make(map[string]*sitter.Node)
	for name, node := range variables {
		joinedVariables[name] = node
	}

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(evaluationState, variables, joinedVariables, filter, rules)
		if matched == nil || !*matched || err != nil {
			return false, nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, joinedVariables, nil
}

func matchEitherFilters(
	evaluationState types.EvaluationState,
	variables,
	joinedVariables map[string]*sitter.Node,
	filters []settings.PatternFilter,
	rules map[string]*settings.Rule,
) (*bool, []*detection.Detection, error) {
	var datatypeDetections []*detection.Detection
	oneMatched := false
	oneNotMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(evaluationState, variables, joinedVariables, subFilter, rules)
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
	evaluationState types.EvaluationState,
	variables,
	joinedVariables map[string]*sitter.Node,
	node *sitter.Node,
	filter settings.PatternFilter,
	rules map[string]*settings.Rule,
) (*bool, []*detection.Detection, error) {
	ruleID := filter.Detection
	sanitizerRuleID := ""
	if rule, ok := rules[ruleID]; ok {
		sanitizerRuleID = rule.SanitizerRuleID
	}

	detections, err := evaluationState.Evaluate(
		evaluationState.NodeFromSitter(node),
		ruleID,
		sanitizerRuleID,
		filter.Scope,
		true,
	)

	if ruleID == "datatype" {
		return boolPointer(len(detections) != 0), detections, err
	}

	var datatypeDetections []*detection.Detection
	ignoredVariables := getIgnoredVariables(detections)
	foundDetection := false

	for _, detection := range detections {
		data, ok := detection.Data.(Data)
		if !ok { // Built-in detector
			foundDetection = true
			continue
		}

		filtersMatch, _, _, err := matchAllFilters(evaluationState, data.VariableNodes, filter.Filters, rules)
		if err != nil {
			return nil, nil, err
		}
		if !filtersMatch {
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
			continue
		}

		foundDetection = true
		for name, node := range data.VariableNodes {
			if _, ignored := ignoredVariables[name]; !ignored {
				joinedVariables[name] = node
			}
		}

		datatypeDetections = append(datatypeDetections, data.Datatypes...)
	}

	return boolPointer(foundDetection), datatypeDetections, err
}

func matchContentFilter(
	evaluationState types.EvaluationState,
	filter settings.PatternFilter,
	sitterNode *sitter.Node,
) (*bool, error) {
	node := evaluationState.NodeFromSitter(sitterNode)
	content := node.Content()

	if len(filter.Values) != 0 {
		return boolPointer(slices.Contains(filter.Values, content)), nil
	}

	if filter.Regex != nil {
		return boolPointer(filter.Regex.MatchString(content)), nil
	}

	if filter.LengthLessThan != nil {
		strValue, _, err := generic.GetStringValue(node, evaluationState)
		if err != nil || strValue == "" {
			return nil, err
		}

		if len(strValue) >= *filter.LengthLessThan {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.StringRegex != nil {
		value, isLiteral, err := generic.GetStringValue(node, evaluationState)
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

func getIgnoredVariables(detections []*detection.Detection) map[string]struct{} {
	ignoredVariables := make(map[string]struct{})
	seenNodes := make(map[string]*sitter.Node)

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
