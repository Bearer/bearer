package ruby

import (
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/ssoroka/slice"

	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/language/types"
)

type rubyLanguage struct {
	language.Base
}

func Get() types.Language {
	return &rubyLanguage{
		Base: language.New(ruby.GetLanguage()),
	}
}

func (lang *rubyLanguage) Parse(input string) (*language.Tree, error) {
	tree, err := lang.Base.Parse(input)
	if err != nil {
		return nil, err
	}

	lang.analyzeFlow(tree.RootNode())

	return tree, nil
}

var variableLookupParents = []string{"pair", "argument_list"}

func (lang *rubyLanguage) analyzeFlow(rootNode *language.Node) {
	scope := make(map[string]*language.Node)

	rootNode.Walk(func(node *language.Node) error {
		switch node.Type() {
		case "method":
			scope = make(map[string]*language.Node)
		case "assignment":
			left := node.ChildByFieldName("left")
			right := node.ChildByFieldName("right")

			if left.Type() == "identifier" {
				scope[left.Content()] = node

				lang.UnifyNodes(node, right)
			}
		case "identifier":
			parent := node.Parent()
			if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
				scopedNode := scope[node.Content()]
				if scopedNode != nil {
					lang.UnifyNodes(node, scopedNode)
				}
			}
		}

		return nil
	})
}
