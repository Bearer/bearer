package common

import "github.com/bearer/curio/new/language/tree"

func GetLiteralKey(keyNode *tree.Node) string {
	switch keyNode.Type() {
	case "hash_key_symbol":
		return keyNode.Content()
	case "simple_symbol":
		return keyNode.Content()[1:]
	case "string":
		if keyNode.NamedChildCount() == 1 && keyNode.Child(1).Type() == "string_content" {
			return keyNode.Child(1).Content()
		}
	}

	return ""
}
