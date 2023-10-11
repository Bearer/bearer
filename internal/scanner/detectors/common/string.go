package common

import (
	"github.com/bearer/bearer/internal/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"

	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

const NonLiteralValue = "\uFFFD" // unicode Replacement character

type String struct {
	Value     string
	IsLiteral bool
}

func GetStringValue(node *tree.Node, detectorContext types.Context) (string, bool, error) {
	detections, err := detectorContext.Scan(node, ruleset.BuiltinStringRule, traversalstrategy.Cursor)
	if err != nil {
		return "", false, err
	}

	switch len(detections) {
	case 0:
		return "", false, nil
	case 1:
		childString := detections[0].Data.(String)

		return childString.Value, childString.IsLiteral, nil
	default:
		literalValue := ""
		for _, detection := range detections {
			childString := detection.Data.(String)
			if childString.IsLiteral && childString.Value != "" {
				if literalValue != "" && childString.Value != literalValue {
					return "", false, nil
				}

				literalValue = childString.Value
			}
		}

		return literalValue, true, nil
	}
}

func ConcatenateChildStrings(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	value := ""
	isLiteral := true

	for _, child := range node.Children() {
		if !child.IsNamed() {
			continue
		}

		childValue, childIsLiteral, err := GetStringValue(child, detectorContext)
		if err != nil {
			return nil, err
		}

		if childValue == "" && !childIsLiteral {
			childValue = NonLiteralValue
		}

		value += childValue

		if !childIsLiteral {
			isLiteral = false
		}
	}

	return []interface{}{String{
		Value:     value,
		IsLiteral: isLiteral,
	}}, nil
}

func ConcatenateAssignEquals(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	left, leftIsLiteral, err := GetStringValue(node.ChildByFieldName("left"), detectorContext)
	if err != nil {
		return nil, err
	}

	right, rightIsLiteral, err := GetStringValue(node.ChildByFieldName("right"), detectorContext)
	if err != nil {
		return nil, err
	}

	if left == "" && !leftIsLiteral {
		left = NonLiteralValue

		// No detection when neither parts are a string
		if right == "" && !rightIsLiteral {
			return nil, nil
		}
	}

	if right == "" && !rightIsLiteral {
		right = NonLiteralValue
	}

	return []interface{}{String{
		Value:     left + right,
		IsLiteral: leftIsLiteral && rightIsLiteral,
	}}, nil
}
