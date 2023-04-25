package java

import (
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

type javaImplementation struct{}

func Get() implementation.Implementation {
	return &javaImplementation{}
}

func (implementation *javaImplementation) SitterLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (*javaImplementation) AnalyzeFlow(rootNode *tree.Node) error {
	scope := implementation.NewScope(nil)

	return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
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
		case "method_declaration", "lambda_expression":
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
		case "field_declaration":
			declarator := node.ChildByFieldName("declarator")
			if declarator != nil {
				name := declarator.ChildByFieldName("name")
				value := declarator.ChildByFieldName("value")
				scope.Assign(name.Content(), value)
				node.UnifyWith(value)
			}
		// String user = "John";
		case "local_variable_declaration":
			declarator := node.ChildByFieldName("declarator")

			name := declarator.ChildByFieldName("name")
			value := declarator.ChildByFieldName("value")

			if name.Type() == "identifier" {
				err := visitChildren()

				scope.Assign(name.Content(), node)
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

			// user.name
			// a = user.name
			// a = user.name()
			// var a = user
			// a = user
			// user["name"] = name;
			// todo: new expression
			if slice.Contains(variableLookupParents, parent.Type()) ||
				(parent.Type() == "scoped_type_identifier" && node.Equal(parent.Child(0))) ||
				(parent.Type() == "method_invocation" && node.Equal(parent.ChildByFieldName("object"))) ||
				(parent.Type() == "field_access" && node.Equal(parent.ChildByFieldName("object"))) ||
				(parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("value"))) ||
				(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("right"))) ||
				(parent.Type() == "array_access" && node.Equal(parent.ChildByFieldName("right"))) {
				if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}

			// user.getName(name, lastName)
			if parent.Type() == "arguments_list" {
				scope.Assign(node.Content(), node)
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

func (javaImplementation *javaImplementation) ShouldSkipNode(node *tree.Node) bool {
	return false
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
	unAnchored := []string{"class_body", "block"}

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

func (implementation *javaImplementation) TranslatePatternContent(fromNodeType, toNodeType, content string) string {
	return content
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
