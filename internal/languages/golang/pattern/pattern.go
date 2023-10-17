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
	patternMatchNodeContainerTypes = []string{"parameter_declaration", "parameter_list", "var_spec"}

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

func (*Pattern) LeafContentTypes() []string {
	return []string{
		"identifier",
		"package_identifier",
		"type_identifier",
		"raw_string_literal",
		"intepreted_string_literal",
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

	if parent.Type() == "function_declaration" {
		// parameters
		if node == parent.ChildByFieldName("parameters") {
			return true, false
		}

		return false, false
	}

	// function declaration_list
	unAnchored := []string{
		"function_declaration",
		"argument_list",
		"var_declaration",
	}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"source_file"}, node.Type()) && !node.IsMissing()
}

func (patternLanguage *Pattern) NodeTypes(node *tree.Node) []string {
	return []string{node.Type()}
}

func (*Pattern) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*Pattern) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}
