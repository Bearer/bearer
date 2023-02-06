package ruby

import (
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

var passThroughMethods = []string{
	"to_a",
	"to_array",
	"to_csv",
	"to_h",
	"to_hash",
	"to_json",
	"to_s",
}

var (
	variableLookupParents = []string{"pair", "argument_list", "interpolation"}

	anonymousPatternNodeParentTypes = []string{"binary"}
	patternMatchNodeContainerTypes  = []string{"argument_list"}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!\.]+)(?::(?P<types>[^>]+))?>`)
	allowedPatternQueryTypes  = []string{"identifier", "constant", "_", "call"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)

	ellipsisRegex = regexp.MustCompile(`\$<\.\.\.>`)
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
		case "method":
			scope = make(map[string]*tree.Node)
		case "assignment":
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
			if parent == nil {
				break
			}

			if slice.Contains(variableLookupParents, parent.Type()) ||
				(parent.Type() == "call" && node.Equal(parent.ChildByFieldName("receiver"))) {
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

func (implementation *rubyImplementation) FindPatternMatchNode(input []byte) [][]int {
	return matchNodeRegex.FindAllIndex(input, -1)
}

func (implementation *rubyImplementation) FindPatternUnanchoredPoints(input []byte) [][]int {
	return ellipsisRegex.FindAllIndex(input, -1)
}

func produceDummyValue(i int, nodeType string) string {
	switch nodeType {
	case "identifier", "call":
		return "curioVar" + fmt.Sprint(i)
	default:
		return "CurioVar" + fmt.Sprint(i)
	}
}

func (implementation *rubyImplementation) AnonymousPatternNodeParentTypes() []string {
	return anonymousPatternNodeParentTypes
}

func (implementation *rubyImplementation) PatternMatchNodeContainerTypes() []string {
	return patternMatchNodeContainerTypes
}

func (implementation *rubyImplementation) PatternIsAnchored(node *tree.Node) bool {
	if node.Type() == "pair" {
		return false
	}

	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Class body
	if parent.Type() == "class" && !node.Equal(parent.ChildByFieldName("name")) {
		return false
	}

	// Block body
	if (parent.Type() == "do_block" || parent.Type() == "block") && !node.Equal(parent.ChildByFieldName("parameters")) {
		return false
	}

	// Method body
	if parent.Type() == "method" && !node.Equal(parent.ChildByFieldName("name")) && !node.Equal(parent.ChildByFieldName("parameters")) {
		return false
	}

	return true
}

func (implementation *rubyImplementation) DescendIntoDetectionNode(node *tree.Node) bool {
	parent := node.Parent()

	if parent != nil && parent.Type() == "call" && node.Equal(parent.ChildByFieldName("receiver")) {
		return slice.Contains(passThroughMethods, parent.ChildByFieldName("method").Content())
	}

	return true
}

func (implementation *rubyImplementation) IsRootOfRuleQuery(node *tree.Node) bool {
	return true
}
