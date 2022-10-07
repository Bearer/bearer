package settings

import (
	"time"
)

type WorkerSettings struct {
	Count                 int           // number of workers to spawn
	Memory                float64       // memory limit per worker in bytes
	FilesToBatch          int           // how many files to process at once
	ProcessOnlineTimeout  time.Duration // how long to wait for process to become available
	TimeoutSecondPerBytes int           // how many bytes produces second of scan before timing out
	TimeoutMinimum        time.Duration // how many seconds is minimum per file scan
	TimeoutMaximum        time.Duration // how many seconds is maximum per file
	MaximumFileSize       int64         // if we find a file bigger than max file size in bytes ignore it
	CustomDetector        CustomDetector
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
