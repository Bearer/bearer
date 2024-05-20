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
	patternQueryVariableRegex      = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex                 = regexp.MustCompile(`\$<!>`)
	ellipsisRegex                  = regexp.MustCompile(`\$<\.\.\.>`)
	unanchoredPatternNodeTypes     = []string{"function_definition"}
	patternMatchNodeContainerTypes = []string{
		"dotted_name",
		"typed_parameter",
		"typed_default_parameter",
		"default_parameter",
		"decorated_definition",
		"module",
	}

	allowedPatternQueryTypes = []string{"_"}
)

type Pattern struct {
	language.PatternBase
}

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
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
	return false
}

func (*Pattern) LeafContentTypes() []string {
	return []string{
		"string",
		"identifier",
		"true",
		"false",
		"float",
		"integer",
		"none",
	}
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
	if slices.Contains(unanchoredPatternNodeTypes, node.Type()) {
		return false, false
	}

	// return type or leading comment
	if node.Type() == "block" {
		return false, true
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	if parent.Type() == "method_declaration" {
		// visibility
		if node == parent.ChildByFieldName("name") {
			return false, true
		}

		// type
		if node == parent.ChildByFieldName("parameters") {
			return true, false
		}

		return false, false
	}

	// type parameters
	if parent.Type() == "function_definition" && node == parent.ChildByFieldName("parameters") {
		return false, true
	}

	// inherited types
	if parent.Type() == "argument_list" && parent.Parent() != nil && parent.Parent().Type() == "class_definition" {
		return false, false
	}

	// Associative array elements are unanchored
	// eg. array("foo" => 42)
	if parent.Type() == "array_creation_expression" &&
		node.Type() == "array_element_initializer" &&
		len(node.NamedChildren()) == 2 {
		return false, false
	}

	if (parent.Type() == "import_statement" || parent.Type() == "import_from_statement" || parent.Type() == "relative_import") &&
		(node.Type() == "dotted_name" || node.Type() == "aliased_import") {
		return false, false
	}

	// Class body declaration_list
	// function/block compound_statement
	unAnchored := []string{}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"module", "expression_statement"}, node.Type()) && !node.IsMissing()
}

func (patternLanguage *Pattern) NodeTypes(node *tree.Node, parentType string) []string {
	if node.Type() == "typed_parameter" {
		return []string{"typed_parameter", "typed_default_parameter"}
	}

	return []string{node.Type()}
}

func (*Pattern) IsContainer(node *tree.Node) bool {
	return slices.Contains(patternMatchNodeContainerTypes, node.Type())
}
