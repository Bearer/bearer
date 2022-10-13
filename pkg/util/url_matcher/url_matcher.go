package url_matcher

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

const prefixPattern = "(?P<match>\\A(?:[^:]+://)?(?:[^/]+\\.)?"
const suffixPattern = "(?:/|\\z)"

type ComparableUrls struct {
	DetectionURL string
	RecipeURL    string
}

func UrlMatcher(urls ComparableUrls) (string, error) {
	parsedURL, err := url.Parse(urls.RecipeURL)
	if err != nil {
		return "", err
	}

	parsedDomain, err := publicsuffix.ParseFromListWithOptions(
		publicsuffix.DefaultList,
		parsedURL.Host,
		&publicsuffix.FindOptions{IgnorePrivate: true, DefaultRule: nil},
	)
	if err != nil {
		return "", err
	}

	matcher, _ := regexp.Compile(prefixPattern + domainPattern(parsedDomain) + pathPattern(parsedURL) + ")" + suffixPattern)
	match := matcher.FindStringSubmatch(urls.DetectionURL)
	if match != nil {
		for i, name := range matcher.SubexpNames() {
			if name == "match" {
				return match[i], nil
			}
		}
	}
	return "", nil
}

func domainPattern(parsedDomain *publicsuffix.DomainName) string {
	var domainPatterns []string
	if parsedDomain.TRD != "" {
		domainPatterns = append(domainPatterns, wildcardPattern(parsedDomain.TRD, "."))
	}
	if parsedDomain.SLD != "" {
		domainPatterns = append(domainPatterns, regexp.QuoteMeta(parsedDomain.SLD))
	}
	if parsedDomain.TLD != "" {
		domainPatterns = append(domainPatterns, regexp.QuoteMeta(parsedDomain.TLD))
	}

	return strings.Join(domainPatterns, "\\.")
}

func pathPattern(parsedURL *url.URL) string {
	return wildcardPattern(strings.TrimSuffix(parsedURL.Path, "/"), "/")
}

func wildcardPattern(myString string, separator string) string {
	var parts []string
	for _, part := range strings.Split(myString, separator) {
		parts = append(parts, "("+strings.ReplaceAll(regexp.QuoteMeta(part), "\\*", ".+")+"|\\*)")
	}
	return strings.Join(parts, regexp.QuoteMeta(separator))
}
