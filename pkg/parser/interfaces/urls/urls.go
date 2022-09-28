package urls

import (
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/report/variables"
	"github.com/bearer/curio/pkg/util/normalize_key"
	"github.com/bearer/curio/pkg/util/regex"
	"golang.org/x/net/publicsuffix"
)

var (
	notEnoughInformationPattern = regexp.MustCompile(`^[\.\*]+$`)
	nonURLCharacterPattern      = regexp.MustCompile(`[^a-zA-Z0-9?/\-._~%=:+*]`)
	urlPattern                  = regexp.MustCompile(`^(?:(https?|\*)://)?([a-z0-9\-._*%]+)(:\d+)?(?:[/?].*)?$`)
	hasUsefulInformation        = regexp.MustCompile(`\w{3,}`)

	keyPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\b(sub)?domain\b`),
		regexp.MustCompile(`\bhost(name)?\b`),
		regexp.MustCompile(`\b(url|uri)\b`),
		regexp.MustCompile(`\bendpoint\b`),
		regexp.MustCompile(`\baddr\b`),
		regexp.MustCompile(`\bsvc\b`),
	}

	allowedDomainPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\.local$`),
		regexp.MustCompile(`\.lan$`),
	}
)

func KeyIsRelevant(key string) bool {
	return regex.AnyMatch(keyPatterns, normalize_key.Normalize(key))
}

func ValueIsRelevant(value *values.Value) bool {
	return textIsRelevant(value) || variablesAreRelevant(value)
}

func textIsRelevant(value *values.Value) bool {
	text := getText(value)
	if text == "" {
		return false
	}

	if nonURLCharacterPattern.MatchString(text) {
		return false
	}

	// Looks like a relative filename
	if strings.Contains(text, "..") {
		return false
	}

	match := urlPattern.FindStringSubmatch(text)
	if len(match) == 0 {
		return false
	}

	if !hasUsefulInformation.Match([]byte(text)) {
		return false
	}

	protocol := match[1]
	domain := match[2]

	// If we have a protocol then assume it's a URL
	if protocol != "" && protocol != "*" {
		return !notEnoughInformationPattern.MatchString(domain)
	}

	return domainValid(domain)
}

func variablesAreRelevant(value *values.Value) bool {
	for _, variableReference := range value.GetVariableReferences() {
		if variableReference.Identifier.Type == variables.VariableName {
			return false
		}

		if KeyIsRelevant(variableReference.Identifier.Name) {
			return true
		}
	}

	return false
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

func domainValid(domain string) bool {
	if domain == "" {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}

	if regex.AnyMatch(allowedDomainPatterns, domain) {
		return true
	}

	eTLD, icann := publicsuffix.PublicSuffix(domain)

	// Private eTLDs from the suffix list always have more than one segment, but
	// unknown ones will only have one segment (no dot)
	return icann || strings.Contains(eTLD, ".")
}
