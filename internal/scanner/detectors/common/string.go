package common

import (
	"fmt"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"

	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

type String struct {
	Value     string
	IsLiteral bool
}

func GetStringValue(node *tree.Node, scanContext types.ScanContext) (string, bool, error) {
	detections, err := scanContext.Scan(node, "string", "", settings.CURSOR_SCOPE)
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
		return "", false, fmt.Errorf(
			"expected single string detection but got %d for %s",
			len(detections),
			node.Debug(true),
		)
	}
}

func ConcatenateChildStrings(node *tree.Node, scanContext types.ScanContext) ([]interface{}, error) {
	value := ""
	isLiteral := true

	for _, child := range node.Children() {
		if !child.SitterNode().IsNamed() {
			continue
		}

		childValue, childIsLiteral, err := GetStringValue(child, scanContext)
		if err != nil {
			return nil, err
		}

		if childValue == "" && !childIsLiteral {
			childValue = "*"
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

func ConcatenateAssignEquals(node *tree.Node, scanContext types.ScanContext) ([]interface{}, error) {
	dataflowSources := node.ChildByFieldName("left").DataflowSources()
	if len(dataflowSources) == 0 {
		return nil, nil
	}
	if len(dataflowSources) != 1 {
		return nil, fmt.Errorf("expected exactly one data source for `+=` node but got %d", len(dataflowSources))
	}

	left, leftIsLiteral, err := GetStringValue(dataflowSources[0], scanContext)
	if err != nil {
		return nil, err
	}

	right, rightIsLiteral, err := GetStringValue(node.ChildByFieldName("right"), scanContext)
	if err != nil {
		return nil, err
	}

	if left == "" && !leftIsLiteral {
		left = "*"

		// No detection when neither parts are a string
		if right == "" && !rightIsLiteral {
			return nil, nil
		}
	}

	if right == "" && !rightIsLiteral {
		right = "*"
	}

	return []interface{}{String{
		Value:     left + right,
		IsLiteral: leftIsLiteral && rightIsLiteral,
	}}, nil
}