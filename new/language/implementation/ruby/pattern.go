package ruby

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

	anonymousPatternNodeParentTypes = []string{"binary"}
	patternMatchNodeContainerTypes  = []string{"argument_list", "keyword_parameter", "optional_parameter"}
	unanchoredPatternNodeTypes      = []string{"pair", "keyword_parameter"}
	allowedPatternQueryTypes        = []string{"identifier", "constant", "_", "call", "simple_symbol"}

	classPatternErrorRegex = regexp.MustCompile(`\Aclass\s*\z`)
)

type patternImplementation struct {
	implementation.PatternBase
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

func (*patternImplementation) FindMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (*patternImplementation) FindUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func produceDummyValue(i int, nodeType string) string {
	return "curioVar" + fmt.Sprint(i)
}

func (*patternImplementation) LeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "constant",
		// datatypes/literals
		"number", "string_content", "integer", "float", "boolean", "nil", "simple_symbol", "hash_key_symbol",
	}
}

func (*patternImplementation) AnonymousParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

func (*patternImplementation) ContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (*patternImplementation) IsAnchored(node *tree.Node) (bool, bool) {
	if slices.Contains(unanchoredPatternNodeTypes, node.Type()) {
		return false, false
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Class body
	if parent.Type() == "class" {
		if node == parent.ChildByFieldName("name") {
			return true, false
		}

		return false, false
	}

	// Block body
	if parent.Type() == "do_block" || parent.Type() == "block" {
		if node == parent.ChildByFieldName("parameters") {
			return true, false
		}

		return false, false
	}

	// Method body
	if parent.Type() == "method" {
		if node == parent.ChildByFieldName("name") || node == parent.ChildByFieldName("parameters") {
			return true, false
		}

		return false, false
	}

	// Conditional body
	if parent.Type() == "then" {
		return false, false
	}

	if (parent.Type() == "if" || parent.Type() == "elsif" || parent.Type() == "unless") &&
		node == parent.ChildByFieldName("condition") {
		return true, false
	}

	return true, true
}

func (*patternImplementation) NodeTypes(node *tree.Node) []string {
	parent := node.Parent()

	// Make these equivalent:
	//   key: value
	//   :key => value
	if parent != nil &&
		parent.Type() == "pair" &&
		node == parent.ChildByFieldName("key") &&
		(node.Type() == "hash_key_symbol" || node.Type() == "simple_symbol") {
		return []string{"hash_key_symbol", "simple_symbol"}
	}

	// Make these equivalent:
	//  call do ... end
	//  call { ... }
	if node.Type() == "block" || node.Type() == "do_block" {
		return []string{"block", "do_block"}
	}

	return []string{node.Type()}
}

func (*patternImplementation) TranslateContent(fromNodeType, toNodeType, content string) string {
	if fromNodeType == "hash_key_symbol" && toNodeType == "simple_symbol" {
		return ":" + content
	}

	if fromNodeType == "simple_symbol" && toNodeType == "hash_key_symbol" {
		return content[1:]
	}

	return content
}

func (*patternImplementation) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	for ancestor := node.Parent(); ancestor != nil; ancestor = ancestor.Parent() {
		if ancestor.Type() != "ERROR" {
			continue
		}

		errorPrefix := input[ancestor.ContentStart.Byte:node.ContentStart.Byte]
		if classPatternErrorRegex.Match(errorPrefix) {
			return strings.ToUpper(string(dummyValue[0])) + dummyValue[1:]
		}
	}

	return dummyValue
}
