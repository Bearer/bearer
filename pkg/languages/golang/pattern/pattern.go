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
	unanchoredPatternNodeTypes     = []string{"import_spec"}
	patternMatchNodeContainerTypes = []string{
		"range_clause",
		"parameter_declaration",
		"argument_list",
		"expression_list",
		"parameter_list",
		"var_spec",
		"import_spec",
		"literal_element", // Can be removed once the tree-sitter-go is updated
	}

	allowedPatternQueryTypes = []string{"_"}
)

type Pattern struct {
	language.PatternBase
}

func (*Pattern) AdjustInput(input string) string {
	return input + "\n"
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
	case "raw_string_literal", "interpreted_string_literal":
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
		"package_identifier",
		"type_identifier",
		"field_identifier",
		"raw_string_literal",
		"interpreted_string_literal",
		"int_literal",
		"float_literal",
		"true",
		"false",
		"nil",
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

	if parent.Type() == "import_spec" {
		if node == parent.ChildByFieldName("path") {
			return false, true
		}
	}

	if slices.Contains([]string{"function_declaration", "method_declaration"}, parent.Type()) {
		// parameters
		if node == parent.ChildByFieldName("parameters") {
			return true, false
		}

		return false, false
	}

	// function declaration_list
	unAnchored := []string{
		"function_declaration",
		"method_declaration",
		"var_declaration",
		"literal_value",
	}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"source_file", "expression_statement"}, node.Type()) && !node.IsMissing()
}

func (patternLanguage *Pattern) NodeTypes(node *tree.Node, parentType string) []string {
	if node.Type() == "identifier" && node.Parent().Type() == "source_file" {
		return []string{"identifier", "package_identifier"}
	}

	return []string{node.Type()}
}

func (*Pattern) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*Pattern) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}
