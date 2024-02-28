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
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex            = regexp.MustCompile(`\$<!>`)
	ellipsisRegex             = regexp.MustCompile(`\$<\.\.\.>`)

	patternMatchNodeContainerTypes = []string{"import_clause", "import_specifier", "required_parameter"}

	allowedPatternQueryTypes = []string{"identifier", "property_identifier", "_", "member_expression", "string", "template_string"}
)

type Pattern struct {
	language.PatternBase
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

func (*Pattern) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (*Pattern) LeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "property_identifier", "shorthand_property_identifier", "type_identifier",
		// datatypes/literals
		"template_string", "string_fragment", "number", "null", "true", "false",
	}
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
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
	unAnchored := []string{"statement_block", "class_body", "object_pattern", "named_imports", "method_definition"}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"expression_statement", "program"}, node.Type())
}

func (*Pattern) NodeTypes(node *tree.Node) []string {
	if node.Type() == "statement_block" && node.Parent().Type() == "program" {
		if len(node.NamedChildren()) == 0 {
			return []string{"object"}
		} else {
			return []string{node.Type(), "program"}
		}
	}

	return []string{node.Type()}
}

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	parent := node.Parent()
	if parent == nil {
		return dummyValue
	}

	if parent.NamedChildren()[0].Type() == "import_clause" {
		return "\"" + dummyValue + "\""
	}

	return dummyValue
}
