package customrule

import (
	"fmt"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule/filters"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/scanner/variableshape"
)

func translateFiltersTop(
	ruleSet *ruleset.Set,
	variableShapeSet *variableshape.Set,
	variableShape *variableshape.Shape,
	sourceFilters []settings.PatternFilter,
) (filters.Filter, error) {
	children, err := translateFilters(ruleSet, variableShapeSet, variableShape, sourceFilters)
	if err != nil {
		return nil, err
	}

	return &filters.All{Children: children}, nil
}

func translateFilters(
	ruleSet *ruleset.Set,
	variableShapeSet *variableshape.Set,
	variableShape *variableshape.Shape,
	sourceFilters []settings.PatternFilter,
) ([]filters.Filter, error) {
	filters := make([]filters.Filter, len(sourceFilters))

	sortFilters(sourceFilters)
	for i, sourceFilter := range sourceFilters {
		filter, err := translateFilter(ruleSet, variableShapeSet, variableShape, &sourceFilter)
		if err != nil {
			return nil, err
		}

		filters[i] = filter
	}

	return filters, nil
}

func translateFilter(
	ruleSet *ruleset.Set,
	variableShapeSet *variableshape.Set,
	variableShape *variableshape.Shape,
	sourceFilter *settings.PatternFilter,
) (filters.Filter, error) {
	if sourceFilter.Not != nil {
		child, err := translateFilter(ruleSet, variableShapeSet, variableShape, sourceFilter.Not)
		if err != nil {
			return nil, err
		}

		return &filters.Not{Child: child}, nil
	}

	if len(sourceFilter.Either) != 0 {
		children, err := translateFilters(ruleSet, variableShapeSet, variableShape, sourceFilter.Either)
		if err != nil {
			return nil, err
		}

		return &filters.Either{Children: children}, nil
	}

	if sourceFilter.FilenameRegex != nil {
		return &filters.FilenameRegex{Regex: sourceFilter.FilenameRegex.Regexp}, nil
	}

	variable, err := variableShape.Variable(sourceFilter.Variable)
	if err != nil {
		return nil, err
	}

	if sourceFilter.Detection != "" {
		rule, err := ruleSet.RuleByID(sourceFilter.Detection)
		if err != nil {
			return nil, err
		}

		ruleFilter, err := translateFiltersTop(
			ruleSet,
			variableShapeSet,
			variableShapeSet.Shape(rule),
			sourceFilter.Filters,
		)
		if err != nil {
			return nil, err
		}

		traversalStrategy, err := traversalstrategy.Get(sourceFilter.Scope)
		if err != nil {
			return nil, err
		}

		childVariableShape := variableShapeSet.Shape(rule)

		importedVariables := make([]filters.ImportedVariable, len(sourceFilter.Imports))
		for i, importedVariable := range sourceFilter.Imports {
			parentVariable, err := variableShape.Variable(importedVariable.As)
			if err != nil {
				return nil, err
			}

			childVariable, err := childVariableShape.Variable(importedVariable.Variable)
			if err != nil {
				return nil, err
			}

			importedVariables[i] = filters.ImportedVariable{
				ParentVariable: parentVariable,
				ChildVariable:  childVariable,
			}
		}

		return &filters.Rule{
			Variable:          variable,
			Rule:              rule,
			TraversalStrategy: traversalStrategy,
			IsDatatypeRule:    sourceFilter.Detection == "datatype",
			Filter:            ruleFilter,
			ImportedVariables: importedVariables,
		}, nil
	}

	if len(sourceFilter.Values) != 0 {
		return &filters.Values{
			Variable: variable,
			Values:   sourceFilter.Values,
		}, nil
	}

	if sourceFilter.Regex != nil {
		return &filters.Regex{
			Variable: variable,
			Regex:    sourceFilter.Regex.Regexp,
		}, nil
	}

	if sourceFilter.LengthLessThan != nil {
		return &filters.StringLengthLessThan{
			Variable: variable,
			Value:    *sourceFilter.LengthLessThan,
		}, nil
	}

	if sourceFilter.StringRegex != nil {
		return &filters.StringRegex{
			Variable: variable,
			Regex:    sourceFilter.StringRegex.Regexp,
		}, nil
	}

	if sourceFilter.LessThan != nil {
		return &filters.IntegerLessThan{
			Variable: variable,
			Value:    *sourceFilter.LessThan,
		}, nil
	}

	if sourceFilter.LessThanOrEqual != nil {
		return &filters.IntegerLessThanOrEqual{
			Variable: variable,
			Value:    *sourceFilter.LessThanOrEqual,
		}, nil
	}

	if sourceFilter.GreaterThan != nil {
		return &filters.IntegerGreaterThan{
			Variable: variable,
			Value:    *sourceFilter.GreaterThan,
		}, nil
	}

	if sourceFilter.GreaterThanOrEqual != nil {
		return &filters.IntegerGreaterThanOrEqual{
			Variable: variable,
			Value:    *sourceFilter.GreaterThanOrEqual,
		}, nil
	}

	log.Debug().Msgf("unknown filter type: %#v", sourceFilter)
	return &filters.Unknown{}, nil
}

func sortFilters(filters []settings.PatternFilter) {
	slices.SortFunc(filters, func(a, b settings.PatternFilter) int {
		return scoreFilter(a) - scoreFilter(b)
	})

	for i := range filters {
		sortFilter(&filters[i])
	}
}

func sortFilter(filter *settings.PatternFilter) {
	switch {
	case len(filter.Either) != 0:
		sortFilters(filter.Either)
	case filter.Not != nil:
		sortFilter(filter.Not)
	}
}

func scoreFilter(filter settings.PatternFilter) int {
	if filter.Regex != nil ||
		len(filter.Values) != 0 ||
		filter.LengthLessThan != nil ||
		filter.LessThan != nil ||
		filter.LessThanOrEqual != nil ||
		filter.GreaterThan != nil ||
		filter.GreaterThanOrEqual != nil ||
		filter.FilenameRegex != nil {
		return 1
	}

	if filter.Detection == "datatype" {
		return 7
	}

	if filter.StringRegex != nil ||
		filter.Detection != "" && filter.Scope == settings.CURSOR_STRICT_SCOPE {
		return 2
	}

	if filter.Detection != "" && filter.Scope == settings.CURSOR_SCOPE {
		return 3
	}

	if filter.Detection != "" && filter.Scope == settings.NESTED_STRICT_SCOPE {
		return 4
	}

	if filter.Detection != "" && filter.Scope == settings.RESULT_SCOPE {
		return 5
	}

	if filter.Detection != "" && filter.Scope == settings.NESTED_SCOPE {
		return 6
	}

	if filter.Not != nil {
		return scoreFilter(*filter.Not)
	}

	if len(filter.Either) != 0 {
		max := 0

		for _, subFilter := range filter.Either {
			if subScore := scoreFilter(subFilter); subScore > max {
				max = subScore
			}
		}

		return max
	}

	panic(fmt.Sprintf("unknown filter %#v", filter))
}
