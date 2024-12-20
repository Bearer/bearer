package filters

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/types"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/scanner/variableshape"
	"github.com/bearer/bearer/pkg/util/entropy"
)

type Result struct {
	matches []Match
}

func NewResult(matches ...Match) *Result {
	return &Result{matches: matches}
}

type Match struct {
	variables          variableshape.Values
	datatypeDetections []*detectortypes.Detection
	value              string
}

func NewMatch(variables variableshape.Values, valueStr string, datatypeDetections []*detectortypes.Detection) Match {
	return Match{variables: variables, value: valueStr, datatypeDetections: datatypeDetections}
}

func (result *Result) Matches() []Match {
	return result.matches
}

func (match *Match) Variables() variableshape.Values {
	return match.variables
}

func (match *Match) Value() string {
	return match.value
}

func (match *Match) DatatypeDetections() []*detectortypes.Detection {
	return match.datatypeDetections
}

type Filter interface {
	Evaluate(
		detectorContext detectortypes.Context,
		patternVariables variableshape.Values,
	) (*Result, error)
}

type Not struct {
	Child Filter
}

func (filter *Not) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	childResult, err := filter.Child.Evaluate(detectorContext, patternVariables)
	if err != nil {
		return nil, err
	}

	if childResult == nil {
		log.Trace().Msg("filters.Not: nil")
		return nil, nil
	}

	result := len(childResult.Matches()) == 0

	if log.Trace().Enabled() {
		log.Trace().Msgf("filters.Not: %t", result)
	}

	return boolResult(patternVariables, result, ""), nil
}

type Either struct {
	Children []Filter
}

func (filter *Either) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	var matches []Match

	unknownResult := true
	for _, child := range filter.Children {
		subResult, err := child.Evaluate(detectorContext, patternVariables)
		if err != nil {
			return nil, err
		}

		if subResult == nil {
			continue
		}

		unknownResult = false
		matches = append(matches, subResult.matches...)
	}

	if unknownResult {
		return nil, nil
	}

	return NewResult(matches...), nil
}

type All struct {
	Children []Filter
}

func (filter *All) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	var matches []Match

	if len(filter.Children) == 0 {
		log.Trace().Msg("filters.All: true (no children)")
		return boolResult(patternVariables, true, ""), nil
	}

	for i, child := range filter.Children {
		subResult, err := child.Evaluate(detectorContext, patternVariables)
		if err != nil {
			return nil, err
		}

		if subResult == nil {
			log.Trace().Msg("filters.All: nil")
			return nil, nil
		}

		if i == 0 {
			matches = subResult.matches
			continue
		}

		matches = filter.joinMatches(matches, subResult.matches)

		if len(matches) == 0 {
			log.Trace().Msg("filters.All: no matches")
			return NewResult(), nil
		}
	}

	log.Trace().Msg("filters.All: matches")
	return NewResult(matches...), nil
}

func (filter *All) joinMatches(matches, childMatches []Match) []Match {
	var result []Match

	for _, match := range matches {
		for _, childMatch := range childMatches {
			if variables, variablesMatch := match.variables.Merge(childMatch.variables); variablesMatch {
				value := match.Value()
				value += childMatch.Value()

				result = append(result, NewMatch(
					variables,
					value,
					// FIXME: this seems like it will create unnecessary duplicates
					append(match.datatypeDetections, childMatch.datatypeDetections...),
				))
			}
		}
	}

	return result
}

type FilenameRegex struct {
	Regex *regexp.Regexp
}

func (filter *FilenameRegex) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	return boolResult(patternVariables, filter.Regex.MatchString(detectorContext.Filename()), ""), nil
}

type ImportedVariable struct {
	ParentVariable,
	ChildVariable *variableshape.Variable
}

type Rule struct {
	Variable          *variableshape.Variable
	Rule              *ruleset.Rule
	TraversalStrategy traversalstrategy.Strategy
	IsDatatypeRule    bool
	Filter            Filter
	ImportedVariables []ImportedVariable
}

func (filter *Rule) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	if node == nil {
		return nil, fmt.Errorf("couldn't find node for var %s", filter.Variable.Name())
	}
	detections, err := detectorContext.Scan(node, filter.Rule, filter.TraversalStrategy)
	if err != nil {
		return nil, err
	}

	if len(detections) == 0 {
		return NewResult(), nil
	}

	if filter.IsDatatypeRule {
		log.Trace().Msg("filters.Rule: match (datatype)")
		return NewResult(NewMatch(patternVariables, "", detections)), nil
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("filters.Rule: %d detections", len(detections))
	}

	var matches []Match
	hasPatternVariableMatch := false

	var datatypeDetections []*detectortypes.Detection

	for _, detection := range detections {
		data, ok := detection.Data.(types.Data)
		if !ok { // Built-in detector
			log.Trace().Msg("filters.Rule: match (built-in)")

			hasPatternVariableMatch = true
			continue
		}

		subResult, err := filter.Filter.Evaluate(detectorContext, data.Variables)
		if err != nil {
			return nil, err
		}

		if subResult == nil {
			log.Trace().Msg("filters.Rule: no match (filter result unknown)")
			continue
		}

		if len(subResult.matches) == 0 {
			log.Trace().Msg("filters.Rule: no match")
			continue
		}

		if len(filter.ImportedVariables) == 0 {
			log.Trace().Msg("filters.Rule: match (no imported vars)")

			hasPatternVariableMatch = true
			datatypeDetections = append(datatypeDetections, data.Datatypes...)

			for _, detectionMatch := range subResult.matches {
				datatypeDetections = append(datatypeDetections, detectionMatch.datatypeDetections...)
			}

			continue
		}

		matched := false
		for _, detectionMatch := range subResult.matches {
			if variables, variablesMatch := filter.importVariables(patternVariables, detectionMatch.variables); variablesMatch {
				matched = true
				matches = append(matches, NewMatch(variables, "", detectionMatch.datatypeDetections))
			}
		}

		if matched {
			log.Trace().Msg("filters.Rule: match")

			if len(data.Datatypes) != 0 {
				hasPatternVariableMatch = true
				datatypeDetections = append(datatypeDetections, data.Datatypes...)
			}
		} else {
			log.Trace().Msg("filters.Rule: no match (variable mismatch)")
		}
	}

	if hasPatternVariableMatch {
		matches = append(matches, NewMatch(patternVariables, "", datatypeDetections))
	}

	return NewResult(matches...), nil
}

func (filter *Rule) importVariables(parentVariables, childVariables variableshape.Values) (variableshape.Values, bool) {
	if len(filter.ImportedVariables) == 0 {
		return parentVariables, true
	}

	variables := parentVariables.Clone()

	for _, importedVariable := range filter.ImportedVariables {
		parentNode := parentVariables.Node(importedVariable.ParentVariable)
		childNode := childVariables.Node(importedVariable.ChildVariable)

		if childNode == nil {
			continue
		}

		if parentNode != nil && parentNode != childNode {
			return nil, false
		}

		variables.Set(importedVariable.ParentVariable, childNode)
	}

	return variables, true
}

type Values struct {
	Variable *variableshape.Variable
	Values   []string
}

func (filter *Values) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	return boolResult(patternVariables, slices.Contains(filter.Values, node.Content()), ""), nil
}

type Regex struct {
	Variable *variableshape.Variable
	Regex    *regexp.Regexp
}

func (filter *Regex) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	result := filter.Regex.MatchString(node.Content())

	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"filters.Regex: %t for pattern %s at %s, content=%s",
			result,
			filter.Regex.String(),
			node.Debug(),
			node.Content(),
		)
	}

	return boolResult(patternVariables, result, ""), nil
}

type StringLengthLessThan struct {
	Variable *variableshape.Variable
	Value    int
}

func (filter *StringLengthLessThan) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isString, err := lookupString(detectorContext, node)
	if err != nil || !isString {
		return nil, err
	}

	return boolResult(patternVariables, len(value) < filter.Value, ""), nil
}

type StringRegex struct {
	Variable *variableshape.Variable
	Regex    *regexp.Regexp
}

func (filter *StringRegex) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isString, err := lookupString(detectorContext, node)
	if err != nil {
		return nil, err
	}

	if !isString {
		if log.Trace().Enabled() {
			log.Trace().Msgf("filters.StringRegex: nil for pattern %s at %s", filter.Regex.String(), node.Debug())
		}

		return nil, nil
	}

	result := filter.Regex.MatchString(value)
	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"filters.StringRegex: %t for pattern %s at %s, content=%s",
			result,
			filter.Regex.String(),
			node.Debug(),
			value,
		)
	}

	return boolResult(patternVariables, result, value), nil
}

type EntropyGreaterThan struct {
	Variable *variableshape.Variable
	Value    float64
}

func (filter *EntropyGreaterThan) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isString, err := lookupString(detectorContext, node)
	if err != nil {
		return nil, err
	}

	if !isString {
		if log.Trace().Enabled() {
			log.Trace().Msgf("filters.EntropyGreaterThan: nil for min %f at %s", filter.Value, node.Debug())
		}

		return nil, nil
	}

	entropy := entropy.Shannon(value)
	result := entropy > filter.Value
	if log.Trace().Enabled() {
		log.Trace().Msgf(
			"filters.EntropyGreaterThan: %t for entropy %f with min %f at %s, content=%s",
			result,
			entropy,
			filter.Value,
			node.Debug(),
			value,
		)
	}

	return boolResult(patternVariables, result, ""), nil
}

type IntegerLessThan struct {
	Variable *variableshape.Variable
	Value    int
}

func (filter *IntegerLessThan) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isInteger, err := parseInteger(node)
	if err != nil || !isInteger {
		return nil, err
	}

	return boolResult(patternVariables, value < filter.Value, ""), nil
}

type IntegerLessThanOrEqual struct {
	Variable *variableshape.Variable
	Value    int
}

func (filter *IntegerLessThanOrEqual) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isInteger, err := parseInteger(node)
	if err != nil || !isInteger {
		return nil, err
	}

	return boolResult(patternVariables, value <= filter.Value, ""), nil
}

type IntegerGreaterThan struct {
	Variable *variableshape.Variable
	Value    int
}

func (filter *IntegerGreaterThan) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isInteger, err := parseInteger(node)
	if err != nil || !isInteger {
		return nil, err
	}

	return boolResult(patternVariables, value > filter.Value, ""), nil
}

type IntegerGreaterThanOrEqual struct {
	Variable *variableshape.Variable
	Value    int
}

func (filter *IntegerGreaterThanOrEqual) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	node := patternVariables.Node(filter.Variable)
	value, isInteger, err := parseInteger(node)
	if err != nil || !isInteger {
		return nil, err
	}

	return boolResult(patternVariables, value >= filter.Value, ""), nil
}

type Unknown struct{}

func (filter *Unknown) Evaluate(
	detectorContext detectortypes.Context,
	patternVariables variableshape.Values,
) (*Result, error) {
	return nil, nil
}

func lookupString(
	detectorContext detectortypes.Context,
	node *tree.Node,
) (string, bool, error) {
	value, isLiteral, err := common.GetStringValue(node, detectorContext)
	if err != nil || (value == "" && !isLiteral) {
		return "", false, err
	}

	return value, true, nil
}

func parseInteger(node *tree.Node) (int, bool, error) {
	value, err := strconv.Atoi(node.Content())
	if err != nil {
		return 0, false, nil
	}

	return value, true, nil
}

func boolResult(patternVariables variableshape.Values, value bool, valueStr string) *Result {
	return NewResult(boolMatches(patternVariables, value, valueStr)...)
}

func boolMatches(patternVariables variableshape.Values, value bool, valueStr string) []Match {
	if value {
		return []Match{NewMatch(patternVariables, valueStr, nil)}
	} else {
		return nil
	}
}
