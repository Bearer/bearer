package interfaces

import (
	"github.com/bearer/bearer/pkg/parser/interfaces/paths"
	"github.com/bearer/bearer/pkg/parser/interfaces/urls"
	"github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/report/values"
)

func KeyIsRelevant(key string) bool {
	return urls.KeyIsRelevant(key)
}

func GetTypeWithKey(key string, value *values.Value) (interfaces.Type, bool) {
	if urls.KeyIsRelevant(key) || urls.ValueIsRelevant(value) {
		return interfaces.TypeURL, true
	}

	return "", false
}

func GetType(value *values.Value, pathAllowed bool) (interfaces.Type, bool) {
	if urls.ValueIsRelevant(value) {
		return interfaces.TypeURL, true
	}

	if pathAllowed {
		if paths.ValueIsRelevant(value) {
			return interfaces.TypePath, true
		}
	}

	return "", false
}
