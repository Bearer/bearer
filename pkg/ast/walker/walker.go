package walker

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

type Walker struct {
	query *sitter.Query
}

type VisitFunction = func(node *sitter.Node, visitChildren func() error) error

type Cursor struct {
	queryCursor *sitter.QueryCursor
	peekCache   *sitter.Node
	onVisit     VisitFunction
}

func NewWalker(language *sitter.Language) *Walker {
	query, err := sitter.NewQuery([]byte("(_) @node"), ruby.GetLanguage())
	if err != nil {
		panic(err)
	}

	return &Walker{query: query}
}

func (walker *Walker) Walk(rootNode *sitter.Node, onVisit VisitFunction) error {
	queryCursor := sitter.NewQueryCursor()
	defer queryCursor.Close()

	queryCursor.Exec(walker.query, rootNode)

	cursor := &Cursor{queryCursor: queryCursor, onVisit: onVisit}

	// Visit the root
	return cursor.visit(cursor.peek())
}

func (walker *Walker) Close() {
	walker.query.Close()
}

func (cursor *Cursor) visit(node *sitter.Node) error {
	cursor.accept()

	visitedChildren := false

	if err := cursor.onVisit(node, func() error {
		visitedChildren = true
		return cursor.visitChildren(node)
	}); err != nil {
		return err
	}

	if !visitedChildren {
		cursor.skipDescendants(node)
	}

	return nil
}

func (cursor *Cursor) visitChildren(node *sitter.Node) error {
	for {
		next := cursor.peek()
		if next == nil || !node.Equal(next.Parent()) {
			return nil
		}

		if err := cursor.visit(next); err != nil {
			return err
		}
	}
}

func (cursor *Cursor) skipDescendants(node *sitter.Node) {
	for {
		if !cursor.isAncestorOfNext(node) {
			return
		}

		cursor.accept()
	}
}

func (cursor *Cursor) peek() *sitter.Node {
	if cursor.peekCache != nil {
		return cursor.peekCache
	}

	match, exists := cursor.queryCursor.NextMatch()
	if !exists {
		return nil
	}

	cursor.peekCache = match.Captures[0].Node
	return cursor.peekCache
}

func (cursor *Cursor) accept() {
	cursor.peekCache = nil
}

func (cursor *Cursor) isAncestorOfNext(ancestor *sitter.Node) bool {
	node := cursor.peek()
	if node == nil {
		return false
	}

	for {
		parent := node.Parent()
		if parent == nil {
			return false
		}

		if parent.Equal(ancestor) {
			return true
		}

		node = parent
	}
}
