package regorule

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/open-policy-agent/opa/rego"
	"golang.org/x/exp/slices"
)

//go:embed prelude.rego
var prelude []byte

type rule struct {
	variables []string
}

type compiler struct {
	idGenerator nodeid.Generator
	regoRules   strings.Builder
	regoQuery   strings.Builder
	sourceRules map[string]*rules.Rule
	rules       map[string]*rule
}

func Compile(rules map[string]*rules.Rule) (*rego.PreparedEvalQuery, error) {
	compiler := &compiler{
		idGenerator: &nodeid.IntGenerator{},
		sourceRules: rules,
		rules:       make(map[string]*rule),
	}

	compiler.regoRules.Write(prelude)

	for _, sourceRule := range rules {
		compiler.compile(sourceRule)
	}

	fmt.Printf("rego code:\n\n%s\n", compiler.regoRules.String())
	fmt.Printf("rego query:\n\n%s\n", compiler.regoQuery.String())

	rego := rego.New(
		rego.Module("rules.rego", compiler.regoRules.String()),
		rego.Query(compiler.regoQuery.String()),
		// rego.Function3(
		// 	&rego.Function{
		// 		Name:             "bearer.builtin.patternMatch",
		// 		Decl:             types.NewFunction([]types.Type{types.N, types.S, types.N}, types.B),
		// 		Memoize:          true,
		// 		Nondeterministic: false,
		// 	}, func(bctx rego.BuiltinContext, nodeTerm, ruleIDTerm, patternIndexTerm *ast.Term) (*ast.Term, error) {
		// 		return ast.BooleanTerm(false), nil
		// 	},
		// ),
	)

	query, err := rego.PrepareForEval(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error preparing rego: %w", err)
	}

	return &query, nil
}

func (compiler *compiler) compile(sourceRule *rules.Rule) (*rule, error) {
	if rule, exists := compiler.rules[sourceRule.Id]; exists {
		return rule, nil
	}

	variables, err := compiler.enumerateVariables(sourceRule)
	if err != nil {
		return nil, err
	}

	var buffer strings.Builder

	for i, pattern := range sourceRule.Patterns {
		patternID := fmt.Sprintf("%s[%d]", sourceRule.Id, i)
		buffer.WriteString(sourceRule.Id)
		buffer.WriteString(" contains [cursor")

		patternVariables := pattern.Query.Variables()
		filterVariables := set.New[string]()

		for _, filter := range pattern.Filters {
			subVariables, err := compiler.enumerateFilterVariables(filter)
			if err != nil {
				return nil, err
			}

			filterVariables.AddAll(subVariables)
		}

		variableUsed := func(variable string) bool {
			return slices.Contains(patternVariables, variable) || filterVariables.Has(variable)
		}

		for _, variable := range variables {
			if variableUsed(variable) {
				buffer.WriteString(", var_")
				buffer.WriteString(variable)
			} else {
				buffer.WriteString(", null")
			}
		}

		buffer.WriteString("] if {\n")
		// buffer.WriteString("  node := input.nodes[_]\n")
		// for _, variable := range variables {
		// 	if !variableUsed(variable) {
		// 		continue
		// 	}

		// 	buffer.WriteString("some var_")
		// 	buffer.WriteString(variable)
		// 	buffer.WriteString("\n")
		// }
		// buffer.WriteString("\n")

		variablesIdentifier := "_"
		if len(patternVariables) != 0 {
			buffer.WriteString("  some variables\n")
			variablesIdentifier = "variables"
		}

		buffer.WriteString(`  pattern_match[[cursor, "`)
		buffer.WriteString(patternID)
		buffer.WriteString(`", `)
		buffer.WriteString(variablesIdentifier)
		buffer.WriteString("]]\n")

		for _, variable := range patternVariables {
			buffer.WriteString("  var_")
			buffer.WriteString(variable)
			buffer.WriteString("  := variables.")
			buffer.WriteString(variable)
			buffer.WriteString("\n")
		}
		buffer.WriteString("\n")

		for _, filter := range pattern.Filters {
			buffer.WriteString("  ")
			compiler.compileFilter(&buffer, filter, variableUsed)
			buffer.WriteString("\n")
		}

		buffer.WriteString("}\n\n")
	}

	compiler.regoRules.WriteString(buffer.String())

	compiler.regoQuery.WriteString(sourceRule.Id)
	compiler.regoQuery.WriteString(" := [ node | data.bearer.rules.")
	compiler.regoQuery.WriteString(sourceRule.Id)

	compiler.regoQuery.WriteString("[[node")

	for range variables {
		compiler.regoQuery.WriteString(", _")
	}

	compiler.regoQuery.WriteString("]] ]\n")

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

		for _, filter := range pattern.Filters {
			filterVariables, err := compiler.enumerateFilterVariables(filter)
			if err != nil {
				return nil, err
			}
			variables.AddAll(filterVariables)
		}
	}

	return variables.Items(), nil
}

func (compiler *compiler) enumerateFilterVariables(filter rules.PatternFilter) ([]string, error) {
	if filter.Detection != "" {
		rule, err := compiler.ruleById(filter.Detection)
		if err != nil {
			return nil, err
		}

		return rule.variables, nil
	}

	if len(filter.Either) != 0 {
		common := make(map[string]int)

		for _, subFilter := range filter.Either {
			subVariables, err := compiler.enumerateFilterVariables(subFilter)
			if err != nil {
				return nil, err
			}

			for _, variable := range subVariables {
				common[variable] += 1
			}
		}

		var variables []string
		for variable, count := range common {
			if count == len(filter.Either) {
				variables = append(variables, variable)
			}
		}

		return variables, nil
	}

	// if filter.Not != nil {
	// 	variables, err := compiler.enumerateFilterVariables(*filter.Not)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return variables, nil
	// }

	return nil, nil
}

func (compiler *compiler) compileFilter(buffer *strings.Builder, filter rules.PatternFilter, variableUsed func(variable string) bool) error {
	if filter.Not != nil {
		subFilter := *filter.Not
		if subFilter.Detection == "" {
			return nil
		}

		rule, err := compiler.ruleById(subFilter.Detection)
		if err != nil {
			return err
		}

		var usedVariables []string
		for _, variable := range rule.variables {
			if variableUsed(variable) {
				usedVariables = append(usedVariables, variable)
			}
		}

		var subBuffer strings.Builder
		if err := compiler.compileFilter(&subBuffer, *filter.Not, variableUsed); err != nil {
			return err
		}

		notRuleId := "not_filter_" + compiler.idGenerator.GenerateId()
		cursorDefined := slices.Contains(usedVariables, subFilter.Variable)
		compiler.regoRules.WriteString(notRuleId)
		compiler.regoRules.WriteString(" contains [")
		if !cursorDefined {
			compiler.regoRules.WriteString("var_")
			compiler.regoRules.WriteString(subFilter.Variable)
		}
		for i, variable := range usedVariables {
			if i != 0 || !cursorDefined {
				compiler.regoRules.WriteString(", ")
			}

			compiler.regoRules.WriteString("var_")
			compiler.regoRules.WriteString(variable)
		}
		compiler.regoRules.WriteString("] if { ")
		compiler.regoRules.WriteString(subBuffer.String())
		compiler.regoRules.WriteString(" }\n")

		buffer.WriteString("not ")
		buffer.WriteString(notRuleId)
	}

	if filter.Detection != "" {
		buffer.WriteString(filter.Detection)
		buffer.WriteString("[[var_")
		buffer.WriteString(filter.Variable)

		rule, err := compiler.ruleById(filter.Detection)
		if err != nil {
			return err
		}

		for _, variable := range rule.variables {
			if variableUsed(variable) {
				buffer.WriteString(", var_")
				buffer.WriteString(variable)
			} else {
				buffer.WriteString(", _")
			}
		}

		buffer.WriteString("]]")
	}

	return nil
}

func (compiler *compiler) ruleById(id string) (*rule, error) {
	sourceRule, exists := compiler.sourceRules[id]
	if !exists {
		return nil, fmt.Errorf("referenced rule %s does not exist", id)
	}

	return compiler.compile(sourceRule)
}
