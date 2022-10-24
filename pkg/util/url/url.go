package url

import (
	"errors"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

const prefixPattern = "(?P<match>\\A(?:[^:]+://)?(?:[^/]+\\.)?"
const suffixPattern = "(?:/|\\z)"

// recipe url matching regexp
var regexpReplaceMatcher = regexp.MustCompile(`<\w+>`)
var regexpVariableMatcher = regexp.MustCompile(`\A[*\/.-:]+\z`)

// url validation regexp
var regexpDependencyFileMatcher = regexp.MustCompile(`Gemfile\.lock|package\.json|yarn\.lock|maven\-dependencies\.json|gemnasium\-maven\-plugin\.json|gradle\-dependencies\.json|Pipfile\.lock|package\-lock\.json|npm\-shrinkwrap\.json|packages\.lock\.json|project\.json|packages\.config|paket\.dependencies|ivy\-report\.xml|composer\.lock|composer\.json|pipdeptree\.json|go\.sum|requirements\.txt|pyproject\.toml|poetry\.lock|pom\.xml|build\.gradle`)
var regexpInvalidFilenameMatcher = regexp.MustCompile(`(trad|/translations?/|locales?|dockerfile|i18n)`)
var regexpValidPathMatcher = regexp.MustCompile(`\A[\w\-.*/?=&\[\]]+\z`)
var regexpInvalidExtensionsInPathMatcher = regexp.MustCompile(` /\.md|\.zip|\.css|\.csv|\.xls|\.sh|\.jpg|\.jpeg|\.png|\.pdf|\.htm|\.html|\.xhtml|\.txt|\.dtd|\.sql|\.xsd|\.gif|\.ico/i`)

var excludedDomains = map[string]struct{}{
	"k8s.io":     {},
	"jquery.com": {},
	"github.com": {},
}

var blocklistTLDs = map[string]struct{}{
	"id":       {},
	"name":     {},
	"country":  {},
	"email":    {},
	"ping":     {},
	"phone":    {},
	"link":     {},
	"menu":     {},
	"zip":      {},
	"domains":  {},
	"search":   {},
	"js":       {},
	"tf":       {},
	"host":     {},
	"dev":      {},
	"local":    {},
	"md":       {},
	"show":     {},
	"video":    {},
	"global":   {},
	"now":      {},
	"py":       {},
	"tab":      {},
	"so":       {},
	"cc":       {},
	"off":      {},
	"services": {},
	"money":    {},
	"org":      {},
}

var subdomainNotAllowed = map[string]struct{}{
	"media":      {},
	"community":  {},
	"example":    {},
	"schemas":    {},
	"static":     {},
	"wiki":       {},
	"dist":       {},
	"support":    {},
	"dl":         {},
	"download":   {},
	"docs":       {},
	"cloudfront": {},
	"img":        {},
	"packages":   {},
	"fonts":      {},
	"releases":   {},
	"msdn":       {},
	"downloads":  {},
	"images":     {},
	"updates":    {},
	"upload":     {},
	"mobile":     {},
	"demo":       {},
	"forum":      {},
	"video":      {},
	"doc":        {},
	"tools":      {},
	"www2":       {},
	"groups":     {},
	"shemas":     {},
	"widgets":    {},
	"feeds":      {},
	"modules":    {},
	"package":    {},
	"blogs":      {},
	"news":       {},
	"www":        {},
	"faq":        {},
	"cdn":        {},
}

var invalidLanguageTypes = map[string]struct{}{
	"markup": {},
}

var allowedFilenameExtensions = map[string]struct{}{
	"twig": {},
	"tpl":  {},
	"ejs":  {},
}

var restrictedDetectorTypes = map[string]struct{}{
	"simple": {},
}

var invalidFilenameExtensions = map[string]struct{}{
	"tf":          {},
	"sql":         {},
	"ipynb":       {},
	"sh":          {},
	"io":          {},
	"j2":          {},
	"feature":     {},
	"xsl":         {},
	"ps1":         {},
	"bzl":         {},
	"cake":        {},
	"dockerfile":  {},
	"hcl":         {},
	"yar":         {},
	"xslt":        {},
	"tfvars":      {},
	"sbt":         {},
	"1":           {},
	"mk":          {},
	"psm1":        {},
	"bat":         {},
	"linq":        {},
	"ftl":         {},
	"rockspec":    {},
	"fs":          {},
	"nomad":       {},
	"es":          {},
	"snippet":     {},
	"pb":          {},
	"bash":        {},
	"m":           {},
	"com":         {},
	"ascs":        {},
	"rtf":         {},
	"genrule_cmd": {},
	"csx":         {},
	"old":         {},
	"tmp":         {},
	"notused":     {},
	"pp":          {},
	"cmd":         {},
	"bundle":      {},
	"purs":        {},
}

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

func ValidateFormat(myURL string, data *detections.Detection) (*classify.ValidationResult, error) {
	ValidationResult := classify.ValidationResult{
		State:  classify.Invalid,
		Reason: "uncertain", // default
	}

	if classify.IsVendored(data.Source.Filename) {
		ValidationResult.Reason = classify.IncludedInVendorFolderReason
		return &ValidationResult, nil
	}

	if classify.IsPotentialDetector(data.DetectorType) {
		ValidationResult.State = classify.Potential
		ValidationResult.Reason = classify.PotentialDetectorReason
		return &ValidationResult, nil
	}

	if myURL == "" {
		ValidationResult.Reason = "blank_url"
		return &ValidationResult, nil
	}

	parsedURL, err := url.Parse(myURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.Host == "" {
		ValidationResult.Reason = "no_host_error"
		return &ValidationResult, nil
	}

	if parsedURL.Host[0:1] == "." {
		ValidationResult.Reason = "tld_error"
		return &ValidationResult, nil
	}

	if net.ParseIP(parsedURL.Host) != nil {
		ValidationResult.Reason = "ip_address_error"
		return &ValidationResult, nil
	}

	if regexpDependencyFileMatcher.MatchString(data.Source.Filename) {
		ValidationResult.Reason = "dependency_file_error"
		return &ValidationResult, nil
	}

	filenameExtension := strings.TrimPrefix(strings.TrimSpace(filepath.Ext(data.Source.Filename)), ".")
	if invalidLanguageType(filenameExtension, data) {
		ValidationResult.Reason = "language_type_error"
		return &ValidationResult, nil
	}

	if invalidLanguage(filenameExtension) {
		ValidationResult.Reason = "invalid_language_error"
		return &ValidationResult, nil
	}

	if regexpInvalidFilenameMatcher.MatchString(data.Source.Filename) {
		ValidationResult.Reason = "filename_error"
		return &ValidationResult, nil
	}

	ValidationResult.State = classify.Potential
	return &ValidationResult, nil
}

func ValidateInternal(myURL string) (*classify.ValidationResult, error) {
	ValidationResult := classify.ValidationResult{
		State:  classify.Invalid,
		Reason: "uncertain", // default
	}

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

	if subdomainIsNotAllowed(parsedDomain.TRD) {
		ValidationResult.Reason = "internal_domain_subdomain_error"
		return &ValidationResult, nil
	}

	if pathError(parsedURL.Path) {
		ValidationResult.Reason = "internal_domain_errors_in_path"
		return &ValidationResult, nil
	}

	if pathContainsAPIorAuth(parsedURL.Path) {
		ValidationResult.State = classify.Valid
		ValidationResult.Reason = "internal_domain_path_contains_api_or_auth"

		return &ValidationResult, nil
	}

	if parsedDomain.TRD == "" {
		ValidationResult.Reason = "internal_domain_but_no_subdomain"
		return &ValidationResult, nil
	}

	if strings.Contains(parsedURL.Host, "*") {
		ValidationResult.State = classify.Potential
		ValidationResult.Reason = "internal_domain_and_subdomain_with_wildcard"
		return &ValidationResult, nil
	}

	ValidationResult.State = classify.Valid
	ValidationResult.Reason = "internal_domain_and_subdomain"
	return &ValidationResult, nil
}

func Validate(myURL string, domainResolver *DomainResolver) (*classify.ValidationResult, error) {
	ValidationResult := classify.ValidationResult{
		State:  classify.Invalid,
		Reason: "uncertain", // default
	}

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

	if parsedDomain.TLD == "" || isBlocklisted(parsedDomain.TLD) {
		ValidationResult.Reason = "tld_error"
		return &ValidationResult, nil
	}

	if !domainResolver.CanReach(parsedDomain.SLD + "." + parsedDomain.TLD) {
		ValidationResult.Reason = "domain_not_reachable"
		return &ValidationResult, nil
	}

	if domainIsExcluded(parsedURL.Host) {
		ValidationResult.Reason = "excluded_domains_error"
		return &ValidationResult, nil
	}

	if subdomainIsNotAllowed(parsedDomain.TRD) {
		ValidationResult.Reason = "subdomain_error"
		return &ValidationResult, nil
	}

	if pathError(parsedURL.Path) {
		ValidationResult.Reason = "errors_in_path"
		return &ValidationResult, nil
	}

	if pathContainsAPIorAuth(parsedURL.Path) {
		ValidationResult.State = classify.Valid
		ValidationResult.Reason = "path_contains_api_or_auth"
		return &ValidationResult, nil
	}

	if parsedDomain.TRD == "" {
		ValidationResult.State = classify.Potential
		ValidationResult.Reason = "no_subdomain"
		return &ValidationResult, nil
	}

	if strings.Contains(parsedDomain.TRD, "api") {
		if strings.Contains(parsedURL.Host, "*") {
			ValidationResult.State = classify.Potential
			ValidationResult.Reason = "subdomain_contains_api_with_wildcard"
			return &ValidationResult, nil
		}

		ValidationResult.State = classify.Valid
		ValidationResult.Reason = "subdomain_contains_api"
		return &ValidationResult, nil
	}

	ValidationResult.State = classify.Potential // uncertain
	return &ValidationResult, nil
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

func domainIsExcluded(host string) bool {
	_, ok := excludedDomains[host]
	return ok
}

func isBlocklisted(tld string) bool {
	_, ok := blocklistTLDs[tld]
	return ok
}

func subdomainIsNotAllowed(trd string) bool {
	_, ok := subdomainNotAllowed[trd]
	return ok
}

func invalidLanguageType(filenameExtension string, data *detections.Detection) bool {
	_, invalidLanguageType := invalidLanguageTypes[data.Source.LanguageType]
	_, validFilenameExtension := allowedFilenameExtensions[filenameExtension]
	_, restrictedDetectorType := restrictedDetectorTypes[string(data.DetectorType)]

	return invalidLanguageType && !validFilenameExtension && restrictedDetectorType
}

func invalidLanguage(filenameExtension string) bool {
	_, ok := invalidFilenameExtensions[filenameExtension]
	return ok
}

func pathError(path string) bool {
	if path == "" {
		return false
	}

	return !regexpValidPathMatcher.MatchString(path) || regexpInvalidExtensionsInPathMatcher.MatchString(path)
}

func pathContainsAPIorAuth(path string) bool {
	if path == "" {
		return false
	}

	return strings.Contains(path, "api") || strings.Contains(path, "auth")
}
