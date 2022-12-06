package settings

import (
	"embed"
	_ "embed"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/rego"
)

type Config struct {
	Worker         flag.WorkerOptions `mapstructure:"worker" json:"worker" yaml:"worker"`
	Scan           flag.ScanOptions   `mapstructure:"scan" json:"scan" yaml:"scan"`
	Report         flag.ReportOptions `mapstructure:"report" json:"report" yaml:"report"`
	CustomDetector map[string]Rule    `mapstructure:"custom_detector" json:"custom_detector" yaml:"custom_detector"`
	Policies       map[string]*Policy `mapstructure:"policies" json:"policies" yaml:"policies"`
	Target         string             `mapstructure:"target" json:"target" yaml:"target"`
}

type PolicyLevel string

var LevelCritical = "critical"
var LevelHigh = "high"
var LevelMedium = "medium"
var LevelLow = "low"

type Modules []*PolicyModule

type Policy struct {
	Query       string      `mapstructure:"query" json:"query" yaml:"query"`
	Id          string      `mapstructure:"id" json:"id" yaml:"id"`
	Name        string      `mapstructure:"name" json:"name" yaml:"name"`
	Description string      `mapstructure:"description" json:"description" yaml:"description"`
	Level       PolicyLevel `mapstructure:"level" json:"level" yaml:"level"`
	Modules     Modules     `mapstructure:"modules" json:"modules" yaml:"modules"`
}

type PolicyModule struct {
	Path    string `mapstructure:"path" json:"path,omitempty" yaml:"path,omitempty"`
	Name    string `mapstructure:"name" json:"name" yaml:"name"`
	Content string `mapstructure:"content" json:"content" yaml:"content"`
}

func (modules Modules) ToRegoModules() (output []rego.Module) {
	for _, module := range modules {
		output = append(output, rego.Module{
			Name:    module.Name,
			Content: module.Content,
		})
	}
	return
}

type PatternFilter struct {
	Variable string
	Values   []string
}

type RulePattern struct {
	Pattern string
	Filters []PatternFilter
}

type Rule struct {
	Disabled       bool          `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Type           string        `mapstructure:"type" json:"type" yaml:"type"`
	Languages      []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting bool          `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Processors     []Processor   `mapstructure:"processors" json:"processors" yaml:"processors"`
	Patterns       []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	RootSingularize bool `mapstructure:"root_singularize" yaml:"root_singularize" `
	RootLowercase   bool `mapstructure:"root_lowercase" yaml:"root_lowercase"`

	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
	Stored         bool               `mapstructure:"stored" json:"stored" yaml:"stored"`
	DetectPresence bool               `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent     bool               `mapstructure:"omit_parent" json:"omit_parent" yaml:"omit_parent"`
}

type Processor struct {
	Query   string  `mapstructure:"query" json:"query" yaml:"query"`
	Modules Modules `mapstructure:"modules" json:"modules" yaml:"modules"`
}

type MetaVar struct {
	Input  string `mapstructure:"input" json:"input" yaml:"input"`
	Output int    `mapstructure:"output" json:"output" yaml:"output"`
	Regex  string `mapstructure:"regex" json:"regex" yaml:"regex"`
}

//go:embed custom_detector.yml
var customDetector []byte

//go:embed policies.yml
var defaultPolicies []byte

//go:embed policies/*
var policiesFs embed.FS

//go:embed processors/*
var processorsFs embed.FS

var CustomDetectorKey string = "scan.custom_detector"
var PoliciesKey string = "scan.policies"

func FromOptions(opts flag.Options) (Config, error) {
	var rules map[string]Rule
	if viper.IsSet(CustomDetectorKey) {
		err := viper.UnmarshalKey(CustomDetectorKey, &rules)
		if err != nil {
			return Config{}, err
		}
	} else {
		rules = DefaultCustomDetector()
	}

	for _, customDetector := range rules {
		for _, processor := range customDetector.Processors {
			for _, module := range processor.Modules {
				if module.Path != "" {
					content, err := processorsFs.ReadFile(module.Path)
					if err != nil {
						return Config{}, err
					}
					module.Content = string(content)
					module.Path = ""
				}
			}
		}
	}

	var policies map[string]*Policy
	if viper.IsSet(PoliciesKey) {
		err := viper.UnmarshalKey(PoliciesKey, &policies)
		if err != nil {
			return Config{}, err
		}
	} else {
		policies = DefaultPolicies()
	}

	for key := range policies {
		policy := policies[key]

		if len(opts.PolicyOptions.OnlyPolicy) > 0 && !opts.PolicyOptions.OnlyPolicy[policy.Id] {
			delete(policies, key)
			continue
		}

		if opts.PolicyOptions.SkipPolicy[policy.Id] {
			delete(policies, key)
			continue
		}

		for _, module := range policy.Modules {
			if module.Path != "" {
				content, err := policiesFs.ReadFile(module.Path)
				if err != nil {
					return Config{}, err
				}
				module.Content = string(content)
			}
		}
	}

	config := Config{
		Worker:         opts.WorkerOptions,
		CustomDetector: rules,
		Scan:           opts.ScanOptions,
		Report:         opts.ReportOptions,
		Policies:       policies,
	}

	return config, nil
}

func (rulePattern *RulePattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to parse as a string
	var pattern string
	if err := unmarshal(&pattern); err == nil {
		rulePattern.Pattern = pattern
		return nil
	}

	// Wasn't a string so it must be the structured format
	type rawRulePattern RulePattern
	return unmarshal((*rawRulePattern)(rulePattern))
}

func DefaultCustomDetector() map[string]Rule {
	var rules map[string]Rule

	err := yaml.Unmarshal(customDetector, &rules)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal database file %e", err)
	}

	return rules
}

func DefaultPolicies() map[string]*Policy {
	var policies map[string]*Policy

	err := yaml.Unmarshal(defaultPolicies, &policies)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal database file %e", err)
	}

	return policies
}
