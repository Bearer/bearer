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
	unanchoredPatternNodeTypes     = []string{}
	patternMatchNodeContainerTypes = []string{"formal_parameters", "simple_parameter", "argument"}

	allowedPatternQueryTypes = []string{"_"}
)

type Pattern struct {
	language.PatternBase
}

// func (*Pattern) AdjustInput(input string) string {
// 	return input
// }

// func (*Pattern) FixupMissing(node *tree.Node) string {
// 	if node.Type() != `";"` {
// 		return ""
// 	}

// 	return ";"
// }

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	if slices.Contains([]string{"named_type"}, node.Parent().Type()) {
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

func (*Pattern) LeafContentTypes() []string {
	return []string{
		"encapsed_string",
		"string",
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

	// Associative array elements are unanchored
	// eg. array("foo" => 42)
	if parent.Type() == "array_creation_expression" &&
		node.Type() == "array_element_initializer" &&
		len(node.NamedChildren()) == 2 {
		return false, false
	}

	// Class body declaration_list
	// function/block compound_statement
	unAnchored := []string{
		"declaration_list",
		"compound_statement",
	}

	isUnanchored := !slices.Contains(unAnchored, parent.Type())
	return isUnanchored, isUnanchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"expression_statement", "program"}, node.Type()) && !node.IsMissing()
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
