package language

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/curio/new/parser"
)

type Language struct {
	sitterLanguage *sitter.Language
}

func Get(name string) (*Language, error) {
	sitterLanguage, err := getSitterLanguage(name)
	if err != nil {
		return nil, err
	}

	return &Language{
		sitterLanguage: sitterLanguage,
	}, nil
}

func getSitterLanguage(name string) (*sitter.Language, error) {
	switch name {
	case "ruby":
		return ruby.GetLanguage(), nil
	default:
		return nil, fmt.Errorf("unsupported language '%s'", name)
	}
}

func (language *Language) Parse(input string) (*parser.Tree, error) {
	return parser.Parse(language.sitterLanguage, input)
}

func (language *Language) CompileQuery(input string) (*parser.Query, error) {
	return parser.CompileQuery(language.sitterLanguage, input)
}
