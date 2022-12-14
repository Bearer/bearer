package ruby

import (
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/language/types"
)

type rubyLanguage struct {
	language.Base
}

func Get() types.Language {
	return &rubyLanguage{
		Base: language.Base{SitterLanguage: ruby.GetLanguage()},
	}
}
