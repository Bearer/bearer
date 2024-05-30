package paths

import (
	"mime"
	"regexp"
	"strings"

	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/util/normalize_key"
	"github.com/bearer/bearer/pkg/util/regex"
)

var (
	onlyValidPattern = regexp.MustCompile(`^[\/\-\w\*\[\]?=&\:]+$`)
	pathPattern      = regexp.MustCompile(`\/`)
	looksLikeUrl     = regexp.MustCompile(`^(.+)?\:\/\/`)
	hasNormalChars   = regexp.MustCompile(`\w{3,}`)

	keyPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\bpath\b`),
	}
)

func KeyIsRelevant(key string) bool {
	return regex.AnyMatch(keyPatterns, normalize_key.Normalize(key))
}

func ValueIsRelevant(value *values.Value) bool {
	text := getText(value)

	_, _, err := mime.ParseMediaType(text)
	if err == nil {
		return false
	}

	if strings.Contains(text, "tmp") {
		return false
	}

	if strings.Contains(text, "..") {
		return false
	}

	if looksLikeUrl.MatchString(text) {
		return false
	}

	if !hasNormalChars.MatchString(text) {
		return false
	}

	return onlyValidPattern.MatchString(text) && pathPattern.MatchString(text)
}

// Replaces non string parts with asterisks
func getText(value *values.Value) string {
	text := ""

	for _, part := range value.Parts {
		if stringPart, ok := part.(*values.String); ok {
			text += stringPart.Value
		} else {
			text += "*"
		}
	}

	return text
}
