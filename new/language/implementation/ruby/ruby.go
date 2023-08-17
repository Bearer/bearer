package ruby

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/ruby/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/ruby/string"
	"github.com/bearer/bearer/new/language/implementation"
)

type rubyImplementation struct {
	pattern patternImplementation
}

func Get() implementation.Implementation {
	return &rubyImplementation{}
}

func (*rubyImplementation) Name() string {
	return "ruby"
}

func (*rubyImplementation) EnryLanguages() []string {
	return []string{"Ruby"}
}

func (*rubyImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorRuby, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*rubyImplementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (implementation *rubyImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*rubyImplementation) AnalyzeTree(ctx context.Context, rootNode *sitter.Node, builder *tree.Builder) error {
	return analyzeNode(ctx, rootNode, builder, implementation.NewScope(nil))
}

// FIXME: extract the common logic across languages
func analyzeNode(ctx context.Context, node *sitter.Node, builder *tree.Builder, scope *implementation.Scope) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func(childScope *implementation.Scope) error {
		childCount := int(node.ChildCount())

		for i := 0; i < childCount; i++ {
			child := node.Child(i)
			if err := analyzeNode(ctx, child, builder, childScope); err != nil {
				return err
			}
		}

		return nil
	}

	lookupVariable := func(variableNode *sitter.Node) {
		if variableNode == nil || variableNode.Type() != "identifier" {
			return
		}

		if scopedNode := scope.Lookup(builder.ContentFor(variableNode)); scopedNode != nil {
			builder.Dataflow(variableNode, scopedNode)
		}
	}

	switch node.Type() {
	case "method":
		return visitChildren(implementation.NewScope(nil))
	case "block", "do_block":
		return visitChildren(implementation.NewScope(scope))
	case "assignment":
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		builder.Dataflow(node, right)
		lookupVariable(right)

		if left.Type() == "identifier" {
			err := visitChildren(scope)
			scope.Assign(builder.ContentFor(left), node)
			return err
		}
	// x += y
	case "operator_assignment":
		builder.SetOperation(node)

		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		builder.Dataflow(node, left, right)
		lookupVariable(left)
		lookupVariable(right)

		if left.Type() == "identifier" {
			err := visitChildren(scope)
			scope.Assign(builder.ContentFor(left), node)
			return err
		}
	// foo.bar(42)
	case "call":
		builder.SetOperation(node)

		if argumentsNode := node.ChildByFieldName("arguments"); argumentsNode != nil {
			builder.Dataflow(node, argumentsNode)
		}

		lookupVariable(node.ChildByFieldName("right"))
	// foo["bar"]
	case "element_reference":
		objectNode := node.ChildByFieldName("object")
		builder.Dataflow(node, objectNode)
		lookupVariable(objectNode)
		lookupVariable(objectNode)
	// case foo
	// ...
	// end
	case "case":
		if valueNode := node.ChildByFieldName("value"); valueNode != nil {
			builder.Dataflow(
				node,
				builder.ChildrenExcept(node, valueNode)...,
			)

			break
		}

		builder.Dataflow(node, builder.ChildrenFor(node)...)
	// case foo
	// when 1
	// end
	case "when":
		if patternNode := node.ChildByFieldName("pattern"); patternNode != nil {
			builder.Dataflow(
				node,
				builder.ChildrenExcept(node, patternNode)...,
			)

			break
		}

		builder.Dataflow(node, builder.ChildrenFor(node)...)
	// if x
	//   expr...
	// end
	// case/if
	//   ...
	// else
	//   expr...
	// end
	case "then", "else":
		if lastChild := builder.LastChild(node); lastChild != nil {
			builder.Dataflow(node, lastChild)
		}
	case "keyword_parameter", "optional_parameter":
		nameNode := node.ChildByFieldName("name")
		builder.Dataflow(node, nameNode)

		if nameNode.Type() == "identifier" {
			scope.Declare(builder.ContentFor(nameNode), nameNode)
		}
	case "method_parameters", "block_parameters":
		children := builder.ChildrenFor(node)
		builder.Dataflow(node, children...)

		for _, child := range children {
			if child.Type() == "identifier" {
				scope.Declare(builder.ContentFor(child), child)
			}
		}
	case "pair", "argument_list":
		children := builder.ChildrenFor(node)
		builder.Dataflow(node, children...)

		for _, child := range children {
			lookupVariable(child)
		}
	case "interpolation", "array", "binary", "unary":
		builder.SetOperation(node)

		children := builder.ChildrenFor(node)
		builder.Dataflow(node, children...)

		for _, child := range children {
			lookupVariable(child)
		}
	default:
		builder.Dataflow(node, builder.ChildrenExcept(node, node.ChildByFieldName("condition"))...)
	}

	return visitChildren(scope)
}
