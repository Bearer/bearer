package javascript

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/util/regex"
)

var (
	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex            = regexp.MustCompile(`\$<!>`)
	ellipsisRegex             = regexp.MustCompile(`\$<\.\.\.>`)

	patternMatchNodeContainerTypes = []string{"import_clause", "import_specifier", "required_parameter"}

	allowedPatternQueryTypes = []string{"identifier", "property_identifier", "_", "member_expression", "string", "template_string"}
)

type patternImplementation struct {
	implementation.PatternBase
}

func (*patternImplementation) IsLeaf(node *tree.Node) bool {
	return node.Type() == "string"
}

func (*patternImplementation) ExtractVariables(input string) (string, []patternquerytypes.Variable, error) {
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

func (*patternImplementation) FindMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (*patternImplementation) FindUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func (*patternImplementation) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (*patternImplementation) LeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "property_identifier", "shorthand_property_identifier", "type_identifier",
		// datatypes/literals
		"template_string", "string_fragment", "number", "null", "true", "false",
	}
}

func (*patternImplementation) IsAnchored(node *tree.Node) (bool, bool) {
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

func (*patternImplementation) IsRoot(node *tree.Node) bool {
	return !(node.Type() == "expression_statement")
}

func (*patternImplementation) NodeTypes(node *tree.Node) []string {
	if node.Type() == "statement_block" && node.Parent().Type() == "program" {
		if len(node.NamedChildren()) == 0 {
			return []string{"object"}
		} else {
			return []string{node.Type(), "program"}
		}
	}

	return []string{node.Type()}
}

func (*patternImplementation) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	parent := node.Parent()
	if parent == nil {
		return dummyValue
	}

	if parent.NamedChildren()[0].Type() == "import_clause" {
		return "\"" + dummyValue + "\""
	}

	return dummyValue
}
