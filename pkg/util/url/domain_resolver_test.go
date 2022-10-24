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
		Name, Domain     string
		MockLookupIPAddr func(ctx context.Context, addr string) ([]net.IPAddr, error)
		MockLookupNS     func(ctx context.Context, name string) ([]*net.NS, error)
		Want             bool
	}{
		{
			Name:   "when we have non-ASCII chars",
			Domain: "meghívó-måradt.com",
			Want:   true,
		},
		{
			Name:   "when the domain contains wildcards and the non wildcard part is not a valid TLD",
			Domain: "_oidc-client-*-partial.tf",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, nil
			},
			Want: false,
		},
		{
			Name:   "when the domain is static and resolves",
			Domain: "example.co.za",
			Want:   true,
		},
		{
			Name:   "when the domain is static and does not resolve",
			Domain: "example.co.za",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, nil
			},
			Want: false,
		},
		{
			Name:   "when the domain is static and the DNS server has an error",
			Domain: "example.co.za",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, errors.New("just being nervous")
			},
			Want: false,
		},
		{
			Name:   "when the domain is static it does not lookup the nameserver",
			Domain: "example.co.za",
			MockLookupNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				panic("we shouldn't reach this")
			},
			Want: true,
		},
		{
			Name:   "when a wildcard domain does not resolve and its static domain has no nameserver",
			Domain: "*-test.example.co.uk",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, nil
			},
			MockLookupNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, nil
			},
			Want: false,
		},
		{
			Name:   "when a wildcard domain resolves and its static domain has no nameserver",
			Domain: "*-test.example.co.uk",
			MockLookupNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, nil
			},
			Want: true,
		},
		{
			Name:   "when a wildcard domain does not resolve but its static domain has a nameserver",
			Domain: "*-test.example.co.uk",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, nil
			},
			Want: true,
		},
		{
			Name:   "when both the DNS and nameserver lookup raises errors for a wildcard domain",
			Domain: "*-test.example.co.uk",
			MockLookupIPAddr: func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{}, errors.New("just being nervous")
			},
			MockLookupNS: func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{}, errors.New("just being nervous")
			},
			Want: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			// reset mocks
			domainResolver := url.NewDomainResolverDefault()
			domainResolver.LookupIPAddr = func(ctx context.Context, addr string) ([]net.IPAddr, error) {
				return []net.IPAddr{{IP: []byte{1}}}, nil
			}
			domainResolver.LookupNS = func(ctx context.Context, name string) ([]*net.NS, error) {
				return []*net.NS{{Host: "example.co.za"}}, nil
			}

			// update mocks if needed
			if testCase.MockLookupIPAddr != nil {
				domainResolver.LookupIPAddr = testCase.MockLookupIPAddr
			}
			if testCase.MockLookupNS != nil {
				domainResolver.LookupNS = testCase.MockLookupNS
			}
			output := domainResolver.CanReach(testCase.Domain)
			assert.Equal(t, testCase.Want, output)
		})
	}
}
