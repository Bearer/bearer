package flag

import (
	"strings"
	"time"
)

type Context string

const (
	Health Context = "health"
)

var (
	SkipPathFlag = Flag{
		Name:       "skip-path",
		ConfigName: "scan.skip-path",
		Value:      []string{},
		Usage:      "specify the comma separated files and directories to skip (supports * syntax), eg. --skip-path users/*.go,users/admin.sql",
	}
	DebugFlag = Flag{
		Name:       "debug",
		ConfigName: "scan.debug",
		Value:      false,
		Usage:      "enable debug logs",
	}
	DisableDomainResolutionFlag = Flag{
		Name:       "disable-domain-resolution",
		ConfigName: "scan.disable-domain-resolution",
		Value:      false,
		Usage:      "do not attempt to resolve detected domains during classification (default false), eg. --disable-domain-resolution=true",
	}
	DomainResolutionTimeoutFlag = Flag{
		Name:       "domain-resolution-timeout",
		ConfigName: "scan.domain-resolution-timeout",
		Value:      3 * time.Second,
		Usage:      "set timeout when attempting to resolve detected domains during classification (default 3 seconds), eg. --domain-resolution-timeout=TODO",
	}
	InternalDomainsFlag = Flag{
		Name:       "internal-domains",
		ConfigName: "scan.internal-domains",
		Value:      []string{},
		Usage:      "define regular expressions for better classification of private or unreachable domains eg. --internal-domains=\"*.my-company.com,private.sh\"",
	}
	ContextFlag = Flag{
		Name:       "context",
		ConfigName: "scan.context",
		Value:      "",
		Usage:      "expand context of schema classification e.g. --context=health to include data types particular to health",
	}
	QuietFlag = Flag{
		Name:       "quiet",
		ConfigName: "scan.quiet",
		Value:      false,
		Usage:      "suppress non-essential messages",
	}
)

type ScanFlagGroup struct {
	SkipPathFlag                *Flag
	DebugFlag                   *Flag
	DisableDomainResolutionFlag *Flag
	DomainResolutionTimeoutFlag *Flag
	InternalDomainsFlag         *Flag
	ContextFlag                 *Flag
	QuietFlag                   *Flag
}

type ScanOptions struct {
	Target                  string        `json:"target"`
	SkipPath                []string      `json:"skip_path"`
	Debug                   bool          `json:"debug"`
	DisableDomainResolution bool          `json:"disable_domain_resolution"`
	DomainResolutionTimeout time.Duration `json:"domain_resolution_timeout"`
	InternalDomains         []string      `json:"internal_domains"`
	Context                 Context       `json:"context"`
	Quiet                   bool          `json:"quiet"`
}

func NewScanFlagGroup() *ScanFlagGroup {
	return &ScanFlagGroup{
		SkipPathFlag:                &SkipPathFlag,
		DebugFlag:                   &DebugFlag,
		DisableDomainResolutionFlag: &DisableDomainResolutionFlag,
		DomainResolutionTimeoutFlag: &DomainResolutionTimeoutFlag,
		InternalDomainsFlag:         &InternalDomainsFlag,
		ContextFlag:                 &ContextFlag,
		QuietFlag:                   &QuietFlag,
	}
}

func (f *ScanFlagGroup) Name() string {
	return "Scan"
}

func (f *ScanFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.SkipPathFlag,
		f.DebugFlag,
		f.DisableDomainResolutionFlag,
		f.DomainResolutionTimeoutFlag,
		f.InternalDomainsFlag,
		f.ContextFlag,
		f.QuietFlag,
	}
}

func (f *ScanFlagGroup) ToOptions(args []string) (ScanOptions, error) {
	var target string
	if len(args) == 1 {
		target = args[0]
	}

	return ScanOptions{
		SkipPath:                getStringSlice(f.SkipPathFlag),
		Debug:                   getBool(f.DebugFlag),
		DisableDomainResolution: getBool(f.DisableDomainResolutionFlag),
		DomainResolutionTimeout: getDuration(f.DomainResolutionTimeoutFlag),
		InternalDomains:         getStringSlice(f.InternalDomainsFlag),
		Context:                 getContext(f.ContextFlag),
		Quiet:                   getBool(f.QuietFlag),
		Target:                  target,
	}, nil
}

func getContext(flag *Flag) Context {
	if flag == nil {
		return ""
	}

	flagStr := strings.ToLower(getString(flag))
	return Context(flagStr)
}
