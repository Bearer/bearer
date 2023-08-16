package javascript

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	"github.com/ssoroka/slice"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/javascript/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/javascript/string"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"
)

var (
	variableLookupParents = []string{
		"pair",
		"arguments",
		"binary_expression",
		"template_substitution",
		"array",
		"spread_element",
		"augmented_assignment_expression",
	}

	passthroughMethods = []string{"JSON.parse", "JSON.stringify"}
)

type javascriptImplementation struct {
	pattern patternImplementation
}

func Get() implementation.Implementation {
	return &javascriptImplementation{}
}

func (*javascriptImplementation) Name() string {
	return "javascript"
}

func (*javascriptImplementation) EnryLanguages() []string {
	return []string{"JavaScript", "TypeScript", "TSX"}
}

func (*javascriptImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJavascript, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*javascriptImplementation) SitterLanguage() *sitter.Language {
	return tsx.GetLanguage()
}

func (*javascriptImplementation) AnalyzeTree(ctx context.Context, rootNode *sitter.Node, builder *tree.Builder) error {
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
	// () => {}
	// function getName() {}
	case "function", "arrow_function", "method_definition":
		return visitChildren(implementation.NewScope(scope))
	// user = ...
	case "assignment_expression":
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")

		if left.Type() == "identifier" {
			err := visitChildren(scope)

			scope.Assign(builder.ContentFor(left), node)
			builder.Dataflow(node, right)

			return err
		}
	// x += y
	case "augmented_assignment_expression":
		err := visitChildren(scope)

		left := node.ChildByFieldName("left")
		if left.Type() == "identifier" {
			scope.Assign(builder.ContentFor(left), node)
		}

		return err
	// const user = ...
	// var user = ...
	// let user = ...
	case "variable_declarator":
		name := node.ChildByFieldName("name")
		value := node.ChildByFieldName("value")

		if name.Type() == "identifier" {
			err := visitChildren(scope)

			scope.Declare(builder.ContentFor(name), node)
			builder.Dataflow(node, value)

			return err
		}
	case "shorthand_property_identifier_pattern":
		scope.Declare(builder.ContentFor(node), node)
	case "identifier":
		if parent == nil {
			break
		}

		if slice.Contains(variableLookupParents, parent.Type()) ||
			(parent.Type() == "assignment_expression" && node == parent.ChildByFieldName("right")) ||
			(parent.Type() == "new_expression" && node == parent.ChildByFieldName("constructor")) ||
			(parent.Type() == "variable_declarator" && node == parent.ChildByFieldName("value")) ||
			(parent.Type() == "member_expression" && node == parent.ChildByFieldName("object")) ||
			(parent.Type() == "call_expression" && node == parent.ChildByFieldName("function")) ||
			(parent.Type() == "subscript_expression" && node == parent.ChildByFieldName("object")) {
			if scopedNode := scope.Lookup(builder.ContentFor(node)); scopedNode != nil {
				builder.Dataflow(node, scopedNode)
			}

			break
		}

		// typescript: different type of identifier
		if parent.Type() == "required_parameter" {
			scope.Declare(builder.ContentFor(node), node)
			break
		}

		if parent.Type() == "arguments" {
			callNode := parent.Parent()
			builder.Dataflow(callNode, node)
			break
		}

		if isImportedIdentifier(node) {
			scope.Declare(builder.ContentFor(node), node)
		}
	case "property_identifier":
		if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
			if scopedNode := scope.Lookup(builder.ContentFor(node)); scopedNode != nil {
				builder.Dataflow(node, scopedNode)
			}
		}
	}

	return visitChildren(scope)
}

func (implementation *javascriptImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

// func (*javascriptImplementation) PassthroughNested(node *tree.Node) bool {
// 	if node.Type() != "arguments" {
// 		return false
// 	}

// 	callNode := node.Parent()
// 	if callNode.Type() != "call_expression" {
// 		return false
// 	}

// 	functionNode := callNode.ChildByFieldName("function")

// 	var method string
// 	var wildcardMethod string
// 	switch functionNode.Type() {
// 	case "identifier":
// 		return slices.Contains(passthroughMethods, functionNode.Content())
// 	case "member_expression":
// 		object := functionNode.ChildByFieldName("object")
// 		if object.Type() == "identifier" {
// 			property := functionNode.ChildByFieldName("property").Content()
// 			method = object.Content() + "." + property
// 			wildcardMethod = "*." + property
// 		}
// 	}

// 	return slices.Contains(passthroughMethods, method) || slices.Contains(passthroughMethods, wildcardMethod)
// }

func contributesToResult(node *sitter.Node) bool {
	// Statements don't have results
	if strings.HasSuffix(node.Type(), "_statement") {
		return false
	}

	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Must not be a ternary condition
	if parent.Type() == "ternary_expression" && node == parent.ChildByFieldName("condition") {
		return false
	}

	// Not the name part of a declaration
	if parent.Type() == "variable_declarator" && node == parent.ChildByFieldName("name") {
		return false
	}

	// Not the left part of an assignment
	if parent.Type() == "assignment_expression" && node == parent.ChildByFieldName("left") {
		return false
	}

	return true
}

func isImportedIdentifier(node *sitter.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return false
	}

	// import x from "library"
	if parent.Type() == "import_clause" {
		return true
	}

	// import * as x from "library"
	if parent.Type() == "namespace_import" {
		return true
	}

	if parent.Type() != "import_specifier" {
		return false
	}

	// import { x } from "library"
	if parent.ChildByFieldName("alias") == nil {
		return true
	}

	// import { a as x } from "library"
	if node == parent.ChildByFieldName("alias") {
		return true
	}

	return false
}
