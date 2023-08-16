package java

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/java/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/java/string"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"
)

var variableLookupParents = []string{
	"argument_list",
	"array_access",
	"array_initializer",
	"binary_expression",
	"field_declaration",
	"ternary_expression",
}

type javaImplementation struct {
	pattern patternImplementation
}

func Get() implementation.Implementation {
	return &javaImplementation{}
}

func (*javaImplementation) Name() string {
	return "java"
}

func (*javaImplementation) EnryLanguages() []string {
	return []string{"Java"}
}

func (*javaImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJava, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*javaImplementation) SitterLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (*javaImplementation) AnalyzeTree(ctx context.Context, rootNode *sitter.Node, builder *tree.Builder) error {
	return analyzeNode(ctx, rootNode, builder, implementation.NewScope(nil))
}

func analyzeNode(ctx context.Context, node *sitter.Node, builder *tree.Builder, scope *implementation.Scope) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	parent := node.Parent()
	if parent != nil && contributesToResult(builder, node) {
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
	// public class Main {
	//
	// }
	case "class_body":
		return visitChildren(implementation.NewScope(scope))
	// public class Main {
	//	// method declaration
	//	static void myMethod() {
	//
	//   }
	// }
	//
	// lambda_expression
	// numbers.forEach( (n) -> { System.out.println(n); } );
	case "method_declaration", "lambda_expression", "for_statement", "enhanced_for_statement", "block":
		return visitChildren(implementation.NewScope(scope))
	// user = ...
	case "assignment_expression":
		err := visitChildren(scope)

		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")

		if builder.ContentFor(node.Child(1)) == "=" {
			builder.Dataflow(node, right)
		}

		if left.Type() == "identifier" {
			scope.Assign(builder.ContentFor(left), node)
		}

		return err
	case "field_declaration":
		err := visitChildren(scope)

		declarator := node.ChildByFieldName("declarator")
		if declarator != nil {
			scope.Declare(builder.ContentFor(declarator.ChildByFieldName("name")), node)

			if value := declarator.ChildByFieldName("value"); value != nil {
				builder.Dataflow(node, value)
			}
		}

		return err
	// String user = "John";
	case "local_variable_declaration":
		declarator := node.ChildByFieldName("declarator")

		name := declarator.ChildByFieldName("name")
		value := declarator.ChildByFieldName("value")

		if name.Type() == "identifier" {
			err := visitChildren(scope)

			scope.Declare(builder.ContentFor(name), node)
			builder.Dataflow(node, value)

			return err
		}
	// // TODO: figure out this one
	// case "shorthand_property_identifier_pattern":
	// 	scope.Assign(node.Content(), node)
	case "identifier":
		if parent == nil {
			break
		}

		if slices.Contains(variableLookupParents, parent.Type()) ||
			(parent.Type() == "scoped_type_identifier" && node == parent.Child(0)) ||
			(parent.Type() == "method_invocation" && node == parent.ChildByFieldName("object")) ||
			(parent.Type() == "field_access" && node == parent.ChildByFieldName("object")) ||
			(parent.Type() == "variable_declarator" && node == parent.ChildByFieldName("value")) ||
			(parent.Type() == "assignment_expression" && node == parent.ChildByFieldName("right")) ||
			(parent.Type() == "assignment_expression" && node == parent.ChildByFieldName("left") && builder.ContentFor(parent.Child(1)) != "=") ||
			(parent.Type() == "enhanced_for_statement" && node == parent.ChildByFieldName("value")) {
			if scopedNode := scope.Lookup(builder.ContentFor(node)); scopedNode != nil {
				builder.Dataflow(node, scopedNode)
			}
		}

		if parent.Type() == "formal_parameter" ||
			parent.Type() == "catch_formal_parameter" ||
			(parent.Type() == "resource" && node == parent.ChildByFieldName("name")) {
			scope.Declare(builder.ContentFor(node), node)
		}

		if parent.Type() == "enhanced_for_statement" && node == parent.ChildByFieldName("name") {
			scope.Declare(builder.ContentFor(node), node)
			builder.Dataflow(node, parent.ChildByFieldName("value"))
		}

		// todo: see what this is
		// case "property_identifier":
		// 	parent := node.Parent()
		// 	if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
		// 		if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
		// 			node.UnifyWith(scopedNode)
		// 		}
		// 	}
		// }
	}

	return visitChildren(scope)
}

func (implementation *javaImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

// func (*javaImplementation) PassthroughNested(node *tree.Node) bool {
// 	return false
// }

func contributesToResult(builder *tree.Builder, node *sitter.Node) bool {
	// Statements don't have results
	if strings.HasSuffix(node.Type(), "_statement") {
		return false
	}

	// Switch case
	if node.Type() == "switch_label" {
		return false
	}

	parent := node.Parent()
	if parent == nil {
		return true
	}

	// Must not be a ternary/switch condition
	if node == parent.ChildByFieldName("condition") {
		return false
	}

	// Not the name part of a declaration
	if parent.Type() == "variable_declarator" && node == parent.ChildByFieldName("name") {
		return false
	}

	// Not the left part of an `=` assignment
	if parent.Type() == "assignment_expression" && node == parent.ChildByFieldName("left") {
		return builder.ContentFor(parent.Child(1)) != "="
	}

	return true
}
