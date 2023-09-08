package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/language"
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
	case "declaration_list",
		"method_declaration",
		"anonymous_function_creation_expression",
		"for_statement",
		"block":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
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
	case "simple_parameter", "variadic_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	case "switch_statement":
		return analyzer.analyzeSwitch(node, visitChildren)
	case "switch_block":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "switch_label":
		return visitChildren()
	case "arguments", "binary_expression", "unary_op_expression":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "while_statement", "do_statement", "if_statement": // statements don't have results
		return visitChildren()
	case "match_expression":
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenExcept(node, node.ChildByFieldName("condition"))...)
		return visitChildren()
	default:
		return visitChildren()
		// return nil
	}
}

// $foo = a
// $foo .= a
func (analyzer *analyzer) analyzeAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	operator := node.Child(1)

	if analyzer.builder.ContentFor(operator) == "=" {
		analyzer.builder.Alias(node, right)
	} else {
		analyzer.lookupVariable(left)
		analyzer.builder.Dataflow(node, left, right)
	}

	analyzer.lookupVariable(right)

	err := visitChildren()

	if left.Type() == "variable_name" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

func (analyzer *analyzer) analyzeParentheses(node *sitter.Node, visitChildren func() error) error {
	child := node.NamedChild(0)
	analyzer.builder.Alias(node, child)
	analyzer.lookupVariable(child)

	return visitChildren()
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

// foo->bar(1, 2);
func (analyzer *analyzer) analyzeMethodInvocation(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.ChildByFieldName("object"))

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

	return visitChildren()
}

func (analyzer *analyzer) analyzeSwitch(node *sitter.Node, visitChildren func() error) error {
	analyzer.builder.Alias(node, node.ChildByFieldName("body"))

	return visitChildren()
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
	if node == nil || node.Type() != "variable_name" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
