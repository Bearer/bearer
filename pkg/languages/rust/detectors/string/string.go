package string

import (
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/util/stringutil"

	"github.com/bearer/bearer/pkg/scanner/detectors/common"
	"github.com/bearer/bearer/pkg/scanner/detectors/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_literal":
		return handleStringLiteral(node)
	case "raw_string_literal":
		return handleRawStringLiteral(node)
	case "binary_expression":
		// String concatenation: a + b
		if len(node.Children()) >= 2 && node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "compound_assignment_expr":
		// String append: s += "value"
		if len(node.Children()) >= 2 && node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}

func handleStringLiteral(node *tree.Node) ([]interface{}, error) {
	value := stringutil.StripQuotes(node.Content())

	return []interface{}{common.String{
		Value:     value,
		IsLiteral: true,
	}}, nil
}

func handleRawStringLiteral(node *tree.Node) ([]interface{}, error) {
	content := node.Content()

	// Raw strings are r"..." or r#"..."# etc.
	// Strip the r prefix and outer delimiters
	if len(content) >= 3 && content[0] == 'r' {
		// Find the start and end of actual content
		start := 1
		for start < len(content) && content[start] == '#' {
			start++
		}
		if start < len(content) && content[start] == '"' {
			start++
		}

		end := len(content) - 1
		for end > start && content[end] == '#' {
			end--
		}
		if end > start && content[end] == '"' {
			end--
		}

		if start <= end {
			content = content[start : end+1]
		} else {
			content = ""
		}
	}

	return []interface{}{common.String{
		Value:     content,
		IsLiteral: true,
	}}, nil
}

