package implementation

import (
	sitter "github.com/smacker/go-tree-sitter"

	patternquerytypes "github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/new/language/tree"
)

type Implementation interface {
	SitterLanguage() *sitter.Language
	// AnalyzeFlow unifies nodes that represent the same value in the tree.
	//
	// eg. given Ruby code like this:
	//   user = { first_name: "" }
	//   some_call(user)
	//   user[:first_name]
	// the `user` identifier node on lines 2 and 3 will be unified with the
	// assignment node
	AnalyzeFlow(rootNode *tree.Node) error
	// ExtractPatternVariables parses variables from a pattern and returns a new
	// pattern with the variables replaced with a dummy value, along with a list
	// of the variables. Dummy values are needed to allow Tree Sitter to parse
	// the pattern without error.
	ExtractPatternVariables(input string) (string, []patternquerytypes.Variable, error)
	// FindPatternMatchNode returns pairs of start and end offsets for the
	// pattern match node. This is to allow different syntax for specifying the
	// match node in different languages. There can only be one match node in a
	// pattern, but multiple are supported here to avoid implementing the error
	// handling in each language.
	//
	// eg. given a Ruby pattern like this (where `$<!>` means the match node)
	//   some_call($<!>$<VAR>)
	// we would return `[[10, 14]]`
	FindPatternMatchNode(input []byte) [][]int
	// FindPatternUnanchoredPoints returns pairs of start and end offsets for the
	// unanchored points in the input. This is to allow different syntax for
	// specifying the unanchored points in different languages.
	//
	// eg. given a Ruby pattern like this (where `$<...>` means an unanchored point):
	//   some_call($<...>$<DATA_TYPE>$<...>)
	// we would return [[10, 16], [29, 35]]
	FindPatternUnanchoredPoints(input []byte) [][]int
	// AnonymousPatternNodeParentTypes returns a list of node types for which
	// anonymous children should be matched against. Generally, we don't want to
	// match anonymous nodes as they make the pattern too restrictive.
	//
	// eg. given Ruby code like this:
	//   a == b
	// you will get a tree like this (where nodes in `"` are anonymous):
	//   (binary (identifier) "==" (identifier))
	// If we don't match the "==" then the pattern would also incorrectly match:
	//   a != b
	AnonymousPatternNodeParentTypes() []string
	// PatternMatchNodeContainerTypes returns a list of node types from which a
	// match node should not be able to escape. There can be multiple nodes in the
	// tree at the same character position, and we want to allow a match node to
	// be the highest position node, terminating at a container node.
	//
	// eg. given the following Ruby pattern:
	//   some_call($<!>key: value)
	// the match node is initially parsed at the `key` node. We want to allow it to
	// expand up to the pair node `key: value`, but not into the argument list. ie.
	// given the following Ruby code matching the pattern:
	//   some_call key: value, other_key: value2
	// we want the content of the match to be `key: value` and not `key: value, other_key: value2`
	PatternMatchNodeContainerTypes() []string
	// PatternIsAnchored returns whether a node in a pattern should be compiled
	// with anchors (`.`) before and after it in the resulting tree sitter query
	//
	// eg. given a Ruby pattern like this:
	//   some_call($<ARG>) do
	//     other_call
	//   end
	// it is natural for `$<ARG>`` to only match the first argument, but
	// we wouldn't expect `other_call` to be the first expression in the block
	PatternIsAnchored(node *tree.Node) (bool, bool)
	// PatternNodeTypes returns the types to use for a given node. This allows us
	// to match using equivalent syntax without having to enumerate all the
	// combinations in rules.
	//
	// eg. given a Ruby pattern like this:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)
	// we want to match both of these code examples, despite differences in the
	// way they parse:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)
	//   call(:verify_mode => OpenSSL::SSL::VERIFY_NONE)
	PatternNodeTypes(node *tree.Node) []string
	// TranslatePatternContent converts the content of a pattern node to a
	// different type. This is used when PatternNodeTypes returns multiple types
	// for a leaf node.
	//
	// eg. given the situation described in the comment for PatternNodeTypes, we
	// must match against the following content for the symbol:
	//   call(verify_mode: OpenSSL::SSL::VERIFY_NONE)    -> verify_mode
	//   call(:verify_mode => OpenSSL::SSL::VERIFY_NONE) -> :verify_mode
	TranslatePatternContent(fromNodeType, toNodeType, content string) string
	// IsRootOfRuleQuery returns whether a node should be ignored or be a root
	// of a custom rule query
	//
	// eg. given a javascript code like this:
	//    const context = {
	//    		email: "foo@domain.com",
	//    }
	//    logger.child(context).info(user.name);
	// if we want to pull both datatypes inside `child()` as well as inside `info()`
	// we want to ignore member_expressions as roots.
	IsRootOfRuleQuery(node *tree.Node) bool
	PatternLeafContentTypes() []string
	// ShouldSkipNode returns wether a node should should be skipped or assigned variable to it
	// it is useful for cases when you have nested nodes to ignore and want to assign variable to its child
	//
	// eg. given the following tree sitter
	// arguments
	//	formal_parameters
	//		required_parameter
	//			identifier
	//
	//  if you want to get only identifier instead of required parameter ShouldSkipNode should return true
	//  for required parameter
	ShouldSkipNode(node *tree.Node) bool

	PassthroughNested(node *tree.Node) bool

	ContributesToValue(node *tree.Node) bool
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

func (scope *Scope) Assign(name string, node *tree.Node) {
	scope.variables[name] = node
}

func (scope *Scope) Lookup(name string) *tree.Node {
	if node, exists := scope.variables[name]; exists {
		return node
	}

	if scope.parent != nil {
		return scope.parent.Lookup(name)
	}

	return nil
}
