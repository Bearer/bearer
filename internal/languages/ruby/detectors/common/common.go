package common

import "github.com/bearer/bearer/internal/scanner/ast/tree"

func GetLiteralKey(keyNode *tree.Node) string {
	switch keyNode.Type() {
	case "hash_key_symbol":
		return keyNode.Content()
	case "simple_symbol":
		return keyNode.Content()[1:]
	case "string":
		if len(keyNode.Children()) == 3 && keyNode.Children()[1].Type() == "string_content" {
			return keyNode.Children()[1].Content()
		}
	}

	return ""
}
