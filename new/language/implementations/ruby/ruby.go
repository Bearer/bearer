package ruby

import (
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/ssoroka/slice"

	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/language/types"
)

var variableLookupParents = []string{"pair", "argument_list"}

type rubyLanguage struct {
	language.Base
}

func Get() types.Language {
	return &rubyLanguage{
		Base: language.New(ruby.GetLanguage(), analyzeFlow),
	}
}

func analyzeFlow(rootNode *language.Node) {
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
