package schema

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/report/variables"
	sitter "github.com/smacker/go-tree-sitter"
)

type Finder struct {
	tree      *parser.Tree
	values    map[parser.NodeID]*Node
	parseNode func(finder *Finder, node *parser.Node, value *Node)
}

type Node struct {
	Terminating bool
	Variables   []*Variable
}

type Variable string

func New(tree *parser.Tree, parseNode func(finder *Finder, node *parser.Node, value *Node)) *Finder {
	return &Finder{
		tree:      tree,
		parseNode: parseNode,
		values:    make(map[parser.NodeID]*Node),
	}
}

func (finder *Finder) Annotate() {
	finder.tree.WalkBottomUp(func(child *parser.Node) error { //nolint:all,errcheck

		value := &Node{}
		finder.parseNode(finder, child, value)
		finder.values[child.ID()] = value

		return nil
	})
}

func (finder *Finder) ToVariableValues() map[*sitter.Node]*values.Value {
	newMap := make(map[*sitter.Node]*values.Value)

	for key, node := range finder.values {
		newValue := values.New()
		for _, value := range node.Variables {
			newValue.AppendVariableReference(variables.VariableName, string(*value))
		}
		newMap[key] = newValue
	}

	return newMap
}

func (finder *Finder) NonTerminatingValues(root *parser.Node) []*Variable {
	variables := []*Variable{}
	for i := 0; i < root.ChildCount(); i++ {
		child := root.Child(i)

		childValue := finder.values[child.ID()]

		if childValue.Terminating {
			continue
		}

		variables = append(variables, childValue.Variables...)
	}
	return variables
}

func (finder *Finder) ReportSchemas(detectorType detectors.Type, output report.Report) {
	finder.tree.WalkBottomUp(func(node *parser.Node) error { //nolint:all,errcheck
		value := finder.values[node.ID()]

		if !value.Terminating {
			return nil
		}

		var previousVariable *Variable
		for _, variable := range value.Variables {

			if previousVariable != nil {
				output.AddSchema(detectorType,
					schema.Schema{
						ObjectName: string(*previousVariable),
						FieldName:  string(*variable),
					},
					node.Source(true),
				)
			}

			previousVariable = variable
		}

		return nil
	})
}
