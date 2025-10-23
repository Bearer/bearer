package types

import (
	"time"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/pkg/util/set"
)

type Flag struct {
	// Name is for CLI flag and environment variable.
	// If this field is empty, it will be available only in config file.
	Name string

	// ConfigName is a key in config file. It is also used as a key of viper.
	ConfigName string

	// Shorthand is a shorthand letter.
	Shorthand string

	// Value is the default value. It must be filled to determine the flag type.
	Value interface{}

	// Usage explains how to use the flag.
	Usage string

	// DisableInConfig represents if flag should be present in config
	DisableInConfig bool

	// Do not show flag in the helper
	Hide bool

	// Deprecated represents if the flag is deprecated
	Deprecated bool

	// Additional environment variables to read the value from, in addition to the default
	EnvironmentVariables []string
}

type FlagGroup interface {
	Name() string
	Flags() []*Flag
	SetOptions(options *Options, args []string) error
}

type Context string

// Options holds all the runtime configuration
type Options struct {
	ReportOptions
	RuleOptions
	ScanOptions
	RepositoryOptions
	GeneralOptions
	IgnoreAddOptions
	IgnoreShowOptions
	IgnoreMigrateOptions
	WorkerOptions
}

type ScanOptions struct {
	Target                  string        `mapstructure:"target" json:"target" yaml:"target"`
	SkipTest                bool          `mapstructure:"skip-test" json:"skip-test" yaml:"skip-test"`
	SkipGitIgnore           bool          `mapstructure:"skip-git-ignore" json:"skip-git-ignore" yaml:"skip-git-ignore"`
	SkipPath                []string      `mapstructure:"skip-path" json:"skip-path" yaml:"skip-path"`
	DisableDomainResolution bool          `mapstructure:"disable-domain-resolution" json:"disable-domain-resolution" yaml:"disable-domain-resolution"`
	DomainResolutionTimeout time.Duration `mapstructure:"domain-resolution-timeout" json:"domain-resolution-timeout" yaml:"domain-resolution-timeout"`
	InternalDomains         []string      `mapstructure:"internal-domains" json:"internal-domains" yaml:"internal-domains"`
	Context                 Context       `mapstructure:"context" json:"context" yaml:"context"`
	DataSubjectMapping      string        `mapstructure:"data_subject_mapping" json:"data_subject_mapping" yaml:"data_subject_mapping"`
	Quiet                   bool          `mapstructure:"quiet" json:"quiet" yaml:"quiet"`
	HideProgressBar         bool          `mapstructure:"hide_progress_bar" json:"hide_progress_bar" yaml:"hide_progress_bar"`
	Force                   bool          `mapstructure:"force" json:"force" yaml:"force"`
	ExternalRuleDir         []string      `mapstructure:"external-rule-dir" json:"external-rule-dir" yaml:"external-rule-dir"`
	Scanner                 []string      `mapstructure:"scanner" json:"scanner" yaml:"scanner"`
	Parallel                int           `mapstructure:"parallel" json:"parallel" yaml:"parallel"`
	ExitCode                int           `mapstructure:"exit-code" json:"exit-code" yaml:"exit-code"`
	Diff                    bool          `mapstructure:"diff" json:"diff" yaml:"diff"`
}

type RuleOptions struct {
	DisableDefaultRules bool            `mapstructure:"disable-default-rules" json:"disable-default-rules" yaml:"disable-default-rules"`
	SkipRule            map[string]bool `mapstructure:"skip-rule" json:"skip-rule" yaml:"skip-rule"`
	OnlyRule            map[string]bool `mapstructure:"only-rule" json:"only-rule" yaml:"only-rule"`
}

type ReportOptions struct {
	Format             string          `mapstructure:"format" json:"format" yaml:"format"`
	Report             string          `mapstructure:"report" json:"report" yaml:"report"`
	Output             string          `mapstructure:"output" json:"output" yaml:"output"`
	Severity           set.Set[string] `mapstructure:"severity" json:"severity" yaml:"severity"`
	FailOnSeverity     set.Set[string] `mapstructure:"fail-on-severity" json:"fail-on-severity" yaml:"fail-on-severity"`
	ExcludeFingerprint map[string]bool `mapstructure:"exclude_fingerprints" json:"exclude_fingerprints" yaml:"exclude_fingerprints"`
	NoExtract          bool            `mapstructure:"no-extract" json:"no-extract" yaml:"no-extract"`
	NoRuleMeta         bool            `mapstructure:"no-rule-meta" json:"no-rule-meta" yaml:"no-rule-meta"`
	IncludeStats       bool            `mapstructure:"include-stats" json:"include-stats" yaml:"include-stats"`
}

type RepositoryOptions struct {
	OriginURL         string
	Branch            string
	Commit            string
	DefaultBranch     string
	DiffBaseBranch    string
	DiffBaseCommit    string
	GithubToken       string
	GithubRepository  string
	GithubAPIURL      string
	PullRequestNumber string
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile          string `json:"config_file" yaml:"config_file"`
	Client              *api.API
	DisableVersionCheck bool
	NoColor             bool   `mapstructure:"no_color" json:"no_color" yaml:"no_color"`
	IgnoreFile          string `mapstructure:"ignore_file" json:"ignore_file" yaml:"ignore_file"`
	Debug               bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	LogLevel            string `mapstructure:"log-level" json:"log-level" yaml:"log-level"`
	DebugProfile        bool
	IgnoreGit           bool `mapstructure:"ignore-git" json:"ignore-git" yaml:"ignore-git"`
}

type IgnoreAddOptions struct {
	Author        string `mapstructure:"author" json:"author" yaml:"author"`
	Comment       string `mapstructure:"comment" json:"comment" yaml:"comment"`
	FalsePositive bool   `mapstructure:"false_positive" json:"false_positive" yaml:"false_positive"`
	Force         bool   `mapstructure:"ignore_add_force" json:"ignore_add_force" yaml:"ignore_add_force"`
}

type IgnoreShowOptions struct {
	All bool `mapstructure:"all" json:"all" yaml:"all"`
}

type IgnoreMigrateOptions struct {
	Force bool `mapstructure:"ignore_migrate_force" json:"ignore_migrate_force" yaml:"ignore_migrate_force"`
}

type WorkerOptions struct {
	ParentProcessID int
	WorkerID        string `mapstructure:"worker-id" json:"worker-id" yaml:"worker-id"`
	Port            string `mapstructure:"port" json:"port" yaml:"port"`
}
