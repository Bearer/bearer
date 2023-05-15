package settings

import (
	"embed"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/rego"
)

var (
	Workers                   = 1                 // The number of processing workers to spawn
	Timeout                   = 10 * time.Minute  // "The maximum time alloted to complete the scan
	TimeoutFileMinimum        = 5 * time.Second   // Minimum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileMaximum        = 30 * time.Second  // Maximum timeout assigned for scanning each file. This config superseeds timeout-second-per-bytes
	TimeoutFileBytesPerSecond = 1 * 1000          // 1 Kb/s minimum number of bytes per second allowed to scan a file
	TimeoutWorkerOnline       = 60 * time.Second  // Maximum time to wait for a worker process to come online
	FileSizeMaximum           = 2 * 1000 * 1000   // 2 MB Ignore files larger than the specified value
	FilesToBatch              = 1                 // Specify the number of files to batch per worker
	MemoryMaximum             = 800 * 1000 * 1000 // 800 MB If the memory needed to scan a file surpasses the specified limit, skip the file.
	ExistingWorker            = ""                // Specify the URL of an existing worker
)

type WorkerOptions struct {
	Workers                   int           `mapstructure:"workers" json:"workers" yaml:"workers"`
	Timeout                   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	TimeoutFileMinimum        time.Duration `mapstructure:"timeout-file-min" json:"timeout-file-min" yaml:"timeout-file-min"`
	TimeoutFileMaximum        time.Duration `mapstructure:"timeout-file-max"  json:"timeout-file-max" yaml:"timeout-file-max"`
	TimeoutFileBytesPerSecond int           `mapstructure:"timeout-file-bytes-per-second" json:"timeout-file-bytes-per-second" yaml:"timeout-file-bytes-per-second"`
	TimeoutWorkerOnline       time.Duration `mapstructure:"timeout-worker-online" json:"timeout-worker-online" yaml:"timeout-worker-online"`
	FileSizeMaximum           int           `mapstructure:"file-size-max" json:"file-size-max" yaml:"file-size-max"`
	FilesToBatch              int           `mapstructure:"files-to-batch" json:"files-to-batch" yaml:"files-to-batch"`
	MemoryMaximum             int           `mapstructure:"memory-max" json:"memory-max" yaml:"memory-max"`
	ExistingWorker            string        `mapstructure:"existing-worker" json:"existing-worker" yaml:"existing-worker"`
}

type Config struct {
	Client             *api.API
	Worker             WorkerOptions          `mapstructure:"worker" json:"worker" yaml:"worker"`
	Scan               flag.ScanOptions       `mapstructure:"scan" json:"scan" yaml:"scan"`
	Report             flag.ReportOptions     `mapstructure:"report" json:"report" yaml:"report"`
	Policies           map[string]*Policy     `mapstructure:"policies" json:"policies" yaml:"policies"`
	Target             string                 `mapstructure:"target" json:"target" yaml:"target"`
	Rules              map[string]*rules.Rule `mapstructure:"rules" json:"rules" yaml:"rules"`
	BuiltInRules       map[string]*rules.Rule `mapstructure:"built_in_rules" json:"built_in_rules" yaml:"built_in_rules"`
	CacheUsed          bool                   `mapstructure:"cache_used" json:"cache_used" yaml:"cache_used"`
	BearerRulesVersion string                 `mapstructure:"bearer_rules_version" json:"bearer_rules_version" yaml:"bearer_rules_version"`
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

type LoadRulesResult struct {
	BuiltInRules       map[string]*rules.Rule
	Rules              map[string]*rules.Rule
	CacheUsed          bool
	BearerRulesVersion string
}

type Processor struct {
	Query   string  `mapstructure:"query" json:"query" yaml:"query"`
	Modules Modules `mapstructure:"modules" json:"modules" yaml:"modules"`
}

//go:embed policies.yml
var defaultPolicies []byte

//go:embed built_in_rules/*
var buildInRulesFs embed.FS

//go:embed policies/*
var policiesFs embed.FS

//go:embed processors/*
var processorsFs embed.FS

func defaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		Workers:                   Workers,
		Timeout:                   Timeout,
		TimeoutFileMinimum:        TimeoutFileMinimum,
		TimeoutFileMaximum:        TimeoutFileMaximum,
		TimeoutFileBytesPerSecond: TimeoutFileBytesPerSecond,
		TimeoutWorkerOnline:       TimeoutWorkerOnline,
		FilesToBatch:              FilesToBatch,
		FileSizeMaximum:           FileSizeMaximum,
		MemoryMaximum:             MemoryMaximum,
		ExistingWorker:            ExistingWorker,
	}
}

func FromOptions(opts flag.Options, foundLanguages []string) (Config, error) {
	policies := DefaultPolicies()
	workerOptions := defaultWorkerOptions()
	result, err := loadRules(
		opts.ExternalRuleDir,
		opts.RuleOptions,
		foundLanguages,
		opts.ScanOptions.Force,
	)
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
		Client:             opts.Client,
		Worker:             workerOptions,
		Scan:               opts.ScanOptions,
		Report:             opts.ReportOptions,
		Policies:           policies,
		Rules:              result.Rules,
		BuiltInRules:       result.BuiltInRules,
		CacheUsed:          result.CacheUsed,
		BearerRulesVersion: result.BearerRulesVersion,
	}

	return config, nil
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
