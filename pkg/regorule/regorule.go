package regorule

import (
	"fmt"
	"strings"

	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/open-policy-agent/opa/rego"
)

type rule struct {
	variables []string
	rego      string
}

type compiler struct {
	regoCode    strings.Builder
	sourceRules map[string]*rules.Rule
	rules       map[string]*rule
}

func Compile(rules map[string]*rules.Rule) (rego.Rego, error) {
	compiler := &compiler{
		sourceRules: rules,
		rules:       make(map[string]*rule),
	}
	// rule contains [..., vars] if {}

	for _, sourceRule := range rules {
		compiler.compile(sourceRule)
		regoSource.WriteString("\n\n")
	}

	return nil, nil
}

func (compiler *compiler) compile(sourceRule *rules.Rule) (*rule, error) {
	if rule, exists := compiler.rules[sourceRule.Id]; exists {
		return rule, nil
	}

	variables, err := compiler.enumerateVariables(sourceRule)
	if err != nil {
		return nil, err
	}

	compiler.regoCode.WriteString()
	compiler.regoCode.WriteString("\n\n")

	rule := &rule{
		variables: variables,
	}

	compiler.rules[sourceRule.Id] = rule

	return rule, err
}

func (compiler *compiler) enumerateVariables(sourceRule *rules.Rule) ([]string, error) {
	variables := set.New[string]()

	for _, pattern := range sourceRule.Patterns {
		variables.AddAll(pattern.Query.Variables())

		filterVariables, err := compiler.enumerateFilterVariables(pattern.Filters)
		if err != nil {
			return nil, err
		}
		variables.AddAll(filterVariables)
	}

	return variables.Items(), nil
}

func (compiler *compiler) enumerateFilterVariables(filters []rules.PatternFilter) ([]string, error) {
	variables := set.New[string]()

	for _, filter := range filters {
		if filter.Detection != "" {
			sourceRule, exists := compiler.sourceRules[filter.Detection]
			if !exists {
				return nil, fmt.Errorf("referenced rule %s does not exist", filter.Detection)
			}

			rule := compiler.compile(sourceRule)
			variables.AddAll(rule.variables)

			continue
		}

	}

	return variables.Items(), nil
}
