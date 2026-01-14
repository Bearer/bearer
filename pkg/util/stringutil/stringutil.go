package stringutil

import (
	"fmt"
	"strconv"
	"strings"
)

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

func Unescape(value string) (string, error) {
	return strconv.Unquote(fmt.Sprintf(`"%s"`, value))
}

// UnescapeJavaScript handles JavaScript escape sequences.
// JavaScript is more permissive than Go - unknown escape sequences like \s, \d
// in string literals are treated as just the character after the backslash.
func UnescapeJavaScript(value string) string {
	return unescapePermissive(value)
}

// UnescapePython handles Python escape sequences.
// Python (historically) treats unknown escape sequences like \s, \d in string
// literals as just the character after the backslash. In Python 3.6+ this
// raises a DeprecationWarning, but for scanning purposes we use the permissive behavior.
func UnescapePython(value string) string {
	return unescapePermissive(value)
}

// unescapePermissive handles escape sequences permissively.
// It first tries standard Go unquoting, and if that fails (for unknown escapes),
// strips the backslash prefix.
func unescapePermissive(value string) string {
	// First try standard Go unquoting (handles \n, \t, \\, \", \xNN, \uNNNN, etc.)
	if result, err := strconv.Unquote(fmt.Sprintf(`"%s"`, value)); err == nil {
		return result
	}

	// For unknown escapes that Go doesn't recognize,
	// just strip the backslash (JavaScript/Python behavior for unknown escapes)
	if len(value) >= 2 && value[0] == '\\' {
		return value[1:]
	}

	// Fallback to raw value
	return value
}
