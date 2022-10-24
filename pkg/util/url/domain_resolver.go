package url

import (
	"context"
	"net"
	"regexp"
	"time"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

// for domain resolution - find anything between * and .
var regexpDomainSplitMatcher = regexp.MustCompile(`\*\s*(.*?)\s*\.`)

type DomainResolver struct {
	Enabled      bool
	Timeout      time.Duration
	LookupIPAddr func(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupNS     func(ctx context.Context, host string) ([]*net.NS, error)
}

func NewDomainResolver(enabled bool, timeout time.Duration) *DomainResolver {
	var resolver = net.Resolver{PreferGo: true}

	return &DomainResolver{
		Enabled:      enabled,
		Timeout:      timeout,
		LookupNS:     resolver.LookupNS,
		LookupIPAddr: resolver.LookupIPAddr,
	}
}

func NewDomainResolverDefault() *DomainResolver {
	var resolver = net.Resolver{PreferGo: true}

	return &DomainResolver{
		Enabled:      true,
		Timeout:      3 * time.Second,
		LookupNS:     resolver.LookupNS,
		LookupIPAddr: resolver.LookupIPAddr,
	}
}

func (domainResolver *DomainResolver) CanReach(domain string) bool {
	if domainResolver == nil || !domainResolver.Enabled {
		// skip disabled check
		return true
	}

	resolverContext, cancel := context.WithTimeout(context.TODO(), domainResolver.Timeout)
	defer cancel()

	staticDomain := getStaticDomain(domain)
	if domain == staticDomain {
		return domainResolver.domainResolves(domain, resolverContext)
	}

	if domainResolver.isNameserver(domain, staticDomain, resolverContext) {
		return true
	}

	return domainResolver.domainResolves(domain, resolverContext)
}

func (domainResolver *DomainResolver) isNameserver(domain string, staticDomain string, resolverContext context.Context) bool {
	_, err := publicsuffix.ParseFromListWithOptions(
		publicsuffix.DefaultList,
		staticDomain,
		&publicsuffix.FindOptions{IgnorePrivate: true, DefaultRule: nil},
	)
	if err != nil {
		// invalid domain
		return false
	}

	nameserver, err := domainResolver.LookupNS(resolverContext, staticDomain)
	if err != nil {
		// return false even for transient errors
		return false
	}

	return len(nameserver) > 0
}

func (domainResolver *DomainResolver) domainResolves(domain string, resolverContext context.Context) bool {
	// handle any special characters
	sanitizedDomain, err := publicsuffix.ToASCII(domain)
	if err != nil {
		return false
	}

	names, err := domainResolver.LookupIPAddr(resolverContext, sanitizedDomain)
	if err != nil {
		// return false even for transient errors
		return false
	}

	return len(names) > 0
}

func getStaticDomain(domain string) string {
	domainSplit := regexpDomainSplitMatcher.Split(domain, -1)
	lenDomainSplit := len(domainSplit)
	if lenDomainSplit <= 1 {
		// single part, no static domain
		return domain
	}

	return domainSplit[lenDomainSplit-1]
}
