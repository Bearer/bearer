package flag

import (
	"errors"
	"os"
	"strings"
	"time"

	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/spf13/viper"
)

const (
	Health flagtypes.Context = "health"
	Empty  flagtypes.Context = ""

	ScannerSAST    = "sast"
	ScannerSecrets = "secrets"
)

var (
	ErrInvalidContext = errors.New("invalid context argument; supported values: health")
	ErrInvalidScanner = errors.New("invalid scanner argument; supported values: sast, secrets")
)

type scanFlagGroup struct{ flagGroupBase }

var ScanFlagGroup = &scanFlagGroup{flagGroupBase{name: "Scan"}}

var (
	SkipPathFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "skip-path",
		ConfigName: "scan.skip-path",
		Value:      []string{},
		Usage:      "Specify the comma separated files and directories to skip. Supports * syntax, e.g. --skip-path users/*.go,users/admin.sql",
	})
	SkipTestFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "skip-test",
		ConfigName: "scan.skip-test",
		Value:      true,
		Usage:      "Disable automatic skipping of test files",
	})
	SkipGitIgnore = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "skip-git-ignore",
		ConfigName: "scan.skip-git-ignore",
		Value:      false,
		Usage:      "Do not automatically skip files that match patterns in .gitignore",
	})
	DisableDomainResolutionFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "disable-domain-resolution",
		ConfigName: "scan.disable-domain-resolution",
		Value:      true,
		Usage:      "Do not attempt to resolve detected domains during classification",
	})
	DomainResolutionTimeoutFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "domain-resolution-timeout",
		ConfigName: "scan.domain-resolution-timeout",
		Value:      3 * time.Second,
		Usage:      "Set timeout when attempting to resolve detected domains during classification, e.g. --domain-resolution-timeout=3s",
	})
	InternalDomainsFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "internal-domains",
		ConfigName: "scan.internal-domains",
		Value:      []string{},
		Usage:      "Define regular expressions for better classification of private or unreachable domains e.g. --internal-domains=\".*.my-company.com,private.sh\"",
	})
	ContextFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "context",
		ConfigName: "scan.context",
		Value:      "",
		Usage:      "Expand context of schema classification e.g., --context=health, to include data types particular to health",
	})
	DataSubjectMappingFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "data-subject-mapping",
		ConfigName: "scan.data_subject_mapping",
		Value:      "",
		Usage:      "Override default data subject mapping by providing a path to a custom mapping JSON file",
	})
	QuietFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "quiet",
		ConfigName: "scan.quiet",
		Value:      false,
		Usage:      "Suppress non-essential messages",
	})
	HideProgressBarFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "hide-progress-bar",
		ConfigName: "scan.hide_progress_bar",
		Value:      false,
		Usage:      "Hide progress bar from output",
	})
	ForceFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "force",
		ConfigName: "scan.force",
		Value:      false,
		Usage:      "Disable the cache and runs the detections again",
	})
	ExternalRuleDirFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "external-rule-dir",
		ConfigName: "scan.external-rule-dir",
		Value:      []string{},
		Usage:      "Specify directories paths that contain .yaml files with external rules configuration",
	})
	ScannerFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:                 "scanner",
		ConfigName:           "scan.scanner",
		Value:                []string{ScannerSAST},
		Usage:                "Specify which scanner to use e.g. --scanner=secrets, --scanner=secrets,sast",
		EnvironmentVariables: []string{"SCANNER"},
	})
	ParallelFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "parallel",
		ConfigName: "scan.parallel",
		Value:      0,
		Usage:      "Specify the amount of parallelism to use during the scan",
	})
	ExitCodeFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:       "exit-code",
		ConfigName: "scan.exit-code",
		Value:      -1,
		Usage:      "Force a given exit code for the scan command. Set this to 0 (success) to always return a success exit code despite any findings from the scan.",
	})
	DiffFlag = ScanFlagGroup.add(flagtypes.Flag{
		Name:            "diff",
		ConfigName:      "scan.diff",
		Value:           false,
		Usage:           "Only report differences in findings relative to a base branch.",
		DisableInConfig: true,
	})
)

type ScanOptions struct {
	Target                  string            `mapstructure:"target" json:"target" yaml:"target"`
	SkipPath                []string          `mapstructure:"skip-path" json:"skip-path" yaml:"skip-path"`
	DisableDomainResolution bool              `mapstructure:"disable-domain-resolution" json:"disable-domain-resolution" yaml:"disable-domain-resolution"`
	DomainResolutionTimeout time.Duration     `mapstructure:"domain-resolution-timeout" json:"domain-resolution-timeout" yaml:"domain-resolution-timeout"`
	InternalDomains         []string          `mapstructure:"internal-domains" json:"internal-domains" yaml:"internal-domains"`
	Context                 flagtypes.Context `mapstructure:"context" json:"context" yaml:"context"`
	DataSubjectMapping      string            `mapstructure:"data_subject_mapping" json:"data_subject_mapping" yaml:"data_subject_mapping"`
	Quiet                   bool              `mapstructure:"quiet" json:"quiet" yaml:"quiet"`
	HideProgressBar         bool              `mapstructure:"hide_progress_bar" json:"hide_progress_bar" yaml:"hide_progress_bar"`
	Force                   bool              `mapstructure:"force" json:"force" yaml:"force"`
	ExternalRuleDir         []string          `mapstructure:"external-rule-dir" json:"external-rule-dir" yaml:"external-rule-dir"`
	Scanner                 []string          `mapstructure:"scanner" json:"scanner" yaml:"scanner"`
	Parallel                int               `mapstructure:"parallel" json:"parallel" yaml:"parallel"`
	ExitCode                int               `mapstructure:"exit-code" json:"exit-code" yaml:"exit-code"`
	Diff                    bool              `mapstructure:"diff" json:"diff" yaml:"diff"`
}

func (scanFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	var target string
	if len(args) == 1 {
		target = args[0]
	}

	context := getContext(ContextFlag)
	switch context {
	case Empty, Health:
	default:
		return ErrInvalidContext
	}

	scanners := getStringSlice(ScannerFlag)
	for _, scanner := range scanners {
		switch scanner {
		case ScannerSAST:
		case ScannerSecrets:
		default:
			return ErrInvalidScanner
		}
	}

	// DIFF_BASE_BRANCH is used for backwards compatibilty
	diff := getBool(DiffFlag) || os.Getenv("DIFF_BASE_BRANCH") != ""

	options.ScanOptions = flagtypes.ScanOptions{
		SkipPath:                getStringSlice(SkipPathFlag),
		SkipTest:                getBool(SkipTestFlag),
		SkipGitIgnore:           getBool(SkipGitIgnore),
		DisableDomainResolution: getBool(DisableDomainResolutionFlag),
		DomainResolutionTimeout: getDuration(DomainResolutionTimeoutFlag),
		InternalDomains:         getStringSlice(InternalDomainsFlag),
		Context:                 context,
		DataSubjectMapping:      getString(DataSubjectMappingFlag),
		Quiet:                   getBool(QuietFlag),
		HideProgressBar:         getBool(HideProgressBarFlag),
		Force:                   getBool(ForceFlag),
		Target:                  target,
		ExternalRuleDir:         getStringSlice(ExternalRuleDirFlag),
		Scanner:                 scanners,
		Parallel:                viper.GetInt(ParallelFlag.ConfigName),
		ExitCode:                viper.GetInt(ExitCodeFlag.ConfigName),
		Diff:                    diff,
	}

	return nil
}

func getContext(flag *flagtypes.Flag) flagtypes.Context {
	if flag == nil {
		return ""
	}

	flagStr := strings.ToLower(getString(flag))
	return flagtypes.Context(flagStr)
}
