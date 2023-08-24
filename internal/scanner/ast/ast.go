package ast

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/language"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
)

func Parse(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
) (*tree.Builder, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language.SitterLanguage())

	sitterTree, err := parser.ParseCtx(ctx, nil, contentBytes)
	if err != nil {
		return nil, err
	}

	return tree.NewBuilder(contentBytes, sitterTree.RootNode()), nil
}

func ParseAndBuild(
	ctx context.Context,
	language language.Language,
	contentBytes []byte,
) (*tree.Tree, error) {
	builder, err := Parse(ctx, language, contentBytes)
	if err != nil {
		return nil, err
	}

	return builder.Build(), nil
}
