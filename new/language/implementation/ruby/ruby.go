package ruby

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
	patternquerybuilder "github.com/bearer/curio/new/language/patternquery/builder"
	"github.com/bearer/curio/new/language/tree"
	"github.com/bearer/curio/pkg/util/regex"
)

var (
	variableLookupParents = []string{"pair", "argument_list", "interpolation"}

	anonymousPatternNodeParentTypes = []string{"binary"}

	// $<name:type> or $<name:type1|type2> or $<name>
	patternQueryVariableRegex = regexp.MustCompile(`\$<(?P<name>[^>:!]+)(?::(?P<types>[^>]+))?>`)
	allowedPatternQueryTypes  = []string{"identifier", "constant", "_", "call"}

	matchNodeRegex = regexp.MustCompile(`\$<!>`)
)

type rubyImplementation struct{}

func Get() implementation.Implementation {
	return &rubyImplementation{}
}

func (implementation *rubyImplementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (implementation *rubyImplementation) AnalyzeFlow(rootNode *tree.Node) {
	scope := make(map[string]*tree.Node)

	rootNode.Walk(func(node *tree.Node) error {
		switch node.Type() {
		case "method":
			scope = make(map[string]*tree.Node)
		case "assignment":
			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if left.Type() == "identifier" {
				scope[left.Content()] = node

				node.UnifyWith(right)
			}
		case "identifier":
			parent := node.Parent()
			if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
				scopedNode := scope[node.Content()]
				if scopedNode != nil {
					node.UnifyWith(scopedNode)
				}
			}
		}

		return nil
	})
}

func (implementation *rubyImplementation) ExtractPatternVariables(input string) (string, []patternquerybuilder.Variable, error) {
	nameIndex := patternQueryVariableRegex.SubexpIndex("name")
	typesIndex := patternQueryVariableRegex.SubexpIndex("types")
	i := 0

	var params []patternquerybuilder.Variable

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

		params = append(params, patternquerybuilder.Variable{
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
