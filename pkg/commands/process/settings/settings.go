package settings

import (
	"embed"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/rego"
)

var (
	Workers                   = 1                 // The number of processing workers to spawn
	Timeout                   = 10 * time.Minute  // "The maximum time alloted to complete the scan
	TimeoutFileMinimum        = 5 * time.Second   // Minimum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileMaximum        = 30 * time.Second  // Maximum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileSecondPerBytes = 10 * 1000         // 10kb/s number of file size bytes producing a second of timeout assigned to scanning a file
	TimeoutWorkerOnline       = 60 * time.Second  // Maximum time to wait for a worker process to come online
	FileSizeMaximum           = 2 * 1000 * 1000   // 2MB Ignore files larger than the specified value
	FilesToBatch              = 1                 // Specify the number of files to batch per worker
	MemoryMaximum             = 800 * 1000 * 1000 // 800 MB If the memory needed to scan a file surpasses the specified limit, skip the file.
	ExistingWorker            = ""                // Specify the URL of an existing worker
)

type WorkerOptions struct {
	Workers                   int           `mapstructure:"workers" json:"workers" yaml:"workers"`
	Timeout                   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	TimeoutFileMinimum        time.Duration `mapstructure:"timeout-file-min" json:"timeout-file-min" yaml:"timeout-file-min"`
	TimeoutFileMaximum        time.Duration `mapstructure:"timeout-file-max"  json:"timeout-file-max" yaml:"timeout-file-max"`
	TimeoutFileSecondPerBytes int           `mapstructure:"timeout-file-second-per-bytes" json:"timeout-file-second-per-bytes" yaml:"timeout-file-second-per-bytes"`
	TimeoutWorkerOnline       time.Duration `mapstructure:"timeout-worker-online" json:"timeout-worker-online" yaml:"timeout-worker-online"`
	FileSizeMaximum           int           `mapstructure:"file-size-max" json:"file-size-max" yaml:"file-size-max"`
	FilesToBatch              int           `mapstructure:"files-to-batch" json:"files-to-batch" yaml:"files-to-batch"`
	MemoryMaximum             int           `mapstructure:"memory-max" json:"memory-max" yaml:"memory-max"`
	ExistingWorker            string        `mapstructure:"existing-worker" json:"existing-worker" yaml:"existing-worker"`
}

type Config struct {
	Worker       WorkerOptions      `mapstructure:"worker" json:"worker" yaml:"worker"`
	Scan         flag.ScanOptions   `mapstructure:"scan" json:"scan" yaml:"scan"`
	Report       flag.ReportOptions `mapstructure:"report" json:"report" yaml:"report"`
	Policies     map[string]*Policy `mapstructure:"policies" json:"policies" yaml:"policies"`
	Target       string             `mapstructure:"target" json:"target" yaml:"target"`
	Rules        map[string]*Rule   `mapstructure:"rules" json:"rules" yaml:"rules"`
	BuiltInRules map[string]*Rule   `mapstructure:"built_in_rules" json:"built_in_rules" yaml:"built_in_rules"`
}

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
	AssociatedRecipe   string `mapstructure:"associated_recipe" json:"associated_recipe" yaml:"associated_recipe"`
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
	Id        string        `mapstructure:"id" json:"id" yaml:"id"`
	Type      string        `mapstructure:"type" json:"type" yaml:"type"`
	Languages []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	Patterns  []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

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
	Id                 string            `mapstructure:"id" json:"id,omitempty" yaml:"id,omitempty"`
	AssociatedRecipe   string            `mapstructure:"associated_recipe" json:"associated_recipe" yaml:"associated_recipe"`
	Type               string            `mapstructure:"type" json:"type,omitempty" yaml:"type,omitempty"`          // TODO: use enum value
	Trigger            string            `mapstructure:"trigger" json:"trigger,omitempty" yaml:"trigger,omitempty"` // TODO: use enum value
	DetailedContext    bool              `mapstructure:"detailed_context" json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	Detectors          []string          `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors         []string          `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	Stored             bool              `mapstructure:"stored" json:"stored,omitempty" yaml:"stored,omitempty"`
	AutoEncrytPrefix   string            `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	OmitParentContent  bool              `mapstructure:"omit_parent_content" json:"omit_parent_content,omitempty" yaml:"omit_parent_content,omitempty"`
	SkipDataTypes      []string          `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes      []string          `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	Severity           map[string]string `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	Description        string            `mapstructure:"description" json:"description" yaml:"description"`
	RemediationMessage string            `mapstructure:"remediation_message" json:"remediation_messafe" yaml:"remediation_messafe"`
	DSRID              string            `mapstructure:"dsr_id" json:"dsr_id" yaml:"dsr_id"`
	Languages          []string          `mapstructure:"languages" json:"languages" yaml:"languages"`
	Patterns           []RulePattern     `mapstructure:"patterns" json:"patterns" yaml:"patterns"`

	// FIXME: remove after refactor of sql
	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
	ParamParenting bool               `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
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

//go:embed built_in_rules/*
var builtInRulesFs embed.FS

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

func defaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		Workers:                   Workers,
		Timeout:                   Timeout,
		TimeoutFileMinimum:        TimeoutFileMinimum,
		TimeoutFileMaximum:        TimeoutFileMaximum,
		TimeoutFileSecondPerBytes: TimeoutFileSecondPerBytes,
		TimeoutWorkerOnline:       TimeoutWorkerOnline,
		FilesToBatch:              FilesToBatch,
		FileSizeMaximum:           FileSizeMaximum,
		MemoryMaximum:             MemoryMaximum,
		ExistingWorker:            ExistingWorker,
	}
}

func FromOptions(opts flag.Options) (Config, error) {
	policies := DefaultPolicies()
	workerOptions := defaultWorkerOptions()

	builtInRules, rules, err := loadRules(opts.ExternalRuleDir, opts.RuleOptions)
	if err != nil {
		return Config{}, err
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

	config := Config{
		Worker:       workerOptions,
		Scan:         opts.ScanOptions,
		Report:       opts.ReportOptions,
		Policies:     policies,
		Rules:        rules,
		BuiltInRules: builtInRules,
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
