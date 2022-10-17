package url

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

const prefixPattern = "(?P<match>\\A(?:[^:]+://)?(?:[^/]+\\.)?"
const suffixPattern = "(?:/|\\z)"

var regexpReplaceMatcher = regexp.MustCompile(`<\w+>`)
var regexpVariableMatcher = regexp.MustCompile(`\A[*\/.-:]+\z`)

func Match(url string, matcher *regexp.Regexp) (string, error) {
	match := matcher.FindStringSubmatch(url)
	if match != nil {
		for i, name := range matcher.SubexpNames() {
			if name == "match" {
				return match[i], nil
			}
		}
	}
	return "", nil
}

func PrepareRegexpMatcher(myURL string) (*regexp.Regexp, error) {
	parsedURL, err := url.Parse(myURL)
	if err != nil {
		return nil, err
	}

	parsedDomain, err := publicsuffix.ParseFromListWithOptions(
		publicsuffix.DefaultList,
		parsedURL.Host,
		&publicsuffix.FindOptions{IgnorePrivate: true, DefaultRule: nil},
	)
	if err != nil {
		return nil, err
	}

	return regexp.Compile(prefixPattern + domainPattern(parsedDomain) + pathPattern(parsedURL) + ")" + suffixPattern)
}

func PrepareURLValue(myURL string) (string, error) {
	if regexpVariableMatcher.MatchString(myURL) {
		return "", errors.New("URL is only made of variables")
	}

	var preparedURL string
	// replace placeholders with wildcard *
	preparedURL = strings.ReplaceAll(myURL, `%d`, "*")
	preparedURL = strings.ReplaceAll(preparedURL, `%s`, "*")
	preparedURL = regexpReplaceMatcher.ReplaceAllString(preparedURL, "*")

	// ensure scheme is present
	if strings.HasPrefix(preparedURL, "http://") || strings.HasPrefix(preparedURL, "https://") {
		return preparedURL, nil
	}
	if strings.Contains(preparedURL, ".") {
		preparedURL = "https://" + preparedURL
	}

	return preparedURL, nil
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
