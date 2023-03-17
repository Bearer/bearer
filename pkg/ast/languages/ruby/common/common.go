package common

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/pkg/util/set"
)

var (
	anonymousPatternNodeParentTypes = []string{"binary"}

	leafNodeTypes = []string{
		"identifier",
		"constant",
		"integer",
		"float",
		"complex",
		"rational",
		"string_content",
		"simple_symbol",
		"hash_key_symbol",
	}

	nodeTypeToFieldNames = map[string][]string{
		"method":               {"name", "parameters"},
		"singleton_method":     {"object", "name", "parameters"},
		"splat_parameter":      {"name"},
		"hash_splat_parameter": {"name"},
		"block_parameter":      {"name"},
		"keyword_parameter":    {"name", "value"},
		"optional_parameter":   {"name", "value"},
		"class":                {"name", "superclass"},
		"singleton_class":      {"value"},
		"module":               {"name"},
		"if_modifier":          {"body", "condition"},
		"unless_modifier":      {"body", "condition"},
		"while_modifier":       {"body", "condition"},
		"until_modifier":       {"body", "condition"},
		"rescue_modifier":      {"body", "handler"},
		"while":                {"condition", "body"},
		"until":                {"condition", "body"},
		"for":                  {"pattern", "value", "body"},
		"case":                 {"value"},
		"when":                 {"pattern", "body"},
		"if":                   {"condition", "consequence", "alternative"},
		"unless":               {"condition", "consequence", "alternative"},
		"elsif":                {"condition", "consequence", "alternative"},
		"rescue":               {"exceptions", "variable", "body"},
		"element_reference":    {"object"},
		"scope_resolution":     {"scope", "name"},
		"call":                 {"receiver", "method", "arguments", "block"},
		"do_block":             {"parameters"},
		"block":                {"parameters"},
		"assignment":           {"left", "right"},
		"operator_assignment":  {"left", "right", "operator"},
		"conditional":          {"condition", "consequence", "alternative"},
		"range":                {"begin", "end", "operator"},
		"binary":               {"left", "right", "operator"},
		"unary":                {"operator", "operand"},
		"setter":               {"name"},
		"alias":                {"name", "alias"},
		"pair":                 {"key", "value"},
		"lambda":               {"parameters", "body"},
	}

	specialCasedNotMissing = map[string][]string{
		"call": {"arguments"},
	}
)

func MatchContent(node *sitter.Node) bool {
	return slices.Contains(leafNodeTypes, node.Type())
}

// FIXME: move to patterns?
func GetMissingFields(node *sitter.Node) set.Set[string] {
	allFields := nodeTypeToFieldNames[node.Type()]
	notMissing := specialCasedNotMissing[node.Type()]

	missingFields := set.New[string]()
	for _, fieldName := range allFields {
		if node.ChildByFieldName(fieldName) == nil && !slices.Contains(notMissing, fieldName) {
			missingFields.Add(fieldName)
		}
	}

	return missingFields
}

func FieldName(node *sitter.Node) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	for _, fieldName := range nodeTypeToFieldNames[parent.Type()] {
		child := parent.ChildByFieldName(fieldName)
		if child != nil && child.Equal(node) {
			return fieldName
		}
	}

	return ""
}

func IncludeNode(node *sitter.Node) bool {
	parent := node.Parent()
	if parent != nil {
		return node.IsNamed() || slices.Contains(anonymousPatternNodeParentTypes, parent.Type())
	}

	return node.IsNamed()
}
