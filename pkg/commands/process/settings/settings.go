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
	Worker         flag.WorkerOptions `json:"worker"`
	Scan           flag.ScanOptions   `json:"scan"`
	Report         flag.ReportOptions `json:"report"`
	CustomDetector map[string]Rule    `json:"custom_detector"`
	Policies       map[string]*Policy `json:"policies"`
}

type policyLevel string

var LevelMedium = "medium"
var LevelWarning = "warning"
var LevelCritical = "critical"

type Modules []*PolicyModule

type Policy struct {
	Query   string
	Message string
	Modules Modules
	Level   policyLevel
}

type PolicyModule struct {
	Path    string `yaml:"path,omitempty"`
	Name    string
	Content string
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

type Rule struct {
	Disabled       bool
	Type           string
	Verifier       bool
	Languages      []string
	Patterns       []string
	ParamParenting bool `yaml:"param_parenting"`
	Processors     []Processor

	Singularilize bool
	Lowercase     bool
	Metavars      map[string]MetaVar
	Stored        bool
}

type Processor struct {
	Query   string
	Modules Modules
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

//go:embed processors/*
var processorsFs embed.FS

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

	policies := DefaultPolicies()
	if viper.IsSet(PoliciesKey) {
		err := viper.UnmarshalKey(PoliciesKey, &policies)
		if err != nil {
			return Config{}, err
		}
	}

	for _, policy := range policies {
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
