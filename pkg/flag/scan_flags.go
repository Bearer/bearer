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
		Usage:      "Specify the comma separated files and directories to skip. Supports * syntax, e.g. --skip-path users/*.go,users/admin.sql",
	}
	DebugFlag = Flag{
		Name:       "debug",
		ConfigName: "scan.debug",
		Value:      false,
		Usage:      "Enable debug logs",
	}
	DisableDomainResolutionFlag = Flag{
		Name:       "disable-domain-resolution",
		ConfigName: "scan.disable-domain-resolution",
		Value:      true,
		Usage:      "Do not attempt to resolve detected domains during classification (default false), e.g. --disable-domain-resolution=true",
	}
	DomainResolutionTimeoutFlag = Flag{
		Name:       "domain-resolution-timeout",
		ConfigName: "scan.domain-resolution-timeout",
		Value:      3 * time.Second,
		Usage:      "Set timeout when attempting to resolve detected domains during classification (default 3 seconds), e.g. --domain-resolution-timeout=3s",
	}
	InternalDomainsFlag = Flag{
		Name:       "internal-domains",
		ConfigName: "scan.internal-domains",
		Value:      []string{},
		Usage:      "Define regular expressions for better classification of private or unreachable domains e.g. --internal-domains=\"*.my-company.com,private.sh\"",
	}
	ContextFlag = Flag{
		Name:       "context",
		ConfigName: "scan.context",
		Value:      "",
		Usage:      "Expand context of schema classification e.g., --context=health, to include data types particular to health",
	}
	QuietFlag = Flag{
		Name:       "quiet",
		ConfigName: "scan.quiet",
		Value:      false,
		Usage:      "Suppress non-essential messages",
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
	Target                  string        `json:"target" yaml:"target"`
	SkipPath                []string      `json:"skip_path" yaml:"skip_path"`
	Debug                   bool          `json:"debug" yaml:"debug"`
	DisableDomainResolution bool          `json:"disable_domain_resolution" yaml:"disable_domain_resolution"`
	DomainResolutionTimeout time.Duration `json:"domain_resolution_timeout" yaml:"domain_resolution_timeout"`
	InternalDomains         []string      `json:"internal_domains" yaml:"internal_domains"`
	Context                 Context       `json:"context" yaml:"context"`
	Quiet                   bool          `json:"quiet" yaml:"quiet"`
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
