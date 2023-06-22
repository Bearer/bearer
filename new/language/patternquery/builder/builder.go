package builder

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/rs/zerolog/log"
)

type InputParams struct {
	Variables         []types.Variable
	MatchNodeOffset   int
	UnanchoredOffsets []int
}

type Result struct {
	Query              string
	ParamToVariable    map[string]string
	EqualParams        [][]string
	ParamToContent     map[string]map[string]string
	SingleVariableName string
}

type builder struct {
	langImplementation implementation.Implementation
	stringBuilder      strings.Builder
	idGenerator        nodeid.Generator
	inputParams        InputParams
	variableToParams   map[string][]string
	paramToContent     map[string]map[string]string
	matchNode          *tree.Node
}

func Build(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	input string,
	focusedVariable string,
) (*Result, error) {
	processedInput, inputParams, err := processInput(langImplementation, input)
	if err != nil {
		return nil, err
	}

	tree, err := lang.Parse(context.TODO(), processedInput)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	if fixedInput, fixed := fixupInput(
		langImplementation,
		processedInput,
		inputParams.Variables,
		tree.RootNode(),
	); fixed {
		tree.Close()
		tree, err = lang.Parse(context.TODO(), fixedInput)
		if err != nil {
			return nil, err
		}
	}

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
		paramToContent:     make(map[string]map[string]string),
	}

	builder.setMatchNode(
		inputParams.MatchNodeOffset,
		focusedVariable,
		langImplementation.PatternMatchNodeContainerTypes(),
		tree.RootNode(),
	)
	if builder.matchNode == nil {
		return nil, fmt.Errorf("match node not found")
	}

	result, err := builder.build(root)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func fixupInput(
	langImplementation implementation.Implementation,
	input string,
	variables []types.Variable,
	rootNode *tree.Node,
) (string, bool) {
	insideError := false
	inputOffset := 0

	byteInput := []byte(input)
	newInput := []byte(input)
	fixed := false

	err := rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		oldInsideError := insideError
		if node.Type() == "ERROR" {
			insideError = true
		}
		if err := visitChildren(); err != nil {
			return err
		}
		insideError = oldInsideError

		if !insideError {
			return nil
		}

		variable := getVariableFor(node, langImplementation, variables)
		if variable == nil {
			return nil
		}

		if log.Trace().Enabled() {
			var parentDebug, grandparentDebug string
			if parent := node.Parent(); parent != nil {
				parentDebug = parent.Debug(true)
				if grandparent := parent.Parent(); grandparent != nil {
					grandparentDebug = grandparent.Debug(true)
				}
			}

			log.Trace().Msgf("attempting pattern fixup. node: %s", node.Debug(true))
			log.Trace().Msgf("fixup parent: %s", parentDebug)
			log.Trace().Msgf("fixup grandparent: %s", grandparentDebug)
		}

		newValue := langImplementation.FixupPatternVariableDummyValue(byteInput, node, variable.DummyValue)
		if newValue == variable.DummyValue {
			return nil
		}

		fixed = true
		valueOffset := len(newValue) - len(variable.DummyValue)
		variable.DummyValue = newValue

		newInput = append(
			append(
				newInput[:node.StartByte()+inputOffset],
				newValue...,
			),
			newInput[node.EndByte()+inputOffset:]...,
		)

		inputOffset += valueOffset

		return nil
	})

	// walk errors are only ones we produce, and we don't make any
	if err != nil {
		panic(err)
	}

	return string(newInput), fixed
}

func (builder *builder) build(rootNode *tree.Node) (*Result, error) {
	if rootNode.ChildCount() == 0 {
		variable := builder.getVariableFor(rootNode)
		if variable != nil {
			return &Result{SingleVariableName: variable.Name}, nil
		}
	}

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
			node.StartLineNumber(),
			node.StartColumnNumber(),
			node.Content(),
		)
	}

	nodeAnchoredBefore, nodeAnchoredAfter := builder.langImplementation.PatternIsAnchored(node)
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
	return getVariableFor(node, builder.langImplementation, builder.inputParams.Variables)
}

func getVariableFor(
	node *tree.Node,
	langImplementation implementation.Implementation,
	variables []types.Variable,
) *types.Variable {
	for i, variable := range variables {
		if langImplementation.ShouldSkipNode(node) {
			continue
		}

		if (node.NamedChildCount() == 0 || langImplementation.IsMatchLeaf(node)) && node.Content() == variable.DummyValue {
			return &variables[i]
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

func (builder *builder) setMatchNode(
	offset int,
	focusedVariable string,
	containerTypes []string,
	node *tree.Node,
) {
	err := node.Walk(func(node *tree.Node, visitChildren func() error) error {
		if focusedVariable != "" {
			if variable := builder.getVariableFor(node); variable != nil && variable.Name == focusedVariable {
				builder.matchNode = node
				return nil
			}
		} else {
			if node.StartByte() == offset && !slices.Contains(containerTypes, node.Type()) {
				builder.matchNode = node
				return nil
			}
		}

		return visitChildren()
	})

	// walk itself shouldn't trigger an error, and we aren't creating any
	if err != nil {
		panic(err)
	}
}
