package analyzer

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/ast/tree"
)

type analyzer struct {
	builder *tree.Builder
	scope   *implementation.Scope
}

func New(builder *tree.Builder) implementation.Analyzer {
	return &analyzer{
		builder: builder,
		scope:   implementation.NewScope(nil),
	}
}

func (analyzer *analyzer) Analyze(node *sitter.Node, visitChildren func() error) error {
	switch node.Type() {
	// () => {}
	// function getName() {}
	case "function", "arrow_function", "method_definition":
		return analyzer.withScope(implementation.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "assignment_expression":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "augmented_assignment_expression":
		return analyzer.analyzeAugmentedAssignment(node, visitChildren)
	case "variable_declarator":
		return analyzer.analyzeVariableDeclarator(node, visitChildren)
	case "shorthand_property_identifier_pattern":
		return analyzer.analyzeShorthandPropertyIdentifierPattern(node, visitChildren)
	case "member_expression":
		return analyzer.analyzeMember(node, visitChildren)
	case "subscript_expression":
		return analyzer.analyzeSubscript(node, visitChildren)
	case "new_expression":
		return analyzer.analyzeNew(node, visitChildren)
	case "call_expression":
		return analyzer.analyzeCall(node, visitChildren)
	case "required_parameter", "optional_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	case "import_clause":
		return analyzer.analyzeImportClause(node, visitChildren)
	case "namespace_import":
		return analyzer.analyzeNamespaceImport(node, visitChildren)
	case "import_specifier":
		return analyzer.analyzeImportSpecifier(node, visitChildren)
	case "ternary_expression":
		return analyzer.analyzeTernary(node, visitChildren)
	case "parenthesized_expression":
		return analyzer.analyzeParentheses(node, visitChildren)
	case "arguments",
		"array",
		"binary_expression",
		"pair",
		"spread_element",
		"template_substitution",
		"unary_expression":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	default:
		// statements don't have results
		if !strings.HasSuffix(node.Type(), "_statement") {
			analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		}

		return visitChildren()
	}
}

// user = ...
func (analyzer *analyzer) analyzeAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	analyzer.builder.Alias(node, right)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "identifier" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

// x += y
func (analyzer *analyzer) analyzeAugmentedAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	analyzer.builder.Dataflow(node, left, right)
	analyzer.lookupVariable(left)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "identifier" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

// const user = ...
// var user = ...
// let user = ...
func (analyzer *analyzer) analyzeVariableDeclarator(node *sitter.Node, visitChildren func() error) error {
	name := node.ChildByFieldName("name")
	value := node.ChildByFieldName("value")
	analyzer.builder.Alias(node, value)
	analyzer.lookupVariable(value)

	err := visitChildren()

	if name.Type() == "identifier" {
		analyzer.scope.Declare(analyzer.builder.ContentFor(name), node)
	}

	return err
}

// const { foo } = ...
func (analyzer *analyzer) analyzeShorthandPropertyIdentifierPattern(
	node *sitter.Node,
	visitChildren func() error,
) error {
	analyzer.scope.Declare(analyzer.builder.ContentFor(node), node)

	return visitChildren()
}

// foo.bar
func (analyzer *analyzer) analyzeMember(node *sitter.Node, visitChildren func() error) error {
	object := node.ChildByFieldName("object")
	analyzer.builder.Dataflow(node, object)
	analyzer.lookupVariable(object)

	return visitChildren()
}

// foo["bar"]
func (analyzer *analyzer) analyzeSubscript(node *sitter.Node, visitChildren func() error) error {
	object := node.ChildByFieldName("object")
	analyzer.builder.Dataflow(node, object)
	analyzer.lookupVariable(object)

	return visitChildren()
}

// new Foo()
func (analyzer *analyzer) analyzeNew(node *sitter.Node, visitChildren func() error) error {
	constructor := node.ChildByFieldName("constructor")
	analyzer.lookupVariable(constructor)

	return visitChildren()
}

// foo.bar(1, 2)
func (analyzer *analyzer) analyzeCall(node *sitter.Node, visitChildren func() error) error {
	if arguments := node.ChildByFieldName("arguments"); arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// parameter definition
// foo(a, b = 1)
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	pattern := node.ChildByFieldName("pattern")
	if pattern.Type() == "identifier" {
		analyzer.scope.Declare(analyzer.builder.ContentFor(pattern), node)
	}

	if value := node.ChildByFieldName("value"); value != nil {
		analyzer.lookupVariable(value)
		analyzer.builder.Alias(node, value)
	}

	return visitChildren()
}

// parts between "import" and "from":
// import a, * as x from "library"
func (analyzer *analyzer) analyzeImportClause(node *sitter.Node, visitChildren func() error) error {
	for _, child := range analyzer.builder.ChildrenFor(node) {
		if child.Type() == "identifier" {
			analyzer.scope.Declare(analyzer.builder.ContentFor(child), child)
		}
	}

	return visitChildren()
}

// "* as x" part from:
// import * as x from "library"
func (analyzer *analyzer) analyzeNamespaceImport(node *sitter.Node, visitChildren func() error) error {
	for _, child := range analyzer.builder.ChildrenFor(node) {
		if child.Type() == "identifier" {
			analyzer.scope.Declare(analyzer.builder.ContentFor(child), child)
		}
	}

	return visitChildren()
}

// individual items inside the {}:
// import { x, y as foo } from "library"
func (analyzer *analyzer) analyzeImportSpecifier(node *sitter.Node, visitChildren func() error) error {
	importedName := node.ChildByFieldName("name")

	if alias := node.ChildByFieldName("alias"); alias != nil {
		importedName = alias
	}

	analyzer.scope.Declare(analyzer.builder.ContentFor(importedName), node)

	return visitChildren()
}

// a ? x : y
func (analyzer *analyzer) analyzeTernary(node *sitter.Node, visitChildren func() error) error {
	condition := node.ChildByFieldName("condition")
	consequence := node.ChildByFieldName("consequence")
	alternative := node.ChildByFieldName("alternative")

	analyzer.lookupVariable(condition)
	analyzer.lookupVariable(consequence)
	analyzer.lookupVariable(alternative)

	analyzer.builder.Alias(node, consequence, alternative)

	return visitChildren()
}

// (foo)
func (analyzer *analyzer) analyzeParentheses(node *sitter.Node, visitChildren func() error) error {
	child := node.NamedChild(0)
	analyzer.builder.Alias(node, child)
	analyzer.lookupVariable(child)

	return visitChildren()
}

// default analysis, where the children are assumed to be data sources
func (analyzer *analyzer) analyzeGenericOperation(node *sitter.Node, visitChildren func() error) error {
	children := analyzer.builder.ChildrenFor(node)
	analyzer.builder.Dataflow(node, children...)

	for _, child := range children {
		analyzer.lookupVariable(child)
	}

	return visitChildren()
}

func (analyzer *analyzer) withScope(newScope *implementation.Scope, body func() error) error {
	oldScope := analyzer.scope

	analyzer.scope = newScope
	err := body()
	analyzer.scope = oldScope

	return err
}

func (analyzer *analyzer) lookupVariable(node *sitter.Node) {
	if node == nil || node.Type() != "identifier" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
