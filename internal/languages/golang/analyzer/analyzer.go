package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"
	"golang.org/x/exp/slices"

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
	case "for_statement", "block", "method_declaration", "function_declaration":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "short_var_declaration":
		return analyzer.analyzeShortVarDeclaration(node, visitChildren)
	case "var_spec":
		return analyzer.analyzeVarSpecDeclaration(node, visitChildren)
	case "assignment_expression":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "call_expression":
		return analyzer.analyzeCallExpression(node, visitChildren)
	case "selector_expression":
		return analyzer.analyzeSelectorExpression(node, visitChildren)
	case "parameter_declaration":
		return analyzer.analyzeParameter(node, visitChildren)
	case "expression_switch_statement":
		return analyzer.analyzeSwitch(node, visitChildren)
	case "expression_case", "default_case":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "argument_list":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "return_statement", "go_statement", "defer_statement", "if_statement": // statements don't have results
		return visitChildren()
	case "identifier":
		return visitChildren()
	case "index_expression":
		return visitChildren()
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		return visitChildren()
	}
}

// foo = a
// foo += a
func (analyzer *analyzer) analyzeAssignment(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left").Child(0)
	right := node.ChildByFieldName("right").Child(0)

	if analyzer.builder.ContentFor(node.Child(1)) == "=" {
		analyzer.builder.Alias(node, right)
	} else { // +=
		analyzer.builder.Dataflow(node, left, right)
		analyzer.lookupVariable(left)
	}

	analyzer.lookupVariable(right)

	err := visitChildren()

	analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)

	return err
}

// foo, err := a
func (analyzer *analyzer) analyzeShortVarDeclaration(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	for _, child := range analyzer.builder.ChildrenFor(left) {
		if !slices.Contains([]string{"_", "err"}, analyzer.builder.ContentFor(child)) {
			analyzer.scope.Declare(analyzer.builder.ContentFor(child), child)
			analyzer.scope.Assign(analyzer.builder.ContentFor(child), node)
		}
	}

	for _, child := range analyzer.builder.ChildrenFor(right) {
		analyzer.builder.Alias(node, child)
		analyzer.lookupVariable(child)
	}

	err := visitChildren()

	return err
}

// var a, b string
func (analyzer *analyzer) analyzeVarSpecDeclaration(node *sitter.Node, visitChildren func() error) error {
	for _, child := range analyzer.builder.ChildrenFor(node) {
		if child.Type() == "identifier" {
			analyzer.scope.Declare(analyzer.builder.ContentFor(child), child)
		}
	}

	err := visitChildren()

	return err
}

// foo(1, 2)
func (analyzer *analyzer) analyzeCallExpression(node *sitter.Node, visitChildren func() error) error {
	if arguments := node.ChildByFieldName("arguments"); arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// foo.bar
func (analyzer *analyzer) analyzeSelectorExpression(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.ChildByFieldName("operand"))

	return visitChildren()
}

// method parameter declaration
//
// fn(a string)
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	name := node.ChildByFieldName("name")
	if name != nil {
		analyzer.builder.Alias(node, name)
		analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)
	}

	return visitChildren()
}

func (analyzer *analyzer) analyzeSwitch(node *sitter.Node, visitChildren func() error) error {
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
	if node == nil || node.Type() != "identifier" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
