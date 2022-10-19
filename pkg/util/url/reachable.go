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

type Reachable struct {
	LookUpAddr func(ctx context.Context, addr string) ([]string, error)
	LookUpNS   func(ctx context.Context, addr string) ([]*net.NS, error)
}

func NewReachable() *Reachable {
	var resolver = net.Resolver{PreferGo: true}

	return &Reachable{
		LookUpAddr: resolver.LookupAddr,
		LookUpNS:   resolver.LookupNS,
	}
}

func (reachable *Reachable) CanReach(myURL string, timeout time.Duration) bool {
	resolverContext, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	staticDomain := getStaticDomain(myURL)
	if myURL == staticDomain {
		return reachable.domainResolves(myURL, resolverContext)
	}

	if reachable.isNameserver(myURL, staticDomain, resolverContext) {
		return true
	}

	return reachable.domainResolves(myURL, resolverContext)
}

func (reachable *Reachable) isNameserver(myURL string, staticDomain string, resolverContext context.Context) bool {
	_, err := publicsuffix.ParseFromListWithOptions(
		publicsuffix.DefaultList,
		staticDomain,
		&publicsuffix.FindOptions{IgnorePrivate: true, DefaultRule: nil},
	)
	if err != nil {
		// invalid domain
		return false
	}

	nameserver, err := reachable.LookUpNS(resolverContext, staticDomain)
	if err != nil {
		// return false even for transient errors
		return false
	}

	return len(nameserver) > 0
}

func (reachable *Reachable) domainResolves(myURL string, resolverContext context.Context) bool {
	// handle any special characters
	sanitizedURL, err := publicsuffix.ToASCII(myURL)
	if err != nil {
		return false
	}

	names, err := reachable.LookUpAddr(resolverContext, sanitizedURL)
	if err != nil {
		// return false even for transient errors
		return false
	}

	return len(names) > 0
}

func getStaticDomain(myURL string) string {
	domainSplit := regexpDomainSplitMatcher.Split(myURL, -1)
	lenDomainSplit := len(domainSplit)
	if lenDomainSplit <= 1 {
		// single part, no static domain
		return myURL
	}

	return domainSplit[lenDomainSplit-1]
}
