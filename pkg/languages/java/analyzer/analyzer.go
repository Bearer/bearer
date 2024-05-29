package analyzer

import (
	"slices"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/language"
)

// methods that use `this` in their result
var reflexiveMethods = []string{
	// String
	"getBytes",
	"replace",
	"replaceAll",
	"split",
	"substring",
	"toCharArray",
	// StringBuilder
	"append",
	"toString",
	// Enumeration
	"nextElement",
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
	case "class_body",
		"method_declaration",
		"lambda_expression",
		"for_statement",
		"block",
		"try_with_resources_statement":
		return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
			return visitChildren()
		})
	case "assignment_expression":
		return analyzer.analyzeAssignment(node, visitChildren)
	case "variable_declarator":
		return analyzer.analyzeVariableDeclarator(node, visitChildren)
	case "parenthesized_expression":
		return analyzer.analyzeParentheses(node, visitChildren)
	case "ternary_expression":
		return analyzer.analyzeTernary(node, visitChildren)
	case "method_invocation":
		return analyzer.analyzeMethodInvocation(node, visitChildren)
	case "field_access":
		return analyzer.analyzeFieldAccess(node, visitChildren)
	case "enhanced_for_statement":
		return analyzer.analyzeEnhancedForStatement(node, visitChildren)
	case "formal_parameter", "catch_formal_parameter":
		return analyzer.analyzeParameter(node, visitChildren)
	case "resource":
		return analyzer.analyzeResource(node, visitChildren)
	case "cast_expression":
		return analyzer.analyzeCastExpression(node, visitChildren)
	case "switch_expression":
		return analyzer.analyzeSwitch(node, visitChildren)
	case "switch_block":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "switch_label":
		return visitChildren()
	case "argument_list", "array_access", "array_initializer", "binary_expression", "unary_expression":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "while_statement", "do_statement", "if_statement": // statements don't have results
		return visitChildren()
	default:
		analyzer.builder.Dataflow(node, analyzer.builder.ChildrenFor(node)...)
		return visitChildren()
	}
}

// foo = a
// foo += a
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

	if left.Type() == "identifier" {
		analyzer.scope.Assign(analyzer.builder.ContentFor(left), node)
	}

	return err
}

func (analyzer *analyzer) analyzeCastExpression(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")

	analyzer.builder.Alias(node, value)

	analyzer.lookupVariable(value)

	return visitChildren()
}

// the "foo = 1" part in:
//
//	class X {
//	  void m() {
//	  	 Integer foo = 1;
//	  }
//	}
func (analyzer *analyzer) analyzeVariableDeclarator(node *sitter.Node, visitChildren func() error) error {
	name := node.ChildByFieldName("name")

	// backwards compatibility with rules. fixup rules to use variable name node,
	// and then remove this
	analyzer.builder.Alias(name, node.Parent())

	if value := node.ChildByFieldName("value"); value != nil {
		analyzer.lookupVariable(value)
		analyzer.builder.Alias(name, value)
	}

	err := visitChildren()

	analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)

	return err
}

// (foo)
func (analyzer *analyzer) analyzeParentheses(node *sitter.Node, visitChildren func() error) error {
	child := node.NamedChild(0)
	analyzer.builder.Alias(node, child)
	analyzer.lookupVariable(child)

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

// foo.bar(1, 2);
func (analyzer *analyzer) analyzeMethodInvocation(node *sitter.Node, visitChildren func() error) error {
	if object := node.ChildByFieldName("object"); object != nil {
		analyzer.lookupVariable(object)

		if slices.Contains(reflexiveMethods, analyzer.builder.ContentFor(node.ChildByFieldName("name"))) {
			analyzer.builder.Dataflow(node, object)
		}
	}

	if arguments := node.ChildByFieldName("arguments"); arguments != nil {
		analyzer.builder.Dataflow(node, arguments)
	}

	return visitChildren()
}

// foo.bar
func (analyzer *analyzer) analyzeFieldAccess(node *sitter.Node, visitChildren func() error) error {
	analyzer.lookupVariable(node.ChildByFieldName("object"))

	return visitChildren()
}

// for (String value : array)
func (analyzer *analyzer) analyzeEnhancedForStatement(node *sitter.Node, visitChildren func() error) error {
	return analyzer.withScope(language.NewScope(analyzer.scope), func() error {
		name := node.ChildByFieldName("name")
		value := node.ChildByFieldName("value")

		analyzer.lookupVariable(value)
		analyzer.builder.Dataflow(name, value)
		analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)

		return visitChildren()
	})
}

// method parameter declaration or catch parameter declaration
//
// void m(String foo) {}
// try {} catch (Exception foo) {}
func (analyzer *analyzer) analyzeParameter(node *sitter.Node, visitChildren func() error) error {
	name := node.ChildByFieldName("name")
	analyzer.builder.Alias(node, name)

	if name.Type() == "identifier" {
		analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)
	}

	return visitChildren()
}

// parts like "foo" and "File f = open()" from:
// try (foo; File f = open(); Other b = ...) {}
func (analyzer *analyzer) analyzeResource(node *sitter.Node, visitChildren func() error) error {
	if name := node.ChildByFieldName("name"); name != nil {
		value := node.ChildByFieldName("value")
		analyzer.builder.Alias(node, value)
		analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)
	}

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
	if node == nil || node.Type() != "identifier" {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
