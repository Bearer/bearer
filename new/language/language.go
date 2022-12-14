package language

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

type Base struct {
	SitterLanguage *sitter.Language
}

func (lang *Base) Parse(input string) (*Tree, error) {
	inputBytes := []byte(input)

	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(lang.SitterLanguage)

	sitterTree, err := parser.ParseCtx(context.Background(), nil, inputBytes)
	if err != nil {
		return nil, err
	}

	return &Tree{
		input:      inputBytes,
		sitterTree: sitterTree,
	}, nil
}

func (lang *Base) CompileQuery(input string) (*Query, error) {
	sitterQuery, err := sitter.NewQuery([]byte(input), lang.SitterLanguage)
	if err != nil {
		return nil, err
	}

	return &Query{sitterQuery: sitterQuery}, nil
}
