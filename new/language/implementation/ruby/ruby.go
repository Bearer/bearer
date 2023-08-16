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
	return nil
	// scope := implementation.NewScope(nil)

	// return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
	// 	if ctx.Err() != nil {
	// 		return ctx.Err()
	// 	}

	// 	switch node.Type() {
	// 	case "method":
	// 		scope = implementation.NewScope(nil)
	// 	case "assignment":
	// 		left := node.ChildByFieldName("left")
	// 		right := node.ChildByFieldName("right")

	// 		if left.Type() == "identifier" {
	// 			err := visitChildren()

	// 			scope.Assign(left.Content(), node)
	// 			node.UnifyWith(right)

	// 			return err
	// 		}
	// 	// x += y
	// 	case "operator_assignment":
	// 		err := visitChildren()

	// 		left := node.ChildByFieldName("left")
	// 		if left.Type() == "identifier" {
	// 			scope.Assign(left.Content(), node)
	// 		}

	// 		return err
	// 	case "identifier":
	// 		parent := node.Parent()
	// 		if parent == nil {
	// 			break
	// 		}

	// 		if slices.Contains(variableLookupParents, parent.Type()) ||
	// 			(parent.Type() == "assignment" && node.Equal(parent.ChildByFieldName("right"))) ||
	// 			(parent.Type() == "call" && node.Equal(parent.ChildByFieldName("receiver"))) ||
	// 			(parent.Type() == "element_reference" && node.Equal(parent.ChildByFieldName("object"))) {
	// 			if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
	// 				node.UnifyWith(scopedNode)
	// 			}
	// 		}

	// 		if parent.Type() == "method_parameters" ||
	// 			parent.Type() == "block_parameters" ||
	// 			(parent.Type() == "keyword_parameter" && node.Equal(parent.ChildByFieldName("name"))) ||
	// 			(parent.Type() == "optional_parameter" && node.Equal(parent.ChildByFieldName("name"))) {
	// 			scope.Declare(node.Content(), node)
	// 		}

	// 		if parent.Type() == "argument_list" {
	// 			callNode := parent.Parent()
	// 			callNode.UnifyWith(node)
	// 		}
	// 	case "block", "do_block":
	// 		previousScope := scope
	// 		scope = implementation.NewScope(scope)
	// 		err := visitChildren()
	// 		scope = previousScope
	// 		return err
	// 	}

	// 	return visitChildren()
	// })
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

func (*rubyImplementation) ContributesToResult(node *tree.Node) bool {
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
		parentChildren := parent.Children()
		if node != parentChildren[len(parentChildren)-1] {
			return false
		}
	}

	return true
}
