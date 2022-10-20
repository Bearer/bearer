package url_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/bearer/curio/pkg/util/url"
	"github.com/stretchr/testify/assert"
)

func TestCanReach(t *testing.T) {
	tests := []struct {
		Name, URL      string
		MockLookUpAddr func(ctx context.Context, addr string) ([]string, error)
		MockLookUpNS   func(ctx context.Context, name string) ([]*net.NS, error)
		Want           bool
	}{
		{
			Name: "when we have non-ASCII chars",
			URL:  "meghívó-måradt.com",
			Want: true,
		},
		{
			Name: "when the domain contains wildcards and the non wildcard part is not a valid TLD",
			URL:  "_oidc-client-*-partial.tf",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, nil
			},
			Want: false,
		},
		{
			Name: "when the domain is static and resolves",
			URL:  "www.example.co.za",
			Want: true,
		},
		{
			Name: "when the domain is static and does not resolve",
			URL:  "www.example.co.za",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, nil
			},
			Want: false,
		},
		{
			Name: "when the domain is static and the DNS server has an error",
			URL:  "www.example.co.za",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, errors.New("just being nervous")
			},
			Want: false,
		},
		{
			Name: "when the domain is static it does not lookup the nameserver",
			URL:  "www.example.co.za",
			MockLookUpNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				panic("we shouldn't reach this")
			},
			Want: true,
		},
		{
			Name: "when a wildcard domain does not resolve and its static domain has no nameserver",
			URL:  "*-test.example.co.uk",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, nil
			},
			MockLookUpNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, nil
			},
			Want: false,
		},
		{
			Name: "when a wildcard domain resolves and its static domain has no nameserver",
			URL:  "*-test.example.co.uk",
			MockLookUpNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, nil
			},
			Want: true,
		},
		{
			Name: "when a wildcard domain does not resolve but its static domain has a nameserver",
			URL:  "*-test.example.co.uk",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, nil
			},
			Want: true,
		},
		{
			Name: "when both the DNS and nameserver lookup raises errors for a wildcard domain",
			URL:  "*-test.example.co.uk",
			MockLookUpAddr: func(ctx context.Context, addr string) ([]string, error) {
				return []string{}, errors.New("just being nervous")
			},
			MockLookUpNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, errors.New("just being nervous")
			},
			Want: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			// reset mocks
			domainResolver := url.NewDomainResolverDefault()
			domainResolver.LookUpAddr = func(ctx context.Context, addr string) ([]string, error) {
				return []string{"www.example.co.za"}, nil
			}
			domainResolver.LookUpNS = func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{{Host: "example.co.za"}}, nil
			}

			// update mocks if needed
			if testCase.MockLookUpAddr != nil {
				domainResolver.LookUpAddr = testCase.MockLookUpAddr
			}
			if testCase.MockLookUpNS != nil {
				domainResolver.LookUpNS = testCase.MockLookUpNS
			}
			output := domainResolver.CanReach(testCase.URL)
			assert.Equal(t, testCase.Want, output)
		})
	}
}
