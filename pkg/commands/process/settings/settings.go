package settings

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
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
	DisplayId   string      `mapstructure:"display_id" json:"display_id" yaml:"display_id"`
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
	Either             []PatternFilter `mapstructure:"either" json:"either" yaml:"either"`
	Variable           string          `mapstructure:"variable" json:"variable" yaml:"variable"`
	Detection          string          `mapstructure:"detection" json:"detection" yaml:"detection"`
	Values             []string        `mapstructure:"values" json:"values" yaml:"values"`
	LessThan           *int            `mapstructure:"less_than" json:"less_than" yaml:"less_than"`
	LessThanOrEqual    *int            `mapstructure:"less_than_or_equal" json:"less_than_or_equal" yaml:"less_than_or_equal"`
	GreaterThan        *int            `mapstructure:"greater_than" json:"greater_than" yaml:"greater_than"`
	GreaterThanOrEqual *int            `mapstructure:"greater_than_or_equal" json:"greater_than_or_equal" yaml:"greater_than_or_equal"`

	// FIXME: remove when refactor is complete
	Minimum        *int `mapstructure:"minimum" json:"minimum" yaml:"minimum"`
	Maximum        *int `mapstructure:"maximum" json:"maximum" yaml:"maximum"`
	MatchViolation bool `mapstructure:"match_violation" json:"match_violation" yaml:"match_violation"`
}

type RulePattern struct {
	Pattern string          `mapstructure:"pattern" json:"pattern" yaml:"pattern"`
	Filters []PatternFilter `mapstructure:"filters" json:"filters" yaml:"filters"`
}

type Rule struct {
	Disabled       bool          `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Type           string        `mapstructure:"type" json:"type" yaml:"type"`
	Languages      []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting bool          `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns       []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	RootSingularize bool `mapstructure:"root_singularize" yaml:"root_singularize" `
	RootLowercase   bool `mapstructure:"root_lowercase" yaml:"root_lowercase"`

	Stored bool `mapstructure:"stored" json:"stored" yaml:"stored"`

	// FIXME: remove after refactor
	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
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

//go:embed policies.yml
var defaultPolicies []byte

//go:embed custom_detectors/*
var customDetectorFS embed.FS

//go:embed policies/*
var policiesFs embed.FS

//go:embed processors/*
var processorsFs embed.FS

var CustomDetectorKey string = "scan.custom_detector"
var PoliciesKey string = "scan.policies"

func FromOptions(opts flag.Options) (Config, error) {
	detectors := DefaultCustomDetector()

	policies := DefaultPolicies()

	// validate detector options
	onlyDetector := opts.DetectorOptions.OnlyDetector
	skipDetector := opts.DetectorOptions.SkipDetector

	validDetectors := make(map[string]bool)
	for key := range detectors {
		validDetectors[key] = true
	}

	var invalidDetectors []string
	for key := range onlyDetector {
		if !validDetectors[key] {
			invalidDetectors = append(invalidDetectors, key)
		}
	}

	for key := range skipDetector {
		if !validDetectors[key] {
			invalidDetectors = append(invalidDetectors, key)
		}
	}

	if len(invalidDetectors) > 0 {
		return Config{}, fmt.Errorf("unknown detectors %s", invalidDetectors)
	}

	// apply detector options
	for key := range detectors {
		if len(onlyDetector) > 0 && !onlyDetector[key] {
			delete(detectors, key)
			continue
		}

		if skipDetector[key] {
			delete(detectors, key)
			continue
		}
	}

	externalDetectors, err := LoadExternalDetectors(opts.ExternalDetectorDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load external detectors %w", err)
	}

	for ruleName, rule := range externalDetectors {
		_, ok := detectors[ruleName]
		if ok {
			return Config{}, fmt.Errorf("tried to overwrite default custom detector %s with external detector", ruleName)
		}

		detectors[ruleName] = rule
	}

	// validate policy options
	onlyPolicy := opts.PolicyOptions.OnlyPolicy
	skipPolicy := opts.PolicyOptions.SkipPolicy

	policyDisplayIds := make(map[string]bool)
	for key := range policies {
		policy := policies[key]
		policyDisplayIds[policy.DisplayId] = true
	}

	var invalidPolicyDisplayIds []string
	for key := range onlyPolicy {
		if !policyDisplayIds[key] {
			invalidPolicyDisplayIds = append(invalidPolicyDisplayIds, key)
		}
	}

	for key := range skipPolicy {
		if !policyDisplayIds[key] {
			invalidPolicyDisplayIds = append(invalidPolicyDisplayIds, key)
		}
	}

	if len(invalidPolicyDisplayIds) > 0 {
		return Config{}, fmt.Errorf("unknown policy ids %s", invalidPolicyDisplayIds)
	}

	// apply policy options
	for key := range policies {
		policy := policies[key]

		if len(onlyPolicy) > 0 && !onlyPolicy[policy.DisplayId] {
			delete(policies, key)
			continue
		}

		if skipPolicy[policy.DisplayId] {
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

	// apply external policies
	externalPolicies, err := LoadExternalPolicies(opts.ExternalPolicyDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load external policies %w", err)
	}

	for policyName, policy := range externalPolicies {
		_, ok := policies[policyName]
		if ok {
			return Config{}, fmt.Errorf("tried to overwrite default policy %s with external detector", policyName)
		}

		policies[policyName] = policy
	}

	config := Config{
		Worker:         opts.WorkerOptions,
		CustomDetector: detectors,
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
	customDetectorsDir := "custom_detectors"
	rules := make(map[string]Rule)

	dirEntries, err := customDetectorFS.ReadDir(customDetectorsDir)
	if err != nil {
		log.Fatal().Msgf("failed to read custom detectors dir %e", err)
	}

	for _, entry := range dirEntries {
		fileName := entry.Name()

		ext := filepath.Ext(fileName)
		ruleName := strings.TrimSuffix(fileName, ext)

		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		fileContent, err := customDetectorFS.ReadFile(customDetectorsDir + "/" + fileName)
		if err != nil {
			log.Fatal().Msgf("failed to read custom detector file %e", err)
		}

		var rule Rule
		err = yaml.Unmarshal(fileContent, &rule)
		if err != nil {
			log.Fatal().Msgf("failed to unmarshal database file %e", err)
		}

		rules[ruleName] = rule
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

func ProcessorRegoModuleText(processorName string) (string, error) {
	processorPath := fmt.Sprintf("processors/%s.rego", processorName)
	data, err := processorsFs.ReadFile(processorPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
