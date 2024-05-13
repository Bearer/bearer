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
	case "class_definition", "block", "function_definition":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "augmented_assignment":
		return analyzer.analyzeAugmentedAssignment(node, visitChildren)
	case "assignment":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "attribute":
		return analyzer.analyzeAttribute(node, visitChildren)
	case "subscript":
		return analyzer.analyzeSubscript(node, visitChildren)
	case "call":
		return analyzer.analyzeCall(node, visitChildren)
	case "argument_list", "expression_statement", "list", "tuple":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "parenthesized_expression":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "keyword_argument":
		return analyzer.analyzeKeywordArgumentOperation(node, visitChildren)
	case "while_statement", "try_statement", "if_statement": // statements don't have results
		return visitChildren()
	case "conditional_expression":
		return analyzer.analyzeConditional(node, visitChildren)
	case "boolean_operator":
		return analyzer.analyzeBoolean(node, visitChildren)
	case "identifier":
		return visitChildren()
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		return visitChildren()
	}
}

// foo += a
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

// foo = a
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

// foo.bar(a, b)
func (analyzer *analyzer) analyzeCall(node *sitter.Node, visitChildren func() error) error {
	if receiver := node.ChildByFieldName("function"); receiver != nil {
		analyzer.lookupVariable(receiver)

		analyzer.builder.Dataflow(node, receiver)
	}

	if argumentsNode := node.ChildByFieldName("arguments"); argumentsNode != nil {
		analyzer.builder.Dataflow(node, argumentsNode)
	}

	return visitChildren()
}

// foo.bar
func (analyzer *analyzer) analyzeAttribute(node *sitter.Node, visitChildren func() error) error {
	if receiver := node.ChildByFieldName("object"); receiver != nil {
		analyzer.lookupVariable(receiver)
		analyzer.builder.Dataflow(node, receiver)
	}

	return visitChildren()
}

// foo["bar"]
func (analyzer *analyzer) analyzeSubscript(node *sitter.Node, visitChildren func() error) error {
	objectNode := node.ChildByFieldName("value")
	analyzer.builder.Dataflow(node, objectNode)
	analyzer.lookupVariable(objectNode)

	return visitChildren()
}

// x if foo else y
func (analyzer *analyzer) analyzeConditional(node *sitter.Node, visitChildren func() error) error {
	condition := node.NamedChild(1)
	consequence := node.NamedChild(0)
	alternative := node.NamedChild(2)

	analyzer.lookupVariable(condition)
	analyzer.lookupVariable(consequence)
	analyzer.lookupVariable(alternative)

	analyzer.builder.Alias(node, consequence, alternative)

	return visitChildren()
}

// a or b
func (analyzer *analyzer) analyzeBoolean(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	analyzer.lookupVariable(left)
	analyzer.lookupVariable(right)

	analyzer.builder.Alias(node, left, right)

	return visitChildren()
}

func (analyzer *analyzer) analyzeKeywordArgumentOperation(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")

	analyzer.builder.Alias(node, value)
	analyzer.lookupVariable(value)

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
	if node == nil || node.Type() != "identifier" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
