package variableshape

import (
	"fmt"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	patternquerybuilder "github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

type Values []*tree.Node

func (values Values) Clone() Values {
	if len(values) == 0 {
		return nil
	}

	result := make(Values, len(values))
	copy(result, values)
	return result
}

func (values Values) Merge(other Values) (Values, bool) {
	if len(values) == 0 {
		return nil, true
	}

	if &values[0] == &other[0] {
		return values, true
	}

	result := make(Values, len(values))
	for i, node := range values {
		otherNode := other[i]

		if node == nil || node == otherNode {
			result[i] = otherNode
			continue
		}

		if otherNode == nil {
			result[i] = node
			continue
		}

		return nil, false
	}

	return result, true
}

func (values Values) Node(variable *Variable) *tree.Node {
	return values[variable.id]
}

func (values Values) Set(variable *Variable, node *tree.Node) {
	values[variable.id] = node
}

type Variable struct {
	id   int
	name string
}

func (variable *Variable) Name() string {
	return variable.name
}

type Builder struct {
	variables []Variable
	nameToID  map[string]int
}

func NewBuilder() *Builder {
	return &Builder{
		nameToID: make(map[string]int),
	}
}

func (builder *Builder) Add(name string) *Builder {
	_, exists := builder.nameToID[name]
	if exists {
		return builder
	}

	id := len(builder.variables)

	builder.variables = append(builder.variables, Variable{
		id:   id,
		name: name,
	})

	builder.nameToID[name] = id

	return builder
}

func (builder *Builder) Build() Shape {
	nameToVariable := make(map[string]*Variable)
	for i := range builder.variables {
		variable := &builder.variables[i]
		nameToVariable[variable.name] = variable
	}

	return Shape{
		variables:      builder.variables,
		nameToVariable: nameToVariable,
	}
}

type Shape struct {
	variables      []Variable
	nameToVariable map[string]*Variable
}

func (shape *Shape) Variable(name string) (*Variable, error) {
	variable, exists := shape.nameToVariable[name]
	if !exists {
		return nil, fmt.Errorf("unknown variable '%s'", name)
	}

	return variable, nil
}

func (shape *Shape) NewValues() Values {
	if len(shape.variables) == 0 {
		return nil
	}

	return make(Values, len(shape.variables))
}

type Set struct {
	shapes []Shape
}

func NewSet(language language.Language, ruleSet *ruleset.Set) (*Set, error) {
	set := &Set{
		shapes: make([]Shape, len(ruleSet.Rules())),
	}

	for _, rule := range ruleSet.Rules() {
		if err := set.add(language, rule); err != nil {
			return nil, err
		}
	}

	return set, nil
}

// FIXME: don't do this!
func (set *Set) add(language language.Language, rule *ruleset.Rule) error {
	builder := NewBuilder()

	for _, pattern := range rule.Patterns() {
		if err := addVariablesFromPattern(language, builder, pattern.Pattern); err != nil {
			return err
		}

		addVariablesFromFilters(builder, pattern.Filters)
	}

	set.shapes[rule.Index()] = builder.Build()
	return nil
}

func (set *Set) Shape(rule *ruleset.Rule) *Shape {
	return &set.shapes[rule.Index()]
}

func addVariablesFromPattern(language language.Language, builder *Builder, pattern string) error {
	result, err := patternquerybuilder.Build(language, pattern, "")
	if err != nil {
		return err
	}

	if result.RootVariable != nil {
		builder.Add(result.RootVariable.Name)
		return nil
	}

	for _, name := range result.VariableNames {
		builder.Add(name)
	}

	return nil
}

func addVariablesFromFilters(builder *Builder, filters []settings.PatternFilter) {
	for _, filter := range filters {
		addVariablesFromFilter(builder, filter)
	}
}

func addVariablesFromFilter(builder *Builder, filter settings.PatternFilter) {
	for _, importedVariable := range filter.Imports {
		builder.Add(importedVariable.As)
	}

	addVariablesFromFilters(builder, filter.Either)
	addVariablesFromFilters(builder, filter.Filters)
}
