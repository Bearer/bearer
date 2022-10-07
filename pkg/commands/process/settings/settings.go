package settings

import (
	"github.com/bearer/curio/pkg/flag"
)

type Config struct {
	Worker         flag.WorkerOptions
	CustomDetector CustomDetector
}

type CustomDetector struct {
	RulesConfig *RulesConfig
}

type RulesConfig struct {
	Rules map[string]Rule
}

type Rule struct {
	Disabled       bool
	Languages      []string
	Patterns       []string
	ParamParenting bool `yaml:"param_parenting"`
	Metavars       map[string]MetaVar
}

type MetaVar struct {
	Input  string
	Output int
	Regex  string
}
