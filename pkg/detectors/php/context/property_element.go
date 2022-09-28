package context

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/variables"
)

func variablePropertyElement(node *parser.Node) *variables.Variable {
	variable := &variables.Variable{}

	foundName := false
	foundIntializer := false

	for i := 0; i < node.ChildCount(); i++ {
		child := node.Child(i)

		if child.Type() == "variable_name" {
			foundName = true
			variable.Name = child.Child(0).Content()
		}

		if child.Type() == "property_initializer" {
			foundIntializer = true
			returnVariable := resolveBaseNode(child.Child(0))
			variable.Complexity = returnVariable.Complexity
			variable.Data = returnVariable.Data
			variable.DataType = returnVariable.DataType
		}
	}

	if foundName && foundIntializer {
		return variable
	}

	return nil
}
