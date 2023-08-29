package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/language"
)

// methods that use `self` in their result
var reflexiveMethods = []string{
	"to_a",
	"to_ary",
	"to_h",
	"to_hash",
	"to_s",
	"to_str",
	"to_i",
	"to_f",
	"to_c",
	"to_d",
	"to_r",
	"to_sym",
	"to_json",
	"sub",
	"sub!",
	"gsub",
	"gsub!",
}

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
	case "method":
		return analyzer.withScope(language.NewScope(nil), func() error {
			return visitChildren()
		})
	case "block", "do_block":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "assignment":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "operator_assignment":
		return analyzer.analyzeOperatorAssignment(node, visitChildren)
	case "call":
		return analyzer.analyzeCall(node, visitChildren)
	case "element_reference":
		return analyzer.analyzeElementReference(node, visitChildren)
	case "case":
		return analyzer.analyzeCase(node, visitChildren)
	case "when":
		return analyzer.analyzeWhen(node, visitChildren)
	case "then", "else":
		return analyzer.analyzeBasicBlock(node, visitChildren)
	case "keyword_parameter", "optional_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	case "method_parameters", "block_parameters":
		return analyzer.analyzeParameterList(node, visitChildren)
	case "parenthesized_statements":
		return analyzer.analyzeParentheses(node, visitChildren)
	case "conditional":
		return analyzer.analyzeConditional(node, visitChildren)
	case "pair", "argument_list", "interpolation", "array", "binary", "unary":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenExcept(node, node.ChildByFieldName("condition"))...)

		return visitChildren()
	}
}

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
func (analyzer *analyzer) analyzeOperatorAssignment(node *sitter.Node, visitChildren func() error) error {
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

// foo.bar(42)
func (analyzer *analyzer) analyzeCall(node *sitter.Node, visitChildren func() error) error {
	if receiver := node.ChildByFieldName("receiver"); receiver != nil {
		analyzer.lookupVariable(receiver)

		if slices.Contains(reflexiveMethods, analyzer.builder.ContentFor(node.ChildByFieldName("method"))) {
			analyzer.builder.Dataflow(node, receiver)
		}
	}

	if argumentsNode := node.ChildByFieldName("arguments"); argumentsNode != nil {
		analyzer.builder.Dataflow(node, argumentsNode)
	}

	return visitChildren()
}

// foo["bar"]
func (analyzer *analyzer) analyzeElementReference(node *sitter.Node, visitChildren func() error) error {
	objectNode := node.ChildByFieldName("object")
	analyzer.builder.Dataflow(node, objectNode)
	analyzer.lookupVariable(objectNode)

	return visitChildren()
}

// case foo
// ...
// end
func (analyzer *analyzer) analyzeCase(node *sitter.Node, visitChildren func() error) error {
	if valueNode := node.ChildByFieldName("value"); valueNode != nil {
		analyzer.builder.Alias(
			node,
			analyzer.builder.ChildrenExcept(node, valueNode)...,
		)
	} else {
		analyzer.builder.Alias(node, analyzer.builder.ChildrenFor(node)...)
	}

	return visitChildren()
}

// Any construct that is just a block of code. eg.
//
//	if x
//	  expr...
//	end
//	case/if
//	  ...
//	else
//	  expr...
//	end
func (analyzer *analyzer) analyzeBasicBlock(node *sitter.Node, visitChildren func() error) error {
	if lastChild := analyzer.builder.LastChild(node); lastChild != nil {
		analyzer.builder.Alias(node, lastChild)
	}

	return visitChildren()
}

// case foo
// when 1
// end
func (analyzer *analyzer) analyzeWhen(node *sitter.Node, visitChildren func() error) error {
	if patternNode := node.ChildByFieldName("pattern"); patternNode != nil {
		analyzer.builder.Alias(
			node,
			analyzer.builder.ChildrenExcept(node, patternNode)...,
		)
	} else {
		analyzer.builder.Alias(node, analyzer.builder.ChildrenFor(node)...)
	}

	return visitChildren()
}

// keyword or default parameter definition
// def m(a = 1, b:)
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	nameNode := node.ChildByFieldName("name")

	if nameNode.Type() == "identifier" {
		analyzer.scope.Declare(analyzer.builder.ContentFor(nameNode), nameNode)
	}

	return visitChildren()
}

// all parameter definitions for a method/block
// def m(a, b = 1)
func (analyzer *analyzer) analyzeParameterList(node *sitter.Node, visitChildren func() error) error {
	children := analyzer.builder.ChildrenFor(node)
	analyzer.builder.Dataflow(node, children...)

	for _, child := range children {
		if child.Type() == "identifier" {
			analyzer.scope.Declare(analyzer.builder.ContentFor(child), child)
		}
	}

	return visitChildren()
}

// (foo)
func (analyzer *analyzer) analyzeParentheses(node *sitter.Node, visitChildren func() error) error {
	child := node.NamedChild(0)
	analyzer.builder.Alias(node, child)
	analyzer.lookupVariable(child)

	return visitChildren()
}

// foo ? x : y
func (analyzer *analyzer) analyzeConditional(node *sitter.Node, visitChildren func() error) error {
	condition := node.ChildByFieldName("condition")
	consequence := node.ChildByFieldName("consequence")
	alternative := node.ChildByFieldName("alternative")

	analyzer.lookupVariable(condition)
	analyzer.lookupVariable(consequence)
	analyzer.lookupVariable(alternative)

	analyzer.builder.Alias(node, consequence, alternative)

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
	if node == nil || node.Type() != "identifier" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
