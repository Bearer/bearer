package rust

import (
	"strings"

	"github.com/smacker/go-tree-sitter/rust"

	"github.com/bearer/bearer/pkg/detectors/rust/datatype"
	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/interfacedetector"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/report/variables"
	"github.com/bearer/bearer/pkg/util/file"
)

var language = rust.GetLanguage()

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Language != "Rust" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	datatype.Discover(report, tree, detector.idGenerator)

	if err := annotate(tree); err != nil {
		return false, err
	}
	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorRust,
		PathAllowed:      false,
	}); err != nil {
		return false, err
	}

	return true, nil
}

func annotate(tree *parser.Tree) error {
	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "binary_expression":
			// Handle string concatenation
			if node.FirstUnnamedChild() != nil && node.FirstUnnamedChild().Content() == "+" {
				left := node.ChildByFieldName("left")
				right := node.ChildByFieldName("right")
				if left != nil {
					value.Append(left.Value())
				}
				if right != nil {
					value.Append(right.Value())
				}
				return
			}
		case "identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())
			return
		case "string_literal":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				value.AppendString(text)
				return nil
			}, func(child *parser.Node) error {
				return nil
			})
			return
		case "raw_string_literal":
			// Raw strings have no escapes, just strip delimiters
			content := stripRawStringDelimiters(node.Content())
			value.AppendString(content)
			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func stripRawStringDelimiters(content string) string {
	if len(content) < 3 {
		return content
	}

	// Raw strings: r"..." or r#"..."# or r##"..."## etc.
	if content[0] != 'r' {
		return content
	}

	// Count the # characters
	hashCount := 0
	for i := 1; i < len(content) && content[i] == '#'; i++ {
		hashCount++
	}

	// Find the opening quote
	start := 1 + hashCount
	if start >= len(content) || content[start] != '"' {
		return content
	}
	start++

	// Find the closing quote and hashes
	end := len(content) - 1 - hashCount
	if end < start || content[end] != '"' {
		return content
	}

	return content[start:end]
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "field_initializer":
			// Struct { "ignored.domain": "..." }
			keyNode := parent.ChildByFieldName("name")
			if keyNode != nil && keyNode.Equal(lastNode) {
				return false
			}
		case "index_expression":
			// something["ignored.domain"]
			return false
		case "use_declaration":
			// use "ignored.domain"
			return false
		case "attribute_item":
			// #[derive(...)] - ignore attributes
			return false
		}

		lastNode = parent
	}

	// Ignore doc comments (they start with //! or ///)
	if node.Type() == "string_literal" {
		content := node.Content()
		if strings.HasPrefix(content, "//!") || strings.HasPrefix(content, "///") {
			return false
		}
	}

	return true
}

