package builder

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/language"
)

type InputParams struct {
	VariableNames     []string
	Variables         []language.PatternVariable
	MatchNodeOffset   int
	UnanchoredOffsets []int
}

type Result struct {
	Query           string
	VariableNames   []string
	ParamToVariable map[string]string
	EqualParams     [][]string
	ParamToContent  map[string]map[string]string
	RootVariable    *language.PatternVariable
}

type builder struct {
	patternLanguage  language.Pattern
	stringBuilder    strings.Builder
	idGenerator      nodeid.Generator
	inputParams      InputParams
	variableToParams map[string][]string
	paramToContent   map[string]map[string]string
	matchNode        *tree.Node
}

func Build(
	language language.Language,
	input string,
	focusedVariable string,
) (*Result, error) {
	patternLanguage := language.Pattern()
	processedInput, inputParams, err := processInput(patternLanguage, input)
	if err != nil {
		return nil, err
	}

	tree, err := ast.Parse(context.TODO(), language, processedInput)
	if err != nil {
		return nil, err
	}

	if fixedInput, fixed := fixupInput(
		patternLanguage,
		processedInput,
		inputParams.Variables,
		tree.RootNode(),
	); fixed {
		tree, err = ast.Parse(context.TODO(), language, fixedInput)
		if err != nil {
			return nil, err
		}
	}

	root := tree.RootNode()

	if len(root.Children()) != 1 {
		return nil, fmt.Errorf("expecting 1 node but got %d", len(root.Children()))
	}

	for {
		root = root.Children()[0]

		if patternLanguage.IsRoot(root) {
			break
		}
	}

	builder := builder{
		patternLanguage:  patternLanguage,
		stringBuilder:    strings.Builder{},
		idGenerator:      &nodeid.IntGenerator{},
		inputParams:      *inputParams,
		variableToParams: make(map[string][]string),
		paramToContent:   make(map[string]map[string]string),
	}

	builder.setMatchNode(
		inputParams.MatchNodeOffset,
		focusedVariable,
		patternLanguage.ContainerTypes(),
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
	patternLanguage language.Pattern,
	byteInput []byte,
	variables []language.PatternVariable,
	rootNode *tree.Node,
) ([]byte, bool) {
	insideError := false
	inputOffset := 0

	newInput := make([]byte, len(byteInput))
	copy(newInput, byteInput)
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

		variable := getVariableFor(node, patternLanguage, variables)
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

		newValue := patternLanguage.FixupVariableDummyValue(byteInput, node, variable.DummyValue)
		if newValue == variable.DummyValue {
			return nil
		}

		fixed = true
		valueOffset := len(newValue) - len(variable.DummyValue)
		variable.DummyValue = newValue

		newInput = append(
			append(
				newInput[:node.ContentStart.Byte+inputOffset],
				newValue...,
			),
			newInput[node.ContentEnd.Byte+inputOffset:]...,
		)

		inputOffset += valueOffset

		return nil
	})

	// walk errors are only ones we produce, and we don't make any
	if err != nil {
		panic(err)
	}

	return newInput, fixed
}

func (builder *builder) build(rootNode *tree.Node) (*Result, error) {
	if len(rootNode.Children()) == 0 {
		variable := builder.getVariableFor(rootNode)
		if variable != nil {
			return &Result{RootVariable: variable}, nil
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
		VariableNames:   builder.inputParams.VariableNames,
		ParamToVariable: paramToVariable,
		EqualParams:     equalParams,
		ParamToContent:  builder.paramToContent,
	}, nil
}

func (builder *builder) compileNode(node *tree.Node, isRoot bool, isLastChild bool) error {
	if node.SitterNode().IsError() {
		return fmt.Errorf(
			"error parsing pattern at %d:%d: %s",
			node.ContentStart.Line,
			node.ContentStart.Column,
			node.Content(),
		)
	}

	nodeAnchoredBefore, nodeAnchoredAfter := builder.patternLanguage.IsAnchored(node)
	anchored := !isRoot && node.SitterNode().IsNamed() && nodeAnchoredBefore

	if anchored && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.ContentStart.Byte) {
		builder.write(". ")
	}

	if variable := builder.getVariableFor(node); variable != nil {
		builder.compileVariableNode(variable)
	} else if !node.SitterNode().IsNamed() {
		builder.compileAnonymousNode(node)
	} else if len(node.NamedChildren()) == 0 {
		builder.compileLeafNode(node)
	} else if err := builder.compileNodeWithChildren(node); err != nil {
		return err
	}

	if node == builder.matchNode {
		builder.write(" @match")
	}

	if anchored && isLastChild && nodeAnchoredAfter && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.ContentEnd.Byte) {
		builder.write(" .")
	}

	return nil
}

// variable nodes match their type and capture their content
func (builder *builder) compileVariableNode(variable *language.PatternVariable) {
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
	if !slices.Contains(builder.patternLanguage.AnonymousParentTypes(), node.Parent().Type()) {
		return
	}

	builder.write(strconv.Quote(node.Content()))
}

// Leaves match their type and content
func (builder *builder) compileLeafNode(node *tree.Node) {
	if !slices.Contains(builder.patternLanguage.LeafContentTypes(), node.Type()) {
		builder.write("[")

		for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
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

	for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
		paramContent[nodeType] = builder.patternLanguage.TranslateContent(
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

	var children []*tree.Node
	if slices.Contains(builder.patternLanguage.AnonymousParentTypes(), node.Type()) {
		children = node.Children()
	} else {
		children = node.NamedChildren()
	}

	lastNode := children[len(children)-1]

	for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
		builder.write("(")
		builder.write(nodeType)

		for _, child := range node.Children() {
			builder.write(" ")

			if err := builder.compileNode(child, false, child == lastNode); err != nil {
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

func (builder *builder) getVariableFor(node *tree.Node) *language.PatternVariable {
	return getVariableFor(node, builder.patternLanguage, builder.inputParams.Variables)
}

func getVariableFor(
	node *tree.Node,
	patternLanguage language.Pattern,
	variables []language.PatternVariable,
) *language.PatternVariable {
	for i, variable := range variables {
		if (len(node.NamedChildren()) == 0 || patternLanguage.IsLeaf(node)) && node.Content() == variable.DummyValue {
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
			if node.ContentStart.Byte == offset && !slices.Contains(containerTypes, node.Type()) {
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
