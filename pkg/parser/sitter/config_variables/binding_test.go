package config_variables_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/parser/sitter/config_variables"
	sitter "github.com/smacker/go-tree-sitter"
)

type NodeContent struct {
	Type    string
	Content string
}

func TestGrammar(t *testing.T) {
	input := []byte("Test {{ my.var }}${{ steps.dockerhub-check.outcome == 'success' }} 123 $VAR ${VAR2}")
	rootNode, err := sitter.ParseCtx(context.Background(), input, config_variables.GetLanguage())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(
		t,
		"(string (literal) (variable) (unknown) (literal) (variable) (literal) (variable))",
		rootNode.String(),
	)

	var childContents []NodeContent
	n := int(rootNode.ChildCount())
	for i := 0; i < n; i++ {
		child := rootNode.Child(i)
		if child.IsNamed() {
			childContents = append(childContents, NodeContent{Type: child.Type(), Content: child.Content(input)})
		}
	}

	expectedChildContents := []NodeContent{
		{Type: "literal", Content: "Test "},
		{Type: "variable", Content: "my.var"},
		{Type: "unknown", Content: "steps.dockerhub-check.outcome == 'success'"},
		{Type: "literal", Content: " 123 "},
		{Type: "variable", Content: "VAR"},
		{Type: "literal", Content: " "},
		{Type: "variable", Content: "VAR2"},
	}

	assert.Equal(t, expectedChildContents, childContents)
}
