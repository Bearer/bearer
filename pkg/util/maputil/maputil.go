package maputil

import (
	"cmp"
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

func ToSortedSlice[mapKey cmp.Ordered, T any](input map[mapKey]T) []T {
	keys := make([]mapKey, 0)
	for key := range input {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	data := make([]T, 0)
	for _, key := range keys {
		data = append(data, input[key])
	}

	return data
}
