package custom

import (
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/implementation/generic"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/rs/zerolog/log"
)

func matchFilter(
	scope settings.RuleReferenceScope,
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	variableNodes map[string]*tree.Node,
	filter settings.PatternFilter,
	rules map[string]*settings.Rule,
) (*bool, []*types.Detection, error) {
	if filter.Not != nil {
		match, _, err := matchFilter(scope, result, evaluator, variableNodes, *filter.Not, rules)
		if match == nil {
			return nil, nil, err
		}
		return boolPointer(!*match), nil, err
	}

	if len(filter.Either) != 0 {
		return matchEitherFilters(scope, result, evaluator, variableNodes, filter.Either, rules)
	}

	if filter.FilenameRegex != nil {
		return boolPointer(filter.FilenameRegex.MatchString(evaluator.FileName())), nil, nil
	}

	node, ok := result.Variables[filter.Variable]
	// shouldn't happen if filters are validated against pattern
	if !ok {
		return nil, nil, nil
	}

	if filter.Detection != "" {
		effectiveScope := filter.Scope
		if effectiveScope == settings.NESTED_SCOPE && scope == settings.RESULT_SCOPE {
			effectiveScope = settings.RESULT_SCOPE
		}

		return matchDetectionFilter(
			result,
			evaluator,
			variableNodes,
			node,
			filter.Detection,
			effectiveScope,
			filter.Contains == nil || *filter.Contains,
			rules,
		)
	}

	matched, err := matchContentFilter(filter, evaluator, node)
	return matched, nil, err
}

func matchAllFilters(
	scope settings.RuleReferenceScope,
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	filters []settings.PatternFilter,
	rules map[string]*settings.Rule,
) (bool, []*types.Detection, map[string]*tree.Node, error) {
	var datatypeDetections []*types.Detection

	variableNodes := make(map[string]*tree.Node)
	for name, node := range result.Variables {
		variableNodes[name] = node
	}

	for _, filter := range filters {
		matched, subDataTypeDetections, err := matchFilter(scope, result, evaluator, variableNodes, filter, rules)
		if matched == nil || !*matched || err != nil {
			return false, nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, variableNodes, nil
}

func matchEitherFilters(
	scope settings.RuleReferenceScope,
	result *languagetypes.PatternQueryResult,
	evaluator types.Evaluator,
	variableNodes map[string]*tree.Node,
	filters []settings.PatternFilter,
	rules map[string]*settings.Rule,
) (*bool, []*types.Detection, error) {
	var datatypeDetections []*types.Detection
	oneMatched := false
	oneNotMatched := false

	for _, subFilter := range filters {
		subMatch, subDatatypeDetections, err := matchFilter(scope, result, evaluator, variableNodes, subFilter, rules)
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
	scope settings.RuleReferenceScope,
	contains bool,
	rules map[string]*settings.Rule,
) (*bool, []*types.Detection, error) {
	sanitizerRuleID := ""
	if rule, ok := rules[detectorType]; ok {
		sanitizerRuleID = rule.SanitizerRuleID
	}

	if detectorType == "datatype" {
		detections, err := evaluator.Evaluate(node, "datatype", sanitizerRuleID, scope, true)

		return boolPointer(len(detections) != 0), detections, err
	}

	detections, err := evaluator.Evaluate(node, detectorType, sanitizerRuleID, scope, true)

	var datatypeDetections []*types.Detection
	ignoredVariables := getIgnoredVariables(detections)
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
			if _, ignored := ignoredVariables[name]; !ignored {
				variableNodes[name] = node
			}
		}

		datatypeDetections = append(datatypeDetections, data.Datatypes...)
	}

	return boolPointer(foundDetection), datatypeDetections, err
}

func matchContentFilter(filter settings.PatternFilter, evaluator types.Evaluator, node *tree.Node) (*bool, error) {
	content := node.Content()

	if len(filter.Values) != 0 {
		return boolPointer(slices.Contains(filter.Values, content)), nil
	}

	if filter.Regex != nil {
		return boolPointer(filter.Regex.MatchString(content)), nil
	}

	if filter.LengthLessThan != nil {
		strValue, _, err := generic.GetStringValue(node, evaluator)
		if err != nil || strValue == "" {
			return nil, err
		}

		if len(strValue) >= *filter.LengthLessThan {
			return boolPointer(false), nil
		}

		return boolPointer(true), nil
	}

	if filter.StringRegex != nil {
		value, _, err := generic.GetStringValue(node, evaluator)
		if err != nil || value == "" {
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
			if seenNode != nil && !seenNode.Equal(node) {
				ignoredVariables[name] = struct{}{}
			}

			seenNodes[name] = node
		}
	}

	return ignoredVariables
}
