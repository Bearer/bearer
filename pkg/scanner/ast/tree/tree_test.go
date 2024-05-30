package tree_test

import (
	"context"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
)

func parseTree(t *testing.T, content string) *tree.Tree {
	contentBytes := []byte(content)
	sitterLanguage := ruby.GetLanguage()

	sitterRootNode, err := sitter.ParseCtx(context.Background(), contentBytes, sitterLanguage)
	if err != nil {
		t.Fatalf("failed to parse input: %s", err)
	}

	return tree.NewBuilder(sitterLanguage, contentBytes, sitterRootNode, 0).Build()
}

func TestTree(t *testing.T) {
	tree := parseTree(t, `
		def m(a)
			a.foo
		end
	`)

	cupaloy.SnapshotT(t, tree.RootNode().Dump())
}
