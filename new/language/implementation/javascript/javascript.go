package javascript

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/ssoroka/slice"
	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/language/implementation"
	patternquerytypes "github.com/bearer/curio/new/language/patternquery/types"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/pkg/util/regex"
)

var (
	variableLookupParents = []string{"pair", "arguments", "binary_expression", "template_string"}

	// anonymousPatternNodeParentTypes = []string{"binary"}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!]+)(?::(?P<types>[^>]+))?>`)
	allowedPatternQueryTypes  = []string{"identifier", "property_identifier", "_", "member_expression"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)
)

type rubyImplementation struct{}

func Get() implementation.Implementation {
	return &rubyImplementation{}
}

func (implementation *rubyImplementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (implementation *rubyImplementation) AnalyzeFlow(rootNode *tree.Node) error {
	scope := make(map[string]*tree.Node)

	return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		switch node.Type() {
		case "function":
			scope = make(map[string]*tree.Node)
		case "arrow_function":
			scope = make(map[string]*tree.Node)
		case "method_definition":
			scope = make(map[string]*tree.Node)
		case "assignment_expression":
			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if left.Type() == "identifier" {
				err := visitChildren()

				scope[left.Content()] = node
				node.UnifyWith(right)

				return err
			}
		case "identifier":
			parent := node.Parent()
			if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
				scopedNode := scope[node.Content()]
				if scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}
		case "property_identifier":
			parent := node.Parent()
			if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
				scopedNode := scope[node.Content()]
				if scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}
		}

		return visitChildren()
	})
}

func (implementation *rubyImplementation) ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error) {
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

func (implementation *rubyImplementation) ExtractPatternMatchNode(input string) (string, int, error) {
	inputBytes := []byte(input)
	matches := matchNodeRegex.FindAllIndex(inputBytes, -1)

	if len(matches) == 0 {
		return input, 0, nil
	}

	if len(matches) > 1 {
		return "", 0, errors.New("pattern must only contain a single match node")
	}

	match := matches[0]
	return string(inputBytes[0:match[0]]) + string(inputBytes[match[1]:]), match[0], nil
}

func produceDummyValue(i int, nodeType string) string {
	// TODO: See if anything needs to be added here
	switch nodeType {
	case "identifier", "call":
		return "curioVar" + fmt.Sprint(i)
	default:
		return "CurioVar" + fmt.Sprint(i)
	}
}

func (implementation *rubyImplementation) AnonymousPatternNodeParentTypes() []string {
	return []string{}
}

func (implementation *rubyImplementation) PatternIsAnchored(node *tree.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return true
	}

	// // Class body
	// if parent.Type() == "class" && !node.Equal(parent.ChildByFieldName("name")) {
	// 	return false
	// }

	// // Block body
	// if (parent.Type() == "do_block" || parent.Type() == "block") && !node.Equal(parent.ChildByFieldName("parameters")) {
	// 	return false
	// }

	// // Method body
	// if parent.Type() == "method" && !node.Equal(parent.ChildByFieldName("name")) && !node.Equal(parent.ChildByFieldName("parameters")) {
	// 	return false
	// }

	return true
}
