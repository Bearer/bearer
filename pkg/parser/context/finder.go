package context

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/variables"
	sitter "github.com/smacker/go-tree-sitter"
)

type FinderRequest struct {
	ContextKeywords  []string
	VariableResolver func(*parser.Node) *variables.Variable
	Tree             *parser.Tree
}

type Finder struct {
	RootContext *Context
	Request     *FinderRequest
}

type Context struct {
	ParentContext  *Context
	ChildContext   []*Context
	LocalVariables map[string]*variables.Variable
	Node           *sitter.Node
}

func NewFinder(req *FinderRequest) *Finder {
	return &Finder{
		Request: req,
		RootContext: &Context{
			Node:           req.Tree.Sitter().RootNode(),
			LocalVariables: make(map[string]*variables.Variable),
		},
	}
}

func (finder *Finder) Find() {
	finder.walk(finder.RootContext.Node, finder.RootContext)
}

func (finder *Finder) walk(node *sitter.Node, currentContext *Context) {
	if node.IsNamed() {
		result := finder.Request.VariableResolver(finder.Request.Tree.Wrap(node))
		if result != nil {
			currentContext.LocalVariables[result.Name] = result
		}
	}

	nextContext := currentContext

	contextBreak := false
	for _, contextKeyword := range finder.Request.ContextKeywords {
		if contextKeyword == node.Type() {
			contextBreak = true
		}
	}

	if contextBreak {
		nextContext = &Context{
			ParentContext:  currentContext,
			Node:           node,
			LocalVariables: make(map[string]*variables.Variable),
		}
		currentContext.ChildContext = append(currentContext.ChildContext, nextContext)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		finder.walk(node.Child(i), nextContext)
	}
}

func (finder *Finder) ToNodeMap() map[parser.NodeID]*Context {
	values := make(map[parser.NodeID]*Context)
	call := func(ctx *Context) {
		values[ctx.Node] = ctx
	}

	walkContext(finder.RootContext, call)

	return values
}

func (finder *Finder) ToResolver() *Resolver {
	nodeMap := finder.ToNodeMap()
	return NewResolver(nodeMap)
}

func (finder *Finder) CtxExport() (returned []*Context) {
	values := finder.ToNodeMap()
	for _, v := range values {
		returned = append(returned, v)
	}

	return
}

func walkContext(ctx *Context, call func(ctx *Context)) {
	call(ctx)

	for i := 0; i < len(ctx.ChildContext); i++ {
		walkContext(ctx.ChildContext[i], call)
	}
}
