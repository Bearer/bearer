package context

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report/variables"
	sitter "github.com/smacker/go-tree-sitter"
)

type Resolver struct {
	Map map[parser.NodeID]*Context
}

func NewResolver(contextMap map[parser.NodeID]*Context) *Resolver {
	return &Resolver{
		Map: contextMap,
	}
}

// ResolveContext gets closest context node belongs to
func (res *Resolver) ResolveClosestContext(node *sitter.Node) *Context {
	currentNode := node
	for {
		if ctx, ok := res.Map[currentNode]; ok {
			return ctx
		}

		if currentNode.Parent() != nil {
			currentNode = currentNode.Parent()
		}
	}
}

func (res *Resolver) ResolveVariable(node *sitter.Node, variableName string) *variables.Variable {
	currentContext := res.ResolveClosestContext(node)
	for currentContext != nil {

		if variable, ok := currentContext.LocalVariables[variableName]; ok {
			return variable
		}

		currentContext = currentContext.ParentContext
	}
	return nil
}

func (res *Resolver) VariableStringValue(node *sitter.Node, variableName string) (string, error) {
	variable := res.ResolveVariable(node, variableName)

	if variable == nil {
		return "", errors.New("variable not found")
	}

	if variable.Complexity == variables.VariableComplexitySimple {
		return fmt.Sprintf("%s", variable.Data), nil
	}

	return "", errors.New("variable is complex")
}
