package implementation

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"

	detectortypes "github.com/bearer/bearer/new/detector/types"
	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
)

type Pattern interface {
	// ExtractVariables parses variables from a pattern and returns a new pattern
	// with the variables replaced with a dummy value, along with a list of the
	// variables. Dummy values are needed to allow Tree Sitter to parse the
	// pattern without error.
	ExtractVariables(input string) (string, []patternquerytypes.Variable, error)
	// FixupVariableDummyValue is used to return a new dummy value to use, when
	// the initial parse of a pattern resulted in errors. This can be used in the
	// case where the default dummy value is not valid in the syntax.
	FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string
	// IsRoot returns whether a node should be ignored or be a root of a pattern
	//
	// eg. given a javascript code like this:
	//    const context = {
	//    		email: "foo@domain.com",
	//    }
	//    logger.child(context).info(user.name);
	// if we want to pull both datatypes inside `child()` as well as inside `info()`
	// we want to ignore member_expressions as roots.
	IsRoot(node *tree.Node) bool
	// IsLeaf returns whether the given node should be treated as a leaf, even if
	// it has children
	IsLeaf(node *tree.Node) bool
	// FindUnanchoredPoints returns pairs of start and end offsets for the
	// unanchored points in the input. This is to allow different syntax for
	// specifying the unanchored points in different languages.
	//
	// eg. given a Ruby pattern like this (where `$<...>` means an unanchored point):
	//   some_call($<...>$<DATA_TYPE>$<...>)
	// we would return [[10, 16], [29, 35]]
	FindUnanchoredPoints(input []byte) [][]int
	// IsAnchored returns whether a node in a pattern should be compiled with
	// anchors (`.`) before and after it in the resulting tree sitter query
	//
	// eg. given a Ruby pattern like this:
	//   some_call($<ARG>) do
	//     other_call
	//   end
	// it is natural for `$<ARG>`` to only match the first argument, but
	// we wouldn't expect `other_call` to be the first expression in the block
	IsAnchored(node *tree.Node) (bool, bool)
	// AnonymousParentTypes returns a list of node types for which anonymous
	// children should be matched against. Generally, we don't want to match
	// anonymous nodes as they make the pattern too restrictive.
	//
	// eg. given Ruby code like this:
	//   a == b
	// you will get a tree like this (where nodes in `"` are anonymous):
	//   (binary (identifier) "==" (identifier))
	// If we don't match the "==" then the pattern would also incorrectly match:
	//   a != b
	AnonymousParentTypes() []string
	// NodeTypes returns the types to use for a given node. This allows us
	// to match using equivalent syntax without having to enumerate all the
	// combinations in rules.
	//
	// eg. given a Ruby pattern like this:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)
	// we want to match both of these code examples, despite differences in the
	// way they parse:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)
	//   call(:verify_mode => OpenSSL::SSL::VERIFY_NONE)
	NodeTypes(node *tree.Node) []string
	// LeafContentTypes returns all the leaf node types which should be matched
	// on their content. eg. strings literals will match their literal values
	LeafContentTypes() []string
	// TranslateContent converts the content of a pattern node to a different
	// type. This is used when NodeTypes returns multiple types for a leaf node.
	//
	// eg. given the situation described in the comment for NodeTypes, we must
	// match against the following content for the symbol:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)    -> verify_mode
	//   call(:verify_mode => OpenSSL::SSL::VERIFY_NONE) -> :verify_mode
	TranslateContent(fromNodeType, toNodeType, content string) string
	// FindMatchNode returns pairs of start and end offsets for the pattern match
	// node. This is to allow different syntax for specifying the match node in
	// different languages. There can only be one match node in a pattern, but
	// multiple are supported here to avoid implementing the error handling in
	// each language.
	//
	// eg. given a Ruby pattern like this (where `$<!>` means the match node)
	//   some_call($<!>$<VAR>)
	// we would return `[[10, 14]]`
	FindMatchNode(input []byte) [][]int
	// ContainerTypes returns a list of node types from which a match node should
	// not be able to escape. There can be multiple nodes in the tree at the same
	// character position, and we want to allow a match node to be the highest
	// position node, terminating at a container node.
	//
	// eg. given the following Ruby pattern:
	//   some_call($<!>key: value)
	// the match node is initially parsed at the `key` node. We want to allow it to
	// expand up to the pair node `key: value`, but not into the argument list. ie.
	// given the following Ruby code matching the pattern:
	//   some_call key: value, other_key: value2
	// we want the content of the match to be `key: value` and not `key: value, other_key: value2`
	ContainerTypes() []string
}

type Implementation interface {
	Name() string
	EnryLanguages() []string
	NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector
	SitterLanguage() *sitter.Language
	// AnalyzeTree populates the dataflow analysis
	AnalyzeTree(ctx context.Context, rootNode *sitter.Node, astBuilder *tree.Builder) error
	Pattern() Pattern
	PassthroughNested(node *tree.Node) bool
	ContributesToResult(node *tree.Node) bool
}

type PatternBase struct{}

func (*PatternBase) IsLeaf(node *tree.Node) bool {
	return false
}

func (*PatternBase) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*PatternBase) IsRoot(node *tree.Node) bool {
	return true
}

func (*PatternBase) ShouldSkipNode(node *tree.Node) bool {
	return false
}

func (*PatternBase) ContainerTypes() []string {
	return nil
}

func (*PatternBase) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	return dummyValue
}

func (*PatternBase) AnonymousParentTypes() []string {
	return nil
}

type Scope struct {
	parent    *Scope
	variables map[string]*tree.Node
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]*tree.Node),
	}
}

func (scope *Scope) Declare(name string, node *tree.Node) {
	scope.variables[name] = node
}

func (scope *Scope) Assign(name string, node *tree.Node) {
	targetScope := scope
	if _, declarationScope := scope.lookupWithScope(name); declarationScope != nil {
		targetScope = declarationScope
	}

	targetScope.variables[name] = node
}

func (scope *Scope) Lookup(name string) *tree.Node {
	node, _ := scope.lookupWithScope(name)
	return node
}

func (scope *Scope) lookupWithScope(name string) (*tree.Node, *Scope) {
	if node, exists := scope.variables[name]; exists {
		return node, scope
	}

	if scope.parent != nil {
		return scope.parent.lookupWithScope(name)
	}

	return nil, nil
}
