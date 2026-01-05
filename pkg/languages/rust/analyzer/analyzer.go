package analyzer

import (
	"slices"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/language"
)

type analyzer struct {
	builder *tree.Builder
	scope   *language.Scope
}

func New(builder *tree.Builder) language.Analyzer {
	return &analyzer{
		builder: builder,
		scope:   language.NewScope(nil),
	}
}

func (analyzer *analyzer) Analyze(node *sitter.Node, visitChildren func() error) error {
	switch node.Type() {
	// Scope-creating nodes
	case "function_item", "impl_item", "mod_item", "block", "closure_expression":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	// Variable declarations
	case "let_declaration":
		return analyzer.analyzeLetDeclaration(node, visitChildren)
	// Assignments
	case "assignment_expression":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "compound_assignment_expr":
		return analyzer.analyzeCompoundAssignment(node, visitChildren)
	// Function calls
	case "call_expression":
		return analyzer.analyzeCall(node, visitChildren)
	// Method calls
	case "method_call_expression":
		return analyzer.analyzeMethodCall(node, visitChildren)
	// Field access
	case "field_expression":
		return analyzer.analyzeFieldExpression(node, visitChildren)
	// Match expressions
	case "match_expression":
		return analyzer.analyzeMatch(node, visitChildren)
	// Macro invocations
	case "macro_invocation":
		return analyzer.analyzeMacro(node, visitChildren)
	// If expressions
	case "if_expression":
		return analyzer.analyzeIf(node, visitChildren)
	// For loops
	case "for_expression":
		return analyzer.analyzeFor(node, visitChildren)
	// While loops
	case "while_expression":
		return visitChildren()
	// Use statements (imports)
	case "use_declaration":
		return analyzer.analyzeUse(node, visitChildren)
	// Parameters
	case "parameter", "self_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	// Tuple expressions, arrays, argument lists, etc.
	case "tuple_expression", "array_expression", "arguments", "binary_expression", "unary_expression":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	// Return, await, try expressions
	case "return_expression", "await_expression", "try_expression":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	// Index expressions
	case "index_expression":
		return analyzer.analyzeIndexExpression(node, visitChildren)
	// Reference expressions
	case "reference_expression":
		return analyzer.analyzeReference(node, visitChildren)
	// Dereference expressions
	case "dereference_expression":
		return analyzer.analyzeDereference(node, visitChildren)
	// Identifier lookup
	case "identifier":
		return visitChildren()
	// Statements without results
	case "expression_statement", "loop_expression":
		return visitChildren()
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		return visitChildren()
	}
}

// let x = value;
// let mut x = value;
// let x: Type = value;
func (analyzer *analyzer) analyzeLetDeclaration(node *sitter.Node, visitChildren func() error) error {
	pattern := node.ChildByFieldName("pattern")
	value := node.ChildByFieldName("value")

	if value != nil {
		analyzer.builder.Alias(node, value)
		analyzer.lookupVariable(value)
	}

	err := visitChildren()

	if pattern != nil && pattern.Type() == "identifier" {
		analyzer.scope.Declare(analyzer.builder.ContentFor(pattern), pattern)
		analyzer.scope.Assign(analyzer.builder.ContentFor(pattern), node)
	}

	return err
}

// x = value;
func (analyzer *analyzer) analyzeAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	analyzer.builder.Alias(node, right)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left != nil && left.Type() == "identifier" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

// x += value;
func (analyzer *analyzer) analyzeCompoundAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	analyzer.builder.Dataflow(node, left, right)
	analyzer.lookupVariable(left)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left != nil && left.Type() == "identifier" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

// function(args)
func (analyzer *analyzer) analyzeCall(node *sitter.Node, visitChildren func() error) error {
	function := node.ChildByFieldName("function")
	arguments := node.ChildByFieldName("arguments")

	analyzer.lookupVariable(function)

	if arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// object.method(args)
func (analyzer *analyzer) analyzeMethodCall(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")
	arguments := node.ChildByFieldName("arguments")

	analyzer.lookupVariable(value)
	analyzer.builder.Dataflow(node, value)

	if arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// object.field
func (analyzer *analyzer) analyzeFieldExpression(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")
	analyzer.lookupVariable(value)
	analyzer.builder.Dataflow(node, value)

	return visitChildren()
}

// match expr { ... }
func (analyzer *analyzer) analyzeMatch(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")
	analyzer.lookupVariable(value)

	return visitChildren()
}

// macro!(args)
func (analyzer *analyzer) analyzeMacro(node *sitter.Node, visitChildren func() error) error {
	// Macros can have token trees as arguments
	// We treat the whole macro as a dataflow node
	analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)

	return visitChildren()
}

// if condition { ... } else { ... }
func (analyzer *analyzer) analyzeIf(node *sitter.Node, visitChildren func() error) error {
	condition := node.ChildByFieldName("condition")
	analyzer.lookupVariable(condition)

	return visitChildren()
}

// for pattern in iterator { ... }
func (analyzer *analyzer) analyzeFor(node *sitter.Node, visitChildren func() error) error {
	pattern := node.ChildByFieldName("pattern")
	value := node.ChildByFieldName("value")

	if value != nil {
		analyzer.lookupVariable(value)
	}

	if pattern != nil && pattern.Type() == "identifier" {
		analyzer.builder.Dataflow(pattern, value)
		analyzer.scope.Declare(analyzer.builder.ContentFor(pattern), pattern)
	}

	return visitChildren()
}

// use path::to::item;
func (analyzer *analyzer) analyzeUse(node *sitter.Node, visitChildren func() error) error {
	// Use declarations bring items into scope
	// For now, we just visit children
	return visitChildren()
}

// fn foo(param: Type)
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	pattern := node.ChildByFieldName("pattern")
	if pattern != nil && pattern.Type() == "identifier" {
		analyzer.builder.Alias(node, pattern)
		analyzer.scope.Declare(analyzer.builder.ContentFor(pattern), pattern)
	}

	return visitChildren()
}

// array[index]
func (analyzer *analyzer) analyzeIndexExpression(node *sitter.Node, visitChildren func() error) error {
	value := node.NamedChild(0)
	analyzer.lookupVariable(value)
	analyzer.builder.Dataflow(node, value)

	return visitChildren()
}

// &value or &mut value
func (analyzer *analyzer) analyzeReference(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")
	analyzer.lookupVariable(value)
	analyzer.builder.Alias(node, value)

	return visitChildren()
}

// *value
func (analyzer *analyzer) analyzeDereference(node *sitter.Node, visitChildren func() error) error {
	value := node.NamedChild(0)
	analyzer.lookupVariable(value)
	analyzer.builder.Alias(node, value)

	return visitChildren()
}

// default analysis, where the children are assumed to be aliases
func (analyzer *analyzer) analyzeGenericConstruct(node *sitter.Node, visitChildren func() error) error {
	children := analyzer.builder.ChildrenFor(node)
	analyzer.builder.Alias(node, children...)

	for _, child := range children {
		analyzer.lookupVariable(child)
	}

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

func (analyzer *analyzer) withScope(newScope *language.Scope, body func() error) error {
	oldScope := analyzer.scope

	analyzer.scope = newScope
	err := body()
	analyzer.scope = oldScope

	return err
}

func (analyzer *analyzer) lookupVariable(node *sitter.Node) {
	if node == nil || !slices.Contains([]string{"identifier", "self"}, node.Type()) {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}

