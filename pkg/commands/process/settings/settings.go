package settings

import (
	_ "embed"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/flag"
)

type Config struct {
	Worker         flag.WorkerOptions `json:"worker"`
	Scan           flag.ScanOptions   `json:"scan"`
	CustomDetector CustomDetector     `json:"custom_detector"`
}

type CustomDetector struct {
	RulesConfig *RulesConfig `json:"rules_config"`
}

type RulesConfig struct {
	Rules map[string]Rule `json:"rules"`
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

//go:embed custom_detector.yml
var customDetector []byte

func DefaultCustomDetector() CustomDetector {
	var rules RulesConfig

	err := yaml.Unmarshal(customDetector, &rules)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal database file %e", err)
	}

	return CustomDetector{
		RulesConfig: &rules,
	}
}
