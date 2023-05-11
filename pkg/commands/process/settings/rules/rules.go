package rules

import (
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type MatchOn string

const (
	PRESENCE          MatchOn = "presence"
	ABSENCE           MatchOn = "absence"
	STORED_DATA_TYPES MatchOn = "stored_data_types"
)

type MetaVar struct {
	Input  string `mapstructure:"input" json:"input" yaml:"input"`
	Output int    `mapstructure:"output" json:"output" yaml:"output"`
	Regex  string `mapstructure:"regex" json:"regex" yaml:"regex"`
}

type RuleTrigger struct {
	MatchOn           MatchOn `mapstructure:"match_on" json:"match_on" yaml:"match_on"`
	DataTypesRequired bool    `mapstructure:"data_types_required" json:"data_types_required" yaml:"data_types_required"`
	RequiredDetection *string `mapstructure:"required_detection" json:"required_detection" yaml:"required_detection"`
}

type RuleDefinitionTrigger struct {
	MatchOn           *MatchOn `mapstructure:"match_on" json:"match_on" yaml:"match_on"`
	RequiredDetection *string  `mapstructure:"required_detection" json:"required_detection" yaml:"required_detection"`
	DataTypesRequired *bool    `mapstructure:"data_types_required" json:"data_types_required" yaml:"data_types_required"`
}

type RuleMetadata struct {
	Description        string   `mapstructure:"description" json:"description" yaml:"description"`
	RemediationMessage string   `mapstructure:"remediation_message" json:"remediation_message" yaml:"remediation_message"`
	CWEIDs             []string `mapstructure:"cwe_id" json:"cwe_id" yaml:"cwe_id"`
	AssociatedRecipe   string   `mapstructure:"associated_recipe" json:"associated_recipe" yaml:"associated_recipe"`
	ID                 string   `mapstructure:"id" json:"id" yaml:"id"`
	DocumentationUrl   string   `mapstructure:"documentation_url" json:"documentation_url" yaml:"documentation_url"`
}

type RuleDefinition struct {
	Disabled           bool                    `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Type               string                  `mapstructure:"type" json:"type" yaml:"type"`
	Languages          []string                `mapstructure:"languages" json:"languages" yaml:"languages"`
	ParamParenting     bool                    `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns           []RuleDefinitionPattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`
	Stored             bool                    `mapstructure:"stored" json:"stored" yaml:"stored"`
	Detectors          []string                `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors         []string                `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	AutoEncrytPrefix   string                  `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	DetectPresence     bool                    `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	Trigger            *RuleDefinitionTrigger  `mapstructure:"trigger" json:"trigger" yaml:"trigger"` // TODO: use enum value
	Severity           string                  `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	SkipDataTypes      []string                `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes      []string                `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	HasDetailedContext bool                    `mapstructure:"has_detailed_context" json:"has_detailed_context,omitempty" yaml:"has_detailed_context,omitempty"`
	Metadata           *RuleMetadata           `mapstructure:"metadata" json:"metadata" yaml:"metadata"`
	Auxiliary          []Auxiliary             `mapstructure:"auxiliary" json:"auxiliary" yaml:"auxiliary"`
}

type Auxiliary struct {
	Id        string                  `mapstructure:"id" json:"id" yaml:"id"`
	Type      string                  `mapstructure:"type" json:"type" yaml:"type"`
	Languages []string                `mapstructure:"languages" json:"languages" yaml:"languages"`
	Patterns  []RuleDefinitionPattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	RootSingularize bool `mapstructure:"root_singularize" yaml:"root_singularize" `
	RootLowercase   bool `mapstructure:"root_lowercase" yaml:"root_lowercase"`

	Stored           bool     `mapstructure:"stored" json:"stored" yaml:"stored"`
	Detectors        []string `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors       []string `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	AutoEncrytPrefix string   `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`

	// FIXME: remove after refactor of sql
	ParamParenting bool `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	DetectPresence bool `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent     bool `mapstructure:"omit_parent" json:"omit_parent,omitempty" yaml:"omit_parent,omitempty"`
}

type Rule struct {
	Id                 string        `mapstructure:"id" json:"id,omitempty" yaml:"id,omitempty"`
	AssociatedRecipe   string        `mapstructure:"associated_recipe" json:"associated_recipe" yaml:"associated_recipe"`
	Type               string        `mapstructure:"type" json:"type,omitempty" yaml:"type,omitempty"` // TODO: use enum value
	Trigger            RuleTrigger   `mapstructure:"trigger" json:"trigger,omitempty" yaml:"trigger,omitempty"`
	IsLocal            bool          `mapstructure:"is_local" json:"is_local,omitempty" yaml:"is_local,omitempty"`
	Detectors          []string      `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors         []string      `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	Stored             bool          `mapstructure:"stored" json:"stored,omitempty" yaml:"stored,omitempty"`
	AutoEncrytPrefix   string        `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	HasDetailedContext bool          `mapstructure:"has_detailed_context" json:"has_detailed_context,omitempty" yaml:"has_detailed_context,omitempty"`
	SkipDataTypes      []string      `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes      []string      `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	Severity           string        `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	Description        string        `mapstructure:"description" json:"description" yaml:"description"`
	RemediationMessage string        `mapstructure:"remediation_message" json:"remediation_messafe" yaml:"remediation_messafe"`
	CWEIDs             []string      `mapstructure:"cwe_ids" json:"cwe_ids" yaml:"cwe_ids"`
	Languages          []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	Patterns           []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`
	DocumentationUrl   string        `mapstructure:"documentation_url" json:"documentation_url" yaml:"documentation_url"`
	IsAuxilary         bool          `mapstructure:"is_auxilary" json:"is_auxilary" yaml:"is_auxilary"`

	// FIXME: remove after refactor of sql
	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
	ParamParenting bool               `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	DetectPresence bool               `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent     bool               `mapstructure:"omit_parent" json:"omit_parent" yaml:"omit_parent"`
}

type PatternFilter struct {
	Not                *PatternFilter  `mapstructure:"not" json:"not" yaml:"not"`
	Either             []PatternFilter `mapstructure:"either" json:"either" yaml:"either"`
	Variable           string          `mapstructure:"variable" json:"variable" yaml:"variable"`
	Detection          string          `mapstructure:"detection" json:"detection" yaml:"detection"`
	Contains           *bool           `mapstructure:"contains" json:"contains" yaml:"contains"`
	Regex              *Regexp         `mapstructure:"regex" json:"regex" yaml:"regex"`
	Values             []string        `mapstructure:"values" json:"values" yaml:"values"`
	LengthLessThan     *int            `mapstructure:"length_less_than" json:"length_less_than" yaml:"length_less_than"`
	LessThan           *int            `mapstructure:"less_than" json:"less_than" yaml:"less_than"`
	LessThanOrEqual    *int            `mapstructure:"less_than_or_equal" json:"less_than_or_equal" yaml:"less_than_or_equal"`
	GreaterThan        *int            `mapstructure:"greater_than" json:"greater_than" yaml:"greater_than"`
	GreaterThanOrEqual *int            `mapstructure:"greater_than_or_equal" json:"greater_than_or_equal" yaml:"greater_than_or_equal"`
	StringRegex        *Regexp         `mapstructure:"string_regex" json:"string_regex" yaml:"string_regex"`
	FilenameRegex      *Regexp         `mapstructure:"filename_regex" json:"filename_regex" yaml:"filename_regex"`
}

type RuleDefinitionPattern struct {
	Pattern string          `mapstructure:"pattern" json:"pattern" yaml:"pattern"`
	Filters []PatternFilter `mapstructure:"filters" json:"filters" yaml:"filters"`
}

type RulePattern struct {
	Pattern string
	Query   languagetypes.PatternQuery
	Filters []PatternFilter
}

func (rulePattern *RuleDefinitionPattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to parse as a string
	var pattern string
	if err := unmarshal(&pattern); err == nil {
		rulePattern.Pattern = pattern
		return nil
	}

	// Wasn't a string so it must be the structured format
	type rawRulePattern RuleDefinitionPattern
	return unmarshal((*rawRulePattern)(rulePattern))
}

func (rule *Rule) PolicyType() bool {
	return rule.Type == "risk"
}

func (rule *Rule) Language() string {
	if rule.Languages == nil {
		return "secret"
	}

	switch rule.Languages[0] {
	case "javascript":
		return "JavaScript"
	case "ruby":
		return "Ruby"
	case "sql":
		return "SQL"
	default:
		return rule.Languages[0]
	}
}
