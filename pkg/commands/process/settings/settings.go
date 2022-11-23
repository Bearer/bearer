package settings

import (
	"embed"
	_ "embed"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/flag"
)

type Config struct {
	Worker         flag.WorkerOptions `json:"worker" yaml:"worker"`
	Scan           flag.ScanOptions   `json:"scan" yaml:"scan"`
	Report         flag.ReportOptions `json:"report" yaml:"report"`
	CustomDetector map[string]Rule    `json:"custom_detector" yaml:"custom_detector"`
	Policies       map[string]*Policy `json:"policies" yaml:"policies"`
}

type PolicyLevel string

var LevelCritical = "critical"
var LevelHigh = "high"
var LevelMedium = "medium"
var LevelLow = "low"

type Policy struct {
	Query       string
	Name        string
	Id          string
	Description string
	Modules     []*PolicyModule
	Level       PolicyLevel
}

type PolicyModule struct {
	Path    string
	Name    string
	Content string
}

type Rule struct {
	Disabled       bool
	Languages      []string
	Patterns       []string
	ParamParenting bool `yaml:"param_parenting"`
	Metavars       map[string]MetaVar
	Stored         bool
}

type MetaVar struct {
	Input  string
	Output int
	Regex  string
}

//go:embed custom_detector.yml
var customDetector []byte

//go:embed policies.yml
var defaultPolicies []byte

//go:embed policies/*
var policiesFs embed.FS

var CustomDetectorKey string = "scan.custom_detector"
var PoliciesKey string = "scan.policies"

func FromOptions(opts flag.Options) (Config, error) {
	rules := DefaultCustomDetector()
	if viper.IsSet(CustomDetectorKey) {
		err := viper.UnmarshalKey(CustomDetectorKey, &rules)
		if err != nil {
			return Config{}, err
		}
	}

	policies := DefaultPolicies()
	if viper.IsSet(PoliciesKey) {
		err := viper.UnmarshalKey(PoliciesKey, &rules)
		if err != nil {
			return Config{}, err
		}
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
