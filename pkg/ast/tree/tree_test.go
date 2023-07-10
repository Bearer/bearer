package tree_test

import (
	"context"
	"testing"

	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bradleyjkemp/cupaloy"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

func parseTree(t *testing.T, content string) *tree.Node {
	sitterLanguage := ruby.GetLanguage()

	sitterRootNode, err := sitter.ParseCtx(context.Background(), []byte(content), sitterLanguage)
	if err != nil {
		t.Fatalf("failed to parse input: %s", err)
	}

	return tree.NewBuilder(sitterRootNode).Build()
}

func TestTree(t *testing.T) {
	rootNode := parseTree(t, `
		def m(a)
			a.foo
		end
	`)

	cupaloy.SnapshotT(t, rootNode.Dump())
}

func TestNodeAndDescendentIDs(t *testing.T) {
	rootNode := parseTree(t, `
		a.foo
		b.bar
	`)

	cupaloy.SnapshotT(
		t,
		rootNode.Child(0).NodeAndDescendentIDs(),
		rootNode.Child(1).NodeAndDescendentIDs(),
	)
}
