package builder

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

type Variable struct {
	NodeTypes  []string
	DummyValue string
	Name       string
}

type Result struct {
	Query           string
	ParamToVariable map[string]string
	EqualParams     [][]string
	ParamToContent  map[string]string
}

type builder struct {
	anonymousParentTypes []string
	stringBuilder        strings.Builder
	idGenerator          nodeid.Generator
	variables            []Variable
	variableToParams     map[string][]string
	paramToContent       map[string]string
	matchNodeOffset      int
	matchNodeFound       bool
}

func Build(
	lang types.Language,
	anonymousParentTypes []string,
	input string,
	variables []Variable,
	matchNodeOffset int,
) (*Result, error) {
	tree, err := lang.Parse(input)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	if tree.RootNode().ChildCount() != 1 {
		return nil, fmt.Errorf("expecting 1 node but got %d", tree.RootNode().ChildCount())
	}

	builder := builder{
		anonymousParentTypes: anonymousParentTypes,
		stringBuilder:        strings.Builder{},
		idGenerator:          &nodeid.IntGenerator{},
		variables:            variables,
		variableToParams:     make(map[string][]string),
		paramToContent:       make(map[string]string),
		matchNodeOffset:      matchNodeOffset,
	}

	result := builder.build(tree.RootNode().Child(0))

	if !builder.matchNodeFound {
		return nil, fmt.Errorf("match node not found")
	}

	return result, nil
}

func (builder *builder) build(rootNode *tree.Node) *Result {
	builder.write("(")
	builder.compileNode(rootNode)
	builder.write(" @root")
	builder.write(")")

	paramToVariable, equalParams := builder.processVariableToParams()

	return &Result{
		Query:           builder.stringBuilder.String(),
		ParamToVariable: paramToVariable,
		EqualParams:     equalParams,
		ParamToContent:  builder.paramToContent,
	}
}

func (builder *builder) compileNode(node *tree.Node) error {
	if node.IsError() {
		return fmt.Errorf(
			"error parsing pattern at %d:%d: %s",
			node.LineNumber(),
			node.ColumnNumber(),
			node.Content(),
		)
	}

	writeMatch := false
	if !builder.matchNodeFound && node.StartByte() == builder.matchNodeOffset {
		builder.matchNodeFound = true
		writeMatch = true
	}

	if variable := builder.getVariableFor(node); variable != nil {
		builder.compileVariableNode(variable)
	} else if !node.IsNamed() {
		builder.compileAnonymousNode(node)
	} else if node.ChildCount() == 0 {
		builder.compileLeafNode(node)
	} else if err := builder.compileNodeWithChildren(node); err != nil {
		return err
	}

	if writeMatch {
		builder.write(" @match")
	}

	return nil
}

// variable nodes match their type and capture their content
func (builder *builder) compileVariableNode(variable *Variable) {
	paramName := builder.newParam()
	builder.variableToParams[variable.Name] = append(builder.variableToParams[variable.Name], paramName)

	builder.write("[")

	for _, nodeType := range variable.NodeTypes {
		builder.write("(")
		builder.write(nodeType)
		builder.write(")")
	}

	builder.write("] @")
	builder.write(paramName)
}

// Anonymous nodes match their content as a literal
func (builder *builder) compileAnonymousNode(node *tree.Node) {
	if !slices.Contains(builder.anonymousParentTypes, node.Parent().Type()) {
		return
	}

	builder.write(strconv.Quote(node.Content()))
}

// Leaves match their type and content
func (builder *builder) compileLeafNode(node *tree.Node) {
	paramName := builder.newParam()
	builder.paramToContent[paramName] = node.Content()

	builder.write("(")
	builder.write(node.Type())
	builder.write(") @")
	builder.write(paramName)
}

// Nodes with children match their type and child nodes
func (builder *builder) compileNodeWithChildren(node *tree.Node) error {
	builder.write("(")
	builder.write(node.Type())

	for i := 0; i < node.ChildCount(); i++ {
		builder.write(" ")

		if err := builder.compileNode(node.Child(i)); err != nil {
			return err
		}
	}

	builder.write(")")

	return nil
}

func (builder *builder) processVariableToParams() (map[string]string, [][]string) {
	paramToVariable := make(map[string]string)
	var equalParams [][]string

	for variableName, paramNames := range builder.variableToParams {
		if len(paramNames) > 1 {
			equalParams = append(equalParams, paramNames)
		}

		paramToVariable[paramNames[0]] = variableName
	}

	return paramToVariable, equalParams
}

func (builder *builder) getVariableFor(node *tree.Node) *Variable {
	for _, variable := range builder.variables {
		if node.Content() == variable.DummyValue {
			return &variable
		}
	}

	return nil
}

func (builder *builder) write(value string) {
	builder.stringBuilder.WriteString(value)
}

func (builder *builder) newParam() string {
	return "param" + builder.idGenerator.GenerateId()
}
