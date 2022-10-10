package settings

import (
	_ "embed"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

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
