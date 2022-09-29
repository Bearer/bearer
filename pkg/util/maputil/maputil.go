package maputil

import (
	"reflect"
	"sort"
)

func SortedStringKeys(mapValue interface{}) []string {
	interfaceKeys := reflect.ValueOf(mapValue).MapKeys()
	keys := make([]string, len(interfaceKeys))

	for i, interfaceKey := range interfaceKeys {
		keys[i] = interfaceKey.Interface().(string)
	}

	sort.Strings(keys)

	return keys
}
