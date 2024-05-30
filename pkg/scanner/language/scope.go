package language

import sitter "github.com/smacker/go-tree-sitter"

type Scope struct {
	parent    *Scope
	variables map[string]*sitter.Node
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]*sitter.Node),
	}
}

func (scope *Scope) Declare(name string, node *sitter.Node) {
	scope.variables[name] = node
}

func (scope *Scope) Assign(name string, node *sitter.Node) {
	targetScope := scope
	if _, declarationScope := scope.lookupWithScope(name); declarationScope != nil {
		targetScope = declarationScope
	}

	targetScope.variables[name] = node
}

func (scope *Scope) Lookup(name string) *sitter.Node {
	node, _ := scope.lookupWithScope(name)
	return node
}

func (scope *Scope) lookupWithScope(name string) (*sitter.Node, *Scope) {
	if node, exists := scope.variables[name]; exists {
		return node, scope
	}

	if scope.parent != nil {
		return scope.parent.lookupWithScope(name)
	}

	return nil, nil
}
