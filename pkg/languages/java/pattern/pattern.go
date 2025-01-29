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
	queryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	matchNodeRegex     = regexp.MustCompile(`\$<!>`)
	ellipsisRegex      = regexp.MustCompile(`\$<\.\.\.>`)

	matchNodeContainerTypes = []string{
		"catch_formal_parameter",
		"catch_type",
		"formal_parameters",
		"modifiers",
		"method_declaration",
		"program",
	}

	// todo: see if it is ok to replace typescripts `member_expression` with javas `field_access` and `method_invocation`
	allowedQueryTypes = []string{"identifier", "type_identifier", "_", "field_access", "method_invocation", "string_literal"}
)

type Pattern struct {
	language.PatternBase
}

func (*Pattern) ExtractVariables(input string) (string, []language.PatternVariable, error) {
	nameIndex := queryVariableRegex.SubexpIndex("name")
	typesIndex := queryVariableRegex.SubexpIndex("types")
	i := 0

	var params []language.PatternVariable

	replaced, err := regex.ReplaceAllWithSubmatches(queryVariableRegex, input, func(submatches []string) (string, error) {
		nodeTypes := strings.Split(submatches[typesIndex], "|")
		if nodeTypes[0] == "" {
			nodeTypes = []string{"_"}
		}

		for _, nodeType := range nodeTypes {
			if !slices.Contains(allowedQueryTypes, nodeType) {
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

func (*Pattern) LeafContentTypes() []string {
	return []string{
		// todo: see if type identifier should be removed from here (User user) `User` is type
		// identifiers
		"identifier", "modifier",
		// types
		// int user, User user, void user function,
		"integral_type", "type_identifier", "void_type",
		// datatypes/literals
		"string_literal", "character_literal", "null_literal", "true", "false", "decimal_integer_literal", "decimal_floating_point_literal",
	}
}

func (*Pattern) IsAnchored(node *tree.Node) (bool, bool) {
	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Class body class_body
	// function block
	// lambda () -> {} block
	// try {} catch () {}
	unAnchored := []string{
		"class_declaration",
		"class_body",
		"block",
		"try_statement",
		"catch_type",
		"resource_specification",
		"modifiers",
	}

	isAnchored := !slices.Contains(unAnchored, parent.Type())
	return isAnchored, isAnchored
}

func (*Pattern) IsRoot(node *tree.Node) bool {
	return !slices.Contains([]string{"expression_statement", "program"}, node.Type())
}

func (*Pattern) FixupMissing(node *tree.Node) string {
	if node.Type() != `";"` {
		return ""
	}

	return ";"
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

func (*Pattern) IsContainer(node *tree.Node) bool {
	if slices.Contains(matchNodeContainerTypes, node.Type()) {
		return true
	}

	if node.Type() == "class_declaration" {
		if children := node.NamedChildren(); len(children) != 0 && children[0].Type() == "modifiers" {
			return true
		}
	}

	return false
}
