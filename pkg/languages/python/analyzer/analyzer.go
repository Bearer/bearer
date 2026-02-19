package analyzer

import (
	"slices"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/language"
)

var reflexiveMethods = []string{
	"decode",
	"encode",
	"format",
	"replace",
	"split",
	"lower",
	"upper",
	"strip",
	"lstrip",
	"rstrip",
	"capitalize",
	"title",
	"swapcase",
	"casefold",
	"expandtabs",
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
	case "class_definition", "function_definition":
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
	case "pair", "argument_list", "expression_statement", "list", "tuple", "unary_operator", "binary_operator":
		return analyzer.analyzeGenericOperation(node, visitChildren)
	case "parenthesized_expression", "interpolation", "return_statement":
		return analyzer.analyzeGenericConstruct(node, visitChildren)
	case "parameters":
		return analyzer.analyzeParameters(node, visitChildren)
	case "keyword_argument":
		return analyzer.analyzeKeywordArgument(node, visitChildren)
	case "for_statement":
		return analyzer.analyzeForStatement(node, visitChildren)
	case "with_item":
		return analyzer.analyzeWithItem(node, visitChildren)
	case "while_statement", "try_statement", "if_statement": // statements don't have results
		return visitChildren()
	case "conditional_expression":
		return analyzer.analyzeConditional(node, visitChildren)
	case "boolean_operator":
		return analyzer.analyzeBoolean(node, visitChildren)
	case "import_statement", "import_from_statement":
		return analyzer.analyzeImport(node, visitChildren)
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
	if function := node.ChildByFieldName("function"); function != nil {
		analyzer.lookupVariable(function)
		analyzer.builder.Dataflow(function)

		object := function.ChildByFieldName("object")
		analyzer.lookupVariable(object)

		if function.Type() == "identifier" {
			analyzer.builder.Dataflow(node, object)
		}

		if function.Type() == "attribute" {
			attribute := function.ChildByFieldName("attribute")
			if attribute.Type() == "identifier" && slices.Contains(reflexiveMethods, analyzer.builder.ContentFor(attribute)) {
				analyzer.builder.Dataflow(node, object)
			}
		}
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
// foo[x]
// globals()[y]
func (analyzer *analyzer) analyzeSubscript(node *sitter.Node, visitChildren func() error) error {
	objectNode := node.ChildByFieldName("value")
	analyzer.builder.Dataflow(node, objectNode)
	analyzer.lookupVariable(objectNode)

	subscriptNode := node.ChildByFieldName("subscript")
	if subscriptNode.Type() == "identifier" {
		analyzer.lookupVariable(subscriptNode)
	}

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

// def f(self, param: Type)
// def f(param: Type = default)
func (analyzer *analyzer) analyzeParameters(node *sitter.Node, visitChildren func() error) error {
	err := visitChildren()

	for _, parameter := range analyzer.builder.ChildrenFor(node) {
		switch parameter.Type() {
		case "typed_parameter", "typed_default_parameter", "default_parameter":
			name := parameter.NamedChild(0)

			analyzer.builder.Alias(parameter, name)
			analyzer.scope.Declare(analyzer.builder.ContentFor(name), name)

			analyzer.lookupVariable(parameter.ChildByFieldName("type"))

			// Create dataflow from default value to parameter name
			if defaultValue := parameter.ChildByFieldName("value"); defaultValue != nil {
				analyzer.builder.Alias(name, defaultValue)
				analyzer.lookupVariable(defaultValue)
			}
		case "identifier":
			analyzer.scope.Declare(analyzer.builder.ContentFor(parameter), parameter)
		}
	}

	return err
}

func (analyzer *analyzer) analyzeKeywordArgument(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")

	analyzer.builder.Alias(node, value)
	analyzer.lookupVariable(value)

	return visitChildren()
}

func (analyzer *analyzer) analyzeForStatement(node *sitter.Node, visitChildren func() error) error {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")

	analyzer.builder.Dataflow(left, right)
	analyzer.scope.Declare(analyzer.builder.ContentFor(left), left)
	analyzer.lookupVariable(right)

	return visitChildren()
}

// with x, y as foo:
func (analyzer *analyzer) analyzeWithItem(node *sitter.Node, visitChildren func() error) error {
	value := node.ChildByFieldName("value")
	analyzer.lookupVariable(value)

	err := visitChildren()

	if value.Type() == "as_pattern" {
		aliasValue := value.NamedChild(0)
		alias := value.ChildByFieldName("alias")
		analyzer.lookupVariable(aliasValue)
		analyzer.builder.Alias(alias, aliasValue)
		analyzer.scope.Declare(analyzer.builder.ContentFor(alias), alias)
	}

	return err
}

// import x
// import a.b
// from z import x
// import x as y (aliased_import)
// from z import x as y (aliased_import)
func (analyzer *analyzer) analyzeImport(node *sitter.Node, visitChildren func() error) error {
	children := analyzer.builder.ChildrenExcept(node, node.ChildByFieldName("module_name"))

	for _, child := range children {
		switch child.Type() {
		case "aliased_import":
			aliasedImportIdentifier := child.ChildByFieldName("alias")
			analyzer.scope.Declare(analyzer.builder.ContentFor(aliasedImportIdentifier), aliasedImportIdentifier)
		case "dotted_name":
			analyzer.scope.Declare(analyzer.builder.ContentFor(child.NamedChild(0)), child.NamedChild(0))
		}
	}

	return nil
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
	if node == nil || !slices.Contains([]string{"identifier", "type"}, node.Type()) {
		return
	}

	if pointsToNode := analyzer.scope.Lookup(analyzer.builder.ContentFor(node)); pointsToNode != nil {
		analyzer.builder.Alias(node, pointsToNode)
	}
}
