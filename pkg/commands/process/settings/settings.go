package settings

import (
	"time"

	"github.com/bearer/bearer/api"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	ignoretypes "github.com/bearer/bearer/pkg/util/ignore/types"
	"github.com/bearer/bearer/pkg/util/regex"
	"github.com/bearer/bearer/pkg/util/rego"

	globaltypes "github.com/bearer/bearer/pkg/types"
)

var (
	Timeout                          = 10 * time.Minute  // "The maximum time alloted to complete the scan
	TimeoutFileMinimum               = 5 * time.Second   // Minimum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileMaximum               = 30 * time.Second  // Maximum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileBytesPerSecond        = 1 * 1000          // 1 Kb/s minimum number of bytes per second allowed to scan a file
	TimeoutWorkerFileGrace           = 5 * time.Second   // Grace period to allow a worker to timeout on it's own
	TimeoutWorkerOnline              = 60 * time.Second  // Maximum time to wait for a worker process to come online
	TimeoutWorkerShutdown            = 5 * time.Second   // Maximum time to wait for a worker process to shut down cleanly
	CodeExtractBuffer                = 3                 // Number of lines allowed before or after the detection
	FileSizeMaximum                  = 2 * 1000 * 1000   // 2 MB Ignore files larger than the specified value
	FilesPerWorker                   = 1000              // By default, start a worker per this many files, up to the number of CPUs
	MemorySoftMaximum         uint64 = 650 * 1000 * 1000 // 650 MB If the memory needed to scan a file surpasses the specified limit, ask the worker to reduce memory usage.
	MemoryMaximum             uint64 = 800 * 1000 * 1000 // 800 MB If the memory needed to scan a file surpasses the specified limit, skip the file.
	ExistingWorker                   = ""                // Specify the URL of an existing worker
)

type WorkerOptions struct {
	Timeout                   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	TimeoutFileMinimum        time.Duration `mapstructure:"timeout-file-min" json:"timeout-file-min" yaml:"timeout-file-min"`
	TimeoutFileMaximum        time.Duration `mapstructure:"timeout-file-max"  json:"timeout-file-max" yaml:"timeout-file-max"`
	TimeoutFileBytesPerSecond int           `mapstructure:"timeout-file-bytes-per-second" json:"timeout-file-bytes-per-second" yaml:"timeout-file-bytes-per-second"`
	TimeoutWorkerOnline       time.Duration `mapstructure:"timeout-worker-online" json:"timeout-worker-online" yaml:"timeout-worker-online"`
	FileSizeMaximum           int           `mapstructure:"file-size-max" json:"file-size-max" yaml:"file-size-max"`
	ExistingWorker            string        `mapstructure:"existing-worker" json:"existing-worker" yaml:"existing-worker"`
}

type Config struct {
	Client                     *api.API
	Worker                     WorkerOptions                             `mapstructure:"worker" json:"worker" yaml:"worker"`
	Scan                       flagtypes.ScanOptions                     `mapstructure:"scan" json:"scan" yaml:"scan"`
	Report                     flagtypes.ReportOptions                   `mapstructure:"report" json:"report" yaml:"report"`
	IgnoredFingerprints        map[string]ignoretypes.IgnoredFingerprint `mapstructure:"ignored_fingerprints" json:"ignored_fingerprints" yaml:"ignored_fingerprints"`
	StaleIgnoredFingerprintIds []string                                  `mapstructure:"stale_ignored_fingerprint_ids" json:"stale_ignored_fingerprint_ids" yaml:"stale_ignored_fingerprint_ids"`
	CloudIgnoresUsed           bool                                      `mapstructure:"cloud_ignores_used" json:"cloud_ignores_used" yaml:"cloud_ignores_used"`
	Policies                   map[string]*Policy                        `mapstructure:"policies" json:"policies" yaml:"policies"`
	Target                     string                                    `mapstructure:"target" json:"target" yaml:"target"`
	IgnoreFile                 string                                    `mapstructure:"ignore_file" json:"ignore_file" yaml:"ignore_file"`
	Rules                      map[string]*Rule                          `mapstructure:"rules" json:"rules" yaml:"rules"`
	LoadedRuleCount            int
	BuiltInRules               map[string]*Rule `mapstructure:"built_in_rules" json:"built_in_rules" yaml:"built_in_rules"`
	CacheUsed                  bool             `mapstructure:"cache_used" json:"cache_used" yaml:"cache_used"`
	BearerRulesVersion         string           `mapstructure:"bearer_rules_version" json:"bearer_rules_version" yaml:"bearer_rules_version"`
	NoColor                    bool             `mapstructure:"no_color" json:"no_color" yaml:"no_color"`
	Debug                      bool             `mapstructure:"debug" json:"debug" yaml:"debug"`
	LogLevel                   string           `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	DebugProfile               bool             `mapstructure:"debug_profile" json:"debug_profile" yaml:"debug_profile"`
	IgnoreGit                  bool             `mapstructure:"ignore_git" json:"ignore_git" yaml:"ignore_git"`
}

type Processor struct {
	Query   string  `mapstructure:"query" json:"query" yaml:"query"`
	Modules Modules `mapstructure:"modules" json:"modules" yaml:"modules"`
}

type Modules []*PolicyModule

func (modules Modules) ToRegoModules() (output []rego.Module) {
	for _, module := range modules {
		output = append(output, rego.Module{
			Name:    module.Name,
			Content: module.Content,
		})
	}
	return
}

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

type MatchOn string

const (
	PRESENCE          MatchOn = "presence"
	ABSENCE           MatchOn = "absence"
	STORED_DATA_TYPES MatchOn = "stored_data_types"
)

type RuleReferenceScope string

const (
	CURSOR_STRICT_SCOPE RuleReferenceScope = "cursor_strict"
	CURSOR_SCOPE        RuleReferenceScope = "cursor"
	NESTED_SCOPE        RuleReferenceScope = "nested"
	NESTED_STRICT_SCOPE RuleReferenceScope = "nested_strict"
	RESULT_SCOPE        RuleReferenceScope = "result"

	DefaultScope = NESTED_SCOPE
)

type RuleTrigger struct {
	MatchOn            MatchOn  `mapstructure:"match_on" json:"match_on" yaml:"match_on"`
	DataTypesRequired  bool     `mapstructure:"data_types_required" json:"data_types_required" yaml:"data_types_required"`
	RequiredDetections []string `mapstructure:"required_detections" json:"required_detections" yaml:"required_detections"`
}

type RuleDefinitionTrigger struct {
	MatchOn            *MatchOn `mapstructure:"match_on" json:"match_on" yaml:"match_on"`
	RequiredDetection  *string  `mapstructure:"required_detection" json:"required_detection" yaml:"required_detection"`
	RequiredDetections []string `mapstructure:"required_detections" json:"required_detections" yaml:"required_detections"`
	DataTypesRequired  *bool    `mapstructure:"data_types_required" json:"data_types_required" yaml:"data_types_required"`
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
	Disabled           bool                   `mapstructure:"disabled" json:"disabled" yaml:"disabled"`
	Type               string                 `mapstructure:"type" json:"type" yaml:"type"`
	Languages          []string               `mapstructure:"languages" json:"languages" yaml:"languages"`
	Imports            []string               `mapstructure:"imports" json:"imports" yaml:"imports"`
	ParamParenting     bool                   `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	Patterns           []RulePattern          `mapstructure:"patterns" json:"patterns" yaml:"patterns"`
	SanitizerRuleID    string                 `mapstructure:"sanitizer" json:"sanitizer,omitempty" yaml:"sanitizer,omitempty"`
	Stored             bool                   `mapstructure:"stored" json:"stored" yaml:"stored"`
	Detectors          []string               `mapstructure:"detectors" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Processors         []string               `mapstructure:"processors" json:"processors,omitempty" yaml:"processors,omitempty"`
	AutoEncrytPrefix   string                 `mapstructure:"auto_encrypt_prefix" json:"auto_encrypt_prefix,omitempty" yaml:"auto_encrypt_prefix,omitempty"`
	DetectPresence     bool                   `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	Trigger            *RuleDefinitionTrigger `mapstructure:"trigger" json:"trigger" yaml:"trigger"` // TODO: use enum value
	Severity           string                 `mapstructure:"severity" json:"severity,omitempty" yaml:"severity,omitempty"`
	SkipDataTypes      []string               `mapstructure:"skip_data_types" json:"skip_data_types,omitempty" yaml:"skip_data_types,omitempty"`
	OnlyDataTypes      []string               `mapstructure:"only_data_types" json:"only_data_types,omitempty" yaml:"only_data_types,omitempty"`
	HasDetailedContext bool                   `mapstructure:"has_detailed_context" json:"has_detailed_context,omitempty" yaml:"has_detailed_context,omitempty"`
	Metadata           *RuleMetadata          `mapstructure:"metadata" json:"metadata" yaml:"metadata"`
	Auxiliary          []Auxiliary            `mapstructure:"auxiliary" json:"auxiliary" yaml:"auxiliary"`
	DependencyCheck    bool                   `mapstructure:"dependency_check" json:"dependency_check" yaml:"dependency_check"`
	Dependency         *Dependency            `mapstructure:"dependency" json:"dependency" yaml:"dependency"`
	Text               string                 `mapstructure:"-" json:"-" yaml:"-"`
}

type Dependency struct {
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	Name       string `mapstructure:"name" json:"name" yaml:"name"`
	MinVersion string `mapstructure:"min_version" json:"min_version" yaml:"min_version"`
}

type Auxiliary struct {
	Id              string        `mapstructure:"id" json:"id" yaml:"id"`
	Type            string        `mapstructure:"type" json:"type" yaml:"type"`
	Languages       []string      `mapstructure:"languages" json:"languages" yaml:"languages"`
	Patterns        []RulePattern `mapstructure:"patterns" json:"patterns" yaml:"patterns"`
	SanitizerRuleID string        `mapstructure:"sanitizer" json:"sanitizer,omitempty" yaml:"sanitizer,omitempty"`

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
	SanitizerRuleID    string        `mapstructure:"sanitizer" json:"sanitizer" yaml:"sanitizer"`
	DocumentationUrl   string        `mapstructure:"documentation_url" json:"documentation_url" yaml:"documentation_url"`
	IsAuxilary         bool          `mapstructure:"is_auxilary" json:"is_auxilary" yaml:"is_auxilary"`
	DependencyCheck    bool          `mapstructure:"dependency_check" json:"dependency_check" yaml:"dependency_check"`
	Dependency         *Dependency   `mapstructure:"dependency" json:"dependency" yaml:"dependency"`

	// FIXME: remove after refactor of sql
	Metavars       map[string]MetaVar `mapstructure:"metavars" json:"metavars" yaml:"metavars"`
	ParamParenting bool               `mapstructure:"param_parenting" json:"param_parenting" yaml:"param_parenting"`
	DetectPresence bool               `mapstructure:"detect_presence" json:"detect_presence" yaml:"detect_presence"`
	OmitParent     bool               `mapstructure:"omit_parent" json:"omit_parent" yaml:"omit_parent"`
}

type RuleReferenceImport struct {
	Variable string `mapstructure:"variable" json:"variable" yaml:"variable"`
	As       string `mapstructure:"as" json:"as" yaml:"as"`
}

type PatternFilter struct {
	Not        *PatternFilter        `mapstructure:"not" json:"not,omitempty" yaml:"not,omitempty"`
	Either     []PatternFilter       `mapstructure:"either" json:"either,omitempty" yaml:"either,omitempty"`
	Variable   string                `mapstructure:"variable" json:"variable,omitempty" yaml:"variable,omitempty"`
	Type       string                `mapstructure:"type" json:"type,omitempty" yaml:"type,omitempty"`
	StaticType string                `mapstructure:"static_type" json:"static_type,omitempty" yaml:"static_type,omitempty"`
	Detection  string                `mapstructure:"detection" json:"detection,omitempty" yaml:"detection,omitempty"`
	Scope      RuleReferenceScope    `mapstructure:"scope" json:"scope,omitempty" yaml:"scope,omitempty"`
	IsSource   bool                  `mapstructure:"is_source" json:"is_source" yaml:"is_source"`
	Filters    []PatternFilter       `mapstructure:"filters" json:"filters,omitempty" yaml:"filters,omitempty"`
	Imports    []RuleReferenceImport `mapstructure:"imports" json:"imports,omitempty" yaml:"imports,omitempty"`
	// Contains is deprecated in favour of Scope
	Contains           *bool                     `mapstructure:"contains" json:"contains,omitempty" yaml:"contains,omitempty"`
	Regex              *regex.SerializableRegexp `mapstructure:"regex" json:"regex,omitempty" yaml:"regex,omitempty"`
	Values             []string                  `mapstructure:"values" json:"values,omitempty" yaml:"values,omitempty"`
	LengthLessThan     *int                      `mapstructure:"length_less_than" json:"length_less_than,omitempty" yaml:"length_less_than,omitempty"`
	LessThan           *int                      `mapstructure:"less_than" json:"less_than,omitempty" yaml:"less_than,omitempty"`
	LessThanOrEqual    *int                      `mapstructure:"less_than_or_equal" json:"less_than_or_equal,omitempty" yaml:"less_than_or_equal,omitempty"`
	GreaterThan        *int                      `mapstructure:"greater_than" json:"greater_than,omitempty" yaml:"greater_than,omitempty"`
	GreaterThanOrEqual *int                      `mapstructure:"greater_than_or_equal" json:"greater_than_or_equal,omitempty" yaml:"greater_than_or_equal,omitempty"`
	StringRegex        *regex.SerializableRegexp `mapstructure:"string_regex" json:"string_regex,omitempty" yaml:"string_regex,omitempty"`
	EntropyGreaterThan *float64                  `mapstructure:"entropy_greater_than" json:"entropy_greater_than,omitempty" yaml:"entropy_greater_than,omitempty"`
	FilenameRegex      *regex.SerializableRegexp `mapstructure:"filename_regex" json:"filename_regex,omitempty" yaml:"filename_regex,omitempty"`
}

type RulePattern struct {
	Pattern string          `mapstructure:"pattern" json:"pattern" yaml:"pattern"`
	Focus   string          `mapstructure:"focus" json:"focus,omitempty" yaml:"focus,omitempty"`
	Filters []PatternFilter `mapstructure:"filters" json:"filters" yaml:"filters"`
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

func (filter *PatternFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type wrapper PatternFilter
	var wrapped wrapper
	if err := unmarshal(&wrapped); err != nil {
		return err
	}

	*filter = PatternFilter(wrapped)

	// Default Scope to "contains" and maintain backwards compatibility with rules
	// using the `contains` flag
	if filter.Detection != "" {
		if filter.Contains != nil {
			if !*filter.Contains {
				filter.Scope = CURSOR_SCOPE
			}
		}
		if filter.Scope == "" {
			filter.Scope = NESTED_SCOPE
		}
	}

	return nil
}

type MetaVar struct {
	Input  string `mapstructure:"input" json:"input" yaml:"input"`
	Output int    `mapstructure:"output" json:"output" yaml:"output"`
	Regex  string `mapstructure:"regex" json:"regex" yaml:"regex"`
}

func (rule *Rule) PolicyType() bool {
	return rule.Type == "risk"
}

func (rule *Rule) GetSeverity() string {
	if rule.Severity == "" {
		return globaltypes.LevelLow
	}

	return rule.Severity
}

func (rule *Rule) IsSecrets() bool {
	return rule.Languages == nil
}
