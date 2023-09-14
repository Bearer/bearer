package stringutil

import "strings"

func SliceContains(slice []string, value string) bool {
	for _, sliceValue := range slice {
		if sliceValue == value {
			return true
		}
	}

	return false
}

func StripQuotes(input string) string {
	output := strings.Trim(input, `"`)
	return strings.Trim(output, `'`)
}
