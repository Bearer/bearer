package builder

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/language/implementation"
	"github.com/bearer/curio/new/language/patternquery/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

type InputParams struct {
	Variables         []types.Variable
	MatchNodeOffset   int
	UnanchoredOffsets []int
}

type Result struct {
	Query           string
	ParamToVariable map[string]string
	EqualParams     [][]string
	ParamToContent  map[string]string
}

type builder struct {
	langImplementation implementation.Implementation
	stringBuilder      strings.Builder
	idGenerator        nodeid.Generator
	inputParams        InputParams
	variableToParams   map[string][]string
	paramToContent     map[string]string
	matchNodeFound     bool
}

func Build(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	input string,
) (*Result, error) {
	processedInput, inputParams, err := processInput(langImplementation, input)
	if err != nil {
		return nil, err
	}

	tree, err := lang.Parse(processedInput)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()

	if root.ChildCount() != 1 {
		return nil, fmt.Errorf("expecting 1 node but got %d", root.ChildCount())
	}

	for {
		root = root.Child(0)

		if langImplementation.IsRootOfRuleQuery(root) {
			break
		}
	}

	builder := builder{
		langImplementation: langImplementation,
		stringBuilder:      strings.Builder{},
		idGenerator:        &nodeid.IntGenerator{},
		inputParams:        *inputParams,
		variableToParams:   make(map[string][]string),
		paramToContent:     make(map[string]string),
	}

	result, err := builder.build(root)
	if err != nil {
		return nil, err
	}

	if !builder.matchNodeFound {
		return nil, fmt.Errorf("match node not found")
	}

	return result, nil
}

func (builder *builder) build(rootNode *tree.Node) (*Result, error) {
	builder.write("(")

	err := builder.compileNode(rootNode, true, false)
	if err != nil {
		return nil, err
	}

	builder.write(" @root")
	builder.write(")")

	paramToVariable, equalParams := builder.processVariableToParams()

	return &Result{
		Query:           builder.stringBuilder.String(),
		ParamToVariable: paramToVariable,
		EqualParams:     equalParams,
		ParamToContent:  builder.paramToContent,
	}, nil
}

func (builder *builder) compileNode(node *tree.Node, isRoot bool, isLastChild bool) error {
	if node.IsError() {
		return fmt.Errorf(
			"error parsing pattern at %d:%d: %s",
			node.LineNumber(),
			node.ColumnNumber(),
			node.Content(),
		)
	}

	writeMatch := false
	if !builder.matchNodeFound &&
		node.StartByte() == builder.inputParams.MatchNodeOffset &&
		!slices.Contains(builder.langImplementation.PatternMatchNodeContainerTypes(), node.Type()) {
		builder.matchNodeFound = true
		writeMatch = true
	}

	anchored := !isRoot && node.IsNamed() && builder.langImplementation.PatternIsAnchored(node)

	if anchored && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.StartByte()) {
		builder.write(". ")
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

	if anchored && isLastChild && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.EndByte()) {
		builder.write(" .")
	}

	if writeMatch {
		builder.write(" @match")
	}

	return nil
}

// variable nodes match their type and capture their content
func (builder *builder) compileVariableNode(variable *types.Variable) {
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
	if !slices.Contains(builder.langImplementation.AnonymousPatternNodeParentTypes(), node.Parent().Type()) {
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

		if err := builder.compileNode(node.Child(i), false, i == node.ChildCount()-1); err != nil {
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

func (builder *builder) getVariableFor(node *tree.Node) *types.Variable {
	for _, variable := range builder.inputParams.Variables {
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
