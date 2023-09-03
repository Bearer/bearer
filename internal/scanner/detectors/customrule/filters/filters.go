package filters

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule/types"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/pointers"
)

type Filter interface {
	Evaluate(
		detectorContext detectortypes.Context,
		variables,
		joinedVariables map[string]*tree.Node,
	) (*bool, []*detectortypes.Detection, error)
}

type Not struct {
	Child Filter
}

func (filter *Not) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	match, _, err := filter.Child.Evaluate(detectorContext, variables, joinedVariables)
	if match == nil {
		return nil, nil, err
	}

	return pointers.Bool(!*match), nil, err
}

type Either struct {
	Children []Filter
}

func (filter *Either) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	var datatypeDetections []*detectortypes.Detection
	oneMatched := false
	oneNotMatched := false

	for _, child := range filter.Children {
		subMatch, subDatatypeDetections, err := child.Evaluate(detectorContext, variables, joinedVariables)
		if err != nil {
			return nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDatatypeDetections...)
		oneMatched = oneMatched || (subMatch != nil && *subMatch)
		oneNotMatched = oneNotMatched || (subMatch != nil && !*subMatch)
	}

	if oneMatched {
		return pointers.Bool(true), datatypeDetections, nil
	}

	if oneNotMatched {
		return pointers.Bool(false), nil, nil
	}

	return nil, nil, nil
}

type FilenameRegex struct {
	Regex *regexp.Regexp
}

func (filter *FilenameRegex) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	return pointers.Bool(filter.Regex.MatchString(detectorContext.Filename())), nil, nil
}

type Rule struct {
	VariableName      string
	Rule              *ruleset.Rule
	TraversalStrategy *traversalstrategy.Strategy
	IsDatatypeRule    bool
	Filters           []Filter
}

func (filter *Rule) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	node, err := lookupVariable(variables, filter.VariableName)
	if err != nil {
		return nil, nil, err
	}

	detections, err := detectorContext.Scan(node, filter.Rule, filter.TraversalStrategy)
	if err != nil {
		return nil, nil, err
	}

	if filter.IsDatatypeRule {
		return pointers.Bool(len(detections) != 0), detections, nil
	}

	var datatypeDetections []*detectortypes.Detection
	ignoredVariables := getIgnoredVariables(detections)
	foundDetection := false

	for _, detection := range detections {
		data, ok := detection.Data.(types.Data)
		if !ok { // Built-in detector
			foundDetection = true
			log.Trace().Msg("detection match (built-in)")
			continue
		}

		filtersMatch, _, _, err := Match(detectorContext, data.VariableNodes, filter.Filters)
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

	return pointers.Bool(foundDetection), datatypeDetections, nil
}

func getIgnoredVariables(detections []*detectortypes.Detection) map[string]struct{} {
	ignoredVariables := make(map[string]struct{})
	seenNodes := make(map[string]*tree.Node)

	for _, detection := range detections {
		data, ok := detection.Data.(types.Data)
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

type Values struct {
	VariableName string
	Values       []string
}

func (filter *Values) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	node, err := lookupVariable(variables, filter.VariableName)
	if err != nil {
		return nil, nil, err
	}

	return pointers.Bool(slices.Contains(filter.Values, node.Content())), nil, nil
}

type Regex struct {
	VariableName string
	Regex        *regexp.Regexp
}

func (filter *Regex) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	node, err := lookupVariable(variables, filter.VariableName)
	if err != nil {
		return nil, nil, err
	}

	return pointers.Bool(filter.Regex.MatchString(node.Content())), nil, nil
}

type StringLengthLessThan struct {
	VariableName string
	Value        int
}

func (filter *StringLengthLessThan) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isString, err := lookupString(detectorContext, variables, filter.VariableName)
	if err != nil || !isString {
		return nil, nil, err
	}

	if len(value) >= filter.Value {
		return pointers.Bool(false), nil, nil
	}

	return pointers.Bool(true), nil, nil
}

type StringRegex struct {
	VariableName string
	Regex        *regexp.Regexp
}

func (filter *StringRegex) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isString, err := lookupString(detectorContext, variables, filter.VariableName)
	if err != nil || !isString {
		return nil, nil, err
	}

	return pointers.Bool(filter.Regex.MatchString(value)), nil, nil
}

type IntegerLessThan struct {
	VariableName string
	Value        int
}

func (filter *IntegerLessThan) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isInteger, err := lookupInteger(variables, filter.VariableName)
	if err != nil || !isInteger {
		return nil, nil, err
	}

	return pointers.Bool(value < filter.Value), nil, nil
}

type IntegerLessThanOrEqual struct {
	VariableName string
	Value        int
}

func (filter *IntegerLessThanOrEqual) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isInteger, err := lookupInteger(variables, filter.VariableName)
	if err != nil || !isInteger {
		return nil, nil, err
	}

	return pointers.Bool(value <= filter.Value), nil, nil
}

type IntegerGreaterThan struct {
	VariableName string
	Value        int
}

func (filter *IntegerGreaterThan) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isInteger, err := lookupInteger(variables, filter.VariableName)
	if err != nil || !isInteger {
		return nil, nil, err
	}

	return pointers.Bool(value > filter.Value), nil, nil
}

type IntegerGreaterThanOrEqual struct {
	VariableName string
	Value        int
}

func (filter *IntegerGreaterThanOrEqual) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	value, isInteger, err := lookupInteger(variables, filter.VariableName)
	if err != nil || !isInteger {
		return nil, nil, err
	}

	return pointers.Bool(value >= filter.Value), nil, nil
}

type Unknown struct{}

func (filter *Unknown) Evaluate(
	detectorContext detectortypes.Context,
	variables,
	joinedVariables map[string]*tree.Node,
) (*bool, []*detectortypes.Detection, error) {
	return nil, nil, nil
}

func Match(
	detectorContext detectortypes.Context,
	variables map[string]*tree.Node,
	filters []Filter,
) (bool, []*detectortypes.Detection, map[string]*tree.Node, error) {
	var datatypeDetections []*detectortypes.Detection

	joinedVariables := make(map[string]*tree.Node)
	for name, node := range variables {
		joinedVariables[name] = node
	}

	for _, filter := range filters {
		matched, subDataTypeDetections, err := filter.Evaluate(detectorContext, variables, joinedVariables)
		if matched == nil || !*matched || err != nil {
			return false, nil, nil, err
		}

		datatypeDetections = append(datatypeDetections, subDataTypeDetections...)
	}

	return true, datatypeDetections, joinedVariables, nil
}

func lookupVariable(variables map[string]*tree.Node, name string) (*tree.Node, error) {
	node, exists := variables[name]
	if !exists {
		return nil, fmt.Errorf("invalid variable '%s'", name)
	}

	return node, nil
}

func lookupString(
	detectorContext detectortypes.Context,
	variables map[string]*tree.Node,
	variableName string,
) (string, bool, error) {
	node, err := lookupVariable(variables, variableName)
	if err != nil {
		return "", false, err
	}

	value, isLiteral, err := common.GetStringValue(node, detectorContext)
	if err != nil || (value == "" && !isLiteral) {
		return "", false, err
	}

	return value, true, nil
}

func lookupInteger(variables map[string]*tree.Node, variableName string) (int, bool, error) {
	node, err := lookupVariable(variables, variableName)
	if err != nil {
		return 0, false, err
	}

	value, err := strconv.Atoi(node.Content())
	if err != nil {
		return 0, false, nil
	}

	return value, true, nil
}
