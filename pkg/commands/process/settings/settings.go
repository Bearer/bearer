package settings

import (
	_ "embed"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/flag"
)

type Config struct {
	Worker         flag.WorkerOptions `json:"worker"`
	Scan           flag.ScanOptions   `json:"scan"`
	Report         flag.ReportOptions `json:"report"`
	CustomDetector map[string]Rule    `json:"custom_detector"`
}

type Rule struct {
	Disabled       bool
	Languages      []string
	Patterns       []string
	ParamParenting bool `yaml:"param_parenting"`
	Metavars       map[string]MetaVar
	Stored         bool
}

type MetaVar struct {
	Input  string
	Output int
	Regex  string
}

//go:embed custom_detector.yml
var customDetector []byte

var CustomDetectorKey string = "scan.custom_detector"

func FromOptions(opts flag.Options) (Config, error) {
	rules := DefaultCustomDetector()
	if viper.IsSet(CustomDetectorKey) {
		err := viper.UnmarshalKey(CustomDetectorKey, &rules)
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		Worker:         opts.WorkerOptions,
		CustomDetector: rules,
		Scan:           opts.ScanOptions,
		Report:         opts.ReportOptions,
	}, nil
}

func DefaultCustomDetector() map[string]Rule {
	var rules map[string]Rule

	err := yaml.Unmarshal(customDetector, &rules)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal database file %e", err)
	}

	return rules
}
