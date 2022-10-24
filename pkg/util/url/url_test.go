package url_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/url"
	"github.com/stretchr/testify/assert"
)

type MatchTestCase struct {
	Name         string
	DetectionURL string
	RecipeURL    string
	Want         string
}

func TestMatch(t *testing.T) {
	tests := []MatchTestCase{
		{
			Name:         "when the urls are just domains and match exactly",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "https://api.bearer.com",
		},
		{
			Name:         "when the urls are just domains and don't match",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://docs.bearer.com",
			Want:         "",
		},
		{
			Name:         "when the detection url has a subdomain prefix",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://test.api.bearer.com",
			Want:         "https://test.api.bearer.com",
		},
		{
			Name:         "when the detection url has a non-subdomain prefix",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://testapi.bearer.com",
			Want:         "",
		},
		{
			Name:         "when the url has a subdomain prefix",
			RecipeURL:    "https://test.api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "",
		},
		{
			Name:         "when the detection url contains a wildcard for an entire domain segment",
			RecipeURL:    "https://s3.eu.amazonaws.com",
			DetectionURL: "https://s3.*.amazonaws.com",
			Want:         "https://s3.*.amazonaws.com",
		},
		{
			Name:         "when the detection url contains a wildcard for the tld",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.bearer.*",
			Want:         "",
		},
		{
			Name:         "when the detection url contains a wildcard for the sld",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.*.com",
			Want:         "",
		},
		{
			Name:         "when the detection url contains a wildcard for part of a domain segment",
			RecipeURL:    "https://s3.eu-west.amazonaws.com",
			DetectionURL: "https://s3.eu-*.amazonaws.com",
			Want:         "",
		},
		{
			Name:         "when the url contains a wildcard for an entire domain segment",
			RecipeURL:    "https://s3.*.amazonaws.com",
			DetectionURL: "https://s3.eu.amazonaws.com",
			Want:         "https://s3.eu.amazonaws.com",
		},
		{
			Name:         "when the url contains a wildcard for part of a domain segment",
			RecipeURL:    "https://s3.eu-*.amazonaws.com",
			DetectionURL: "https://s3.eu-west.amazonaws.com",
			Want:         "https://s3.eu-west.amazonaws.com",
		},
		{
			Name:         "when the url contains a wildcard matching across detection url segments",
			RecipeURL:    "https://api.*.amazonaws.com",
			DetectionURL: "https://api.s3.eu-west.amazonaws.com",
			Want:         "https://api.s3.eu-west.amazonaws.com",
		},
		{
			Name:         "when the urls include a path and match exactly",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "https://api.bearer.com/path",
		},
		{
			Name:         "when the urls include a path and the path doesn't match",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/otherpath",
			Want:         "",
		},
		{
			Name:         "when the urls match but the url has a trailing slash",
			RecipeURL:    "https://api.bearer.com/path/",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "https://api.bearer.com/path",
		},
		{
			Name:         "when the urls match but the detection url has a trailing slash",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path/",
			Want:         "https://api.bearer.com/path",
		},
		{
			Name:         "when the detection url has a path segment suffix",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path/other",
			Want:         "https://api.bearer.com/path",
		},
		{
			Name:         "when the url has a path segment suffix",
			RecipeURL:    "https://api.bearer.com/path/other",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "",
		},
		{
			Name:         "when the detection url has a non-segment path suffix",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/pathother",
			Want:         "",
		},
		{
			Name:         "when the detection url contains a wildcard for an entire path segment",
			RecipeURL:    "https://bearer.com/api/v2",
			DetectionURL: "https://bearer.com/api/*",
			Want:         "https://bearer.com/api/*",
		},
		{
			Name:         "when the detection url contains a wildcard for part of a path segment",
			RecipeURL:    "https://bearer.com/api/v2",
			DetectionURL: "https://bearer.com/api/v*",
			Want:         "",
		},
		{
			Name:         "when the url contains a wildcard for an entire path segment",
			RecipeURL:    "https://bearer.com/api/*",
			DetectionURL: "https://bearer.com/api/v2",
			Want:         "https://bearer.com/api/v2",
		},
		{
			Name:         "when the url contains a wildcard for part of a path segment",
			RecipeURL:    "https://bearer.com/api/v*",
			DetectionURL: "https://bearer.com/api/v2",
			Want:         "https://bearer.com/api/v2",
		},
		{
			Name:         "when the url contains a wildcard matching multiple path segments",
			RecipeURL:    "https://bearer.com/*/api",
			DetectionURL: "https://bearer.com/eu/west/api/fred",
			Want:         "https://bearer.com/eu/west/api",
		},
		{
			Name:         "when the urls match exactly except the scheme",
			RecipeURL:    "http://api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "https://api.bearer.com",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			regexpMatcher, _ := url.PrepareRegexpMatcher(testCase.RecipeURL)
			output, err := url.Match(
				testCase.DetectionURL,
				regexpMatcher,
			)
			if err != nil {
				t.Errorf("UrlMatcher returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}

func TestPrepareURLValue(t *testing.T) {
	tests := []struct {
		Name, Input string
		Want        string
	}{
		{
			Name:  "No change",
			Input: "http://my.example.com",
			Want:  "http://my.example.com",
		},
		{
			Name:  "Missing scheme",
			Input: "my.example.com",
			Want:  "https://my.example.com",
		},
		{
			Name:  `Wildcard replacement of %d`,
			Input: `my.%d.example.com`,
			Want:  "https://my.*.example.com",
		},
		{
			Name:  `Wildcard replacement of %s`,
			Input: `my.%s.example.com`,
			Want:  "https://my.*.example.com",
		},
		{
			Name:  "Wildcard replacement of <>",
			Input: "my.<variable>.example.com",
			Want:  "https://my.*.example.com",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := url.PrepareURLValue(testCase.Input)
			if err != nil {
				t.Errorf("PrepareURLValue returned unexpected error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}

	errorTestCases := []struct {
		Name, Input string
		Want        error
	}{
		{
			Name:  "Variables only",
			Input: "**",
			Want:  errors.New("URL is only made of variables"),
		},
	}

	for _, errorTestCase := range errorTestCases {
		t.Run(errorTestCase.Name, func(t *testing.T) {
			output, err := url.PrepareURLValue(errorTestCase.Input)
			if err == nil {
				t.Errorf("PreparedURLValue returned unexpected result %s", output)
			}

			assert.Equal(t, errorTestCase.Want, err)
		})
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Data *detections.Detection
		Want *classify.ValidationResult
	}{
		{
			Name: "when a detection is from a potential detector",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				DetectorType: detectors.DetectorEnvFile,
			},
			Want: &classify.ValidationResult{
				State:  classify.Potential,
				Reason: "potential_detectors",
			},
		},
		{
			Name: "when a detection is included inside a vendor folder - case 1",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				Source: source.Source{
					Filename: "foo/vendor/symfony/symfony/src/Symfony/Component/Validator/Mapping/MemberMetadata.php",
				},
			},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "included_in_vendor_folder",
			},
		},
		{
			Name: "when a detection is included inside a vendor folder - case 2",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				Source: source.Source{
					Filename: "rancher-powerdns4/vendor/github.com/prasmussen/gandi-api/client/client.go",
				},
			},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "included_in_vendor_folder",
			},
		},
		{
			Name: "when there's not data to make a strong decision",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{},
			Want: &classify.ValidationResult{
				State:  classify.Potential,
				Reason: "uncertain",
			},
		},
		{
			Name: "when the URL is blank",
			URL:  "",
			Data: &detections.Detection{},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "blank_url",
			},
		},
		{
			Name: "when an IP address is given",
			URL:  "https://127.0.0.1",
			Data: &detections.Detection{},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "ip_address_error",
			},
		},
		{
			Name: "when a dependency file is provided",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				Source: source.Source{
					Filename: "Gemfile.lock",
				},
			},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "dependency_file_error",
			},
		},
		{
			Name: "when markup is detected with a simple detector type",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				DetectorType: detectors.DetectorSimple,
				Source: source.Source{
					LanguageType: "markup",
				},
			},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "language_type_error",
			},
		},
		{
			Name: "when the filename isn't accepted",
			URL:  "https://eu.example.com/path/*",
			Data: &detections.Detection{
				Source: source.Source{
					Filename: "config/translations/en.js",
				},
			},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "filename_error",
			},
		},
		{
			Name: "missing TLD",
			URL:  "https://.",
			Data: &detections.Detection{},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "tld_error",
			},
		},
		{
			Name: "domain is empty",
			URL:  "/nothing",
			Data: &detections.Detection{},
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "no_host_error",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := url.ValidateFormat(testCase.URL, testCase.Data)
			if err != nil {
				t.Errorf("ValidateFormat returned unexpected error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}

func TestValidateInternal(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want *classify.ValidationResult
	}{
		{
			Name: "internal domain with restricted subdomain",
			URL:  "https://cdn.mish-company.com",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "internal_domain_subdomain_error",
			},
		},
		{
			Name: "internal domain with no subdomain",
			URL:  "https://mish-company.com",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "internal_domain_but_no_subdomain",
			},
		},
		{
			Name: "internal domain with invalid chars in path",
			URL:  "https://eu.mish-company.com/%20",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "internal_domain_errors_in_path",
			},
		},
		{
			Name: "internal domain with file extension in path",
			URL:  "https://eu.mish-company.com/photo.jpeg",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "internal_domain_errors_in_path",
			},
		},
		{
			Name: "internal domain with 'api' in path",
			URL:  "https://eu.mish-company.com/api/v2/shipment?debug=true",
			Want: &classify.ValidationResult{
				State:  classify.Valid,
				Reason: "internal_domain_path_contains_api_or_auth",
			},
		},
		{
			Name: "internal domain with a wildcard",
			URL:  "https://cdn*.mish-company.com/",
			Want: &classify.ValidationResult{
				State:  classify.Potential,
				Reason: "internal_domain_and_subdomain_with_wildcard",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := url.ValidateInternal(testCase.URL)
			if err != nil {
				t.Errorf("ValidateInternal returned unexpected error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}
func TestValidate(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want *classify.ValidationResult
	}{
		{
			Name: "TLD is not allowed",
			URL:  "https://example.id",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "tld_error",
			},
		},
		{
			Name: "blacklisted domain",
			URL:  "https://github.com",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "excluded_domains_error",
			},
		},
		{
			Name: "blacklisted subdomain",
			URL:  "https://wiki.example.com",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "subdomain_error",
			},
		},
		{
			Name: "path with invalid characters",
			URL:  "https://eu.example.com/%20",
			Want: &classify.ValidationResult{
				State:  classify.Invalid,
				Reason: "errors_in_path",
			},
		},
		{
			Name: "path contains 'api'",
			URL:  "https://eu.example.com/api/v2",
			Want: &classify.ValidationResult{
				State:  classify.Valid,
				Reason: "path_contains_api_or_auth",
			},
		},
		{
			Name: "subdomain not provided",
			URL:  "https://example.com",
			Want: &classify.ValidationResult{
				State:  classify.Potential,
				Reason: "no_subdomain",
			},
		},
		{
			Name: "subdomain contains 'api'",
			URL:  "https://api.example.com",
			Want: &classify.ValidationResult{
				State:  classify.Valid,
				Reason: "subdomain_contains_api",
			},
		},
		{
			Name: "subdomain contains a wildcard",
			URL:  "https://api.*.example.com",
			Want: &classify.ValidationResult{
				State:  classify.Potential,
				Reason: "subdomain_contains_api_with_wildcard",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			domainResolver := url.NewDomainResolverDefault()
			domainResolver.Enabled = false // disable domain resolution checks to avoid mocking

			output, err := url.Validate(testCase.URL, domainResolver)
			if err != nil {
				t.Errorf("Validate returned unexpected error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}

	// test domain reachability
	t.Run("unreachable domain", func(t *testing.T) {
		domainResolver := url.NewDomainResolverDefault()
		domainResolver.LookupIPAddr = func(ctx context.Context, host string) ([]net.IPAddr, error) {
			return []net.IPAddr{}, nil // no DNS address found
		}
		domainResolver.LookupNS = func(ctx context.Context, name string) ([]*net.NS, error) {
			return []*net.NS{}, nil // no NS found
		}

		output, err := url.Validate("https://unreachable.com", domainResolver)
		if err != nil {
			t.Errorf("Validate returned unexpected error %s", err)
		}

		assert.Equal(t, &classify.ValidationResult{
			State:  classify.Invalid,
			Reason: "domain_not_reachable",
		}, output)
	})
}
