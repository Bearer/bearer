package ruby

import (
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/ssoroka/slice"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/language/implementation"
	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/util/regex"
)

var (
	variableLookupParents = []string{"pair", "argument_list", "interpolation", "array", "binary"}

	anonymousPatternNodeParentTypes = []string{"binary"}
	patternMatchNodeContainerTypes  = []string{"argument_list", "keyword_parameter", "optional_parameter"}
	unanchoredPatternNodeTypes      = []string{"pair", "keyword_parameter"}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	allowedPatternQueryTypes  = []string{"identifier", "constant", "_", "call", "simple_symbol"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)

	ellipsisRegex = regexp.MustCompile(`\$<\.\.\.>`)

	passthroughMethods = []string{"JSON.parse", "JSON.parse!", "*.to_json"}
)

type rubyImplementation struct{}

func Get() implementation.Implementation {
	return &rubyImplementation{}
}

func (*rubyImplementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (*rubyImplementation) AnalyzeFlow(rootNode *tree.Node) error {
	scope := implementation.NewScope(nil)

	return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		switch node.Type() {
		case "method":
			scope = implementation.NewScope(nil)
		case "assignment":
			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if left.Type() == "identifier" {
				err := visitChildren()

				scope.Assign(left.Content(), node)
				node.UnifyWith(right)

				return err
			}
		case "identifier":
			parent := node.Parent()
			if parent == nil {
				break
			}

			if slice.Contains(variableLookupParents, parent.Type()) ||
				(parent.Type() == "assignment" && node.Equal(parent.ChildByFieldName("right"))) ||
				(parent.Type() == "call" && node.Equal(parent.ChildByFieldName("receiver"))) ||
				(parent.Type() == "element_reference" && node.Equal(parent.ChildByFieldName("object"))) {
				if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}

			if parent.Type() == "method_parameters" ||
				parent.Type() == "block_parameters" ||
				(parent.Type() == "keyword_parameter" && node.Equal(parent.ChildByFieldName("name"))) ||
				(parent.Type() == "optional_parameter" && node.Equal(parent.ChildByFieldName("name"))) {
				scope.Assign(node.Content(), node)
			}

			if parent.Type() == "argument_list" {
				callNode := parent.Parent()
				callNode.UnifyWith(node)
			}
		case "block", "do_block":
			previousScope := scope
			scope = implementation.NewScope(scope)
			err := visitChildren()
			scope = previousScope
			return err
		}

		return visitChildren()
	})
}

func (*rubyImplementation) ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error) {
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

func (*rubyImplementation) FindPatternMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (*rubyImplementation) FindPatternUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func produceDummyValue(i int, nodeType string) string {
	switch nodeType {
	case "identifier", "call":
		return "curioVar" + fmt.Sprint(i)
	case "simple_symbol":
		return ":curioVar" + fmt.Sprint(i)
	default:
		return "CurioVar" + fmt.Sprint(i)
	}
}

func (*rubyImplementation) PatternLeafContentTypes() []string {
	return []string{
		// identifiers
		"identifier", "constant",
		// datatypes/literals
		"number", "string_content", "integer", "float", "boolean", "nil", "simple_symbol", "hash_key_symbol",
	}
}

func (*rubyImplementation) AnonymousPatternNodeParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

func (*rubyImplementation) ShouldSkipNode(node *tree.Node) bool {
	return false
}

func (*rubyImplementation) PatternMatchNodeContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (*rubyImplementation) PatternIsAnchored(node *tree.Node) (bool, bool) {
	if slices.Contains(unanchoredPatternNodeTypes, node.Type()) {
		return false, false
	}

	parent := node.Parent()
	if parent == nil {
		return true, true
	}

	// Class body
	if parent.Type() == "class" {
		if node.Equal(parent.ChildByFieldName("name")) {
			return true, false
		}

		return false, false
	}

	// Block body
	if parent.Type() == "do_block" || parent.Type() == "block" {
		if node.Equal(parent.ChildByFieldName("parameters")) {
			return true, false
		}

		return false, false
	}

	// Method body
	if parent.Type() == "method" {
		if node.Equal(parent.ChildByFieldName("name")) || node.Equal(parent.ChildByFieldName("parameters")) {
			return true, false
		}

		return false, false
	}

	// Conditional body
	if parent.Type() == "then" {
		return false, false
	}

	if (parent.Type() == "if" || parent.Type() == "elsif" || parent.Type() == "unless") &&
		node.Equal(parent.ChildByFieldName("condition")) {
		return true, false
	}

	return true, true
}

func (*rubyImplementation) PatternNodeTypes(node *tree.Node) []string {
	parent := node.Parent()

	// Make these equivalent:
	//   key: value
	//   :key => value
	if parent != nil &&
		parent.Type() == "pair" &&
		node.Equal(parent.ChildByFieldName("key")) &&
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

func (*rubyImplementation) TranslatePatternContent(fromNodeType, toNodeType, content string) string {
	if fromNodeType == "hash_key_symbol" && toNodeType == "simple_symbol" {
		return ":" + content
	}

	if fromNodeType == "simple_symbol" && toNodeType == "hash_key_symbol" {
		return content[1:]
	}

	return content
}

func (implementation *rubyImplementation) IsRootOfRuleQuery(node *tree.Node) bool {
	return true
}

func (*rubyImplementation) PassthroughNested(node *tree.Node) bool {
	callNode := node.Parent()
	if callNode.Type() != "call" {
		return false
	}

	receiverNode := callNode.ChildByFieldName("receiver")

	if node.Type() != "arguments_list" && (receiverNode == nil || !node.Equal(receiverNode)) {
		return false
	}

	var receiverMethod string
	var wildcardMethod string

	if receiverNode != nil {
		methodName := callNode.ChildByFieldName("method").Content()

		if receiverNode.Type() == "identifier" {
			receiverMethod = receiverNode.Content() + "." + methodName
		}

		wildcardMethod = "*." + methodName
	}

	return slices.Contains(passthroughMethods, receiverMethod) || slices.Contains(passthroughMethods, wildcardMethod)
}
