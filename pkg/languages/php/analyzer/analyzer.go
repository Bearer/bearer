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
	case "declaration_list", "class_declaration", "anonymous_function_creation_expression", "for_statement", "block", "method_declaration":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "augmented_assignment_expression":
		return analyzer.analyzeAugmentedAssignment(node, visitChildren)
	case "assignment_expression":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "parenthesized_expression":
		return analyzer.analyzeParentheses(node, visitChildren)
	case "conditional_expression":
		return analyzer.analyzeConditional(node, visitChildren)
	case "function_call_expression", "member_call_expression":
		return analyzer.analyzeMethodInvocation(node, visitChildren)
	case "member_access_expression":
		return analyzer.analyzeFieldAccess(node, visitChildren)
	case "simple_parameter", "variadic_parameter", "property_promotion_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	case "switch_statement":
		return analyzer.analyzeSwitch(node, visitChildren)
	case "switch_block":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "switch_label":
		return visitChildren()
	case "const_declaration":
		return analyzer.analyzeConstDeclaration(node, visitChildren)
	case "dynamic_variable_name":
		return analyzer.analyzeDynamicVariableName(node, visitChildren)
	case "subscript_expression":
		return analyzer.analyzeSubscript(node, visitChildren)
	case "catch_clause":
		return analyzer.analyzeCatchClause(node, visitChildren)
	case "foreach_statement":
		return analyzer.analyzeForeach(node, visitChildren)
	case "binary_expression",
		"unary_op_expression",
		"argument",
		"encapsed_string",
		"sequence_expression",
		"array_element_initializer",
		"formal_parameters",
		"include_expression",
		"include_once_expression",
		"require_expression",
		"require_once_expression",
		"echo_statement",
		"print_intrinsic":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "while_statement", "do_statement", "if_statement", "expression_statement", "compound_statement": // statements don't have results
		return visitChildren()
	case "variable_name":
		return visitChildren()
	case "match_expression":
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenExcept(node, node.ChildByFieldName("condition"))...)
		return visitChildren()
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		return visitChildren()
	}
}

// $foo = a
func (analyzer *analyzer) analyzeAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	analyzer.builder.Alias(node, right)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "variable_name" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

// $foo .= a
func (analyzer *analyzer) analyzeAugmentedAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	analyzer.builder.Dataflow(node, left, right)
	analyzer.lookupVariable(left)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "variable_name" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

func (analyzer *analyzer) analyzeParentheses(node *sitter.Node, visitChildren func() error) error {
	analyzer.builder.Alias(node, node.NamedChild(0))
	analyzer.lookupVariable(node.NamedChild(0))
	err := visitChildren()

	return err
}

// a ? x : y
// a ?: x
func (analyzer *analyzer) analyzeConditional(node *sitter.Node, visitChildren func() error) error {
	condition := node.ChildByFieldName("condition")
	consequence := node.ChildByFieldName("body")
	alternative := node.ChildByFieldName("alternative")

	analyzer.lookupVariable(condition)
	analyzer.lookupVariable(consequence)
	analyzer.lookupVariable(alternative)

	if consequence != nil {
		analyzer.builder.Alias(node, consequence, alternative)
	} else {
		analyzer.builder.Alias(node, condition, alternative)
	}

	return visitChildren()
}

// foo(1, 2);
// foo->bar(1, 2);
func (analyzer *analyzer) analyzeMethodInvocation(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.ChildByFieldName("object"))   // method
	analyzer.lookupVariable(node.ChildByFieldName("function")) // function

	if arguments := node.ChildByFieldName("arguments"); arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// foo->bar
func (analyzer *analyzer) analyzeFieldAccess(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.ChildByFieldName("object"))

	return visitChildren()
}

// method parameter declaration
//
// fn(bool $a) => $a;
// fn($x = 42) => $x;
// fn($x, ...$rest) => $rest;
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	name := node.ChildByFieldName("name")
	analyzer.builder.Alias(node, name)
	analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)

	return visitChildren()
}

func (analyzer *analyzer) analyzeSwitch(node *sitter.Node, visitChildren func() error) error {
	analyzer.builder.Alias(node, node.ChildByFieldName("body"))

	return visitChildren()
}

func (analyzer *analyzer) analyzeDynamicVariableName(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.NamedChild(0))

	return visitChildren()
}

// foo["bar"]
func (analyzer *analyzer) analyzeSubscript(node *sitter.Node, visitChildren func() error) error {
	object := node.NamedChild(0)
	analyzer.builder.Dataflow(node, object)
	analyzer.lookupVariable(object)

	analyzer.lookupVariable(node.NamedChild(1))

	return visitChildren()
}

func (analyzer *analyzer) analyzeConstDeclaration(node *sitter.Node, visitChildren func() error) error {
	var child *sitter.Node

	for i := 0; i < int(node.ChildCount()); i++ {
		child = node.NamedChild(i)

		if child.Type() == "const_element" {
			break
		}
	}

	left := child.NamedChild(0)
	right := child.NamedChild(1)
	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "name" {
		analyzer.builder.Alias(left, right)
		analyzer.scope.Declare("self::"+analyzer.builder.ContentFor(left), right)
	}

	return err
}

// catch(FooException | BarException $e) {}
// catch(FooException) {} // PHP 8+ allows catch without variable
func (analyzer *analyzer) analyzeCatchClause(node *sitter.Node, visitChildren func() error) error {
	return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
		name := node.ChildByFieldName("name")
		if name != nil {
			analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)
		}

		return visitChildren()
	})
}

// foreach ($array as $value) {}
// foreach ($array as $key => $value) {}
func (analyzer *analyzer) analyzeForeach(node *sitter.Node, visitChildren func() error) error {
	return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
		array := node.NamedChild(0)
		analyzer.lookupVariable(array)

		value := node.NamedChild(1)
		if value.Type() == "pair" {
			key := value.NamedChild(0)
			analyzer.scope.Declare(analyzer.builder.ContentFor(key), key)
			value = value.NamedChild(1)
		}

		analyzer.scope.Declare(analyzer.builder.ContentFor(value), value)
		analyzer.builder.Dataflow(value, array)

		return visitChildren()
	})
}

// default analysis, where the children are assumed to be aliases
func (analyzer *analyzer) analyzeGenericConstruct(node *sitter.Node, visitChildren func() error) error {
	analyzer.builder.Alias(node, analyzer.builder.ChildrenFor(node)...)

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
	if node == nil || !slices.Contains([]string{"variable_name", "class_constant_access_expression"}, node.Type()) {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
