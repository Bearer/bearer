package javascript

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ssoroka/slice"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/regex"

	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
)

var (
	variableLookupParents = []string{
		"pair",
		"arguments",
		"binary_expression",
		"template_substitution",
		"array",
		"spread_element",
		"augmented_assignment_expression",
	}

	anonymousPatternNodeParentTypes = []string{}
	patternMatchNodeContainerTypes  = []string{"import_clause", "import_specifier", "required_parameter"}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	allowedPatternQueryTypes  = []string{"identifier", "property_identifier", "_", "member_expression", "string", "template_string"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)

	ellipsisRegex = regexp.MustCompile(`\$<\.\.\.>`)

	passthroughMethods = []string{"JSON.parse", "JSON.stringify"}
)

type javascriptImplementation struct {
	implementation.Base
}

func Get() implementation.Implementation {
	return &javascriptImplementation{}
}

func (implementation *javascriptImplementation) SitterLanguage() *sitter.Language {
	return tsx.GetLanguage()
}

func (*javascriptImplementation) AnalyzeFlow(ctx context.Context, rootNode *tree.Node) error {
	scope := implementation.NewScope(nil)

	return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		switch node.Type() {
		// () => {}
		// function getName() {}
		case "function", "arrow_function", "method_definition":
			previousScope := scope
			scope = implementation.NewScope(previousScope)
			err := visitChildren()
			scope = previousScope
			return err
		// user = ...
		case "assignment_expression":
			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if left.Type() == "identifier" {
				err := visitChildren()

				scope.Assign(left.Content(), node)
				node.UnifyWith(right)

				return err
			}
		// x += y
		case "augmented_assignment_expression":
			err := visitChildren()

			left := node.ChildByFieldName("left")
			if left.Type() == "identifier" {
				scope.Assign(left.Content(), node)
			}

			return err
		// const user = ...
		// var user = ...
		// let user = ...
		case "variable_declarator":
			name := node.ChildByFieldName("name")
			value := node.ChildByFieldName("value")

			if name.Type() == "identifier" {
				err := visitChildren()

				scope.Assign(name.Content(), node)
				node.UnifyWith(value)

				return err
			}
		case "shorthand_property_identifier_pattern":
			scope.Assign(node.Content(), node)
		case "identifier":
			parent := node.Parent()
			if parent == nil {
				break
			}

			if slice.Contains(variableLookupParents, parent.Type()) ||
				(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("right"))) ||
				(parent.Type() == "new_expression" && node.Equal(parent.ChildByFieldName("constructor"))) ||
				(parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("value"))) ||
				(parent.Type() == "member_expression" && node.Equal(parent.ChildByFieldName("object"))) ||
				(parent.Type() == "call_expression" && node.Equal(parent.ChildByFieldName("function"))) ||
				(parent.Type() == "subscript_expression" && node.Equal(parent.ChildByFieldName("object"))) {
				if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
					node.UnifyWith(scopedNode)
				}

				break
			}

			// typescript: different type of identifier
			if parent.Type() == "required_parameter" {
				scope.Assign(node.Content(), node)
				break
			}

			if parent.Type() == "arguments" {
				callNode := parent.Parent()
				callNode.UnifyWith(node)
				break
			}

			if isImportedIdentifier(node) {
				scope.Assign(node.Content(), node)
			}
		case "property_identifier":
			parent := node.Parent()
			if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
				if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}
		}

		return visitChildren()
	})
}

func (implementation *javascriptImplementation) IsMatchLeaf(node *tree.Node) bool {
	return node.Type() == "string"
}

func (implementation *javascriptImplementation) ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error) {
	nameIndex := patternQueryVariableRegex.SubexpIndex("name")
	typesIndex := patternQueryVariableRegex.SubexpIndex("types")
	i := 0

	var params []patternquerytypes.Variable

	replaced, err := regex.ReplaceAllWithSubmatches(patternQueryVariableRegex, input, func(submatches []string) (string, error) {
		nodeTypes := strings.Split(submatches[typesIndex], "|")
		if nodeTypes[0] == "" {
			nodeTypes = []string{"_"}
		}

		for _, nodeType := range nodeTypes {
			if !slices.Contains(allowedPatternQueryTypes, nodeType) {
				return "", fmt.Errorf("invalid node type '%s' in pattern query", nodeType)
			}
		}

		dummyValue := produceDummyValue(i, nodeTypes[0])

		params = append(params, patternquerytypes.Variable{
			Name:       submatches[nameIndex],
			NodeTypes:  nodeTypes,
			DummyValue: dummyValue,
		})

		i += 1

		return dummyValue, nil
	})

	if err != nil {
		return "", nil, err
	}

	return replaced, params, nil
}

func produceDummyValue(i int, nodeType string) string {
	return "CurioVar" + fmt.Sprint(i)
}

func (implementation *javascriptImplementation) AnonymousPatternNodeParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

func (implementation *javascriptImplementation) FindPatternMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (implementation *javascriptImplementation) FindPatternUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func (implementation *javascriptImplementation) PatternMatchNodeContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (javascriptImplementation *javascriptImplementation) ShouldSkipNode(node *tree.Node) bool {
	return node.Type() == "required_parameter"
}

func (*javascriptImplementation) PatternLeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "property_identifier", "shorthand_property_identifier", "type_identifier",
		// datatypes/literals
		"template_string", "string_fragment", "number", "null", "true", "false",
	}
}

func (implementation *javascriptImplementation) PatternIsAnchored(node *tree.Node) (bool, bool) {
	if node.Type() == "pair" {
		return false, false
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Class body class_body
	// arrow functions statement_block
	// function statement_block
	// method statement_block
	unAnchored := []string{"statement_block", "class_body", "object_pattern", "named_imports"}

	isUnanchored := !slices.Contains(unAnchored, parent.Type())
	return isUnanchored, isUnanchored
}

func (implementation *javascriptImplementation) IsRootOfRuleQuery(node *tree.Node) bool {
	return !(node.Type() == "expression_statement")
}

func (implementation *javascriptImplementation) PatternNodeTypes(node *tree.Node) []string {
	if node.Type() == "statement_block" && node.Parent().Type() == "program" {
		if node.NamedChildCount() == 0 {
			return []string{"object"}
		} else {
			return []string{node.Type(), "program"}
		}
	}

	return []string{node.Type()}
}

func (*javascriptImplementation) PassthroughNested(node *tree.Node) bool {
	if node.Type() != "arguments" {
		return false
	}

	callNode := node.Parent()
	if callNode.Type() != "call_expression" {
		return false
	}

	functionNode := callNode.ChildByFieldName("function")

	var method string
	var wildcardMethod string
	switch functionNode.Type() {
	case "identifier":
		return slices.Contains(passthroughMethods, functionNode.Content())
	case "member_expression":
		object := functionNode.ChildByFieldName("object")
		if object.Type() == "identifier" {
			property := functionNode.ChildByFieldName("property").Content()
			method = object.Content() + "." + property
			wildcardMethod = "*." + property
		}
	}

	return slices.Contains(passthroughMethods, method) || slices.Contains(passthroughMethods, wildcardMethod)
}

func (*javascriptImplementation) ContributesToResult(node *tree.Node) bool {
	// Statements don't have results
	if strings.HasSuffix(node.Type(), "_statement") {
		return false
	}

	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Must not be a ternary condition
	if parent.Type() == "ternary_expression" && node.Equal(parent.ChildByFieldName("condition")) {
		return false
	}

	// Not the name part of a declaration
	if parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("name")) {
		return false
	}

	// Not the left part of an assignment
	if parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("left")) {
		return false
	}

	return true
}

func isImportedIdentifier(node *tree.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return false
	}

	// import x from "library"
	if parent.Type() == "import_clause" {
		return true
	}

	// import * as x from "library"
	if parent.Type() == "namespace_import" {
		return true
	}

	if parent.Type() != "import_specifier" {
		return false
	}

	// import { x } from "library"
	if parent.ChildByFieldName("alias") == nil {
		return true
	}

	// import { a as x } from "library"
	if node.Equal(parent.ChildByFieldName("alias")) {
		return true
	}

	return false
}
