package dependencies

import "github.com/bearer/bearer/internal/report/detectors"

type Provider string

const DetectorGemFileLock detectors.Type = "gemfile-lock"
const LanguageRuby string = "ruby"

// Dependency is a dependency that keeps the name and version
type Dependency struct {
	PackageManager string `json:"package_manager" yaml:"package_manager"`
	Group          string `json:"group" yaml:"group"`
	Name           string `json:"name" yaml:"name"`
	Version        string `json:"version" yaml:"version"`
}
