package pattern

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/util/regex"
)

var (
	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex      = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex                 = regexp.MustCompile(`\$<!>`)
	ellipsisRegex                  = regexp.MustCompile(`\$<\.\.\.>`)
	unanchoredPatternNodeTypes     = []string{"catch_clause"}
	patternMatchNodeContainerTypes = []string{"formal_parameters", "simple_parameter", "argument", "type_list"}

	allowedPatternQueryTypes = []string{"_"}

	functionRegex      = regexp.MustCompile(`\bfunction\b`)
	parameterTypeRegex = regexp.MustCompile(`[,(]\s*(public|private|protected|var)?\s*\z`)
)

type Pattern struct {
	language.PatternBase
}

func (*Pattern) AdjustInput(input string) string {
	return "<?php " + input
}

func (*Pattern) FixupMissing(node *tree.Node) string {
	if node.Type() != `";"` {
		return ""
	}

	return ";"
}

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	addDollar := false

	if parent := node.Parent(); parent != nil {
		if parent.Type() == "named_type" {
			addDollar = true
		}

		if parent.Type() == "ERROR" && parent.Parent() != nil && parent.Parent().Type() == "declaration_list" {
			parentContent := []byte(parent.Content())
			parentPrefix := string(parentContent[:node.ContentStart.Byte-parent.ContentStart.Byte])

			isFunctionName := functionRegex.MatchString(parentPrefix) && !strings.Contains(parentPrefix, "(")
			addDollar = !isFunctionName && !parameterTypeRegex.MatchString(parentPrefix)
		}
	}

	if addDollar {
		return "$" + dummyValue
	}

	return dummyValue
}

func (*Pattern) ExtractVariables(input string) (string, []language.PatternVariable, error) {
	nameIndex := patternQueryVariableRegex.SubexpIndex("name")
	typesIndex := patternQueryVariableRegex.SubexpIndex("types")
	i := 0

	var params []language.PatternVariable

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

		params = append(params, language.PatternVariable{
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
	return "BearerVar" + fmt.Sprint(i)
}

func (*Pattern) FindMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (*Pattern) FindUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func (*Pattern) IsLeaf(node *tree.Node) bool {
	// Encapsed string literal
	switch node.Type() {
	case "encapsed_string":
		namedChildren := node.NamedChildren()
		if len(namedChildren) == 1 && namedChildren[0].Type() == "string" {
			return true
		}
	}
	return false
}

func (*Pattern) AnonymousParentTypes() []string {
	return []string{
		"binary_expression",
		"unary_op_expression",
	}
}

func (*Pattern) LeafContentTypes() []string {
	return []string{
		"string_value",
		"name",
		"integer",
		"float",
		"boolean",
	}
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
	if slices.Contains(unanchoredPatternNodeTypes, node.Type()) {
		return false, false
	}

	// Named arguments are unanchored
	// eg. f(x: 42)
	if node.Type() == "argument" && node.ChildByFieldName("name") != nil {
		return false, false
	}

	if node.Type() == "property_element" {
		return false, true
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// optional type on parameters
	if slices.Contains([]string{"property_promotion_parameter", "simple_parameter"}, parent.Type()) &&
		node == parent.ChildByFieldName("name") {
		return false, true
	}

	if slices.Contains([]string{
		"method_declaration",
		"function_definition",
		"anonymous_function_creation_expression",
	}, parent.Type()) {
		// visibility
		if node == parent.ChildByFieldName("name") {
			return false, true
		}

		// type
		if node == parent.ChildByFieldName("body") {
			return false, true
		}

		return false, false
	}

	// Associative array elements are unanchored
	// eg. array("foo" => 42)
	if parent.Type() == "array_creation_expression" &&
		node.Type() == "array_element_initializer" &&
		len(node.NamedChildren()) == 2 {
		return false, false
	}

	// `new Foo` should match `new Foo()`
	if parent.Type() == "object_creation_expression" {
		if node == parent.NamedChildren()[0] {
			return true, false
		}
	}

	// Class body declaration_list
	// function/block compound_statement
	// Type1 | Type2
	unAnchored := []string{
		"declaration_list",
		"compound_statement",
		"type_list",
	}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"expression_statement", "php_tag", "program"}, node.Type()) && !node.IsMissing()
}

func (patternLanguage *Pattern) NodeTypes(node *tree.Node) []string {
	parent := node.Parent()
	if parent == nil {
		return []string{node.Type()}
	}

	if (node.Type() == "string" && parent.Type() != "encapsed_string") ||
		(node.Type() == "encapsed_string" && patternLanguage.IsLeaf(node)) {
		return []string{"encapsed_string", "string"}
	}

	return []string{node.Type()}
}

func (*Pattern) TranslateContent(fromNodeType, toNodeType, content string) string {
	if fromNodeType == "string" && toNodeType == "encapsed_string" {
		return fmt.Sprintf(`"%s"`, content[1:len(content)-1])
	}
	if fromNodeType == "encapsed_string" && toNodeType == "string" {
		return fmt.Sprintf("'%s'", content[1:len(content)-1])
	}

	return content
}

func (*Pattern) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}
