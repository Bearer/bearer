package tree_test

import (
	"context"
	"testing"

	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bradleyjkemp/cupaloy"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

func parseTree(t *testing.T, content string) *tree.Tree {
	sitterLanguage := ruby.GetLanguage()

	sitterRootNode, err := sitter.ParseCtx(context.Background(), []byte(content), sitterLanguage)
	if err != nil {
		t.Fatalf("failed to parse input: %s", err)
	}

	return tree.NewBuilder(content, sitterRootNode).Build()
}

func TestTree(t *testing.T) {
	tree := parseTree(t, `
		def m(a)
			a.foo
		end
	`)

	cupaloy.SnapshotT(t, tree.RootNode().Dump())
}

func TestNodeAndDescendentIDs(t *testing.T) {
	tree := parseTree(t, `
		a.foo
		b.bar
	`)

	children := tree.RootNode().Children()

	cupaloy.SnapshotT(
		t,
		children[0].NodeAndDescendentIDs(),
		children[1].NodeAndDescendentIDs(),
	)
}
