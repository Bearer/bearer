package parser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"

	sitter "github.com/smacker/go-tree-sitter"
)

type Tree struct {
	language   *sitter.Language
	fileInfo   *file.FileInfo
	file       *file.Path
	input      []byte
	sitter     *sitter.Tree
	lineOffset int
	values     map[*sitter.Node]*values.Value
}

type Node struct {
	tree   *Tree
	sitter *sitter.Node
}

type NodeID *sitter.Node

type Captures map[string]*Node

func ParseBytes(fileInfo *file.FileInfo, file *file.Path, input []byte, language *sitter.Language, lineOffset int) (*Tree, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(language)

	sitterTree, err := parser.ParseCtx(context.Background(), nil, input)
	if err != nil {
		return nil, err
	}

	return &Tree{
		language:   language,
		fileInfo:   fileInfo,
		file:       file,
		input:      input,
		sitter:     sitterTree,
		lineOffset: lineOffset,
		values:     make(map[*sitter.Node]*values.Value),
	}, nil
}

func ParseFile(fileInfo *file.FileInfo, file *file.Path, language *sitter.Language) (*Tree, error) {
	input, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil, err
	}

	return ParseBytes(fileInfo, file, input, language, 0)
}

func (tree *Tree) Debug() string {
	return tree.sitter.RootNode().String()
}

func (tree *Tree) File() *file.Path {
	return tree.file
}

func (tree *Tree) Close() {
	tree.sitter.Close()
}

func (tree *Tree) Query(query *sitter.Query, onMatch func(captures Captures) error) error {
	return tree.wrap(tree.sitter.RootNode()).Query(query, onMatch)
}

func (tree *Tree) RootNode() *Node {
	return tree.Wrap(tree.sitter.RootNode())
}

func (tree *Tree) Sitter() *sitter.Tree {
	return tree.sitter
}

func (tree *Tree) QueryMustPass(query *sitter.Query) []Captures {
	return tree.wrap(tree.sitter.RootNode()).QueryMustPass(query)
}

func (tree *Tree) SetValues(newValues map[*sitter.Node]*values.Value) {
	tree.values = newValues
}

// QueryConventional obides certain convetions and returns results accordingly
// capture matching helper_{randomString}
//
//	must have content matching {randomString} once quotes are stripped
func (tree *Tree) QueryConventional(query *sitter.Query) []Captures {
	filteredCaptures := []Captures{}
	captures := tree.wrap(tree.sitter.RootNode()).QueryMustPass(query)

	for _, capture := range captures {
		passesFilter := true
		for i := 0; i < int(query.CaptureCount()); i++ {
			name := query.CaptureNameForId(uint32(i))
			if strings.HasPrefix(name, "helper_") {
				wantMatch := strings.TrimPrefix(name, "helper_")
				content := stringutil.StripQuotes(capture[name].Content())
				if content != wantMatch {
					passesFilter = false
					break
				}
			}
		}

		if !passesFilter {
			continue
		}

		filteredCaptures = append(filteredCaptures, capture)
	}

	return filteredCaptures
}

func (tree *Tree) walkBottomUp(onNode func(child *Node) error) error {
	cursor := sitter.NewTreeCursor(tree.sitter.RootNode())
	defer cursor.Close()
	descending := true

	for {
		if descending && cursor.GoToFirstChild() {
			continue
		}

		if cursor.CurrentNode().IsNamed() {
			err := onNode(tree.wrap(cursor.CurrentNode()))
			if err != nil {
				return err
			}
		}

		if cursor.GoToNextSibling() {
			descending = true
			continue
		}

		if cursor.GoToParent() {
			descending = false
			continue
		}

		break
	}

	return nil
}

func (tree *Tree) WalkBottomUp(onNode func(child *Node) error) error {
	return tree.walkBottomUp(onNode)
}

func (tree *Tree) Annotate(populateValue func(node *Node, value *values.Value)) error {
	return tree.walkBottomUp(func(node *Node) error {
		// Already annotated previously (eg. with a specialist annotator)
		if node.Value() != nil {
			return nil
		}

		value := values.New()
		populateValue(node, value)
		node.SetValue(value)

		return nil
	})
}

func (tree *Tree) walkRootValues(onValue func(node *Node)) error {
	return tree.walkBottomUp(func(node *Node) error {
		// Value itself is unknown (not relevant)
		if node.Value().IsUnknown() {
			return nil
		}

		// Parent that isn't unknown (not a root)
		for parent := node.sitter.Parent(); parent != nil; parent = parent.Parent() {
			_, ok := node.tree.values[parent]
			if !ok {
				return nil
			}

			if !node.tree.values[parent].IsUnknown() {
				return nil
			}
		}

		onValue(node)

		return nil
	})
}

func (tree *Tree) WalkRootValues(onValue func(node *Node)) error {
	return tree.walkRootValues(onValue)
}

func (tree *Tree) wrap(node *sitter.Node) *Node {
	if node == nil {
		return nil
	}

	return &Node{tree: tree, sitter: node}
}

func (tree *Tree) Wrap(node *sitter.Node) *Node {
	return tree.wrap(node)
}

func (node *Node) Debug() string {
	return node.sitter.String()
}

func (node *Node) LineNumber() int {
	return int(node.sitter.StartPoint().Row+1) + node.tree.lineOffset
}

func (node *Node) Source(includeText bool) source.Source {
	text := ""
	if includeText {
		text = strings.TrimSpace(node.Content())
	}

	return source.New(
		node.tree.fileInfo,
		node.tree.File(),
		node.LineNumber(),
		int(node.sitter.StartPoint().Column+1),
		text,
	)
}

func (node *Node) Content() string {
	return node.sitter.Content(node.tree.input)
}

func (node *Node) SetValue(value *values.Value) {
	node.tree.values[node.sitter] = value
}

func (node *Node) ID() NodeID {
	return node.sitter
}

func (node *Node) Value() *values.Value {
	return node.tree.values[node.sitter]
}

func (node *Node) Type() string {
	return node.sitter.Type()
}

func (node *Node) IsNamed() bool {
	return node.sitter.IsNamed()
}

func (node *Node) ChildByFieldName(name string) *Node {
	return node.tree.wrap(node.sitter.ChildByFieldName(name))
}

func (node *Node) Parent() *Node {
	return node.tree.wrap(node.sitter.Parent())
}

var ErrParentNotFound = errors.New("parent not found")
var ErrInChain = errors.New("there is an error in chain")

func (node *Node) FindParent(parentType string) (*Node, error) {
	parent := node.Parent()
	for {
		if parent == nil {
			return nil, ErrParentNotFound
		}

		if parent.Type() == parentType {
			return parent, nil
		}

		if parent.Type() == "error" {
			return nil, ErrInChain
		}

		parent = parent.Parent()
	}
}

func (node *Node) FirstChild() *Node {
	return node.tree.wrap(node.sitter.NamedChild(0))
}

func (node *Node) Child(i int) *Node {
	return node.tree.wrap(node.sitter.NamedChild(i))
}

// latest library version 7621c203ae43fe58c0fc4d18e4dcf3caa1985888 supports this but it is buggy so its of no use
// func (node *Node) FieldNameForChild(i int) string {
// 	return node.sitter.FieldNameForChild(i)
// }

func (node *Node) NamedChildCount() int {
	return int(node.sitter.NamedChildCount())
}

func (node *Node) ChildCount() int {
	return int(node.sitter.NamedChildCount())
}

func (node *Node) Sitter() *sitter.Node {
	return node.sitter
}

func (node *Node) Equal(other *Node) bool {
	if other == nil {
		return false
	}

	return node.sitter.Equal(other.sitter)
}

func (node *Node) FirstUnnamedChild() *Node {
	n := int(node.sitter.ChildCount())

	for i := 0; i < n; i++ {
		child := node.sitter.Child(i)
		if !child.IsNamed() {
			return node.tree.wrap(child)
		}
	}

	return nil
}

func (node *Node) ChildValueParts() []values.Part {
	var result []values.Part

	n := int(node.sitter.ChildCount())
	for i := 0; i < n; i++ {
		child := node.sitter.Child(i)
		if child.IsNamed() {
			result = append(result, node.tree.values[child].GetParts()...)
		}
	}

	return result
}

// EachPart calls the supplied functions for each text span or child node.
// Unnamed children (eg. quotes for string literals) are skipped
func (node *Node) EachPart(onText func(text string) error, onChild func(child *Node) error) error {
	n := int(node.sitter.ChildCount())

	start := node.sitter.StartByte()
	end := start

	emit := func() error {
		if end <= start {
			return nil
		}

		return onText(string(node.tree.input[start:end]))
	}

	for i := 0; i < n; i++ {
		child := node.sitter.Child(i)
		end = child.StartByte()

		if err := emit(); err != nil {
			return err
		}

		if child.IsNamed() {
			if err := onChild(node.tree.wrap(child)); err != nil {
				return err
			}
		}

		start = child.EndByte()
		end = start
	}

	if err := emit(); err != nil {
		return err
	}

	return nil
}

func (node *Node) TextParts() []string {
	var result []string
	node.EachPart(func(text string) error { //nolint:all,errcheck
		result = append(result, text)
		return nil
	}, func(node *Node) error {
		return nil
	})

	return result
}

func (node *Node) Query(query *sitter.Query, onMatch func(captures Captures) error) error {
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(query, node.sitter)

	for {
		match, hasMatch := cursor.NextMatch()
		if !hasMatch {
			break
		}

		captures := make(Captures)
		for _, capture := range match.Captures {
			captures[query.CaptureNameForId(capture.Index)] = node.tree.wrap(capture.Node)
		}

		err := onMatch(captures)
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *Node) QueryConventional(query *sitter.Query) []Captures {
	filteredCaptures := []Captures{}
	captures := node.QueryMustPass(query)

	for _, capture := range captures {
		passesFilter := true
		for i := 0; i < int(query.CaptureCount()); i++ {
			name := query.CaptureNameForId(uint32(i))
			if strings.HasPrefix(name, "helper_") {
				wantMatch := strings.TrimPrefix(name, "helper_")
				content := stringutil.StripQuotes(capture[name].Content())
				if content != wantMatch {
					passesFilter = false
					break
				}
			}
		}

		if !passesFilter {
			continue
		}

		filteredCaptures = append(filteredCaptures, capture)
	}

	return filteredCaptures
}

func (node *Node) QueryMustPass(query *sitter.Query) (captures []Captures) {
	var returningCaptures []Captures
	onMatch := func(captures Captures) error {
		returningCaptures = append(returningCaptures, captures)
		return nil
	}
	err := node.Query(query, onMatch)
	if err != nil {
		log.Fatal().Msgf("invalid query %s", err)
	}
	return returningCaptures
}

func (node *Node) Tree() *Tree {
	return node.tree
}

func (node *Node) Column() uint32 {
	return node.sitter.StartPoint().Column
}

func IsDescendant(node *Node, parent *Node) bool {
	if node.Parent() == nil {
		return false
	}
	if node.Parent().ID() == parent.ID() {
		return true
	}

	return IsDescendant(node.Parent(), parent)
}

func QueryMustCompile(language *sitter.Language, text string) *sitter.Query {
	query, err := sitter.NewQuery([]byte(text), language)

	if err != nil {
		panic(fmt.Sprintf("unexpected error compiling tree sitter query:\n\nquery:\n\n%s\n\nerror: %s", text, err))
	}

	return query
}
