package java

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"

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
	return nil
	// scope := implementation.NewScope(nil)

	// return rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
	// 	if ctx.Err() != nil {
	// 		return ctx.Err()
	// 	}

	// 	switch node.Type() {
	// 	// public class Main {
	// 	//
	// 	// }
	// 	case "class_body":
	// 		previousScope := scope
	// 		scope = implementation.NewScope(previousScope)
	// 		err := visitChildren()
	// 		scope = previousScope
	// 		return err
	// 	// public class Main {
	// 	//	// method declaration
	// 	//	static void myMethod() {
	// 	//
	// 	//   }
	// 	// }
	// 	//
	// 	// lambda_expression
	// 	// numbers.forEach( (n) -> { System.out.println(n); } );
	// 	case "method_declaration", "lambda_expression", "for_statement", "enhanced_for_statement", "block":
	// 		previousScope := scope
	// 		scope = implementation.NewScope(previousScope)
	// 		err := visitChildren()
	// 		scope = previousScope
	// 		return err
	// 	// user = ...
	// 	case "assignment_expression":
	// 		err := visitChildren()

	// 		left := node.ChildByFieldName("left")
	// 		right := node.ChildByFieldName("right")

	// 		if node.AnonymousChild(0).Content() == "=" {
	// 			node.UnifyWith(right)
	// 		}

	// 		if left.Type() == "identifier" {
	// 			scope.Assign(left.Content(), node)
	// 		}

	// 		return err
	// 	case "field_declaration":
	// 		err := visitChildren()

	// 		declarator := node.ChildByFieldName("declarator")
	// 		if declarator != nil {
	// 			scope.Declare(declarator.ChildByFieldName("name").Content(), node)

	// 			if value := declarator.ChildByFieldName("value"); value != nil {
	// 				node.UnifyWith(value)
	// 			}
	// 		}

	// 		return err
	// 	// String user = "John";
	// 	case "local_variable_declaration":
	// 		declarator := node.ChildByFieldName("declarator")

	// 		name := declarator.ChildByFieldName("name")
	// 		value := declarator.ChildByFieldName("value")

	// 		if name.Type() == "identifier" {
	// 			err := visitChildren()

	// 			scope.Declare(name.Content(), node)
	// 			node.UnifyWith(value)

	// 			return err
	// 		}
	// 	// // TODO: figure out this one
	// 	// case "shorthand_property_identifier_pattern":
	// 	// 	scope.Assign(node.Content(), node)
	// 	case "identifier":
	// 		parent := node.Parent()
	// 		if parent == nil {
	// 			break
	// 		}

	// 		if slices.Contains(variableLookupParents, parent.Type()) ||
	// 			(parent.Type() == "scoped_type_identifier" && node.Equal(parent.Child(0))) ||
	// 			(parent.Type() == "method_invocation" && node.Equal(parent.ChildByFieldName("object"))) ||
	// 			(parent.Type() == "field_access" && node.Equal(parent.ChildByFieldName("object"))) ||
	// 			(parent.Type() == "variable_declarator" && node.Equal(parent.ChildByFieldName("value"))) ||
	// 			(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("right"))) ||
	// 			(parent.Type() == "assignment_expression" && node.Equal(parent.ChildByFieldName("left")) && parent.AnonymousChild(0).Content() != "=") ||
	// 			(parent.Type() == "enhanced_for_statement" && node.Equal(parent.ChildByFieldName("value"))) {
	// 			if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
	// 				node.UnifyWith(scopedNode)
	// 			}
	// 		}

	// 		if parent.Type() == "formal_parameter" ||
	// 			parent.Type() == "catch_formal_parameter" ||
	// 			(parent.Type() == "resource" && node.Equal(parent.ChildByFieldName("name"))) {
	// 			scope.Declare(node.Content(), node)
	// 		}

	// 		if parent.Type() == "enhanced_for_statement" && node.Equal(parent.ChildByFieldName("name")) {
	// 			scope.Declare(node.Content(), node)
	// 			node.UnifyWith(parent.ChildByFieldName("value"))
	// 		}

	// 		// todo: see what this is
	// 		// case "property_identifier":
	// 		// 	parent := node.Parent()
	// 		// 	if parent != nil && slice.Contains(variableLookupParents, parent.Type()) {
	// 		// 		if scopedNode := scope.Lookup(node.Content()); scopedNode != nil {
	// 		// 			node.UnifyWith(scopedNode)
	// 		// 		}
	// 		// 	}
	// 		// }
	// 	}
	// 	return visitChildren()
	// })
}

func (implementation *javaImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*javaImplementation) PassthroughNested(node *tree.Node) bool {
	return false
}

func (*javaImplementation) ContributesToResult(node *tree.Node) bool {
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
		return parent.Children()[1].Content() != "="
	}

	return true
}
