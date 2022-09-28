package regex

import "regexp"

// AnyMatch returns true if any of the regexes match the given string
func AnyMatch(exprs []*regexp.Regexp, str string) bool {
	for _, expr := range exprs {
		if expr.MatchString(str) {
			return true
		}
	}

	return false
}
