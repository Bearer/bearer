package url_test

import (
	"testing"

	"github.com/bearer/curio/pkg/util/url"
	"github.com/stretchr/testify/assert"
)

type MatchTestCase struct {
	Name         string
	DetectionURL string
	RecipeURL    string
	Want         string
	Skip         bool
}

func TestMatch(t *testing.T) {
	tests := []MatchTestCase{
		{
			Name:         "when the urls are just domains and match exactly",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "https://api.bearer.com",
			Skip:         false,
		},
		{
			Name:         "when the urls are just domains and don't match",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://docs.bearer.com",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url has a subdomain prefix",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://test.api.bearer.com",
			Want:         "https://test.api.bearer.com",
			Skip:         false,
		},
		{
			Name:         "when the detection url has a non-subdomain prefix",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://testapi.bearer.com",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the url has a subdomain prefix",
			RecipeURL:    "https://test.api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for an entire domain segment",
			RecipeURL:    "https://s3.eu.amazonaws.com",
			DetectionURL: "https://s3.*.amazonaws.com",
			Want:         "https://s3.*.amazonaws.com",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for the tld",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.bearer.*",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for the sld",
			RecipeURL:    "https://api.bearer.com",
			DetectionURL: "https://api.*.com",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for part of a domain segment",
			RecipeURL:    "https://s3.eu-west.amazonaws.com",
			DetectionURL: "https://s3.eu-*.amazonaws.com",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard for an entire domain segment",
			RecipeURL:    "https://s3.*.amazonaws.com",
			DetectionURL: "https://s3.eu.amazonaws.com",
			Want:         "https://s3.eu.amazonaws.com",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard for part of a domain segment",
			RecipeURL:    "https://s3.eu-*.amazonaws.com",
			DetectionURL: "https://s3.eu-west.amazonaws.com",
			Want:         "https://s3.eu-west.amazonaws.com",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard matching across detection url segments",
			RecipeURL:    "https://api.*.amazonaws.com",
			DetectionURL: "https://api.s3.eu-west.amazonaws.com",
			Want:         "https://api.s3.eu-west.amazonaws.com",
			Skip:         false,
		},
		{
			Name:         "when the urls include a path and match exactly",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "https://api.bearer.com/path",
			Skip:         false,
		},
		{
			Name:         "when the urls include a path and the path doesn't match",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/otherpath",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the urls match but the url has a trailing slash",
			RecipeURL:    "https://api.bearer.com/path/",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "https://api.bearer.com/path",
			Skip:         false,
		},
		{
			Name:         "when the urls match but the detection url has a trailing slash",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path/",
			Want:         "https://api.bearer.com/path",
			Skip:         false,
		},
		{
			Name:         "when the detection url has a path segment suffix",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/path/other",
			Want:         "https://api.bearer.com/path",
			Skip:         false,
		},
		{
			Name:         "when the url has a path segment suffix",
			RecipeURL:    "https://api.bearer.com/path/other",
			DetectionURL: "https://api.bearer.com/path",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url has a non-segment path suffix",
			RecipeURL:    "https://api.bearer.com/path",
			DetectionURL: "https://api.bearer.com/pathother",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for an entire path segment",
			RecipeURL:    "https://bearer.com/api/v2",
			DetectionURL: "https://bearer.com/api/*",
			Want:         "https://bearer.com/api/*",
			Skip:         false,
		},
		{
			Name:         "when the detection url contains a wildcard for part of a path segment",
			RecipeURL:    "https://bearer.com/api/v2",
			DetectionURL: "https://bearer.com/api/v*",
			Want:         "",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard for an entire path segment",
			RecipeURL:    "https://bearer.com/api/*",
			DetectionURL: "https://bearer.com/api/v2",
			Want:         "https://bearer.com/api/v2",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard for part of a path segment",
			RecipeURL:    "https://bearer.com/api/v*",
			DetectionURL: "https://bearer.com/api/v2",
			Want:         "https://bearer.com/api/v2",
			Skip:         false,
		},
		{
			Name:         "when the url contains a wildcard matching multiple path segments",
			RecipeURL:    "https://bearer.com/*/api",
			DetectionURL: "https://bearer.com/eu/west/api/fred",
			Want:         "https://bearer.com/eu/west/api",
			Skip:         false,
		},
		{
			Name:         "when the urls match exactly except the scheme",
			RecipeURL:    "http://api.bearer.com",
			DetectionURL: "https://api.bearer.com",
			Want:         "https://api.bearer.com",
			Skip:         false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Skip {
				t.Skip("interfaces not implemented")
			}

			output, err := url.Match(
				url.ComparableUrls{
					RecipeURL:    testCase.RecipeURL,
					DetectionURL: testCase.DetectionURL,
				},
			)
			if err != nil {
				t.Errorf("UrlMatcher returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}
