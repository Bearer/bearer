package pattern

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/util/regex"
)

var (
	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex  = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex             = regexp.MustCompile(`\$<!>`)
	ellipsisRegex              = regexp.MustCompile(`\$<\.\.\.>`)
	unanchoredPatternNodeTypes = []string{
		"function_item",
		"impl_item",
		"struct_item",
		"enum_item",
		"trait_item",
		"mod_item",
	}
	patternMatchNodeContainerTypes = []string{
		"scoped_identifier",
		"field_expression",
		"call_expression",
		"arguments",
		"parameters",
		"use_declaration",
		"source_file",
	}

	allowedPatternQueryTypes = []string{"_"}
)

type Pattern struct {
	language.PatternBase
}

func (*Pattern) AdjustInput(input string) string {
	// Rust requires semicolons for expression statements, but patterns
	// often don't include them. Add a semicolon if the pattern doesn't
	// already end with one (or with braces for blocks).
	trimmed := strings.TrimRight(input, " \t\n\r")
	if len(trimmed) > 0 {
		lastChar := trimmed[len(trimmed)-1]
		if lastChar != ';' && lastChar != '}' && lastChar != '{' {
			return input + ";\n"
		}
	}
	return input + "\n"
}

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	return dummyValue
}

func (*Pattern) FixupMissing(node *tree.Node) string {
	// Rust requires semicolons to terminate expression statements
	if node.Type() == `";"` {
		return ";"
	}

	return ""
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
	switch node.Type() {
	case "string_literal", "raw_string_literal":
		return true
	}
	return false
}

func (*Pattern) AnonymousParentTypes() []string {
	return []string{
		"unary_expression",
		"binary_expression",
	}
}

func (*Pattern) LeafContentTypes() []string {
	return []string{
		"identifier",
		"field_identifier",
		"type_identifier",
		"string_literal",
		"raw_string_literal",
		"char_literal",
		"integer_literal",
		"float_literal",
		"boolean_literal",
		"self",
		"crate",
		"super",
	}
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
	if slices.Contains(unanchoredPatternNodeTypes, node.Type()) {
		return false, false
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Function parameters
	if parent.Type() == "function_item" {
		if node == parent.ChildByFieldName("parameters") {
			return true, false
		}
		return false, false
	}

	// Impl items
	if parent.Type() == "impl_item" {
		return false, false
	}

	// Match arms are unanchored
	if parent.Type() == "match_block" {
		return false, false
	}

	// Use declarations
	if parent.Type() == "use_declaration" || parent.Type() == "use_list" || parent.Type() == "scoped_use_list" {
		return false, false
	}

	// Struct field declarations
	if parent.Type() == "field_declaration_list" {
		return false, false
	}

	unAnchored := []string{
		"function_item",
		"impl_item",
		"struct_item",
		"enum_item",
		"trait_item",
		"declaration_list",
	}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"source_file", "expression_statement"}, node.Type()) && !node.IsMissing()
}

func (patternLanguage *Pattern) NodeTypes(node *tree.Node) []string {
	return []string{node.Type()}
}

func (*Pattern) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*Pattern) IsContainer(node *tree.Node) bool {
	return slices.Contains(patternMatchNodeContainerTypes, node.Type())
}

