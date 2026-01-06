package pattern

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/util/regex"

	"github.com/bearer/bearer/pkg/scanner/language"
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

func (*Pattern) FindMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (*Pattern) FindUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func produceDummyValue(i int, nodeType string) string {
	return "bearerVar" + fmt.Sprint(i)
}

func (*Pattern) LeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "constant",
		// datatypes/literals
		"number", "string_content", "integer", "float", "boolean", "nil", "simple_symbol", "hash_key_symbol",
	}
}

func (*Pattern) AnonymousParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

func (*Pattern) IsContainer(node *tree.Node) bool {
	return slices.Contains(patternMatchNodeContainerTypes, node.Type())
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
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

	// Class/module body - don't anchor methods because there may be
	// visibility modifiers (private, protected, public) as siblings
	if parent.Type() == "body_statement" {
		grandparent := parent.Parent()
		if grandparent != nil && (grandparent.Type() == "class" || grandparent.Type() == "module" || grandparent.Type() == "singleton_class") {
			return false, false
		}
	}

	return true, true
}

func (*Pattern) NodeTypes(node *tree.Node, parentType string) []string {
	parent := node.Parent()

	// Make these equivalent:
	//   key: value
	//   :key => value
	if parentType == "pair" &&
		node == parent.ChildByFieldName("key") &&
		(node.Type() == "hash_key_symbol" || node.Type() == "simple_symbol") {
		return []string{"hash_key_symbol", "simple_symbol"}
	}

	// Make these equivalent:
	//  call do ... end
	//  call { ... }
	blockTypes := []string{"block", "do_block"}
	if slices.Contains(blockTypes, node.Type()) {
		return blockTypes
	}

	// The block types use different bodies. This is to cope with matching both
	// block types as equivalent
	if parentType == "block" && node.Type() == "body_statement" {
		return []string{"block_body"}
	}
	if parentType == "do_block" && node.Type() == "block_body" {
		return []string{"body_statement"}
	}

	return []string{node.Type()}
}

func (*Pattern) TranslateContent(fromNodeType, toNodeType, content string) string {
	if fromNodeType == "hash_key_symbol" && toNodeType == "simple_symbol" {
		return ":" + content
	}

	if fromNodeType == "simple_symbol" && toNodeType == "hash_key_symbol" {
		return content[1:]
	}

	return content
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"program"}, node.Type())
}

func (*Pattern) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
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
