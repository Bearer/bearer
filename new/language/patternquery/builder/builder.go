package builder

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	builderinput "github.com/bearer/bearer/new/language/patternquery/builder/input"
	"github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/parser/nodeid"
)

type Result struct {
	Query           string
	ParamToVariable map[string]string
	EqualParams     [][]string
	ParamToContent  map[string]map[string]string
}

type builder struct {
	langImplementation implementation.Implementation
	stringBuilder      strings.Builder
	idGenerator        nodeid.Generator
	inputParams        builderinput.InputParams
	variableToParams   map[string][]string
	paramToContent     map[string]map[string]string
	matchNode          *tree.Node
}

func Build(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	input string,
) (*Result, error) {
	processedInput, inputParams, err := builderinput.Process(langImplementation, input)
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

	matchNode := findMatchNode(
		inputParams.MatchNodeOffset,
		langImplementation.PatternMatchNodeContainerTypes(),
		tree.RootNode(),
	)
	if matchNode == nil {
		return nil, fmt.Errorf("match node not found")
	}

	builder := builder{
		langImplementation: langImplementation,
		stringBuilder:      strings.Builder{},
		idGenerator:        &nodeid.IntGenerator{},
		inputParams:        *inputParams,
		variableToParams:   make(map[string][]string),
		paramToContent:     make(map[string]map[string]string),
		matchNode:          matchNode,
	}

	result, err := builder.build(root)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (builder *builder) build(rootNode *tree.Node) (*Result, error) {
	builder.write("(")

	if err := builder.compileNode(rootNode, true, false); err != nil {
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

	nodeAnchoredBefore, nodeAnchoredAfter := builder.langImplementation.PatternIsAnchored(node.SitterNode())
	anchored := !isRoot && node.IsNamed() && nodeAnchoredBefore

	if anchored && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.StartByte()) {
		builder.write(". ")
	}

	if variable := builder.getVariableFor(node); variable != nil {
		builder.compileVariableNode(variable)
	} else if !node.IsNamed() {
		builder.compileAnonymousNode(node)
	} else if node.NamedChildCount() == 0 {
		builder.compileLeafNode(node)
	} else if err := builder.compileNodeWithChildren(node); err != nil {
		return err
	}

	if node.Equal(builder.matchNode) {
		builder.write(" @match")
	}

	if anchored && isLastChild && nodeAnchoredAfter && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.EndByte()) {
		builder.write(" .")
	}

	return nil
}

// variable nodes match their type and capture their content
func (builder *builder) compileVariableNode(variable *types.Variable) {
	if variable.Name == "_" {
		builder.write("(_)")
		return
	}

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
	if !slices.Contains(builder.langImplementation.PatternLeafContentTypes(), node.Type()) {
		builder.write("[")

		for _, nodeType := range builder.langImplementation.PatternNodeTypes(node) {
			builder.write(" (")
			builder.write(nodeType)
			builder.write(" )")
		}

		builder.write("]")
		return
	}

	paramName := builder.newParam()
	paramContent := make(map[string]string)
	builder.paramToContent[paramName] = paramContent

	builder.write("[")

	for _, nodeType := range builder.langImplementation.PatternNodeTypes(node) {
		paramContent[nodeType] = builder.langImplementation.TranslatePatternContent(
			node.Type(),
			nodeType, node.Content(),
		)

		builder.write(" (")
		builder.write(nodeType)
		builder.write(" )")
	}

	builder.write("] @")
	builder.write(paramName)
}

// Nodes with children match their type and child nodes
func (builder *builder) compileNodeWithChildren(node *tree.Node) error {
	builder.write("[")

	var lastNode *tree.Node
	if slices.Contains(builder.langImplementation.AnonymousPatternNodeParentTypes(), node.Type()) {
		lastNode = node.Child(node.ChildCount() - 1)
	} else {
		lastNode = node.NamedChild(node.NamedChildCount() - 1)
	}

	for _, nodeType := range builder.langImplementation.PatternNodeTypes(node) {
		builder.write("(")
		builder.write(nodeType)

		for i := 0; i < node.ChildCount(); i++ {
			builder.write(" ")

			child := node.Child(i)

			if err := builder.compileNode(child, false, child.Equal(lastNode)); err != nil {
				return err
			}
		}

		builder.write(")")
	}

	builder.write("]")

	return nil
}

func (builder *builder) processVariableToParams() (map[string]string, [][]string) {
	paramToVariable := make(map[string]string)
	var equalParams [][]string

	for variableName, paramNames := range builder.variableToParams {
		if len(paramNames) > 1 {
			equalParams = append(equalParams, paramNames)
		}

		for _, paramName := range paramNames {
			paramToVariable[paramName] = variableName
		}
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

func findMatchNode(offset int, containerTypes []string, node *tree.Node) (matchNode *tree.Node) {
	err := node.Walk(func(node *tree.Node, visitChildren func() error) error {
		if node.StartByte() == offset && !slices.Contains(containerTypes, node.Type()) {
			matchNode = node
			return nil
		}

		return visitChildren()
	})

	// walk itself shouldn't trigger an error, and we aren't creating any
	if err != nil {
		panic(err)
	}

	return
}
