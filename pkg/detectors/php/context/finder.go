package context

import (
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/context"
	"github.com/bearer/bearer/pkg/report/variables"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

func FindContext(tree *parser.Tree) *context.Resolver {
	finder := context.NewFinder(&context.FinderRequest{
		ContextKeywords: []string{
			"function_definition",
			"class_declaration",
			"method_declaration",
			"anonymous_function_creation_expression",
			"object_creation_expression",
		},
		Tree:             tree,
		VariableResolver: variableResolver,
	})
	finder.Find()
	return finder.ToResolver()
}

func variableResolver(node *parser.Node) *variables.Variable {
	switch node.Type() {
	case "assignment_expression":
		return variableAssignmentExpression(node)
	case "property_element":
		return variablePropertyElement(node)
	}

	return nil
}

func resolveBaseNode(node *parser.Node) *variables.Variable {
	var dataType variables.DataType
	var data interface{}
	complexity := variables.VariableComplexityObject

	if node.Type() == "encapsed_string" {
		dataType = variables.VariableDataTypeString
		data = stringutil.StripQuotes(node.Content())
		complexity = variables.VariableComplexitySimple
	}

	if node.Type() == "string" {
		dataType = variables.VariableDataTypeString
		data = stringutil.StripQuotes(node.Content())
		complexity = variables.VariableComplexitySimple
	}

	if node.Type() == "integer" {
		dataType = variables.VariableDataTypeNumber
		data = node.Content()
		complexity = variables.VariableComplexitySimple
	}

	if node.Type() == "float" {
		dataType = variables.VariableDataTypeNumber
		data = node.Content()
		complexity = variables.VariableComplexitySimple
	}

	if node.Type() == "boolean" {
		dataType = variables.VariableDataTypeBoolean
		data = node.Content()
		complexity = variables.VariableComplexitySimple
	}

	return &variables.Variable{
		DataType:   dataType,
		Complexity: complexity,
		Data:       data,
	}
}
