package java

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/regex"
	"github.com/ssoroka/slice"

	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
)

var (
	variableLookupParents = []string{
		"field_declaration",
		"argument_list",
		"binary_expression",
		"array_access",
		"array_initializer",
	}

	anonymousPatternNodeParentTypes = []string{}
	patternMatchNodeContainerTypes  = []string{}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)

	// todo: see if it is ok to replace typescripts `member_expression` with javas `field_access` and `method_invocation`
	allowedPatternQueryTypes = []string{"identifier", "type_identifier", "_", "field_access", "method_invocation", "string_literal"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)

	ellipsisRegex = regexp.MustCompile(`\$<\.\.\.>`)

	passthroughMethods = []string{}
)

type javaImplementation struct {
	implementation.Base
}

func Get() implementation.Implementation {
	return &javaImplementation{}
}

func (implementation *javaImplementation) SitterLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (*javaImplementation) AnalyzeFlow(ctx context.Context, rootNode *tree.Node) error {
	scope := implementation.NewScope(nil)

	return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		switch node.Type() {
		// public class Main {
		//
		// }
		case "class_body":
			previousScope := scope
			scope = implementation.NewScope(previousScope)
			err := visitChildren()
			scope = previousScope
			return err
		// public class Main {
		//	// method declaration
		//	static void myMethod() {
		//
		//   }
		// }
		//
		// lambda_expression
		// numbers.forEach( (n) -> { System.out.println(n); } );
		case "method_declaration", "lambda_expression", "for_statement", "enhanced_for_statement", "block":
			previousScope := scope
			scope = implementation.NewScope(previousScope)
			err := visitChildren()
			scope = previousScope
			return err
		// user = ...
		case "assignment_expression":
			err := visitChildren()

			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if node.AnonymousChild(0).Content() == "=" {
				node.UnifyWith(right)
			}

			if left.Type() == "identifier" {
				scope.Assign(left.Content(), node)
			}

			return err
		case "field_declaration":
			err := visitChildren()

			declarator := node.ChildByFieldName("declarator")
			if declarator != nil {
				scope.Declare(declarator.ChildByFieldName("name").Content(), node)

				if value := declarator.ChildByFieldName("value"); value != nil {
					node.UnifyWith(value)
				}
			}

			return err
		// String user = "John";
		case "local_variable_declaration":
			declarator := node.ChildByFieldName("declarator")

			name := declarator.ChildByFieldName("name")
			value := declarator.ChildByFieldName("value")

			if name.Type() == "identifier" {
				err := visitChildren()

				scope.Declare(name.Content(), node)
				node.UnifyWith(value)

				return err
			}
		// // TODO: figure out this one
		// case "shorthand_property_identifier_pattern":
		// 	scope.Assign(node.Content(), node)
		case "identifier":
			parent := node.Parent()
			if parent == nil {
				break
			}

			if slice.Contains(variableLookupParents, parent.Type()) ||
				(parent.Type() == "scoped_type_identifier" && node.Equal(parent.Child(0))) ||
				(parent.Type() == "method_invocation" && node.Equal(parent.ChildByFieldName("object"))) ||
				(parent.Type() == "field_access" && node.Equal(parent.ChildByFieldName("object"))) ||
				(parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("value"))) ||
				(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("right"))) ||
				(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("left")) && parent.AnonymousChild(0).Content() != "=") ||
				(parent.Type() == "enhanced_for_statement" && node.Equal(parent.ChildByFieldName("value"))) {
				if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}

			if parent.Type() == "formal_parameter" ||
				parent.Type() == "catch_formal_parameter" ||
				(parent.Type() == "resource" && node.Equal(parent.ChildByFieldName("name"))) {
				scope.Declare(node.Content(), node)
			}

			if parent.Type() == "enhanced_for_statement" && node.Equal(parent.ChildByFieldName("name")) {
				scope.Declare(node.Content(), node)
				node.UnifyWith(parent.ChildByFieldName("value"))
			}

			// todo: see what this is
			// case "property_identifier":
			// 	parent := node.Parent()
			// 	if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
			// 		if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
			// 			node.UnifyWith(scopedNode)
			// 		}
			// 	}
			// }
		}
		return visitChildren()
	})
}

// TODO: See if anything needs to be added here
func (implementation *javaImplementation) ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error) {
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

// TODO: See if anything needs to be added here
func (implementation *javaImplementation) AnonymousPatternNodeParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

// TODO: See if anything needs to be added here
func (implementation *javaImplementation) FindPatternMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

// TODO: See if anything needs to be added here
func (implementation *javaImplementation) FindPatternUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func (implementation *javaImplementation) PatternMatchNodeContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (*javaImplementation) PatternLeafContentTypes() []string {
	return []string{
		// todo: see if type identifier should be removed from here (User user) `User` is type
		// identifiers
		"identifier", "modifier",
		// types
		// int user, User user, void user function,
		"integral_type", "type_identifier", "void_type",
		// datatypes/literals
		"string_literal", "character_literal", "null_literal", "true", "false", "decimal_integer_literal", "decimal_floating_point_literal",
	}
}

func (implementation *javaImplementation) PatternIsAnchored(node *tree.Node) (bool, bool) {
	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Class body class_body
	// function block
	// lambda () -> {} block
	// try {} catch () {}
	unAnchored := []string{"class_body", "block", "try_statement", "catch_type", "resource_specification"}

	isUnanchored := !slices.Contains(unAnchored, parent.Type())
	return isUnanchored, isUnanchored
}

func (implementation *javaImplementation) IsRootOfRuleQuery(node *tree.Node) bool {
	return !(node.Type() == "expression_statement")
}

func (implementation *javaImplementation) PatternNodeTypes(node *tree.Node) []string {
	if node.Type() == "statement_block" && node.Parent().Type() == "program" {
		if node.NamedChildCount() == 0 {
			return []string{"object"}
		} else {
			return []string{node.Type(), "program"}
		}
	}

	return []string{node.Type()}
}

func (*javaImplementation) PassthroughNested(node *tree.Node) bool {
	if node.Type() != "arguments" {
		return false
	}

	callNode := node.Parent()
	if callNode.Type() != "field_access" {
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

func (*javaImplementation) ContributesToResult(node *tree.Node) bool {
	// Statements don't have results
	if strings.HasSuffix(node.Type(), "_statement") {
		return false
	}

	// Switch case
	if node.Type() == "switch_label" {
		return false
	}

	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Must not be a ternary/switch condition
	if node.Equal(parent.ChildByFieldName("condition")) {
		return false
	}

	// Not the name part of a declaration
	if parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("name")) {
		return false
	}

	// Not the left part of an `=` assignment
	if parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("left")) {
		return parent.AnonymousChild(0).Content() != "="
	}

	return true
}

func (*javaImplementation) FixupPatternVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	return dummyValue
}
