package tree_test

import (
	"context"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func parseTree(t *testing.T, sitterLanguage *sitter.Language, content string) *tree.Tree {
	contentBytes := []byte(content)

	sitterRootNode, err := sitter.ParseCtx(context.Background(), contentBytes, sitterLanguage)
	if err != nil {
		t.Fatalf("failed to parse input: %s", err)
	}

	return tree.NewBuilder(sitterLanguage, contentBytes, sitterRootNode, 0).Build()
}

func TestTree(t *testing.T) {
	tree := parseTree(t, ruby.GetLanguage(), `
		def m(a)
			a.foo
		end
	`)

	cupaloy.SnapshotT(t, tree.RootNode().Dump())
}

func TestContentParts(t *testing.T) {
	for _, test := range []struct{ expression, expected string }{
		{"`abc`", "abc"},
		{"`a${b}c`", "a*c"},
		{"`${b}c`", "*c"},
		{"`a${b}`", "a*"},
	} {
		t.Run(test.expression, func(tt *testing.T) {
			ast := parseTree(tt, javascript.GetLanguage(), test.expression)
			stringNode := ast.RootNode().NamedChildren()[0].NamedChildren()[0]
			assert.Equal(tt, "template_string", stringNode.Type())

			var result string
			err := stringNode.EachContentPart(func(text string) error {
				result += text
				return nil
			}, func(child *tree.Node) error {
				result += "*"
				return nil
			})
			assert.NoError(tt, err)

			assert.Equal(tt, test.expected, result)
		})
	}
}
