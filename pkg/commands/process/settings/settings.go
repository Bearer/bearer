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
	Worker   flag.WorkerOptions `mapstructure:"worker" json:"worker" yaml:"worker"`
	Scan     flag.ScanOptions   `mapstructure:"scan" json:"scan" yaml:"scan"`
	Report   flag.ReportOptions `mapstructure:"report" json:"report" yaml:"report"`
	Policies map[string]*Policy `mapstructure:"policies" json:"policies" yaml:"policies"`
	Target   string             `mapstructure:"target" json:"target" yaml:"target"`
	Rules    map[string]*Rule   `mapstructure:"rules" json:"rules" yaml:"rules"`
}

type PolicyLevel string

var LevelCritical = "critical"
var LevelHigh = "high"
var LevelMedium = "medium"
var LevelLow = "low"

type Modules []*PolicyModule

type Policy struct {
	Type    string  `mapstructure:"type" json:"type" yaml:"type"`
	Query   string  `mapstructure:"query" json:"query" yaml:"query"`
	Modules Modules `mapstructure:"modules" json:"modules" yaml:"modules"`
}

type PolicyModule struct {
	Path    string `mapstructure:"path" json:"path,omitempty" yaml:"path,omitempty"`
	Name    string `mapstructure:"name" json:"name" yaml:"name"`
	Content string `mapstructure:"content" json:"content" yaml:"content"`
}

type RuleMetadata struct {
	Description        string `mapstructure:"description" json:"description" yaml:"description"`
	RemediationMessage string `mapstructure:"remediation_message" json:"remediation_messafe" yaml:"remediation_messafe"`
	DSRID              string `mapstructure:"dsr_id" json:"dsr_id" yaml:"dsr_id"`
	ID                 string `mapstructure:"id" json:"id" yaml:"id"`
}

type RuleDefinition struct {
	Disabled          bool              `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Type              string            `mapstructure:"type" json:"type" yaml:"type"`
	Languages         []string          `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting    bool              `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns          []RulePattern     `mapstructure:"patterns" json:"patterns" yaml:"patterns"`
	Stored            bool              `mapstructure:"stored" json:"stored" yaml:"stored"`
	Detectors         []string          `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors        []string          `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	AutoEncrytPrefix  string            `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	DetectPresence    bool              `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	Trigger           string            `mapstructure:"trigger" json:"trigger" yaml:"trigger"` // TODO: use enum value
	Severity          map[string]string `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	SkipDataTypes     []string          `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes     []string          `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	OmitParentContent bool              `mapstructure:"omit_parent_content" json:"omit_parent_content,omitempty" yaml:"omit_parent_content,omitempty"`
	DetailedContext   bool              `mapstructure:"detailed_context" json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	Metadata          RuleMetadata      `mapstructure:"metadata" json:"metadata" yaml:"metadata"`
	Auxiliary         []Auxiliary       `mapstructure:"auxiliary" json:"auxiliary" yaml:"auxiliary"`
}

type Auxiliary struct {
	Id             string        `mapstructure:"id" json:"id" yaml:"id"`
	Type           string        `mapstructure:"type" json:"type" yaml:"type"`
	Languages      []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting bool          `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns       []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	RootSingularize bool `mapstructure:"root_singularize" yaml:"root_singularize" `
	RootLowercase   bool `mapstructure:"root_lowercase" yaml:"root_lowercase"`

	Stored           bool     `mapstructure:"stored" json:"stored" yaml:"stored"`
	Detectors        []string `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors       []string `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	AutoEncrytPrefix string   `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	DetectPresence   bool     `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent       bool     `mapstructure:"omit_parent" json:"omit_parent,omitempty" yaml:"omit_parent,omitempty"`
}

type Rule struct {
	Id                 string            `mapstructure:"id" json:"id,omitempty" yaml:"id,omitempty"`
	Type               string            `mapstructure:"type" json:"type,omitempty" yaml:"type,omitempty"`          // TODO: use enum value
	Trigger            string            `mapstructure:"trigger" json:"trigger,omitempty" yaml:"trigger,omitempty"` // TODO: use enum value
	DetailedContext    bool              `mapstructure:"detailed_context" json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	Detectors          []string          `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors         []string          `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	Stored             bool              `mapstructure:"stored" json:"stored,omitempty" yaml:"stored,omitempty"`
	AutoEncrytPrefix   string            `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	Auxiliary          []Auxiliary       `mapstructure:"auxiliary" json:"auxiliary" yaml:"auxiliary"`
	OmitParentContent  bool              `mapstructure:"omit_parent_content" json:"omit_parent_content,omitempty" yaml:"omit_parent_content,omitempty"`
	SkipDataTypes      []string          `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes      []string          `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	Severity           map[string]string `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	Description        string            `mapstructure:"description" json:"description" yaml:"description"`
	RemediationMessage string            `mapstructure:"remediation_message" json:"remediation_messafe" yaml:"remediation_messafe"`
	DSRID              string            `mapstructure:"dsr_id" json:"dsr_id" yaml:"dsr_id"`
	Disabled           bool              `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Languages          []string          `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting     bool              `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns           []RulePattern     `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	// FIXME: remove after refactor of sql
	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
	DetectPresence bool               `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent     bool               `mapstructure:"omit_parent" json:"omit_parent" yaml:"omit_parent"`
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
}

type RulePattern struct {
	Pattern string          `mapstructure:"pattern" json:"pattern" yaml:"pattern"`
	Filters []PatternFilter `mapstructure:"filters" json:"filters" yaml:"filters"`
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

//go:embed rules/*
var rulesFs embed.FS

//go:embed policies/*
var policiesFs embed.FS

//go:embed processors/*
var processorsFs embed.FS

func (rule *Rule) PolicyType() bool {
	if rule.Type == "data_type" || rule.Type == "verifier" {
		return false
	}

	return true
}

func FromOptions(opts flag.Options) (Config, error) {
	policies := DefaultPolicies()
	rules := defaultRules()

	externalRules, err := LoadExternalRules(opts.ExternalRuleDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load external rules %w", err)
	}

	for ruleName, rule := range externalRules {
		_, ok := rules[ruleName]
		if ok {
			return Config{}, fmt.Errorf("tried to overwrite default rules %s with external rule", ruleName)
		}

		rules[ruleName] = &rule
	}

	for key := range policies {
		policy := policies[key]

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

	// Rule options
	onlyRule := opts.RuleOptions.OnlyRule
	skipRule := opts.RuleOptions.SkipRule

	// validate policy options - raise error if invalid DSW code given
	var invalidRuleIds []string
	for key := range onlyRule {
		if rules[key] == nil {
			invalidRuleIds = append(invalidRuleIds, key)
		}
	}

	for key := range skipRule {
		if rules[key] == nil {
			invalidRuleIds = append(invalidRuleIds, key)
		}
	}

	if len(invalidRuleIds) > 0 {
		return Config{}, fmt.Errorf("unknown rule IDs %s", invalidRuleIds)
	}

	// apply policy options
	for key := range rules {
		rule := rules[key]
		if len(onlyRule) > 0 && !onlyRule[rule.Id] {
			delete(rules, key)
			continue
		}

		if skipRule[rule.Id] {
			delete(rules, key)
			continue
		}
	}

	config := Config{
		Worker:   opts.WorkerOptions,
		Scan:     opts.ScanOptions,
		Report:   opts.ReportOptions,
		Policies: policies,
		Rules:    rules,
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

func DefaultPolicies() map[string]*Policy {
	policies := make(map[string]*Policy)
	var policy []*Policy

	err := yaml.Unmarshal(defaultPolicies, &policy)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal policy file %s", err)
	}

	for _, policy := range policy {
		policies[policy.Type] = policy
	}

	return policies
}

func defaultRules() (rules map[string]*Rule) {
	rules = make(map[string]*Rule)

	// loop through rules langs
	langDirs, err := rulesFs.ReadDir("rules")
	if err != nil {
		log.Fatal().Msgf("failed to read rules dir %e", err)
	}

	for _, langDir := range langDirs {
		lang := langDir.Name()

		if filepath.Ext(langDir.Name()) != "" {
			// not a directory; skip it
			continue
		}

		subLangDirs, err := rulesFs.ReadDir("rules/" + lang)
		if err != nil {
			log.Fatal().Msgf("failed to read rules/%s dir %e", lang, err)
		}

		for _, subLangDir := range subLangDirs {
			subLang := subLangDir.Name()
			dirEntries, err := rulesFs.ReadDir("rules/" + lang + "/" + subLang)
			if err != nil {
				log.Fatal().Msgf("failed to read rules/%s/%s dir %e", lang, subLang, err)
			}

			for _, dirEntry := range dirEntries {
				filename := dirEntry.Name()
				ext := filepath.Ext(filename)

				if ext != ".yaml" && ext != ".yml" {
					continue
				}

				entry, err := rulesFs.ReadFile("rules/" + lang + "/" + subLang + "/" + filename)
				if err != nil {
					log.Fatal().Msgf("failed to read rules/%s/%s/%s file %s", lang, subLang, filename, err)
				}

				var ruleDefinition *RuleDefinition
				err = yaml.Unmarshal(entry, &ruleDefinition)
				if err != nil {
					log.Fatal().Msgf("failed to unmarshal rules/%s/%s/%s %s", lang, subLang, filename, err)
				}

				ruleId := ruleDefinition.Metadata.ID
				if subLang == "internal" {
					// overwrite rule id
					ruleId = strings.TrimSuffix(filename, ext)
				} else {
					// add rule
					rule := Rule{
						Id:                 ruleId,
						Type:               ruleDefinition.Type,
						Trigger:            ruleDefinition.Trigger,
						OmitParentContent:  ruleDefinition.OmitParentContent,
						SkipDataTypes:      ruleDefinition.SkipDataTypes,
						OnlyDataTypes:      ruleDefinition.OnlyDataTypes,
						Severity:           ruleDefinition.Severity,
						Description:        ruleDefinition.Metadata.Description,
						RemediationMessage: ruleDefinition.Metadata.RemediationMessage,
						Stored:             ruleDefinition.Stored,
						Detectors:          ruleDefinition.Detectors,
						Processors:         ruleDefinition.Processors,
						AutoEncrytPrefix:   ruleDefinition.AutoEncrytPrefix,
						DSRID:              ruleDefinition.Metadata.DSRID,
						Disabled:           ruleDefinition.Disabled,
						Languages:          ruleDefinition.Languages,
						ParamParenting:     ruleDefinition.ParamParenting,
						Patterns:           ruleDefinition.Patterns,
						DetectPresence:     ruleDefinition.DetectPresence,
					}

					for _, auxiliaryRuleDefinition := range ruleDefinition.Auxiliary {
						auxiliaryRule := &Rule{
							Type:           auxiliaryRuleDefinition.Type,
							Languages:      auxiliaryRuleDefinition.Languages,
							ParamParenting: auxiliaryRuleDefinition.ParamParenting,
							Patterns:       auxiliaryRuleDefinition.Patterns,
							Stored:         auxiliaryRuleDefinition.Stored,
							DetectPresence: auxiliaryRuleDefinition.DetectPresence,
							OmitParent:     auxiliaryRuleDefinition.OmitParent,
						}

						rules[auxiliaryRuleDefinition.Id] = auxiliaryRule
					}

					rules[ruleId] = &rule
				}
			}
		}
	}

	return rules
}

func ProcessorRegoModuleText(processorName string) (string, error) {
	processorPath := fmt.Sprintf("processors/%s.rego", processorName)
	data, err := processorsFs.ReadFile(processorPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
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
