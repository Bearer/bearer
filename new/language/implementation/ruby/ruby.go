package ruby

import (
	"context"

	"golang.org/x/exp/slices"

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

var (
	variableLookupParents = []string{"pair", "argument_list", "interpolation", "array", "binary", "operator_assignment"}

	passthroughMethods = []string{"JSON.parse", "JSON.parse!", "*.to_json"}
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

func (*rubyImplementation) AnalyzeTree(ctx context.Context, rootNode *sitter.Node, builder *tree.Builder) error {
	return analyzeNode(ctx, rootNode, builder, implementation.NewScope(nil))
}

func analyzeNode(ctx context.Context, node *sitter.Node, builder *tree.Builder, scope *implementation.Scope) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	parent := node.Parent()
	if parent != nil && contributesToResult(node) {
		builder.Dataflow(parent, node)
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

	switch node.Type() {
	case "method":
		return visitChildren(implementation.NewScope(nil))
	case "assignment":
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")

		if left.Type() == "identifier" {
			err := visitChildren(scope)

			scope.Assign(builder.ContentFor(left), node)
			builder.Dataflow(node, right)

			return err
		}
	// x += y
	case "operator_assignment":
		err := visitChildren(scope)

		left := node.ChildByFieldName("left")
		if left.Type() == "identifier" {
			scope.Assign(builder.ContentFor(left), node)
		}

		return err
	case "identifier":
		if parent == nil {
			break
		}

		if slices.Contains(variableLookupParents, parent.Type()) ||
			(parent.Type() == "assignment" && node == parent.ChildByFieldName("right")) ||
			(parent.Type() == "call" && node == parent.ChildByFieldName("receiver")) ||
			(parent.Type() == "element_reference" && node == parent.ChildByFieldName("object")) {
			if scopedNode := scope.Lookup(builder.ContentFor(node)); scopedNode != nil {
				builder.Dataflow(node, scopedNode)
			}
		}

		if parent.Type() == "method_parameters" ||
			parent.Type() == "block_parameters" ||
			(parent.Type() == "keyword_parameter" && node == parent.ChildByFieldName("name")) ||
			(parent.Type() == "optional_parameter" && node == parent.ChildByFieldName("name")) {
			scope.Declare(builder.ContentFor(node), node)
		}

		if parent.Type() == "argument_list" {
			callNode := parent.Parent()
			builder.Dataflow(callNode, node)
		}
	case "block", "do_block":
		return visitChildren(implementation.NewScope(scope))
	}

	return visitChildren(scope)
}

func (implementation *rubyImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*rubyImplementation) PassthroughNested(node *tree.Node) bool {
	callNode := node.Parent()
	if callNode.Type() != "call" {
		return false
	}

	receiverNode := callNode.ChildByFieldName("receiver")

	if node.Type() != "arguments_list" && (receiverNode == nil || node != receiverNode) {
		return false
	}

	var receiverMethod string
	var wildcardMethod string

	if receiverNode != nil {
		methodName := callNode.ChildByFieldName("method").Content()

		if receiverNode.Type() == "identifier" {
			receiverMethod = receiverNode.Content() + "." + methodName
		}

		wildcardMethod = "*." + methodName
	}

	return slices.Contains(passthroughMethods, receiverMethod) || slices.Contains(passthroughMethods, wildcardMethod)
}

func contributesToResult(node *sitter.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Must not be a condition
	if node == parent.ChildByFieldName("condition") {
		return false
	}

	// Must not be a case value
	if parent.Type() == "case" && node == parent.ChildByFieldName("value") {
		return false
	}

	// Must not be a case-when pattern
	if parent.Type() == "when" && node == parent.ChildByFieldName("pattern") {
		return false
	}

	// Not the left part of an assignment
	if parent.Type() == "assignment" && node == parent.ChildByFieldName("left") {
		return false
	}

	// Must be the last expression in an expression block
	if slices.Contains([]string{"then", "else"}, parent.Type()) {
		if node != parent.Child(int(parent.ChildCount())-1) {
			return false
		}
	}

	return true
}
