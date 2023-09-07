package builder

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

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
	// FIXME: Can this ever be valid for languages like PHP?
	// if len(root.Children()) != 1 {
	// 	return nil, fmt.Errorf("expecting 1 node but got %d", len(root.Children()))
	// }

	for {
		for _, children := range root.Children() {
			if !patternLanguage.ShouldSkip(children) {
				root = children
				break
			}
		}

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
		if node.IsError() {
			insideError = true
		}
		if err := visitChildren(); err != nil {
			return err
		}
		insideError = oldInsideError

		if !insideError && !node.IsMissing() {
			return nil
		}

		var newValue string
		var originalValue string

		if insideError {
			variable := getVariableFor(node, patternLanguage, variables)
			if variable == nil {
				return nil
			}

			if log.Trace().Enabled() {
				log.Trace().Msgf("attempting pattern fixup. node: %s", node.Debug())
			}

			newValue = patternLanguage.FixupVariableDummyValue(byteInput, node, variable.DummyValue)
			if newValue == variable.DummyValue {
				return nil
			}
			variable.DummyValue = newValue
			originalValue = variable.DummyValue
		} else {
			if log.Trace().Enabled() {
				log.Trace().Msgf("attempting pattern fixup (missing node). node: %s", node.Debug())
			}

			newValue = patternLanguage.FixupMissing(node)
			if newValue == "" {
				return nil
			}
		}

		fixed = true
		valueOffset := len(newValue) - len(originalValue)

		prefix := newInput[:node.ContentStart.Byte+inputOffset]
		suffix := newInput[node.ContentEnd.Byte+inputOffset:]
		// FIXME: We need to append suffix before
		// newInput seems to be sharing memory in some circumstances
		// suffix before and after the first append are not equal
		appendedInput := appendByte([]byte(newValue), suffix...)
		newInput = appendByte(prefix, appendedInput...)

		inputOffset += valueOffset

		return nil
	})

	// walk errors are only ones we produce, and we don't make any
	if err != nil {
		panic(err)
	}

	return newInput, fixed
}

func appendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
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
	anchored := !isRoot && node.IsNamed() && nodeAnchoredBefore

	if anchored && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.ContentStart.Byte) {
		builder.write(". ")
	}

	if variable := builder.getVariableFor(node); variable != nil {
		builder.compileVariableNode(variable)
	} else if !node.IsNamed() {
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
