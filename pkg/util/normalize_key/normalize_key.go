package normalize_key

import (
	"regexp"
	"slices"
	"strings"
)

var (
	normalizeCaseRegexp      = regexp.MustCompile(`[A-Z][A-Z][a-z]|[a-z][A-Z]`) // Matches "AP(INa)me" or "firs(tN)ame"
	normalizeSeparatorRegexp = regexp.MustCompile(`[$_\-.,\s:0-9]+`)
)

func Normalize(key string) string {
	start := 0
	var pieces []string
	matches := normalizeCaseRegexp.FindAllStringSubmatchIndex(key, -1)
	for _, match := range matches {
		splitPoint := match[0] + 1
		pieces = append(pieces, normalizeKeyPiece(key[start:splitPoint]))
		start = splitPoint
	}

	pieces = append(pieces, normalizeKeyPiece(key[start:]))

	if len(pieces) != 0 && (pieces[0] == "get" || pieces[0] == "set") {
		pieces = slices.Delete(pieces, 0, 1)
	}

	return strings.Join(pieces, " ")
}

func NormalizeAll(keys []string) []string {
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = Normalize(key)
	}
	return result
}

func normalizeKeyPiece(piece string) string {
	return strings.TrimSpace(strings.ToLower(normalizeSeparatorRegexp.ReplaceAllString(piece, " ")))
}
